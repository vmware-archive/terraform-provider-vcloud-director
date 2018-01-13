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

export TF_VAR_CATALOG_NAME="test_acc_cata1" #DESTINATION CATALOG
export TF_VAR_CATALOG_DESCRIPTION="accc_ description1" #DESTINATION



export TF_VAR_OVA_PATH="/Users/srinarayana/vmws/tiny.ova"


go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/ -v -run TestAccResourceCatalogItemOva | grep --line-buffered -vE 'TRACE|terraform|^$'
status=${PIPESTATUS[0]} 


#go test -v 2>&1 | go-junit-report > report.xml
echo "test_ova.sh EXIT STATUS = " $status

exit $status