#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

export TF_ACC=1
export TF_LOG=TRACE
export VCD_ALLOW_UNVERIFIED_SSL=true
export VCD_IP="10.172.158.132"
export VCD_USER="acmeadmin"
export VCD_PASSWORD="VMware1!"
export VCD_ORG="acme"
export TF_VAR_TARGET_VAPP_VDC="ACME_PAYG"
export TF_VAR_TARGET_VM_NAME="pcp_hi_09"
export TF_VAR_NETWORK="External-VM-Network"
export TF_VAR_VAPP_IP_ALLOCATION_MODE="dhcp"
export TF_VAR_HOST_NAME="ubuntu"
export TF_VAR_PASSWORD="abc"
export TF_VAR_PASSWORD_AUTO=false
export TF_VAR_PASSWORD_RESET=false
export TF_VAR_MEMORY=64
export TF_VAR_CORES_PER_SOCKET=2
export TF_VAR_VIRTUAL_CPUS=2
export TF_VAR_POWER_ON=true
export TF_VAR_ALL_EULAS_ACCEPTED=true

# create from vapp
# export TF_VAR_TARGET_VAPP_NAME="test1"
# export TF_VAR_SOURCE_VM_NAME="MS Machine"
# export TF_VAR_SOURCE_VAPP="test2"

# create from catalog
export TF_VAR_TARGET_VAPP_NAME="test2"
export TF_VAR_SOURCE_VM_NAME="Tiny Linux template"
export TF_VAR_SOURCE_CATALOG_NAME="ACME"
export TF_VAR_TEMPLATE_NAME="Tiny Linux VM.ova"

go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/ -v -run TestAccResourceVappVm | grep --line-buffered -vE 'DEBUG|TRACE|terraform|^$'

# go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/ -v -run TestAccResourceVappVmFromVapp | grep --line-buffered -vE 'DEBUG|TRACE|terraform|^$'

