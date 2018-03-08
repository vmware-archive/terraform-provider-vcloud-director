#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************
from proto import vapp_vm_pb2_grpc as vapp_vm_pb2_grpc
from proto import vapp_vm_pb2 as vapp_vm_pb2

from pyvcloud.vcd.vdc import VDC
from pyvcloud.vcd.system import System
from pyvcloud.vcd.org import Org
from pyvcloud.vcd.client import TaskStatus
from pyvcloud.vcd.vapp import VApp


import logging
import grpc

from vcd_client_ref import VCDClientRef
import errors
from lxml import objectify, etree
import lxml


class VappVmServicer(vapp_vm_pb2_grpc.VappVmServicer):
    def __init__(self, pyPluginServer):
        self.py_plugin_server = pyPluginServer
        vref = VCDClientRef()
        self.client = vref.get_ref()

    def Create(self, request, context):
        logging.basicConfig(level=logging.DEBUG)
        logging.info("__INIT__Create[VappVmServicer]")
        
        target_vm_name = request.target_vm_name
        target_vapp = request.target_vapp
        target_vdc = request.target_vdc
        source_vapp = request.source_vapp
        source_vm_name = request.source_vm_name
        hostname = request.hostname
        password = request.password
        password_auto = request.password_auto
        password_reset = request.password_reset
        cust_script = request.cust_script
        network = request.network
        storage_profile = request.storage_profile
        power_on = request.power_on
        all_eulas_accepted  = request.all_eulas_accepted
        
        source_catalog_name = request.source_catalog_name
        source_template_name = request.source_template_name

        ip_allocation_mode=request.ip_allocation_mode

        vapp_vm = VappVm(target_vm_name, power_on, context)
        res = vapp_vm.create(
                target_vapp,
                target_vdc,
                source_vapp,
                source_vm_name,
                source_catalog_name,
                source_template_name,
                hostname,
                password,
                password_auto,
                password_reset,
                cust_script,
                network,
                storage_profile,
                power_on,
                all_eulas_accepted,
                ip_allocation_mode)

        logging.info("__DONE__Create[VappVmServicer]")
        return res

    def Delete(self, request, context):
        logging.info("__DONE__Delete[VappVmServicer]")

        name = request.name
        vapp_vm = VappVm(name=name, context=context)
        res = vapp_vm.delete()
        logging.info("__DONE__Delete[VappVmServicer]")
        return res

    def Update(self, request, context):
        logging.info("__DONE__Update[VappVmServicer]")

        name = request.name
        is_enabled = request.is_enabled
        vapp_vm = VappVm(name=name, is_enabled=is_enabled, context=context)
        res = vapp_vm.update()
        logging.info("__DONE__Update[VappVmServicer]")
        return res

    def Read(self, request, context):
        logging.info("__INIT__Read[VappVmServicer]")
        name = request.name
        vapp_vm = VappVm(name=name, context=context)
        res = vapp_vm.read()
        logging.info("__DONE__Read[VappVmServicer]")
        return res


class VappVm:
    def __repr__(self):
        message = 'VappVm [name ={0} '.format(self.name)
        return message

    def __init__(self,
                 target_vm_name,
                 
                 power_on,
                 context=None):
        vref = VCDClientRef()
        self.client = vref.get_ref()
        self.target_vm_name = target_vm_name
        
        self.context=context


    def create(self,
            target_vapp,
            target_vdc,
            source_vapp,
            source_vm_name,
            source_catalog_name,
            source_template_name,
            hostname,
            password,
            password_auto,
            password_reset,
            cust_script,
            network,
            storage_profile,
            power_on,
            all_eulas_accepted,
            ip_allocation_mode):
        logging.info("__INIT__create[VappVm] source_catalog_name[%s]" , source_catalog_name)
        res = vapp_vm_pb2.CreateVappVmResult()
        res.created = False

        context = self.context
        target_vm_name = self.target_vm_name

        logged_in_org = self.client.get_org()
        org = Org(self.client, resource=logged_in_org)

        try:
            vdc_resource = org.get_vdc(target_vdc)
            vdc = VDC(self.client, name=target_vdc, resource=vdc_resource)
            
            vapp_resource = vdc.get_vapp(target_vapp)
            vapp = VApp(self.client, name=target_vapp, resource=vapp_resource)
            
            catalog_item = org.get_catalog_item(source_catalog_name, source_template_name)
            source_vapp_resource = self.client.get_resource(catalog_item.Entity.get('href'))
            spec = {
                        'source_vm_name': source_vm_name,
                        'vapp': source_vapp_resource,
                        'target_vm_name' : target_vm_name,
                        'hostname' : hostname,
                        'network' : network,
                        'ip_allocation_mode' : ip_allocation_mode
                    }
            # spec['storage_profile'] = storage_profile
            
            vms = [spec]
            create_vapp_vm_resp = vapp.add_vms(vms)
            task = client.get_task_monitor().wait_for_status(
                            task=create_vapp_vm_resp,
                            timeout=60,
                            poll_frequency=2,
                            fail_on_statuses=None,
                            expected_target_statuses=[
                                TaskStatus.SUCCESS,
                                TaskStatus.ABORTED,
                                TaskStatus.ERROR,
                                TaskStatus.CANCELED],
                            callback=None)
        
            st = task.get('status')
            if st == TaskStatus.SUCCESS.value:
                message = 'status : {0} '.format(st)
                logging.info(message)
                res.created = True
            else:
                raise errors.VappVmCreateError(
                    etree.tostring(task, pretty_print=True))
        except Exception as e:
            error_message = '__ERROR_create[VappVm] failed for vm {0}. __ErrorMessage__ {1}'.format(
                self.name, str(e))
            logging.warn(error_message)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
            return res

        logging.info("__DONE__create[VappVm]")
        return res

    def update(self):
        logging.info("__INIT__update[VappVm]")
        res = user_pb2.UpdateUserResult()
        res.updated = False

        context = self.context

        logged_in_org = self.client.get_org()
        org = Org(self.client, resource=logged_in_org)
        name = self.name
        is_enabled = self.is_enabled

        try:
            result = org.update_user(name, is_enabled)
            res.updated = True
        except Exception as e:
            error_message = '__ERROR_update[VappVm] failed for VappVm {0}. __ErrorMessage__ {1}'.format(
                self.name, str(e))
            logging.warn(error_message)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
            return res

        logging.info("__DONE__update[VappVm]")
        return res

    def read(self):
        logging.info("__INIT__read[VappVm]")
        res = user_pb2.ReadUserResult()
        res.present = False
        context = self.context

        logged_in_org = self.client.get_org()
        org = Org(self.client, resource=logged_in_org)
        try:
            result = org.get_user(self.name)

            res.present = True
            res.is_enabled = result.IsEnabled
            res.name = str(result.NameInSource)
            res.role_name = str(result.Role.get('name'))

            res.full_name = str(result.FullName)

            #description not available
            #res.description = result.DESCRIPTION

            res.email = str(result.EmailAddress)
            res.telephone = str(result.Telephone)
            res.im = str(result.IM)

            #alert_email not available
            #res.alert_email = result.Role.ALERT_EMAIL

            #alert_email_prefix not available
            #res.alert_email_prefix = result.ALERT_EMAIL_PREFIX

            res.stored_vm_quota = result.StoredVmQuota
            res.deployed_vm_quota = result.DeployedVmQuota
            res.is_group_role = result.IsGroupRole

            #is_default_cached not available
            #res.is_default_cached = result.Role.get('name')

            res.is_external = result.IsExternal

            #is_alert_enabled not available
            #res.is_alert_enabled = result.Role.get('name')

            res.is_enabled = result.IsEnabled
        except Exception as e:
            error_message = '__ERROR_read[VappVm] failed for VappVm {0}. __ErrorMessage__ {1}'.format(
                self.name, str(e))
            logging.warn(error_message)
            #context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            #context.set_details(error_message)
            return res
        logging.info("__DONE__read[VappVm]")
        return res

    def delete(self):
        logging.info("__INIT__delete[VappVm]")
        res = user_pb2.DeleteUserResult()
        res.deleted = False

        context = self.context

        logged_in_org = self.client.get_org()
        org = Org(self.client, resource=logged_in_org)
        try:
            result = org.delete_user(self.name)
            res.deleted = True
        except Exception as e:
            res.deleted = False
            error_message = '__ERROR_delete[VappVm] failed for VappVm {0}. __ErrorMessage__ {1}'.format(
                self.name, str(e))
            logging.warn(error_message)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
            return res

        logging.info("__DONE__delete[VappVm]")
        return res
