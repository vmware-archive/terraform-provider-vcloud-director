#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

export TF_ACC=1
export TF_LOG=TRACE

. ./test_setlogin.sh

export TF_VAR_CATALOG_NAME_OLD="pcp_test_calatog_1009"
export TF_VAR_CATALOG_NAME_NEW="pcp_test_calatog_10010"

export TF_VAR_CATALOG_DESCRIPTION_OLD="pcp_description9"
export TF_VAR_CATALOG_DESCRIPTION_NEW="pcp_description10"

export TF_VAR_CATALOG_SHARED="true"

go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/  -v -run TestAccResourceCatalogAdvance | grep --line-buffered -vE 'DEBUG|TRACE|terraform.|^$'

#SEE EG BELOW
#go test -run Foo     # Run top-level tests matching "Foo", such as "TestFooBar". 

