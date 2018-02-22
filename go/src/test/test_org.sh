#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

export TF_ACC=1
export TF_LOG=TRACE

export VCD_ALLOW_UNVERIFIED_SSL=true
export VCD_IP="255.255.255.255"


export VCD_USER="administrator"
export VCD_PASSWORD="******"
export VCD_ORG="SYSTEM"

export TF_VAR_ORG_NAME="pcp_org_600"
export TF_VAR_FULL_NAME="test_pcp_org_full_name"
export TF_VAR_IS_ENABLED=true
#export TF_VAR_IS_DESCRIPTION = "desc"
export TF_VAR_IS_DISABLED=false
export TF_VAR_FORCE=true
export TF_VAR_RECURSIVE=true

go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/  -v -run TestAccResourceOrg | grep --line-buffered -vE 'DEBUG|TRACE|terraform.|^$'