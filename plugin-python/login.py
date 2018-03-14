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

import requests
import logging
import catalog
import traceback
import grpc
from vcd_cli.profiles import Profiles

from vcd_client_ref import VCDClientRef

import logging


def vcd_login(context, lc):
    logging.info("__INIT__vcd_login [%s]", lc)
    login_path = ""
    client = Client(
        lc.ip,
        api_version="27.0",
        verify_ssl_certs=not lc.allow_insecure_flag,
        log_file='vcd.log',
        log_requests=True,
        log_headers=True,
        log_bodies=True)
    try:
        if lc.use_vcd_cli_profile:
            login_path = "rehydrate_from_token"
            profiles = Profiles.load()
            token = profiles.get('token')
            client.rehydrate_from_token(token)

        else:
            login_path = "set_credentials"
            client.set_credentials(
                BasicLoginCredentials(lc.username, lc.org, lc.password))

        vref = VCDClientRef()
        vref.set_ref(client)
        logging.debug('LOGIN VIA :[{0}] ARG :[{1}]'.format(
            login_path, str(lc)))

        return client
    except Exception as e:
        traceback.print_exc()
        error_message = 'ERROR IN LOGIN .. Exception [{0}] LOGIN VIA :[{1}] ARG :[{2}]'.format(
            str(e), login_path, str(lc))
        error_message = error_message.replace('\n', ' ')
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)
        raise
