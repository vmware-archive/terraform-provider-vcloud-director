#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

from concurrent import futures
import sys
import time

import grpc


import pyvcloudprovider_pb2_grpc
import pyvcloudprovider_pb2

from grpc_health.v1.health import HealthServicer
from grpc_health.v1 import health_pb2, health_pb2_grpc


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

class PyVcloudProviderServicer(pyvcloudprovider_pb2_grpc.PyVcloudProviderServicer):
    """Implementation of PyVcloudProviderServicer service."""
    def __init__(self,pyPluginServer):
            self.py_plugin_server=pyPluginServer

    def isPresentCatalog(self, request, context):
        return catalog.isPresent(self.client,request.name)

    def Login(self, request, context):
        print("Login Called")
        logging.info("Login Called")
        resp = "GOT LOGIN CRED = "+request.username
        #resp = resp +" "+ request.password
        resp = resp +" "+ request.org + " URL "+ request.ip +"  "
        result = pyvcloudprovider_pb2.LoginResult()
        result.token = resp
        self.client=login.vcdlogin( request.ip,request.username,request.password,request.org)
        return result
    
    def CreateCatalog(self, request, context):
        return catalog.create(self.client,request.name,request.description,request.shared)
        
    def DeleteCatalog(self, request, context):
        return catalog.delete(self.client,request.name)

    def CatalogUploadMedia(self, request, context):
	#here the request object is CatalogUploadMediaInfo of the protoc / python definition
        return catalog_item.upload_media(self.client,request.catalog_name,request.file_path,item_name=request.item_name)

    def CatalogUploadOva(self, request, context):
    #here the request object is CatalogUploadOvaInfo of the protoc / python definition
        return catalog_item.upload_ova(self.client,request.catalog_name,request.file_path,item_name=request.item_name)

    def OvaCheckResolved(self, request, context):
    #here the request object is CatalogUploadOvaInfo of the protoc / python definition
        ovaInfo=pyvcloudprovider_pb2.CatalogUploadOvaInfo()
        ovaInfo.item_name=request.item_name
        ovaInfo.catalog_name=request.catalog_name
        ovaInfo.file_path=request.file_path
        return catalog_item.ova_check_resolved(self.client,ovaInfo)

    
    def DeleteCatalogItem(self, request, context):
    #here the request object is CatalogUploadOvaInfo of the protoc / python definition
        return catalog_item.delete(self.client,request.catalog_name,item_name=request.item_name)


    def isPresentCatalogItem(self, request, context):
    #here the request object is CatalogUploadOvaInfo of the protoc / python definition
        return catalog_item.isPresent(self.client,request.catalog_name,item_name=request.item_name)

    # VAPP Impl
    def CreateVApp(self, request, context):

        vappInfo=pyvcloudprovider_pb2.CreateVAppInfo()
        vappInfo.name=request.name
        vappInfo.catalog_name=request.catalog_name
        vappInfo.template_name=request.template_name
        return vapp.create(self.client,vappInfo)
        
    def DeleteVApp(self, request, context):
        delVappInfo=pyvcloudprovider_pb2.DeleteVAppInfo()
        delVappInfo.name=request.name
        
        return vapp.delete(self.client,delVappInfo)




    def StopPlugin(self, request, context):
         self.py_plugin_server.stop()
         return pyvcloudprovider_pb2.StopResult()

#client,"c3","/home/ws2/tiny.ova",item_name="item2"

class PyPluginServer:
    """My class"""
    def serve(self):


        # We need to build a health service to work with go-plugin
        logging.basicConfig(level=logging.DEBUG)

        if self.isRunning():
            print("1|1|tcp|127.0.0.1:1234|grpc")    #DO NOT REMOTE THIS LINE
            sys.stdout.flush()                      #DO NOT REMOTE THIS LINE
            logging.info("Already running / Not starting new server")
            return 

        health = HealthServicer()
        health.set("plugin", health_pb2.HealthCheckResponse.ServingStatus.Value('SERVING'))

        # Start the server.
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        
        self.server=server
        self.e = threading.Event()

        #print (inspect.getsourcefile(server.start))

        pyvcloudprovider_pb2_grpc.add_PyVcloudProviderServicer_to_server(PyVcloudProviderServicer(self), server)
        health_pb2_grpc.add_HealthServicer_to_server(health, server)
        server.add_insecure_port('127.0.0.1:1234')
        ## check before start if server is running at 1234
        ## if server running connect to the server and execute stop
        server.start()
        
        # Output information
        print("1|1|tcp|127.0.0.1:1234|grpc")
        sys.stdout.flush()

        self.e.wait() # This event if expected from the pluger_stop functionality
        self.server.stop(0)
       
    def stop(self):
            logging.info(" Received STOP event ...Initiating Server stop ")
            self.e.set()

    def isRunning(self):

            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.settimeout(2)                                      #2 Second Timeout
            result = sock.connect_ex(('127.0.0.1',1234))
            if result == 0:
                
                return True
                #cl=plugin_client.PyPluginClient()
                #cl.stopRemote() 
                #self.check_and_kill()# this call is just a precaution in case 
            else:
                
                return False


if __name__ == '__main__':
    pysserver=PyPluginServer()
    pysserver.serve()


