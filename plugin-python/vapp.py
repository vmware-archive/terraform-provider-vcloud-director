#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

from pyvcloud.vcd.client import _WellKnownEndpoint
from pyvcloud.vcd.client import API_CURRENT_VERSIONS
from pyvcloud.vcd.client import BasicLoginCredentials
from pyvcloud.vcd.client import Client
from pyvcloud.vcd.client import EntityType
from pyvcloud.vcd.client import get_links
from pyvcloud.vcd.org import Org
from pyvcloud.vcd.vdc import VDC
from pyvcloud.vcd.vapp import VApp

from proto import pyvcloudprovider_pb2 as pyvcloudprovider_pb2
from proto import pyvcloudprovider_pb2_grpc as pyvcloudprovider_pb2_grpc

from proto import vapp_pb2 as vapp_pb2

import requests
import logging

from pyvcloud.vcd.client import TaskStatus
import grpc
import pdb
import errors

from lxml import objectify, etree
import lxml


def task_callback(task):
    message = '{0}: {1}, status: {2}'.format(
        task.get('operationName'), task.get('operation'), task.get('status'))
    logging.info(message)


def create(client, context, vappInfo):
    logging.debug('__INIT__vapp_create [ {0} ]'.format(vappInfo))
    cresult = vapp_pb2.CreateVAppResult()
    cresult.created = False
    #cresult.in_vapp_info=vapp_pb2.CreateVAppInfo()
    cresult.in_vapp_info.CopyFrom(vappInfo)
    try:
        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)
        v = org.get_vdc(vappInfo.vdc)
        if v is None:
            raise errors.VDCNotFoundError(vappInfo.vdc)
        vdc = VDC(client, href=v.get('href'))

        logging.info("__LOG__vapp_create_calling instantiate_vapp............")

        network = None
        memory = None
        storage_profile = None
        if vappInfo.network:
            network = vappInfo.network
        if vappInfo.memory:
            memory = vappInfo.memory
        if vappInfo.storage_profile:
            storage_profile = vappInfo.storage_profile

        logging.info(
            "__LOG__ CREATE VAPP Params - MEMORY = [%s] NETWORK = [%s] storage_profile =[%s] ",
            memory, network, storage_profile)

        result = vdc.instantiate_vapp(
            name=vappInfo.name,
            catalog=vappInfo.catalog_name,
            template=vappInfo.template_name,
            network=network,
            memory=memory,
            cpu=vappInfo.cpu,
            storage_profile=storage_profile)

        task = client.get_task_monitor().wait_for_status(
            task=result.Tasks.Task[0],
            #TODO timeout configurable
            timeout=60,
            poll_frequency=2,
            fail_on_statuses=None,
            expected_target_statuses=[
                TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                TaskStatus.CANCELED
            ],
            callback=task_callback)

        st = task.get('status')
        if st == TaskStatus.SUCCESS.value:
            cresult.created = True
        else:
            raise errors.VCDVappCreationError(
                etree.tostring(task, pretty_print=True))

    except Exception as e:

        error_message = 'ERROR.. Not Created VApp {0}  {1}'.format(
            vappInfo.name, str(e))
        logging.warn(error_message, e)
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)
    return cresult


def read(client, context, vapp_info):
    logging.debug('__INIT__vapp_read [ {0} ]'.format(vapp_info))
    try:
        cresult = vapp_pb2.ReadVAppResult()
        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)
        v = org.get_vdc(vapp_info.vdc)
        vdc = VDC(client, href=v.get('href'))
        try:
            vapp = vdc.get_vapp(vapp_info.name)
        except Exception as e:
            logging.error("error occured in vapp", e)
            cresult.present = False
            return cresult

        logging.debug('vapp read  [ %s ]', vapp)

        if vapp:

            nconfig_section = vapp.NetworkConfigSection
            nconfig = nconfig_section.findall(
                "{http://www.vmware.com/vcloud/v1.5}NetworkConfig")
            for i in range(0, len(nconfig)):
                print(nconfig[i].get('networkName'))
                cresult.network.append(nconfig[i].get('networkName'))

            cresult.name = vapp_info.name
            cresult.present = True

        else:
            cresult.present = False
            logging.debug('__LOG__ vapp not found  [ {0} ] '.format(vapp_info))

        logging.debug('__DONE__vapp_read [ {0} ][{1}] '.format(
            cresult, vapp_info))
        return cresult
    except Exception as e:
        error_message = 'ERROR........ {0}'.format(e)
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)


def delete(client, vappInfo):
    logging.debug("===INIT delete vapp called === \n")

    try:
        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)
        v = org.get_vdc(vappInfo.vdc)

        vdc = VDC(client, href=v.get('href'))

        cresult = vapp_pb2.DeleteVAppResult()
        cresult.deleted = False
        try:
            logging.info(" calling delete_vapp............")
            result = vdc.delete_vapp(name=vappInfo.name, force=True)

            logging.info(" delete proceeding............")

            task = client.get_task_monitor().wait_for_status(
                task=result,
                timeout=60,
                poll_frequency=2,
                fail_on_statuses=None,
                expected_target_statuses=[
                    TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
                    TaskStatus.CANCELED
                ],
                callback=task_callback)

            if task.get('status') == TaskStatus.SUCCESS.value:
                cresult.deleted = True

        except Exception as e:
            logging.exception("\n Not Deleted VApp [" + vappInfo.name + "]")

        return cresult

    except Exception as e:
        logging.exception("error occured create vapp ", e)


def update(client, context, vappInfo):
    logging.info("__INIT__update[Vapp]")

    cresult = vapp_pb2.UpdateVAppResult()
    cresult.updated = False

    print("Vapp[ {0} ], Vdc[ {1} ]".format(vappInfo.name, vappInfo.vdc))

    org_resource = client.get_org()
    org = Org(client, resource=org_resource)
    try:
        vdc_resource = org.get_vdc(vappInfo.vdc)
        vdc = VDC(client, name=vappInfo.vdc, resource=vdc_resource)
        vapp_resource = vdc.get_vapp(vappInfo.name)
        vapp = VApp(client, name=vappInfo.name, resource=vapp_resource)
        resp = None
        if (vappInfo.power_on == True):
            logging.info("Powering on [Vapp %v]".format(vappInfo.name))
            resp = vapp.power_on()
        else:
            logging.info("Powering off [Vapp %v]".format(vappInfo.name))
            resp = vapp.undeploy()

        task = client.get_task_monitor().wait_for_status(
            task=resp,
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
            message = 'status : {0} '.format(st)
            logging.info(message)
            cresult.updated = True
        else:
            raise errors.VappUpdateError(
                etree.tostring(task, pretty_print=True))
    except Exception as e:
        error_message = '__ERROR_updating [Vapp] power_on={0} for Vapp {1} . __ErrorMessage__ {2}'.format(
            vappInfo.power_on, vappInfo.name, str(e))
        logging.warn(error_message)
        cresult.updated = False
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)
    logging.info("__DONE__update[Vapp]")
    return cresult
