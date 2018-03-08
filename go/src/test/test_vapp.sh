#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

export TF_ACC=1
export TF_LOG=TRACE

. ./test_setlogin.sh


export TF_VAR_VAPP_NAME="ACC_VAPP"
export TF_VAR_VAPP_TEMPLATE_NAME="CentOS7"
export TF_VAR_VAPP_CATALOG_NAME="sri"
export TF_VAR_VAPP_VDC="Solution Engineering"

export TF_VAR_VAPP_NETWORK="Intranet"
export TF_VAR_VAPP_MEMORY="8000"
export TF_VAR_VAPP_CPU="4"
export TF_VAR_VAPP_IP_ALLOCATION_MODE="pool"

go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/ -v -run TestAccResourceVApp | grep --line-buffered -vE 'DEBUG|TRACE|terraform|^$'


#MAX INT GO 18446744073709551615