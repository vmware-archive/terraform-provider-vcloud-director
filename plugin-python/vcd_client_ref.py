#!/bin/bash
#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

class VCDClientRef:

    instance = None

    def set_ref(self, arg):
        VCDClientRef.instance = arg

    def get_ref(self):
        return VCDClientRef.instance
