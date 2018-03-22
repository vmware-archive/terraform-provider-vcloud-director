#!/bin/bash
#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************
export GOPATH=/home/parkash/workspace/vmware/Github/terraform-provider-vcloud-director/go
export PY_PLUGIN='python3 /home/parkash/workspace/vmware/Github/terraform-provider-vcloud-director/plugin-python/plugin.py'
export PY_PLUGIN_STOP='python3 /home/parkash/workspace/vmware/Github/terraform-provider-vcloud-director/plugin-python/plugin_stop.py'
export PATH=$PATH:$GOPATH/bin
