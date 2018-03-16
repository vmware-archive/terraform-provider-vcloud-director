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

export TF_VAR_CATALOG_NAME_OLD="pcp_test_calatog_1009"
export TF_VAR_CATALOG_NAME_NEW="pcp_test_calatog_10010"

export TF_VAR_CATALOG_DESCRIPTION_1="pcp_description9"
export TF_VAR_CATALOG_DESCRIPTION_2="pcp_description10"

export TF_VAR_CATALOG_SHARED="true"

#go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/ -v -run TestAccResourceCatalog | grep --line-buffered -vE 'TRACE|terraform.|^$'

go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/  -v -run TestAccResourceCatalogBasic | grep --line-buffered -vE 'DEBUG|TRACE|terraform.|^$'


#SEE EG BELOW
#go test -run Foo     # Run top-level tests matching "Foo", such as "TestFooBar". 

