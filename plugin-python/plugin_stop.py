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

from grpc_health.v1.health import HealthServicer
from grpc_health.v1 import health_pb2, health_pb2_grpc

import login
import logging
import catalog
from pyvcloud.vcd.org import Org
import inspect


class PyPluginClient:
    def stopRemote(self):
        # We need to build a health service to work with go-plugin
        print("stopRemote")
        try:
            channel = grpc.insecure_channel('127.0.0.1:1234')
            stub = pyvcloudprovider_pb2_grpc.PyVcloudProviderStub(channel)
            si = pyvcloudprovider_pb2.StopInfo()

            print("stopping")

            stub.StopPlugin(si)
            print("stop Remote OK")
        except:
            print("Error occured")
            raise


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG)
    pysserver = PyPluginClient()
    pysserver.stopRemote()
