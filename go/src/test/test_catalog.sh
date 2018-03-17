#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

export TF_ACC=1
export TF_LOG=TRACE

. ./test_setlogin.sh

export TF_VAR_CATALOG_NAME="test_acc_cata1"
export TF_VAR_CATALOG_DESCRIPTION="accc_ description1"

#go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/ -v -run TestAccResourceCatalog | grep --line-buffered -vE 'TRACE|terraform.|^$'

go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/  -v -run TestAccResourceCatalogBasic | grep --line-buffered -vE 'DEBUG|TRACE|terraform.|^$'


#SEE EG BELOW
#go test -run Foo     # Run top-level tests matching "Foo", such as "TestFooBar".