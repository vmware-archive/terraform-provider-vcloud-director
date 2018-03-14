#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

import grpc
import errors
import logging
from lxml import etree
from pyvcloud.vcd.vm import VM
from pyvcloud.vcd.org import Org
from pyvcloud.vcd.vdc import VDC
from pyvcloud.vcd.vapp import VApp
from vcd_client_ref import VCDClientRef
from pyvcloud.vcd.client import TaskStatus
from pyvcloud.vcd.client import EntityType
from proto import vapp_vm_pb2 as vapp_vm_pb2
from proto import vapp_vm_pb2_grpc as vapp_vm_pb2_grpc


class VappVmServicer(vapp_vm_pb2_grpc.VappVmServicer):
    def __init__(self, pyPluginServer):
        self.py_plugin_server = pyPluginServer
        vref = VCDClientRef()
        self.client = vref.get_ref()

    def Create(self, request, context):
        logging.basicConfig(level=logging.DEBUG)
        logging.info("__INIT__Create[VappVmServicer]")
        source_catalog_name = request.source_catalog_name

        if len(source_catalog_name) > 0:
            res = self.CreateFromCatalog(request, context)
        else:
            res = self.CreateFromVapp(request, context)

        return res

    def CreateFromCatalog(self, request, context):
        vapp_vm = VappVm(context)
        res = vapp_vm.create_from_catalog(request)
        logging.info("__DONE__Create[VappVmServicer]")

        return res

    def CreateFromVapp(self, request, context):
        vapp_vm = VappVm(context)
        res = vapp_vm.create_from_vapp(request)
        logging.info("__DONE__Create[VappVmServicer]")

        return res

    def Delete(self, request, context):
        logging.info("__INIT__Delete[VappVmServicer]")
        vapp_vm = VappVm(context)
        res = vapp_vm.delete(request)
        logging.info("__DONE__Delete[VappVmServicer]")

        return res

    # def Update(self, request, context):
    #     logging.info("__INIT__Update[VappVmServicer]")
    #     target_vm_name = request.target_vm_name
    #     is_enabled = request.is_enabled
    #     vapp_vm = VappVm(target_vm_name=target_vm_name)
    #     res = vapp_vm.update()
    #     logging.info("__DONE__Update[VappVmServicer]")

    #     return res

    def Read(self, request, context):
        logging.info("__INIT__Read[VappVmServicer]")
        vapp_vm = VappVm(context)
        res = vapp_vm.read(request)
        logging.info("__DONE__Read[VappVmServicer]")

        return res

    def ModifyCPU(self, request, context):
        logging.info("__INIT__ModifyCPU[VappVmServicer]")
        vapp_vm = VappVm(context)
        res = vapp_vm.modify_cpu(request)
        logging.info("__DONE__ModifyCPU[VappVmServicer]")

        return res

    def ModifyMemory(self, request, context):
        logging.info("__INIT__ModifyMemory[VappVmServicer]")
        vapp_vm = VappVm(context)
        res = vapp_vm.modify_memory(request)
        logging.info("__DONE__ModifyMemory[VappVmServicer]")

        return res

    def PowerOn(self, request, context):
        logging.info("__INIT__PowerOn[VappVmServicer]")
        vapp_vm = VappVm(context)
        res = vapp_vm.power_on(request)
        logging.info("__DONE__PowerOn[VappVmServicer]")

        return res

    def PowerOff(self, request, context):
        logging.info("__INIT__PowerOff[VappVmServicer]")
        vapp_vm = VappVm(context)
        res = vapp_vm.power_off(request)

        logging.info("__DONE__PowerOff[VappVmServicer]")

        return res


class VappVm:
    def __init__(self, context):
        vref = VCDClientRef()
        self.client = vref.get_ref()
        self.context = context

    def get_vapp_resource(self, vdc_name, vapp_name):
        org_resource = Org(self.client, resource=self.client.get_org())
        vdc_resource = VDC(
            self.client, resource=org_resource.get_vdc(vdc_name))
        vapp_resource_href = vdc_resource.get_resource_href(
            name=vapp_name, entity_type=EntityType.VAPP)

        return self.client.get_resource(vapp_resource_href)

    def create_from_catalog(self, request):
        logging.info("__INIT__create[VappVm] source_catalog_name[%s]",
                     request.source_catalog_name)
        res = vapp_vm_pb2.CreateVappVmResult()
        res.created = False

        logged_in_org = self.client.get_org()
        org = Org(self.client, resource=logged_in_org)

        try:
            vdc_resource = org.get_vdc(request.target_vdc)
            vdc = VDC(
                self.client, name=request.target_vdc, resource=vdc_resource)

            vapp_resource = vdc.get_vapp(request.target_vapp)
            vapp = VApp(
                self.client, name=request.target_vapp, resource=vapp_resource)

            catalog_item = org.get_catalog_item(request.source_catalog_name,
                                                request.source_template_name)
            source_vapp_resource = self.client.get_resource(
                catalog_item.Entity.get('href'))
            specs = [{
                'source_vm_name': request.source_vm_name,
                'vapp': source_vapp_resource,
                'target_vm_name': request.target_vm_name,
                'hostname': request.hostname,
                'network': request.network,
                'ip_allocation_mode': request.ip_allocation_mode,
                # 'storage_profile': request.storage_profile
            }]
            create_vapp_vm_resp = vapp.add_vms(
                specs,
                power_on=request.power_on,
                all_eulas_accepted=request.all_eulas_accepted)
            task_monitor = self.client.get_task_monitor()
            task = task_monitor.wait_for_status(
                task=create_vapp_vm_resp,
                timeout=60,
                poll_frequency=2,
                fail_on_statuses=None,
                expected_target_statuses=[
                    TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                    TaskStatus.CANCELED
                ],
                callback=None)

            st = task.get('status')
            if st != TaskStatus.SUCCESS.value:
                raise errors.VappVmCreateError(
                    etree.tostring(task, pretty_print=True))

            message = 'status : {0} '.format(st)
            logging.info(message)
            res.created = True

        except Exception as e:
            errmsg = '''__ERROR_create[VappVm] failed for vm {0}. __ErrorMessage__ {1}'''
            logging.warn(errmsg.format(request.target_vm_name, str(e)))
            self.context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            self.context.set_details(errmsg)

            return res

        logging.info("__DONE__create[VappVm]")
        return res

    def create_from_vapp(self, request):
        logging.info("__INIT__create[VappVm] source_catalog_name[%s]",
                     request.source_vapp)
        res = vapp_vm_pb2.CreateVappVmResult()
        res.created = False
        source_vapp_resource = self.get_vapp_resource(
            request.target_vdc, vapp_name=request.source_vapp)
        target_vapp_resource = self.get_vapp_resource(
            request.target_vdc, vapp_name=request.target_vapp)

        specs = [{
            'vapp': source_vapp_resource,
            'source_vm_name': request.source_vm_name,
            'target_vm_name': request.target_vm_name,
            'hostname': request.hostname,
            'password': request.password,
            'password_auto': request.password_auto,
            'password_reset': request.password_reset,
            'cust_script': request.cust_script,
            'network': request.network,
            # 'storage_profile': request.storage_profile
        }]

        try:
            vapp = VApp(self.client, resource=target_vapp_resource)
            create_vapp_vm_resp = vapp.add_vms(
                specs,
                power_on=request.power_on,
                all_eulas_accepted=request.all_eulas_accepted)
            task_monitor = self.client.get_task_monitor()
            task = task_monitor.wait_for_status(
                task=create_vapp_vm_resp,
                timeout=60,
                poll_frequency=2,
                fail_on_statuses=None,
                expected_target_statuses=[
                    TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                    TaskStatus.CANCELED
                ],
                callback=None)

            st = task.get('status')
            if st != TaskStatus.SUCCESS.value:
                raise errors.VappVmCreateError(
                    etree.tostring(task, pretty_print=True))

            message = 'status : {0} '.format(st)
            logging.info(message)
            res.created = True

        except Exception as e:
            errmsg = '''__ERROR_create[VappVm] failed for vm {0}. __ErrorMessage__ {1}'''
            logging.warn(errmsg.format(request.target_vm_name, str(e)))
            self.context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            self.context.set_details(errmsg)

            return res

        logging.info("__DONE__create[VappVm]")
        return res

    # def update(self):
    #     logging.info("__INIT__update[VappVm]")
    #     res = vapp_vm_pb2.UpdateUserResult()
    #     res.updated = False
    #     logged_in_org = self.client.get_org()
    #     org = Org(self.client, resource=logged_in_org)
    #     name = self.name
    #     is_enabled = self.is_enabled

    #     try:
    #         result = org.update_user(name, is_enabled)
    #         res.updated = True
    #     except Exception as e:
    #         error_message = '__ERROR_update[VappVm] failed for VappVm {0}. __ErrorMessage__ {1}'.format(
    #             self.name, str(e))
    #         logging.warn(error_message)
    #         self.context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
    #         self.context.set_details(error_message)

    #         return res

    #     logging.info("__DONE__update[VappVm]")

    #     return res

    def read(self, request):
        logging.info("__INIT__read[VappVm]")
        res = vapp_vm_pb2.ReadVappVmResult()
        res.present = False
        org_resource = self.client.get_org()
        org = Org(self.client, resource=org_resource)
        try:
            vdc_resource = org.get_vdc(request.target_vdc)
            vdc = VDC(
                self.client, name=request.target_vdc, resource=vdc_resource)

            vapp_resource = vdc.get_vapp(request.target_vapp)
            vapp = VApp(
                self.client, name=request.target_vapp, resource=vapp_resource)
            read_vapp_vm_resp = vapp.get_vm(request.target_vm_name)
            vm = VM(client=self.client, href=None, resource=read_vapp_vm_resp)

            res.present = True
        except Exception as e:
            errmsg = '__ERROR_read[VappVm] failed for VappVm {0}. __ErrorMessage__ {1}'
            logging.warn(errmsg.format(request.target_vm_name, str(e)))

            return res

        logging.info("__DONE__read[VappVm]")

        return res

    def delete(self, request):
        logging.info("__INIT__delete[VappVm]")
        res = vapp_vm_pb2.DeleteVappVmResult()
        res.deleted = False
        org_resource = self.client.get_org()
        org = Org(self.client, resource=org_resource)
        try:
            vdc_resource = org.get_vdc(request.target_vdc)
            vdc = VDC(
                self.client, name=request.target_vdc, resource=vdc_resource)

            vapp_resource = vdc.get_vapp(request.target_vapp)
            vapp = VApp(
                self.client, name=request.target_vapp, resource=vapp_resource)

            # Before deleting power_off vm
            # self.power_off(request.target_vdc, request.target_vapp)

            # Before deleting undeploy vm
            self.undeploy(request)

            vms = [request.target_vm_name]
            delete_vapp_vm_resp = vapp.delete_vms(vms)
            task = self.client.get_task_monitor().wait_for_status(
                task=delete_vapp_vm_resp,
                timeout=60,
                poll_frequency=2,
                fail_on_statuses=None,
                expected_target_statuses=[
                    TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                    TaskStatus.CANCELED
                ],
                callback=None)

            st = task.get('status')
            if st != TaskStatus.SUCCESS.value:
                raise errors.VappVmDeleteError(
                    etree.tostring(task, pretty_print=True))

            message = 'delete vapp_vm status : {0} '.format(st)
            logging.info(message)
            res.deleted = True

        except Exception as e:
            res.deleted = False
            errmsg = '__ERROR_delete[VappVm] failed for VappVm {0}. __ErrorMessage__ {1}'
            logging.warn(errmsg.format(request.target_vm_name, str(e)))
            self.context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            self.context.set_details(errmsg)

            return res

        logging.info("__DONE__delete[VappVm]")

        return res

    def power_off(self, request):
        logging.info("__INIT__power_off[VappVm]")
        res = vapp_vm_pb2.PowerOffVappVmResult()
        res.powered_off = False
        org_resource = self.client.get_org()
        org = Org(self.client, resource=org_resource)
        try:
            vdc_resource = org.get_vdc(request.target_vdc)
            vdc = VDC(
                self.client, name=request.target_vdc, resource=vdc_resource)

            vapp_resource = vdc.get_vapp(request.target_vapp)
            vapp = VApp(
                self.client, name=request.target_vapp, resource=vapp_resource)
            vapp_vm_resource = vapp.get_vm(request.target_vm_name)
            vm = VM(self.client, resource=vapp_vm_resource)
            power_off_response = vm.undeploy()

            task = self.client.get_task_monitor().wait_for_status(
                task=power_off_response,
                timeout=60,
                poll_frequency=2,
                fail_on_statuses=None,
                expected_target_statuses=[
                    TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                    TaskStatus.CANCELED
                ],
                callback=None)

            st = task.get('status')
            if st != TaskStatus.SUCCESS.value:
                raise errors.VappVmCreateError(
                    etree.tostring(task, pretty_print=True))

            message = 'status : {0} '.format(st)
            logging.info(message)
            powered_off = True

        except Exception as e:
            errmsg = '__ERROR_power_off[VappVm] failed for VappVm {0}. __ErrorMessage__ {1}'
            logging.warn(errmsg.format(request.target_vm_name, str(e)))

        logging.info("__DONE__power_off[VappVm]")

        return powered_off

    def power_on(self, request):
        logging.info("__INIT__power_on[VappVm]")
        res = vapp_vm_pb2.PowerOnVappVmResult()
        res.powered_on = False

        org_resource = self.client.get_org()
        org = Org(self.client, resource=org_resource)
        try:
            vdc_resource = org.get_vdc(request.target_vdc)
            vdc = VDC(
                self.client, name=request.target_vdc, resource=vdc_resource)

            vapp_resource = vdc.get_vapp(request.target_vapp)
            vapp = VApp(
                self.client, name=request.target_vapp, resource=vapp_resource)
            vapp_vm_resource = vapp.get_vm(request.target_vm_name)
            vm = VM(self.client, resource=vapp_vm_resource)
            power_on_response = vm.power_on()

            task = self.client.get_task_monitor().wait_for_status(
                task=power_on_response,
                timeout=60,
                poll_frequency=2,
                fail_on_statuses=None,
                expected_target_statuses=[
                    TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                    TaskStatus.CANCELED
                ],
                callback=None)

            st = task.get('status')
            if st != TaskStatus.SUCCESS.value:
                raise errors.VappVmPowerOnError(
                    etree.tostring(task, pretty_print=True))

            message = 'status : {0} '.format(st)
            logging.info(message)
            res.powered_on = True

        except Exception as e:
            errmsg = '__ERROR_power_off[VappVm] failed for VappVm {0}. __ErrorMessage__ {1}'
            logging.warn(errmsg.format(request.target_vm_name, str(e)))
            self.context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            self.context.set_details(errmsg)

        logging.info("__DONE__power_on[VappVm]")

        return res

    def undeploy(self, request):
        logging.info("__INIT__undeploy[VappVm]")
        undeploy = False
        org_resource = self.client.get_org()
        org = Org(self.client, resource=org_resource)
        try:
            vdc_resource = org.get_vdc(request.target_vdc)
            vdc = VDC(
                self.client, name=request.target_vdc, resource=vdc_resource)

            vapp_resource = vdc.get_vapp(request.target_vapp)
            vapp = VApp(
                self.client, name=request.target_vapp, resource=vapp_resource)
            vapp_vm_resource = vapp.get_vm(request.target_vm_name)
            vm = VM(self.client, resource=vapp_vm_resource)
            undeploy_response = vm.undeploy()

            task = self.client.get_task_monitor().wait_for_status(
                task=undeploy_response,
                timeout=60,
                poll_frequency=2,
                fail_on_statuses=None,
                expected_target_statuses=[
                    TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                    TaskStatus.CANCELED
                ],
                callback=None)

            st = task.get('status')
            if st != TaskStatus.SUCCESS.value:
                raise errors.VappVmCreateError(
                    etree.tostring(task, pretty_print=True))

            message = 'status : {0} '.format(st)
            logging.info(message)
            undeploy = True

        except Exception as e:
            errmsg = '__ERROR_undeploy[VappVm] failed for VappVm {0}. __ErrorMessage__ {1}'
            logging.warn(errmsg.format(request.target_vm_name, str(e)))

        logging.info("__DONE__undeploy[VappVm]")

        return undeploy

    def modify_cpu(self, request):
        logging.info("__INIT__modify_cpu[VappVm]")
        res = vapp_vm_pb2.ModifyVappVmCPUResult()
        res.modified = False

        org_resource = self.client.get_org()
        org = Org(self.client, resource=org_resource)
        try:
            vdc_resource = org.get_vdc(request.target_vdc)
            vdc = VDC(
                self.client, name=request.target_vdc, resource=vdc_resource)

            vapp_resource = vdc.get_vapp(request.target_vapp)
            vapp = VApp(
                self.client, name=request.target_vapp, resource=vapp_resource)
            vapp_vm_resource = vapp.get_vm(request.target_vm_name)
            vm = VM(self.client, resource=vapp_vm_resource)

            self.undeploy(request)
            modify_cpu_response = vm.modify_cpu(request.virtual_cpus,
                                                request.cores_per_socket)

            task = self.client.get_task_monitor().wait_for_status(
                task=modify_cpu_response,
                timeout=60,
                poll_frequency=2,
                fail_on_statuses=None,
                expected_target_statuses=[
                    TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                    TaskStatus.CANCELED
                ],
                callback=None)
            st = task.get('status')
            if st != TaskStatus.SUCCESS.value:
                raise errors.VappVmModifyCPUError(
                    etree.tostring(task, pretty_print=True))

            message = 'status : {0} '.format(st)
            logging.info(message)

            self.power_on(request)
            res.modified = True

        except Exception as e:
            errmsg = '__ERROR_modify_cpu[VappVm] failed for VappVm {0}. __ErrorMessage__ {1}'
            logging.warn(errmsg.format(request.target_vm_name, str(e)))
            self.context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            self.context.set_details(errmsg)

        logging.info("__DONE__modify_cpu[VappVm]")

        return res

    def modify_memory(self, request):
        logging.info("__INIT__modify_memory[VappVm]")
        res = vapp_vm_pb2.ModifyVappVmMemoryResult()
        res.modified = False
        org_resource = self.client.get_org()
        org = Org(self.client, resource=org_resource)
        try:
            vdc_resource = org.get_vdc(request.target_vdc)
            vdc = VDC(
                self.client, name=request.target_vdc, resource=vdc_resource)

            vapp_resource = vdc.get_vapp(request.target_vapp)
            vapp = VApp(
                self.client, name=request.target_vapp, resource=vapp_resource)
            vapp_vm_resource = vapp.get_vm(request.target_vm_name)
            vm = VM(self.client, resource=vapp_vm_resource)

            self.undeploy(request)
            modify_memory_response = vm.modify_memory(request.memory)

            task = self.client.get_task_monitor().wait_for_status(
                task=modify_memory_response,
                timeout=60,
                poll_frequency=2,
                fail_on_statuses=None,
                expected_target_statuses=[
                    TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                    TaskStatus.CANCELED
                ],
                callback=None)
            st = task.get('status')
            if st != TaskStatus.SUCCESS.value:
                raise errors.VappVmModifyMemoryError(
                    etree.tostring(task, pretty_print=True))
            message = 'status : {0} '.format(st)
            logging.info(message)

            self.power_on(request)
            res.modified = True

        except Exception as e:
            errmsg = '__ERROR_modify_memory[VappVm] failed for VappVm {0}. __ErrorMessage__ {1}'
            logging.warn(errmsg.format(request.target_vm_name, str(e)))
            self.context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            self.context.set_details(errmsg)

        logging.info("__DONE__modify_memory[VappVm]")

        return res

    def reload_vm(self, request):
        logging.info("__INIT__reload_vm[VappVm]")
        org_resource = self.client.get_org()
        org = Org(self.client, resource=org_resource)
        try:
            vdc_resource = org.get_vdc(request.target_vdc)
            vdc = VDC(
                self.client, name=request.target_vdc, resource=vdc_resource)

            vapp_resource = vdc.get_vapp(request.target_vapp)
            vapp = VApp(
                self.client, name=request.target_vapp, resource=vapp_resource)
            vapp_vm_resource = vapp.get_vm(request.target_vm_name)
            vm = VM(self.client, resource=vapp_vm_resource)

            reload_vm_response = vm.reload()

            task = self.client.get_task_monitor().wait_for_status(
                task=reload_vm_response,
                timeout=60,
                poll_frequency=2,
                fail_on_statuses=None,
                expected_target_statuses=[
                    TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                    TaskStatus.CANCELED
                ],
                callback=None)
            st = task.get('status')
            if st != TaskStatus.SUCCESS.value:
                raise errors.VappVmReloadError(
                    etree.tostring(task, pretty_print=True))

            message = 'status : {0} '.format(st)
            logging.info(message)

        except Exception as e:
            errmsg = '__ERROR_reload_vm[VappVm] failed for VappVm {0}. __ErrorMessage__ {1}'
            logging.warn(errmsg.format(request.target_vm_name, str(e)))
            raise errors.VappVmReloadError(
                etree.tostring(errmsg, pretty_print=True))

        logging.info("__DONE__reload_vm[VappVm]")
