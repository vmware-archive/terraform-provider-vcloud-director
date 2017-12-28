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
from pyvcloud.vcd.client import QueryResultFormat

import pyvcloudprovider_pb2_grpc
import pyvcloudprovider_pb2
import requests
import logging
import time



def isPresent(client,catalog_name, item_name):
        logging.basicConfig(level=logging.DEBUG)
        logging.debug("=== isPresent catalog item called === \n")

        try:
                logged_in_org = client.get_org()
                org = Org(client, resource=logged_in_org)
                result = pyvcloudprovider_pb2.IsPresentCatalogItemResult()
                result.present = False
                try:
                        catalog = org.get_catalog_item(catalog_name,item_name)
                        result.present = True
                except Exception as e:
                        logging.info("catalog is not present")
                return result

        except Exception as e:
                logging.warn("=====>>>>>>>  error occured", e)

def upload_media(client, catalog_name, file_name, item_name):
        logging.basicConfig(level=logging.DEBUG)
        logging.debug("===== upload_media to  catalog called === \n")
        result=pyvcloudprovider_pb2.CatalogUploadMediaResult()
        result.created=False

        try:
                logged_in_org = client.get_org()
                org = Org(client, resource=logged_in_org)
                org.upload_media(catalog_name=catalog_name,
                                 file_name=file_name, item_name=item_name)
                result.created=True
                logging.debug("===== return created true  ")
                return result;
        except Exception as e:
                logging.warn("error occured", e)
                return result;

def upload_ova(client, catalog_name, file_name, item_name):
        logging.debug("===== INIT upload ova to catalog %s  called === %s item %s \n",catalog_name,file_name,item_name)
        cresult=pyvcloudprovider_pb2.CatalogUploadOvaResult()
        cresult.created=False

        try:
                logged_in_org = client.get_org()
                org = Org(client, resource=logged_in_org)
                result=org.upload_ovf(catalog_name=catalog_name,
                                 file_name=file_name, item_name=item_name)

                logging.info("result ---- "+str(result))
                cresult.created = True
                
                logging.debug("===== DONE upload ova to  catalog %s   === \n",catalog_name)
                return cresult;
        except Exception as e:
                logging.warn("error occured", e)
                return cresult;

def ova_check_resolved(client, ovaInfo):
        logging.debug("===== INIT ova_check_resolved %s  called === \n",ovaInfo.item_name)
        cresult=pyvcloudprovider_pb2.CheckResolvedResult()
        cresult.resolved=False

        try:
                logged_in_org = client.get_org()
                org = Org(client, resource=logged_in_org)
                source_ova_item = org.get_catalog_item(ovaInfo.catalog_name,ovaInfo.item_name)
                if source_ova_item is None:
                    return None
                

                check_resolved(client,source_ova_item)
                
                cresult.resolved = True
                
                logging.debug("===== DONE upload ova to  catalog %s   === \n",ovaInfo.catalog_name)
                return cresult;
        except Exception as e:
                logging.warn("error occured", e)
                raise
                




def check_resolved(client,source_ova_item):
        logging.info("INIT check_resolved")
        item_id = source_ova_item.get('id')
        logging.info("[LOG] check_resolved item_id %s ", item_id)
        while True:
                q = client.get_typed_query(
                        'catalogItem',
                        query_result_format=QueryResultFormat.ID_RECORDS,
                        qfilter='id==%s' % item_id)
                records = list(q.execute())
                if records[0].get('status') == 'RESOLVED':
                    
                        logging.info("Completed item import")
                        break;
                    
                else:
                        logging.info('Waiting for upload to complete...')
                                    
                logging.info('===_LOG_===sleeping 5 secs...\n\n')        
                time.sleep(5)

        logging.info("DONE check_resolved")





def delete(client, catalog_name, item_name):
        logging.debug("===== delete ova to  %s catalog called === \n",catalog_name)
        result=pyvcloudprovider_pb2.DeleteCatalogItemResult()
        result.deleted=False

        try:
                logged_in_org = client.get_org()
                org = Org(client, resource=logged_in_org)
                org.delete_catalog_item(name=catalog_name,
                                  item_name=item_name)
                result.deleted=True
                return result;
        except Exception as e:
                logging.warn("error occured", e)
                return result;


##DELETE /api/catalogItem/


