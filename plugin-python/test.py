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
import pyvcloudprovider_pb2_grpc
import pyvcloudprovider_pb2
import requests
import logging
import catalog
import vapp
import catalog_item


def vcdlogin(  host, user, password, org):
    logging.basicConfig(level=logging.DEBUG)  
    logging.info("login called")
    client = Client(host,api_version="27.0",verify_ssl_certs=False,log_file='vcd.log',log_requests=True, log_headers=True, log_bodies=True)
    try:
        client.set_credentials(BasicLoginCredentials(user, org, password))
        x=client._session.headers['x-vcloud-authorization']
        logging.info(" =====  X VCloud ========\n  \n"+x + "\n \n")


		#pRes= catalog.isPresent(client,"c1")
       # res=catalog.create(client,"c44","c44",True)

		#logging.info(" is create ===== \n \n "+ str(res.created)+ "\n \n ")
        #logging.info("\n\n=====callling upload ova===\n")
        #catalog_item.upload_ova(client,"c44","/home/iso/tiny.ova",item_name="tiny3.ova")
        #logging.info("\n\n=====upload done ===\n ")
		#logging.info(" Delete  ===== \n \n "+ str(catalog.delete(client,"c3").deleted)+ "\n \n ")
		#logging.info(" Delete  ===== \n \n "+ str(catalog.delete(client,"c4").deleted)+ "\n \n ")

        #vappInfo=pyvcloudprovider_pb2.CreateVAppInfo()
        #vappInfo.name="vappacc2"
        #vappInfo.catalog_name="c1"
        #vappInfo.template_name="cento_tmpl1"
        #vapp.create(client,vappInfo)


        #delVappInfo=pyvcloudprovider_pb2.DeleteVAppInfo()
        #delVappInfo.name="vappacc2"
        
        #vapp.delete(client,delVappInfo)
        ovaInfo=pyvcloudprovider_pb2.CatalogUploadOvaInfo()
        ovaInfo.item_name="item1"
        ovaInfo.catalog_name="testcata1"
        ovaInfo.file_path="/Users/srinarayana/vmws/tiny.ova"

        catalog_item.ova_check_resolved(client,ovaInfo)

        return client;
    except Exception as e:
        print('error occured',e)


if __name__ == '__main__':
    vcdlogin("10.112.83.27","user1","Admin!23","O1");
