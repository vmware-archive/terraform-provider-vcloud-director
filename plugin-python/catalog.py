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


def isPresent(client, name):
        logging.basicConfig(level=logging.DEBUG)
        logging.debug("=== isPresent called === \n")

        try:
                logged_in_org = client.get_org()
                org = Org(client, resource=logged_in_org)
                result = pyvcloudprovider_pb2.IsPresentCatalogResult()
                result.present = False
                try:
                        catalog = org.get_catalog(name)
                        result.present = True
                except Exception as e:
                        logging.info("catalog is not present")
                return result

        except Exception as e:
                logging.warn("=====>>>>>>>  error occured", e)


def create(client, name, description,shared):
        logging.debug("=== create catalog called === \n")

        try:
                logged_in_org = client.get_org()
                org = Org(client, resource=logged_in_org)
                result = pyvcloudprovider_pb2.CreateCatalogResult()
                result.created = False
                try:
                        catalog = org.create_catalog(name=name, description=description)
                        result.created = True
                except Exception as e:
                        logging.info("\n Not Created catalog [" + name + "]")
                return result

        except Exception as e:
                print(e)
                print("\n")
                logging.warn("error occured", e)


def delete(client, name):
        logging.debug("=== delete catalog called === \n")

        try:
                logged_in_org = client.get_org()
                org = Org(client, resource=logged_in_org)
                result = pyvcloudprovider_pb2.DeleteCatalogResult()
                result.deleted = False
                try:
                        catalog = org.delete_catalog(name)
                        result.deleted = True
                except Exception as e:
                        logging.info("\n Not Deleted  catalog [" + name + "]")
                return result

        except Exception as e:
                logging.warn("error occured", e)


def upload_media(client, catalog_name, file_name, item_name):
        logging.debug("===== upload_media to ++catalog called === \n")
        result=pyvcloudprovider_pb2.CatalogUploadMediaResult()
        result.created=False

        try:
                logged_in_org = client.get_org()
                org = Org(client, resource=logged_in_org)
                org.upload_media(catalog_name=catalog_name,
                                 file_name=file_name, item_name=item_name)
                result.created=True
                return result;
        except Exception as e:
                logging.warn("error occured", e)
                return result;
