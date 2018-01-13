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

from pyvcloud.vcd.client import QueryResultFormat

from proto import pyvcloudprovider_pb2 as pyvcloudprovider_pb2
from proto import pyvcloudprovider_pb2_grpc as pyvcloudprovider_pb2_grpc

from proto import catalog_item_pb2 as catalog_item_pb2
import errors

import requests
import logging
import time
import grpc


def is_present(client, catalog_name, item_name):
    logging.basicConfig(level=logging.DEBUG)
    logging.debug("=== isPresent catalog item called === \n")

    try:
        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)
        result = catalog_item_pb2.IsPresentCatalogItemResult()
        result.present = False
        try:
            catalog = org.get_catalog_item(catalog_name, item_name)
            result.present = True
        except Exception as e:
            logging.info("catalog item is not present")
        return result

    except Exception as e:
        logging.warn("__ERROR__ occured", e)


def upload_media(client, context, catalog_name, file_name, item_name):
    logging.basicConfig(level=logging.DEBUG)
    logging.debug("===== upload_media to  catalog called === \n")
    result = catalog_item_pb2.CatalogUploadMediaResult()
    result.created = False

    try:
        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)
        org.upload_media(
            catalog_name=catalog_name,
            file_name=file_name,
            item_name=item_name)
        result.created = True
        logging.debug("===== return created true  ")
        return result
    except Exception as e:
        error_message = 'ERROR........ Upload Media..... {0}'.format(str(e))
        logging.warn(error_message, e)
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)
        return result


def upload_ova(client, context, catalog_name, file_name, item_name):
    logging.debug(
        "===== INIT upload ova to catalog %s  called === %s item %s \n",
        catalog_name, file_name, item_name)
    cresult = catalog_item_pb2.CatalogUploadOvaResult()
    cresult.created = False

    try:
        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)
        result = org.upload_ovf(
            catalog_name=catalog_name,
            file_name=file_name,
            item_name=item_name)

        logging.info("result ---- " + str(result))
        cresult.created = True

        logging.debug("===== DONE upload ova to  catalog %s   === \n",
                      catalog_name)
        return cresult
    except Exception as e:
        error_message = 'ERROR.. upload_ova failed  {0} '.format(catalog_name)
        logging.warn(error_message, e)
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)


def ova_check_resolved(client, context, ovaInfo):
    logging.debug("===== INIT ova_check_resolved %s  called === \n",
                  ovaInfo.item_name)
    cresult = pyvcloudprovider_pb2.CheckResolvedResult()
    cresult.resolved = False

    try:
        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)
        source_ova_item = org.get_catalog_item(ovaInfo.catalog_name,
                                               ovaInfo.item_name)
        if source_ova_item is None:  #TODO WRITE A TEST CASE
            raise error.ItemFoundError('Item = {0}'.format(ovaInfo.item_name))

        check_resolved(client, source_ova_item)

        cresult.resolved = True

        logging.debug(
            "===== DONE upload/capture ova to  catalog %s  %s === \n",
            ovaInfo.catalog_name, ovaInfo.item_name)
        return cresult
    except Exception as e:

        error_message = 'ERROR.. ova_check_resolved {0} '.format(
            ovaInfo.item_name)
        logging.warn(error_message, e)
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)
        raise


def check_resolved(client, source_ova_item):
    logging.info("__INIT__ check_resolved")
    item_id = source_ova_item.get('id')
    logging.info("[LOG] check_resolved item_id %s ", item_id)
    while True:
        q = client.get_typed_query(
            'catalogItem',
            query_result_format=QueryResultFormat.ID_RECORDS,
            qfilter='id==%s' % item_id)
        records = list(q.execute())
        logging.info("__LOG__status = %s", records[0].get('status'))
        if records[0].get('status') == 'RESOLVED':

            logging.info("Completed item import")
            break

        else:
            logging.info('Waiting for upload/capture to complete...')

        logging.info('__LOG__sleeping 5 secs...\n\n')
        time.sleep(5)
        #TODO might have to check when status goes to other state than resolved

    logging.info("__DONE__ check_resolved")


def delete(client, catalog_name, item_name):
    logging.debug("===== delete ova to  %s catalog called === \n",
                  catalog_name)
    result = catalog_item_pb2.DeleteCatalogItemResult()
    result.deleted = False

    try:
        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)
        org.delete_catalog_item(name=catalog_name, item_name=item_name)
        result.deleted = True
        return result
    except Exception as e:
        logging.warn("error occured", e)
        return result


def capture_vapp(client, context, capture_info):
    logging.debug("__INIT__capture_vapp [%s]", capture_info)
    vapp_name = capture_info.vapp_name
    vdc_name = capture_info.vdc_name
    catalog_name = capture_info.catalog_name
    item_name = capture_info.item_name
    desc = capture_info.description
    customize_on_instantiate = capture_info.customize_on_instantiate

    result = catalog_item_pb2.CaptureVAppResult()
    try:
        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)
        v = org.get_vdc(vdc_name)
        if v is None:
            raise errors.VDCNotFoundError(vappInfo.vdc)
        vdc = VDC(client, href=v.get('href'))
        vapp = vdc.get_vapp(vapp_name)
        logging.info(vapp.get('href'))
        catalog = org.get_catalog(catalog_name)
        org.capture_vapp(
            catalog_resource=catalog,
            vapp_href=vapp.get('href'),
            catalog_item_name=item_name,
            description=desc,
            customize_on_instantiate=customize_on_instantiate)
        result.captured = True
        logging.debug("__DONE__capture_vapp [%s] , [%s] , [%s]  ", vapp_name,
                      vdc_name, catalog_name)
        return result
    except Exception as e:
        error_message = 'ERROR.. capture vapp failed  {0} {1} {2}'.format(
            vapp_name, vdc_name, catalog_name)
        logging.warn(error_message, e)
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)
