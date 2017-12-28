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
    except Exception as e:
        print('error occured', e)


def vcd_login(host, user, password, org):
    logging.basicConfig(level=logging.DEBUG)
    logging.info("login called")
    client = Client(
        host,
        api_version="27.0",
        verify_ssl_certs=False,
        log_file='vcd.log',
        log_requests=True,
        log_headers=True,
        log_bodies=True)
    try:
        client.set_credentials(BasicLoginCredentials(user, org, password))
        x = client._session.headers['x-vcloud-authorization']
        logging.info(" =====  X VCloud ========\n  \n" + x + "\n \n")

        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)
        v = org.get_vdc('OVD2')
        vdc = VDC(client, href=v.get('href'))
        vapp = vdc.get_vapp('cento_vapp11_2')
        nconfig_section = vapp.NetworkConfigSection
        nconfig = nconfig_section.findall(
            "{http://www.vmware.com/vcloud/v1.5}NetworkConfig")
        for i in range(0, len(nconfig)):
            print(nconfig[i].get('networkName'))

        ##print(etree.tostring(nconfig[0], pretty_print=True))

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
        #ovaInfo=pyvcloudprovider_pb2.CatalogUploadOvaInfo()
        #ovaInfo.item_name="item1"
        #ovaInfo.catalog_name="testcata1"
        #ovaInfo.file_path="/Users/srinarayana/vmws/tiny.ova"

        # catalog_item.ova_check_resolved(client,ovaInfo)

        #vappInfo.template_name="cento_tmpl1"
        #vapp.create(client,vappInfo)
        return client
    except Exception as e:
        print('error occured', e)


def print_element(root, i=0):
    print("\n")
    for j in range(0, i):
        print("\t", end='')
    for e in root.getchildren():

        print(e.tag, end='')
        #print (e.tag,end='')
        #print (type(e),end='')
        print_element(e, i + 1)


if __name__ == '__main__':
    vcd_login("10.112.83.27", "user1", "Admin!23", "O1")
    #vcd_login_with_token("10.112.83.27")
