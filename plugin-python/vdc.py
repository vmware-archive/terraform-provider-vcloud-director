#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************
import logging

import errors
import grpc
from lxml import objectify, etree
import lxml
from proto import vdc_pb2 as vdc_pb2
from proto import vdc_pb2_grpc as vdc_pb2_grpc
from pyvcloud.vcd.client import TaskStatus
from pyvcloud.vcd.org import Org
from pyvcloud.vcd.system import System
from vcd_client_ref import VCDClientRef
import json
from pyvcloud.vcd.vdc import VDC


class VdcServicer(vdc_pb2_grpc.VdcServicer):
    def __init__(self, pyPluginServer):
        self.py_plugin_server = pyPluginServer
        vref = VCDClientRef()
        self.client = vref.get_ref()

    def Create(self, request, context):
        logging.basicConfig(level=logging.DEBUG)
        logging.info("__INIT__Create[VdcServicer]")

        name = request.name
        is_enabled = request.is_enabled
        provider_vdc_name = request.provider_vdc
        description = request.description
        allocation_model = request.allocation_model
        storage_profiles = request.storage_profiles
        cpu_units = request.cpu_units
        cpu_allocated = request.cpu_allocated
        cpu_limit = request.cpu_limit
        mem_units = request.mem_units
        mem_allocated = request.mem_allocated
        mem_limit = request.mem_limit
        nic_quota = request.nic_quota
        network_quota = request.network_quota
        vm_quota = request.vm_quota
        resource_guaranteed_memory = request.resource_guaranteed_memory
        resource_guaranteed_cpu = request.resource_guaranteed_cpu
        vcpu_in_mhz = request.vcpu_in_mhz
        is_thin_provision = request.is_thin_provision
        network_pool_name = request.network_pool_name
        uses_fast_provisioning = request.uses_fast_provisioning
        over_commit_allowed = request.over_commit_allowed
        vm_discovery_enabled = request.vm_discovery_enabled

        vdc = Vdc(name, is_enabled, context)
        res = vdc.create(
            provider_vdc_name=provider_vdc_name,
            description=description,
            allocation_model=allocation_model,
            storage_profiles=storage_profiles,
            cpu_units=cpu_units,
            cpu_allocated=cpu_allocated,
            cpu_limit=cpu_limit,
            mem_units=mem_units,
            mem_allocated=mem_allocated,
            mem_limit=mem_limit,
            nic_quota=nic_quota,
            network_quota=network_quota,
            vm_quota=vm_quota,
            resource_guaranteed_memory=resource_guaranteed_memory,
            resource_guaranteed_cpu=resource_guaranteed_cpu,
            vcpu_in_mhz=vcpu_in_mhz,
            is_thin_provision=is_thin_provision,
            network_pool_name=network_pool_name,
            uses_fast_provisioning=uses_fast_provisioning,
            over_commit_allowed=over_commit_allowed,
            vm_discovery_enabled=vm_discovery_enabled)

        logging.info("__DONE__Create[VdcServicer]")
        return res

    def Delete(self, request, context):
        logging.info("__INIT__Delete[VdcServicer]")

        name = request.name
        vdc = Vdc(name, None, context)
        res = vdc.delete()
        logging.info("__DONE__Delete[VdcServicer]")
        return res

    def Update(self, request, context):
        logging.info("__INIT__Update[VdcServicer]")

        name = request.name
        is_enabled = request.is_enabled
        vdc = Vdc(name=name, is_enabled=is_enabled, context=context)
        res = vdc.update()
        logging.info("__DONE__Update[VdcServicer]")
        return res

    def Read(self, request, context):
        logging.info("__INIT__Read[VdcServicer]")
        name = request.name
        vdc = Vdc(name=name, context=context)
        res = vdc.read()
        logging.info("__DONE__Read[VdcServicer]")
        return res


class Vdc:
    def __repr__(self):
        message = 'Vdc [name ={0}, is_enabled = {1}'.format(
            self.name, self.is_enabled)
        return message

    def __init__(self, name, is_enabled=True, context=None):
        vref = VCDClientRef()
        self.client = vref.get_ref()
        self.name = name
        self.is_enabled = is_enabled
        self.context = context

    def create(self, provider_vdc_name, description, allocation_model,
               storage_profiles, cpu_units, cpu_allocated, cpu_limit,
               mem_units, mem_allocated, mem_limit, nic_quota, network_quota,
               vm_quota, resource_guaranteed_memory, resource_guaranteed_cpu,
               vcpu_in_mhz, is_thin_provision, network_pool_name,
               uses_fast_provisioning, over_commit_allowed,
               vm_discovery_enabled):
        logging.info("__INIT__create[Vdc]")

        res = vdc_pb2.CreateVdcResult()
        res.created = False

        context = self.context
        logged_in_org = self.client.get_org()
        org = Org(self.client, resource=logged_in_org)

        #Vdc details
        name = self.name
        is_enabled = self.is_enabled

        try:
            storage_profiles = json.loads(storage_profiles)
            if not network_pool_name.strip():
                network_pool_name = None

            create_vdc_resp = org.create_org_vdc(
                vdc_name=name,
                provider_vdc_name=provider_vdc_name,
                description=description,
                allocation_model=allocation_model,
                storage_profiles=storage_profiles,
                cpu_units=cpu_units,
                cpu_allocated=cpu_allocated,
                cpu_limit=cpu_limit,
                mem_units=mem_units,
                mem_allocated=mem_allocated,
                mem_limit=mem_limit,
                nic_quota=nic_quota,
                network_quota=network_quota,
                vm_quota=vm_quota,
                resource_guaranteed_memory=resource_guaranteed_memory,
                resource_guaranteed_cpu=resource_guaranteed_cpu,
                vcpu_in_mhz=vcpu_in_mhz,
                is_thin_provision=is_thin_provision,
                network_pool_name=network_pool_name,
                uses_fast_provisioning=uses_fast_provisioning,
                over_commit_allowed=over_commit_allowed,
                vm_discovery_enabled=vm_discovery_enabled,
                is_enabled=is_enabled)

            task = self.client.get_task_monitor().wait_for_status(
                task=create_vdc_resp.Tasks.Task[0],
                timeout=60,
                poll_frequency=2,
                fail_on_statuses=None,
                expected_target_statuses=[
                    TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                    TaskStatus.CANCELED
                ],
                callback=None)

            st = task.get('status')
            if st == TaskStatus.SUCCESS.value:
                message = 'delete vdc status : {0} '.format(st)
                logging.info(message)
            else:
                raise errors.VCDVdcCreateError(
                    etree.tostring(task, pretty_print=True))

            res.created = True
        except Exception as e:
            error_message = '__ERROR_create[Vdc] failed for Vdc {0}. __ErrorMessage__ {1}'.format(
                self.name, str(e))
            logging.warn(error_message)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
            return res

        logging.info("__DONE__create[Vdc]")
        return res

    def update(self):
        logging.info("__INIT__update[Vdc]")
        res = vdc_pb2.UpdateVdcResult()
        res.updated = False
        context = self.context
        logged_in_org = self.client.get_org()
        org = Org(self.client, resource=logged_in_org)

        #Vdc details
        name = self.name
        is_enabled = self.is_enabled

        try:
            vdc_resource = org.get_vdc(name)
            vdc = VDC(self.client, name=name, resource=vdc_resource)
            result = vdc.enable_vdc(is_enabled)
            res.updated = True
        except Exception as e:
            error_message = '__ERROR_update[Vdc] failed for Vdc {0}. __ErrorMessage__ {1}'.format(
                self.name, str(e))
            logging.warn(error_message)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
            return res

        logging.info("__DONE__update[Vdc]")
        return res

    def read(self):
        logging.info("__INIT__read[Vdc]")
        res = vdc_pb2.ReadVdcResult()
        res.present = False
        context = self.context
        logged_in_org = self.client.get_org()
        org = Org(self.client, resource=logged_in_org)

        #Vdc details
        name = self.name

        try:
            vdc_resource = org.get_vdc(name)

            res.name = name
            #res.provider_vdc = str(result.provider_vdc)
            res.description = str(vdc_resource.Description)
            res.allocation_model = str(vdc_resource.AllocationModel)

            res.cpu_units = str(vdc_resource.ComputeCapacity.Cpu.Units)
            res.cpu_allocated = int(vdc_resource.ComputeCapacity.Cpu.Allocated)
            res.cpu_limit = int(vdc_resource.ComputeCapacity.Cpu.Limit)

            res.mem_units = str(vdc_resource.ComputeCapacity.Memory.Units)
            res.mem_allocated = int(
                vdc_resource.ComputeCapacity.Memory.Allocated)
            res.mem_limit = int(vdc_resource.ComputeCapacity.Memory.Limit)

            res.nic_quota = int(vdc_resource.NicQuota)
            res.network_quota = int(vdc_resource.NetworkQuota)
            res.vm_quota = int(vdc_resource.VmQuota)

            #res.storage_profiles = str(vdc_resource.VdcStorageProfiles.VdcStorageProfile.get('name'))

            #res.resource_guaranteed_memory = str(result.resource_guaranteed_memory)
            #res.resource_guaranteed_cpu = str(result.resource_guaranteed_cpu)

            res.vcpu_in_mhz = int(vdc_resource.VCpuInMhz2)

            #res.is_thin_provision = str(result.is_thin_provision)
            #res.network_pool_name = str(result.network_pool_name)
            #res.uses_fast_provisioning = str(result.uses_fast_provisioning)
            #res.over_commit_allowed = str(result.over_commit_allowed)
            #res.vm_discovery_enabled = str(result.vm_discovery_enabled)

            res.is_enabled = vdc_resource.IsEnabled
            res.present = True
        except Exception as e:
            error_message = '__ERROR_read[Vdc] failed for Vdc {0}. __ErrorMessage__ {1}'.format(
                self.name, str(e))
            logging.warn(error_message)
            #context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            #context.set_details(error_message)
            return res
        logging.info("__DONE__read[Vdc]")
        return res

    def delete(self):
        logging.info("__INIT__delete[Vdc]")
        res = vdc_pb2.DeleteVdcResult()
        res.deleted = False

        context = self.context
        logged_in_org = self.client.get_org()
        org = Org(self.client, resource=logged_in_org)

        #Vdc details
        name = self.name

        try:
            vdc_resource = org.get_vdc(name)
            vdc = VDC(self.client, name=name, resource=vdc_resource)
            delete_vdc_resp = vdc.delete_vdc()

            task = self.client.get_task_monitor().wait_for_status(
                task=delete_vdc_resp,
                timeout=60,
                poll_frequency=2,
                fail_on_statuses=None,
                expected_target_statuses=[
                    TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                    TaskStatus.CANCELED
                ],
                callback=None)

            st = task.get('status')
            if st == TaskStatus.SUCCESS.value:
                message = 'delete vdc status : {0} '.format(st)
                logging.info(message)
            else:
                raise errors.VCDVdcDeleteError(
                    etree.tostring(task, pretty_print=True))

            res.deleted = True
        except Exception as e:
            res.deleted = False
            error_message = '__ERROR_delete[Vdc] failed for Vdc {0}. __ErrorMessage__ {1}'.format(
                self.name, str(e))
            logging.warn(error_message)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
            return res

        logging.info("__DONE__delete[Vdc]")
        return res
