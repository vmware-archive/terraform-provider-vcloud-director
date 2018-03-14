#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

export TF_ACC=1
export TF_LOG=TRACE

export VCD_ALLOW_UNVERIFIED_SSL=true
export VCD_IP="10.172.158.127"

export VCD_USER="acmeadmin"
export VCD_PASSWORD="VMware1!"
export VCD_ORG="ACME"

#. ./test_setlogin.sh


export TF_VAR_VAPP_NAME="pcp_vapp_100"
export TF_VAR_VAPP_TEMPLATE_NAME="tinyova"
export TF_VAR_VAPP_CATALOG_NAME="ACME"
export TF_VAR_VAPP_VDC="ACME_PAYG"

export TF_VAR_VAPP_NETWORK="172.16.1.0"
export TF_VAR_VAPP_MEMORY="64"
export TF_VAR_VAPP_CPU="1"
export TF_VAR_VAPP_IP_ALLOCATION_MODE="dhcp"

go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/ -v -run TestAccResourceVApp | grep --line-buffered -vE 'DEBUG|TRACE|terraform|^$'


#MAX INT GO 18446744073709551615