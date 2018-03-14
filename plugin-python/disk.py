#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

from proto import disk_pb2_grpc as disk_pb2_grpc
from proto import disk_pb2 as disk_pb2

import logging
import grpc

from vcd_client_ref import VCDClientRef

from pyimports import *


class IndependentDiskServicer(disk_pb2_grpc.IndependentDiskServicer):
    """Implementation of PyVcloudProviderServicer service."""

    def __init__(self, pyPluginServer):
        self.py_plugin_server = pyPluginServer

    def Create(self, request, context):
        logging.info("__INIT__Create[disk_plugin]")
        error_message = "dummy"
        #context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        #context.set_details(error_message)

        d = Disk(
            name=request.name,
            size=request.size,
            vdc=request.vdc,
            description=request.description,
        )
        res = d.create()

        logging.info("__DONOE_Create[disk_plugin]")
        return res

    def Read(self, request, context):
        logging.info("__INIT__Read[disk_plugin]")
        error_message = "dummy"
        #context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        #context.set_details(error_message)

        d = Disk(name=request.name, vdc=request.vdc, disk_id=request.disk_id)
        res = d.read()

        logging.info("__DONE_Read[disk_plugin] %s", res)
        return res

    def Delete(self, request, context):
        logging.info("__INIT_Delete[disk_plugin]")
        res = disk_pb2.DeleteDiskResult()
        d = Disk(request.name, vdc=request.vdc, disk_id=request.disk_id)
        d.delete()
        logging.info("__INIT_Delete[disk_plugin]")
        return res


class Disk:
    def __repr__(self):

        message = 'Disk [name ={0} ,id={1}] '.format(self.name, self.disk_id)
        return message

    def disk_callback(self, task):
        message = '{0}: {1}, status: {2}'.format(
            task.get('operationName'), task.get('operation'),
            task.get('status'))
        logging.info(message)

    def __init__(self,
                 name,
                 vdc,
                 size=None,
                 description=None,
                 storage_profile=None,
                 disk_id=None):
        vref = VCDClientRef()
        self.client = vref.get_ref()
        self.name = name
        self.vdc = vdc
        self.size = size
        self.description = description
        self.storage_profile = storage_profile
        self.disk_id = disk_id

    def create(self):
        logging.info('INIT create_disk %s ', self)
        res = disk_pb2.CreateDiskResult()
        res.created = False

        logged_in_org = self.client.get_org()
        org = Org(self.client, resource=logged_in_org)
        v = org.get_vdc(self.vdc)
        if v is None:
            raise errors.VDCNotFoundError(vappInfo.vdc)
        vdc = VDC(self.client, href=v.get('href'))

        #task=vdc.delete_disk('disk01')

        result = vdc.add_disk(
            self.name,
            size=self.size,
            storage_profile_name=self.storage_profile,
            description=self.description,
            customize_on_instantiate=self.customize_on_instantiate)
        logging.info('DONE create_disk id  [%s]', result.get('id'))

        task = self.client.get_task_monitor().wait_for_status(
            task=result.Tasks.Task[0],
            timeout=60,
            poll_frequency=2,
            fail_on_status=None,
            expected_target_statuses=[
                TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                TaskStatus.CANCELED
            ],
            callback=self.disk_callback)
        st = task.get('status')
        if st == TaskStatus.SUCCESS.value:
            logging.info("__LOG__ created DISK")
            #cresult.created = True
        else:
            raise errors.VCDDiskCreationError(
                etree.tostring(task, pretty_print=True))
        res.created = True
        res.disk_id = result.get('id')
        return res

    def read(self):
        logging.info('__INIT__ read_disk [%s] ', self)

        res = disk_pb2.ReadDiskResult()
        res.present = False

        logged_in_org = self.client.get_org()
        org = Org(self.client, resource=logged_in_org)
        v = org.get_vdc(self.vdc)
        if v is None:
            raise errors.VDCNotFoundError(vappInfo.vdc)
        vdc = VDC(self.client, href=v.get('href'))

        #task=vdc.delete_disk('disk01')

        disk_id = None
        if self.disk_id:
            disk_id = self.disk_id

        result = vdc.get_disk(self.name, disk_id=disk_id)
        if result:
            logging.info("got disk = %s", result.get('id'))
            res.present = True
            res.disk_id = result.get('id')

        logging.info('__DONE__ read_disk [%s] id = [%s] ', self, res)
        return res

    def delete(self):
        logging.info('INIT delete_disk [%s]', self)

        logged_in_org = self.client.get_org()
        org = Org(self.client, resource=logged_in_org)
        v = org.get_vdc(self.vdc)
        if v is None:
            raise errors.VDCNotFoundError(vappInfo.vdc)
        vdc = VDC(self.client, href=v.get('href'))

        disk_id = None
        if self.disk_id:
            disk_id = self.disk_id
        result = vdc.delete_disk(self.name, disk_id=disk_id)

        task = self.client.get_task_monitor().wait_for_status(
            task=result,
            timeout=60,
            poll_frequency=2,
            fail_on_status=None,
            expected_target_statuses=[
                TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                TaskStatus.CANCELED
            ],
            callback=self.disk_callback)
        st = task.get('status')
        if st == TaskStatus.SUCCESS.value:
            logging.info("__LOG__ created DISK")
        else:
            raise errors.VCDDiskDeletetionError(
                etree.tostring(task, pretty_print=True))
        logging.info('DONE Delete_disk %s', result)
