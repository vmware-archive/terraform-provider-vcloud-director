#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

from concurrent import futures
import sys
import time

import grpc

from proto import pyvcloudprovider_pb2 as pyvcloudprovider_pb2
from proto import pyvcloudprovider_pb2_grpc as pyvcloudprovider_pb2_grpc

from proto import org_pb2 as org_pb2
from proto import org_pb2_grpc as org_pb2_grpc

from proto import user_pb2 as user_pb2
from proto import user_pb2_grpc as user_pb2_grpc

from proto import vdc_pb2 as vdc_pb2
from proto import vdc_pb2_grpc as vdc_pb2_grpc

from proto import vapp_vm_pb2 as vapp_vm_pb2
from proto import vapp_vm_pb2_grpc as vapp_vm_pb2_grpc

from proto import vapp_pb2 as vapp_pb2
from proto import catalog_item_pb2 as catalog_item_pb2

from grpc_health.v1.health import HealthServicer
from grpc_health.v1 import health_pb2, health_pb2_grpc
from proto import disk_pb2_grpc as disk_pb2_grpc

# local impls
import login
import logging
import catalog
import catalog_item
import vapp

#pyvcloud import
from pyvcloud.vcd.org import Org

#debug
import inspect

#Packages imported for remote terminating
import threading
import socket
import plugin_stop
from disk import IndependentDiskServicer
from org import OrgServicer
from user import UserServicer
from vdc import VdcServicer

from vapp_vm import VappVmServicer


class PyVcloudProviderServicer(
        pyvcloudprovider_pb2_grpc.PyVcloudProviderServicer):
    """Implementation of PyVcloudProviderServicer service."""

    def __init__(self, pyPluginServer):
        self.py_plugin_server = pyPluginServer

    def ReadCatalog(self, request, context):
        return catalog.read(self.client, context, request.name)

    def Login(self, request, context):
        lc = pyvcloudprovider_pb2.LoginCredentials()
        lc.ip = request.ip

        lc.username = request.username
        lc.password = request.password
        lc.org = request.org

        lc.use_vcd_cli_profile = request.use_vcd_cli_profile
        lc.allow_insecure_flag = request.allow_insecure_flag

        result = pyvcloudprovider_pb2.LoginResult()

        self.client = login.vcd_login(context, lc)

        x = self.client._session.headers['x-vcloud-authorization']
        result.token = x

        return result

    def CreateCatalog(self, request, context):
        logging.info("=========================[%s]", request.description)
        logging.info("=========================[%s]", request.name)
        return catalog.create(self.client, context, request.name,
                              request.description, request.shared)

    def DeleteCatalog(self, request, context):
        return catalog.delete(self.client, context, request.name)

    def UpdateCatalog(self, request, context):
        logging.info("=========================[%s]", request.description)
        logging.info("=========================[%s]", request.name)
        return catalog.update(self.client, context, request.old_name,
                              request.name, request.description,
                              request.shared)

    def CatalogUploadMedia(self, request, context):
        #here the request object is CatalogUploadMediaInfo of the protoc / python definition
        return catalog_item.upload_media(
            self.client,
            context,
            request.catalog_name,
            request.file_path,
            item_name=request.item_name)

    def CatalogUploadOva(self, request, context):
        #here the request object is CatalogUploadOvaInfo of the protoc / python definition
        #client, context,catalog_name, file_name, item_name
        return catalog_item.upload_ova(
            client=self.client,
            context=context,
            catalog_name=request.catalog_name,
            file_name=request.file_path,
            item_name=request.item_name)

    def CaptureVapp(self, request, context):
        #here the request object is CatalogUploadOvaInfo of the protoc / python definition
        #ccapture_vapp(client,vapp_name,vdc_name,catalog_name, item_name,desc,customize_on_instantiate=False):
        logging.info("__INIT__CaptureVapp_plugin.py")
        capture_info = catalog_item_pb2.CaptureVAppInfo()
        capture_info.catalog_name = request.catalog_name
        capture_info.item_name = request.item_name
        capture_info.vdc_name = request.vdc_name
        capture_info.vapp_name = request.vapp_name

        capture_info.description = request.description
        capture_info.customize_on_instantiate = request.customize_on_instantiate
        logging.debug("__LOG__CaptureVapp [%s]", capture_info)
        return catalog_item.capture_vapp(
            client=self.client, context=context, capture_info=capture_info)

    def OvaCheckResolved(self, request, context):
        #here the request object is CatalogUploadOvaInfo of the protoc / python definition
        ovaInfo = catalog_item_pb2.CatalogCheckResolvedInfo()
        ovaInfo.item_name = request.item_name
        ovaInfo.catalog_name = request.catalog_name

        return catalog_item.ova_check_resolved(
            client=self.client, context=context, ovaInfo=ovaInfo)

    def DeleteCatalogItem(self, request, context):
        #here the request object is CatalogUploadOvaInfo of the protoc / python definition
        return catalog_item.delete(
            self.client, request.catalog_name, item_name=request.item_name)

    def isPresentCatalogItem(self, request, context):
        #here the request object is CatalogUploadOvaInfo of the protoc / python definition
        return catalog_item.is_present(
            self.client, request.catalog_name, item_name=request.item_name)

    # VAPP Impl
    def CreateVApp(self, request, context):

        vapp_info = vapp_pb2.CreateVAppInfo()
        vapp_info.name = request.name
        vapp_info.catalog_name = request.catalog_name
        vapp_info.template_name = request.template_name
        vapp_info.vdc = request.vdc
        vapp_info.network = request.network

        vapp_info.ip_allocation_mode = request.ip_allocation_mode
        vapp_info.memory = request.memory
        vapp_info.cpu = request.cpu
        vapp_info.storage_profile = request.storage_profile
        logging.debug("__LOG__ [%s]", vapp_info)
        return vapp.create(self.client, context, vapp_info)

    def DeleteVApp(self, request, context):
        del_app_info = vapp_pb2.DeleteVAppInfo()
        del_app_info.name = request.name
        del_app_info.vdc = request.vdc
        return vapp.delete(self.client, del_app_info)

    def ReadVApp(self, request, context):
        read_app_info = vapp_pb2.ReadVAppInfo()
        read_app_info.name = request.name
        read_app_info.vdc = request.vdc
        return vapp.read(self.client, context, read_app_info)

    def UpdateVApp(self, request, context):
        update_vapp_info = vapp_pb2.UpdateVAppInfo()
        update_vapp_info.name = request.name
        update_vapp_info.vdc = request.vdc
        update_vapp_info.power_on = request.power_on
        return vapp.update(self.client, context, update_vapp_info)

    def StopPlugin(self, request, context):
        self.py_plugin_server.stop()
        return pyvcloudprovider_pb2.StopResult()


#client,"c3","/home/ws2/tiny.ova",item_name="item2"
global_vcd_login = 0


class PyPluginServer:
    """My class"""

    def serve(self):

        # We need to build a health service to work with go-plugin
        logging.basicConfig(level=logging.DEBUG)

        if self.isRunning():
            print("1|1|tcp|127.0.0.1:1234|grpc")  #DO NOT REMOTE THIS LINE
            sys.stdout.flush()  #DO NOT REMOTE THIS LINE
            logging.info("Already running / Not starting new server")
            return

        health = HealthServicer()
        health.set(
            "plugin",
            health_pb2.HealthCheckResponse.ServingStatus.Value('SERVING'))

        # Start the server.
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

        self.server = server
        self.e = threading.Event()

        #print (inspect.getsourcefile(server.start))

        pyvcloudprovider_pb2_grpc.add_PyVcloudProviderServicer_to_server(
            PyVcloudProviderServicer(self), server)

        disk_pb2_grpc.add_IndependentDiskServicer_to_server(
            IndependentDiskServicer(self), server)

        org_pb2_grpc.add_OrgServicer_to_server(OrgServicer(self), server)

        user_pb2_grpc.add_UserServicer_to_server(UserServicer(self), server)

        vdc_pb2_grpc.add_VdcServicer_to_server(VdcServicer(self), server)

        vapp_vm_pb2_grpc.add_VappVmServicer_to_server(
            VappVmServicer(self), server)

        health_pb2_grpc.add_HealthServicer_to_server(health, server)
        server.add_insecure_port('127.0.0.1:1234')
        ## check before start if server is running at 1234
        ## if server running connect to the server and execute stop
        server.start()

        # Output information
        print("1|1|tcp|127.0.0.1:1234|grpc")
        sys.stdout.flush()

        self.e.wait(
        )  # This event if expected from the pluger_stop functionality
        self.server.stop(0)

    def stop(self):
        logging.info(" Received STOP event ...Initiating Server stop ")
        self.e.set()

    def isRunning(self):

        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(2)  #2 Second Timeout
        result = sock.connect_ex(('127.0.0.1', 1234))
        if result == 0:

            return True
            #cl=plugin_client.PyPluginClient()
            #cl.stopRemote()
            #self.check_and_kill()# this call is just a precaution in case
        else:

            return False


if __name__ == '__main__':
    pysserver = PyPluginServer()
    pysserver.serve()
