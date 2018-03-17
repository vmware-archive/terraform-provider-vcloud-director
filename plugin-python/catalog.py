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

from proto import pyvcloudprovider_pb2 as pyvcloudprovider_pb2
from proto import pyvcloudprovider_pb2_grpc as pyvcloudprovider_pb2_grpc
import grpc
import requests
import logging
import inspect


def read(client, context, name):
    #logging.basicConfig(level=logging.DEBUG)
    logging.debug("__INIT__read[catalog]__ %s", name)
    logging.debug("catalog_name  %s", name)

    try:
        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)

        result = pyvcloudprovider_pb2.ReadCatalogResult()
        result.present = False
        try:
            catalog = org.get_catalog(name)
            logging.info(vars(catalog))
            logging.info("\n==desc=[%s]", catalog.Description)

            result.name = str(catalog.get("name"))
            result.description = str(catalog.Description)
            result.shared = catalog.IsPublished
            result.present = True
        except Exception as e:
            logging.warning(
                "__ERROR__ while reading catalog[{0}]. __ERROR_MESSAGE__[{1}]".
                format(name, str(e)))

        return result

    except Exception as e:
        logging.error(
            "__ERROR__ while reading catalog[{0}]. __ERROR_MESSAGE__[{1}]".
            format(name, str(e)))
        raise


def create(client, context, name, description):
    logging.debug("_INIT_create_catalog__")
    logging.debug("name = [%s]  desc = [%s] ", name, description)

    result = pyvcloudprovider_pb2.CreateCatalogResult()
    result.created = False
    try:
        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)

        try:
            catalog = org.create_catalog(name=name, description=description)
            result.created = True
        except Exception as e:
            error_message = "__ ERROR Creating catalog [{0}] __ERROR_MESSAGE__ [{1}]".format(
                name, str(e))
            logging.warn(error_message)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
    except Exception as e:
        error_message = "__ ERROR Creating catalog [{0}] __ERROR_MESSAGE__ [{1}]".format(
            name, str(e))
        logging.warn(error_message)
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)

    logging.debug("_DONE_create_catalog__")
    return result


def update(client, context, old_name, new_name, description):
    logging.info("\n __INIT_update_catalog__")
    logging.debug(
        "\n old_name = [%s] new_name = [%s]  desc = [%s]",
        old_name, new_name, description)
    result = pyvcloudprovider_pb2.UpdateCatalogResult()
    result.updated = False

    try:
        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)
        try:
            catalog = org.update_catalog(
                old_catalog_name=old_name,
                new_catalog_name=new_name,
                description=description)
            result.updated = True
        except Exception as e:
            error_message = "Error occured while updating catalog {0} __ERROR_MESSAGE__ {1}".format(
                old_name, str(e))
            logging.error(error_message)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
    except Exception as e:
        error_message = "Error occured while updating catalog {0} __ERROR_MESSAGE__ {1}".format(
            old_name, str(e))
        logging.error(error_message)
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)

    logging.debug("__DONE_update_catalog__")
    return result


def share_catalog(client, context, name, shared):
    logging.info("\n __INIT_shared_catalog__")
    logging.debug("\n name = [%s] shared = [%s]  ", name, shared)
    result = pyvcloudprovider_pb2.ShareCatalogResult()
    result.success = False
    try:
        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)
        try:
            catalog = org.share_catalog(name=name, share=shared)
            result.success = True
        except Exception as e:
            error_message = "Error occured while sharing catalog[{0}[ to shared = [{1}], __ERROR_MESSAGE__ [{2}[".format(name, shared, str(e))
            logging.error(error_message)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
    except Exception as e:
        error_message = "Error occured while sharing catalog[{0}[ to shared = [{1}], __ERROR_MESSAGE__ [{2}[".format(name, shared, str(e))
        logging.error(error_message)
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)

    logging.info("__DONE_share_catalog__")
    return result


def delete(client, context, name):
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
            error_message = "Error occured while deleting catalog {0} __ERROR_MESSAGE__ {1}".format(
                name, str(e))
            logging.error(error_message)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(error_message)
        return result

    except Exception as e:
        error_message = "Error occured while deleting catalog {0} __ERROR_MESSAGE__ {1}".format(
            name, str(e))
        logging.error(error_message)
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)


def upload_media(client, catalog_name, file_name, item_name):
    logging.debug("===== upload_media to ++catalog called === \n")
    result = pyvcloudprovider_pb2.CatalogUploadMediaResult()
    result.created = False

    try:
        logged_in_org = client.get_org()
        org = Org(client, resource=logged_in_org)
        org.upload_media(
            catalog_name=catalog_name,
            file_name=file_name,
            item_name=item_name)
        result.created = True
        return result
    except Exception as e:
        logging.warn("error occured", e)
        return result
