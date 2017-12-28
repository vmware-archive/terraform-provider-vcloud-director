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
import pyvcloudprovider_pb2_grpc
import pyvcloudprovider_pb2
import requests
import logging
from pyvcloud.vcd.vdc import VDC
from pyvcloud.vcd.client import TaskStatus
import click

def spinning_cursor():
    while True:
        for cursor in '|/-\\':
            yield cursor


spinner = spinning_cursor()
def task_callback(task):
        message = '\x1b[2K\r{0}: {1}, status: {2}'.format(
        task.get('operationName'), task.get('operation'), task.get('status')
        )
        if hasattr(task, 'Progress'):
                message += ', progress: %s%%' % task.Progress
        if task.get('status').lower() in [TaskStatus.QUEUED.value,
                                      TaskStatus.PENDING.value,
                                      TaskStatus.PRE_RUNNING.value,
                                      TaskStatus.RUNNING.value]:
                message += ' %s ' % next(spinner)
        click.secho(message, nl=False)
        ##PRINT()
        ##LOGGING.INFO()  REMOVE THE CLICK DEPENDENCY

def create(client, vappInfo):
        logging.debug("===INIT create vapp called === \n")

        try:
                logged_in_org = client.get_org()
                org = Org(client, resource=logged_in_org)
                v = org.get_vdc("OVD2")

                vdc = VDC(client, href=v.get('href'))

                ##WARN FIX the Hard coding
                logging.warn("VDC hard coded !!!!! FIX THIS!!!")



                cresult = pyvcloudprovider_pb2.CreateVAppResult()
                cresult.created = False
                try:
                        logging.info(" calling instantiate_vapp............")
                        result=vdc.instantiate_vapp(
                                        name=vappInfo.name,
                                        catalog=vappInfo.catalog_name,
                                        template=vappInfo.template_name)

                        logging.info(" ............")


                        task = client.get_task_monitor().wait_for_status(
                            task=result.Tasks.Task[0],
                            timeout=60,
                            poll_frequency=2,
                            fail_on_status=None,
                            expected_target_statuses=[
                                TaskStatus.SUCCESS,
                                TaskStatus.ABORTED,
                                TaskStatus.ERROR,
                                TaskStatus.CANCELED],
                                callback=task_callback)

                        if  task.get('status') == TaskStatus.SUCCESS.value :
                                cresult.created = True

                except Exception as e:
                        logging.exception("\n Not Created VApp [" +  vappInfo.name + "]")
                        print(e)
                return cresult

        except Exception as e:
                logging.exception("error occured create vapp ", e)




def delete(client, vappInfo):
        logging.debug("===INIT delete vapp called === \n")

        try:
                logged_in_org = client.get_org()
                org = Org(client, resource=logged_in_org)
                v = org.get_vdc("OVD2")

                vdc = VDC(client, href=v.get('href'))

                ##WARN FIX the Hard coding
                logging.warn("VDC hard coded !!!!! FIX THIS!!!")



                cresult = pyvcloudprovider_pb2.DeleteVAppResult()
                cresult.deleted = False
                try:
                        logging.info(" calling delete_vapp............")
                        result=vdc.delete_vapp(
                                        name=vappInfo.name,
                                        force=True)

                        logging.info(" delete proceeding............")


                        task = client.get_task_monitor().wait_for_status(
                            task=result,
                            timeout=60,
                            poll_frequency=2,
                            fail_on_status=None,
                            expected_target_statuses=[
                                TaskStatus.SUCCESS,
                                TaskStatus.ABORTED,
                                TaskStatus.ERROR,
                                TaskStatus.CANCELED],
                                callback=task_callback)

                        if  task.get('status') == TaskStatus.SUCCESS.value :
                                cresult.deleted = True

                except Exception as e:
                        logging.exception("\n Not Deleted VApp [" +  vappInfo.name + "]")
                        
                return cresult

        except Exception as e:
                logging.exception("error occured create vapp ", e)