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

from vcd_cli.profiles import Profiles

from lxml import objectify, etree
import lxml
import errors


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
    # We need to build a health service to work with go-plugin
    logging.basicConfig(level=logging.DEBUG)
    file = open('task.xml')
    xslt_content = file.read()
    file.close()

    logging.info("parsing %s ", xslt_content)

    #tree = etree.XML(xslt_content)

    tree = etree.XML(xslt_content)
    error_node = tree.findall("{http://www.vmware.com/vcloud/v1.5}Error")
    msg = error_node[0].get('message')
    logging.info(msg)

    logging.info(etree.tostring(tree, pretty_print=True))
    e = errors.VCDVappCreationError("l")
    raise e
