#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

export TF_ACC=1
export TF_LOG=TRACE

export VCD_ALLOW_UNVERIFIED_SSL=true
export VCD_IP="10.112.83.27"


export VCD_USER="user1"
export VCD_PASSWORD="Admin!23"
export VCD_ORG="O1"


export TF_VAR_VAPP_NAME="acctes_vapp02"
export TF_VAR_VAPP_TEMPLATE_NAME="cento_tmpl1"
export TF_VAR_VAPP_CATALOG_NAME="c1"
export TF_VAR_VAPP_VDC="OVD2"

export TF_VAR_VAPP_NETWORK="ORGNET1"
export TF_VAR_VAPP_MEMORY="8000"
export TF_VAR_VAPP_CPU="4"
export TF_VAR_IP_ALLOCATION_MODE="static"

go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/ -v -run TestAccResourceVApp | grep --line-buffered -vE 'DEBUG|TRACE|terraform|^$'


#MAX INT GO 18446744073709551615