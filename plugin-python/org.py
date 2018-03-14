#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

from proto import org_pb2_grpc as org_pb2_grpc
from proto import org_pb2 as org_pb2

import logging
import grpc
from pyvcloud.vcd.system import System

from vcd_client_ref import VCDClientRef
import errors
from pyvcloud.vcd.org import Org
from pyvcloud.vcd.client import TaskStatus
from lxml import objectify, etree
import lxml

#from __builtin__ import str

#from pyimports import *


class OrgServicer(org_pb2_grpc.OrgServicer):
    def __init__(self, pyPluginServer):
        self.py_plugin_server = pyPluginServer

    def Create(self, request, context):
        """
        # API used https://github.com/vmware/pyvcloud/blob/master/pyvcloud/vcd/system.py#L43
        """

        logging.basicConfig(level=logging.DEBUG)
        logging.info("__INIT__Create[org_plugin]")
        res = org_pb2.CreateOrgResult()
        res.created = False
        try:
            vref = VCDClientRef()
            client = vref.get_ref()
            sys_admin = client.get_admin()
            system = System(client, admin_resource=sys_admin)
            system.create_org(request.name, request.org_full_name,
                              request.is_enabled)
            logging.info("__DONE_Create[org_plugin]")
            res.created = True
            return res
        except Exception as e:
            error_message = '__ERROR_Create[org_plugin] failed  {0} '.format(
                request.name)
            logging.warn(error_message, e)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
            return res

    def Read(self, request, context):
        """Check if an organization exist.
        :param client: client
        :param org_name: (string): name of organisation
        :return: (bool) : true/false based on the api execution success
        """
        logging.info("__INIT__Read[org_plugin]")

        org_name = request.name
        res = org_pb2.ReadOrgResult()
        res.present = False

        try:
            vref = VCDClientRef()
            client = vref.get_ref()
            resource = client.get_org_by_name(org_name)
            org = Org(client, resource=resource)
            org_admin_resource = org.client.get_resource(org.href_admin)

            is_enabled = org_admin_resource['IsEnabled']
            org_full_name = str(org_admin_resource['FullName'])

            res.present = True
            res.name = org_name
            res.org_full_name = org_full_name
            res.is_enabled = is_enabled
            logging.info("__DONE_Read[org_plugin] %s", res)
            return res
        except Exception as e:
            error_message = '__ERROR_Read[org_plugin] failed for org {0} '.format(
                org_name)
            logging.warn(error_message, e)
            #context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            #context.set_details(error_message)
            return res

    def Delete(self, request, context):
        """
        # API used https://github.com/vmware/pyvcloud/blob/master/pyvcloud/vcd/system.py#L62
        """
        logging.info("__INIT_Delete[org_plugin]")
        res = org_pb2.DeleteOrgResult()
        res.deleted = False
        try:
            vref = VCDClientRef()
            client = vref.get_ref()
            sys_admin = client.get_admin()
            system = System(client, admin_resource=sys_admin)
            delete_org_resp = system.delete_org(request.name, request.force,
                                                request.recursive)

            task = client.get_task_monitor().wait_for_status(
                task=delete_org_resp,
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
                message = 'delete org status : {0} '.format(st)
                logging.info(message)
            else:
                raise errors.VCDOrgDeleteError(
                    etree.tostring(task, pretty_print=True))

            logging.info("__DONE_Delete[org_plugin]")
            res.deleted = True
            return res
        except Exception as e:
            error_message = '__ERROR_Delete[org_plugin] failed  {0} '.format(
                request.name)
            logging.warn(error_message, e)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
            return res

    def Update(self, request, context):
        """
        Update an organization. 
        Expected request params: 
            :param org_name: (string): name of organisation
            :param is_enabled: (bool): enable/disable the organization
            :return (bool): true/false based on the api execution success
        """
        logging.info("__INIT_update_org_[org_plugin]")
        org_name = request.name
        is_enabled = request.is_enabled

        res = org_pb2.UpdateOrgResult()
        res.updated = False

        try:
            vref = VCDClientRef()
            client = vref.get_ref()
            org_resource = client.get_org_by_name(org_name)
            org = Org(client, resource=org_resource)
            org.update_org(is_enabled)
            logging.info("__DONE_update_org_[org_plugin]")
            res.updated = True
            return res
        except Exception as e:
            error_message = '__ERROR_update_org_[org_plugin] failed for org {0} '.format(
                org_name)
            logging.warn(error_message, e)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
            return res
