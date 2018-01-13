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
export VCD_ORG="O1"

export TF_VAR_CATALOG_NAME="test_acc_cata1" #DESTINATION CATALOG
export TF_VAR_CATALOG_DESCRIPTION="accc_ description1" #DESTINATION




export TF_VAR_SOURCE_VAPP_NAME="vappacc8" #VAPP THAT IS CAPTURED
export TF_VAR_SOURCE_VDC_NAME="OVD2" #VDC OF VAPP THAT IS CAPTURED

echo ' TEST CAPTURE VAPP'

go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/ -v -run TestAccResourceCatalogItemCaptureVapp | grep --line-buffered -vE 'TRACE|terraform|^$'
