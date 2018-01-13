#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

export TF_ACC=1
export TF_LOG=DEBUG

export VCD_ALLOW_UNVERIFIED_SSL=true
export VCD_IP="10.112.83.27"


 export VCD_USER="user1"
export VCD_PASSWORD="Admin!23"


 #export VCD_USER=""
 #export VCD_PASSWORD=""

export VCD_ORG="O1"

export VCD_USE_VCD_CLI_PROFILE=""

#export VCD_USE_VCD_CLI_PROFILE=true

export TF_VAR_VAPP_CATALOG_NAME="c1"
export TF_VAR_VAPP_VDC="OVD4"

export TF_VAR_VAPP_NETWORK="ORGNET0"
export TF_VAR_VAPP_MEMORY="200"
export TF_VAR_VAPP_CPU="2"
export TF_VAR_VAPP_IP_ALLOCATION_MODE="static"




export TF_VAR_SOURCE_VAPP_NAME="vappacc8"
export TF_VAR_SOURCE_VDC_NAME="OVD2"

export TF_VAR_OVA_PATH="/Users/srinarayana/vmws/tiny.ova"


#MAX INT GO 18446744073709551615

