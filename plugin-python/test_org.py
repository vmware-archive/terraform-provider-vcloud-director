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
from pyvcloud.vcd.system import System

from proto import pyvcloudprovider_pb2 as pyvcloudprovider_pb2
from proto import pyvcloudprovider_pb2_grpc as pyvcloudprovider_pb2_grpc

import requests
import logging
import catalog
import vapp
import catalog_item

from vcd_cli.profiles import Profiles

from lxml import objectify, etree
import lxml
from pyvcloud.vcd.client import TaskStatus


def vcd_login_with_token(host):
    logging.basicConfig(level=logging.DEBUG)
    logging.info("INIT vcd_login_with_token")
    profiles = Profiles.load()
    token = profiles.get('token')
    client = Client(
        host,
        api_version="27.0",
        verify_ssl_certs=False,
        log_file='vcd.log',
        log_requests=True,
        log_headers=True,
        log_bodies=True)
    try:
        client.rehydrate_from_token(token)
        x = client._session.headers['x-vcloud-authorization']
        logging.info(" =====  X VCloud ========\n  \n" + x + "\n \n")
        return client
    except Exception as e:
        print('error occured', e)
        raise


def capture_vapp(client):
    logging.info('ok')

    logged_in_org = client.get_org()
    org = Org(client, resource=logged_in_org)
    v = org.get_vdc('OVD2')
    if v is None:
        raise errors.VDCNotFoundError(vappInfo.vdc)
    vdc = VDC(client, href=v.get('href'))
    vapp = vdc.get_vapp('vappacc8')
    logging.info(vapp.get('href'))
    catalog = org.get_catalog('c1')
    org.capture_vapp(
        catalog_resource=catalog,
        vapp_href=vapp.get('href'),
        catalog_item_name='captured',
        description='captured',
        customize_on_instantiate=False)


def task_callback(task):
    message = '{0}: {1}, status: {2}'.format(
        task.get('operationName'), task.get('operation'), task.get('status'))
    logging.info(message)


def create_disk(client):
    logging.info('INIT create_disk')

    logged_in_org = client.get_org()
    org = Org(client, resource=logged_in_org)
    v = org.get_vdc('OVD2')
    if v is None:
        raise errors.VDCNotFoundError(vappInfo.vdc)
    vdc = VDC(client, href=v.get('href'))

    #task=vdc.delete_disk('disk01')

    result = vdc.add_disk('disk02', "100")
    logging.info('DONE create_disk %s', result)

    task = client.get_task_monitor().wait_for_status(
        task=result.Tasks.Task[0],
        timeout=60,
        poll_frequency=2,
        fail_on_status=None,
        expected_target_statuses=[
            TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
            TaskStatus.CANCELED
        ],
        callback=task_callback)
    st = task.get('status')
    if st == TaskStatus.SUCCESS.value:
        logging.info("__LOG__ created DISK")
        #cresult.created = True
    else:
        raise errors.VCDDiskCreationError(
            etree.tostring(task, pretty_print=True))


def list_disks(client):
    logging.info('INIT create_disk')

    logged_in_org = client.get_org()
    org = Org(client, resource=logged_in_org)
    v = org.get_vdc('OVD2')
    if v is None:
        raise errors.VDCNotFoundError(vappInfo.vdc)
    vdc = VDC(client, href=v.get('href'))
    #disks = vdc.get_disk('disk02')

    #print(etree.tostring(disks, pretty_print=True))
    logging.info(disks.get('href'))

    vdc.delete_disk('disk02')


def delete_disks(client):
    logging.info('INIT create_disk')

    logged_in_org = client.get_org()
    org = Org(client, resource=logged_in_org)
    v = org.get_vdc('OVD2')
    if v is None:
        raise errors.VDCNotFoundError(vappInfo.vdc)
    vdc = VDC(client, href=v.get('href'))
    disks = vdc.get_disks()

    print(etree.tostring(disks[0], pretty_print=True))
    #logging.info(disks[].get('href'))
    result = vdc.delete_disk('disk02')
    task = client.get_task_monitor().wait_for_status(
        task=result,
        timeout=60,
        poll_frequency=2,
        fail_on_status=None,
        expected_target_statuses=[
            TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
            TaskStatus.CANCELED
        ],
        callback=task_callback)
    st = task.get('status')
    if st == TaskStatus.SUCCESS.value:
        logging.info("__LOG__ created DISK")
        #cresult.created = True
    else:
        raise errors.VCDDiskDeletionError(
            etree.tostring(task, pretty_print=True))


def create_org(client):
    logging.info("create org %s", str(client))
    sys_admin = client.get_admin()
    system = System(client, admin_resource=sys_admin)
    task = system.create_org("O2", "O2 ORG")
    print(type(task))


if __name__ == '__main__':

    client = vcd_login_with_token("10.112.83.27")
    #create_org(client)
    #capture_vapp(client)
    # create_disk(client)
    delete_disks(client)
