#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

export TF_ACC=1
export TF_LOG=TRACE



. ./test_setlogin.sh


export TF_VAR_VAPP_NAME="acvapp"
export TF_VAR_VAPP_TEMPLATE_NAME="CentOS7"
export TF_VAR_VAPP_CATALOG_NAME="Components"
export TF_VAR_VAPP_VDC="Terraform_VDC"

export TF_VAR_VAPP_NETWORK="192.168.10.0/247"
export TF_VAR_VAPP_MEMORY="64"
export TF_VAR_VAPP_CPU="1"
export TF_VAR_VAPP_IP_ALLOCATION_MODE="dhcp"

go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/ -v -run TestAccResourceVApp | grep --line-buffered -vE 'DEBUG|TRACE|terraform|^$'


#MAX INT GO 18446744073709551615