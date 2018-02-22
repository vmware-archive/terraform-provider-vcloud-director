#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************
from proto import user_pb2_grpc as user_pb2_grpc
from proto import user_pb2 as user_pb2

import logging
import grpc
from pyvcloud.vcd.system import System

from vcd_client_ref import VCDClientRef
import errors
from pyvcloud.vcd.org import Org
from pyvcloud.vcd.client import TaskStatus
from lxml import objectify, etree
import lxml


class UserServicer(user_pb2_grpc.UserServicer):
    def __init__(self, pyPluginServer):
        self.py_plugin_server = pyPluginServer
        vref = VCDClientRef()
        self.client = vref.get_ref()

    def Create(self, request, context):
        logging.basicConfig(level=logging.DEBUG)
        logging.info("__INIT__Create[UserServicer]")

        name = request.name
        password = request.password
        role_name = request.role_name
        full_name = request.full_name
        description = request.description
        email = request.email
        telephone = request.telephone
        im = request.im
        alert_email = request.alert_email
        alert_email_prefix = request.alert_email_prefix
        stored_vm_quota = request.stored_vm_quota
        deployed_vm_quota = request.deployed_vm_quota
        is_group_role = request.is_group_role
        is_default_cached = request.is_default_cached
        is_external = request.is_external
        is_alert_enabled = request.is_alert_enabled
        is_enabled = request.is_enabled

        user = User(name, password, role_name, is_enabled, context)
        res = user.create(description, full_name, email, telephone, im,
                          alert_email, alert_email_prefix, stored_vm_quota,
                          deployed_vm_quota, is_group_role, is_default_cached,
                          is_external, is_alert_enabled)

        logging.info("__DONE__Create[UserServicer]")
        return res

    def Delete(self, request, context):
        logging.info("__DONE__Delete[UserServicer]")

        name = request.name
        user = User(name=name, context=context)
        res = user.delete()
        logging.info("__DONE__Delete[UserServicer]")
        return res

    def Update(self, request, context):
        logging.info("__DONE__Update[UserServicer]")

        name = request.name
        is_enabled = request.is_enabled
        user = User(name=name, is_enabled=is_enabled, context=context)
        res = user.update()
        logging.info("__DONE__Update[UserServicer]")
        return res

    def Read(self, request, context):
        logging.info("__INIT__Read[UserServicer]")
        name = request.name
        user = User(name=name, context=context)
        res = user.read()
        logging.info("__DONE__Read[UserServicer]")
        return res


class User:
    def __repr__(self):
        message = 'User [name ={0} '.format(self.name)
        return message

    def __init__(self,
                 name,
                 password=None,
                 role_name="",
                 is_enabled=True,
                 context=None):
        vref = VCDClientRef()
        self.client = vref.get_ref()
        self.name = name
        self.password = password
        self.role_name = role_name
        self.is_enabled = is_enabled
        self.context = context

    def create(self,
               description='',
               full_name='',
               email='',
               telephone='',
               im='',
               alert_email='',
               alert_email_prefix='',
               stored_vm_quota=0,
               deployed_vm_quota=0,
               is_group_role=False,
               is_default_cached=False,
               is_external=False,
               is_alert_enabled=False):
        logging.info("__INIT__create[User]")
        res = user_pb2.CreateUserResult()
        res.created = False

        context = self.context

        logged_in_org = self.client.get_org()
        org = Org(self.client, resource=logged_in_org)
        logging.info("__role_name__ %s org[%s]", self.role_name, org)
        role = org.get_role_record(self.role_name)
        role_href = role.get('href')

        try:
            result = org.create_user(
                self.name, self.password, role_href, full_name, description,
                email, telephone, im, alert_email, alert_email_prefix,
                stored_vm_quota, deployed_vm_quota, is_group_role,
                is_default_cached, is_external, is_alert_enabled,
                self.is_enabled)
            res.created = True
        except Exception as e:
            error_message = '__ERROR_create[user] failed for user {0}. __ErrorMessage__ {1}'.format(
                self.name, str(e))
            logging.warn(error_message)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
            return res

        logging.info("__DONE__create[User]")
        return res

    def update(self):
        logging.info("__INIT__update[User]")
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
            error_message = '__ERROR_update[user] failed for user {0}. __ErrorMessage__ {1}'.format(
                self.name, str(e))
            logging.warn(error_message)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
            return res

        logging.info("__DONE__update[User]")
        return res

    def read(self):
        logging.info("__INIT__read[User]")
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
            error_message = '__ERROR_read[user] failed for user {0}. __ErrorMessage__ {1}'.format(
                self.name, str(e))
            logging.warn(error_message)
            #context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            #context.set_details(error_message)
            return res
        logging.info("__DONE__read[User]")
        return res

    def delete(self):
        logging.info("__INIT__delete[User]")
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
            error_message = '__ERROR_delete[user] failed for user {0}. __ErrorMessage__ {1}'.format(
                self.name, str(e))
            logging.warn(error_message)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
            return res

        logging.info("__DONE__delete[User]")
        return res
