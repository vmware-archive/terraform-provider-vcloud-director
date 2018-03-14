#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

export TF_ACC=1
export TF_LOG=TRACE

export VCD_ALLOW_UNVERIFIED_SSL=true
export VCD_IP="10.172.158.127"

export VCD_USER="administrator"
export VCD_PASSWORD="VMware1!"
export VCD_ORG="SYSTEM"

export TF_VAR_VDC_NAME="pcp_vdc_4000"
export TF_VAR_PROVIDER_VDC="PVDC1"
export TF_VAR_DESCRIPTION="PCP desc VDC"
export TF_VAR_ALLOCATION_MODEL="AllocationVApp"

export TF_VAR_CPU_UNITS="MHz"
export TF_VAR_CPU_ALLOCATED=0
export TF_VAR_CPU_LIMIT=0

export TF_VAR_MEM_UNITS="MB"
export TF_VAR_MEM_ALLOCATED=0
export TF_VAR_MEM_LIMIT=0

export TF_VAR_NIC_QUOTA=0
export TF_VAR_NETWORK_QUOTA=0
export TF_VAR_VM_QUOTA=0

export TF_VAR_STORAGE_PROFILE_1="{ \"name\" : \"Performance\",\"enabled\"  : true, \"units\" : \"MB\", \"limit\" : 0, \"default\" : true }"
#export TF_VAR_STORAGE_PROFILE_NAME="Performance"
#export TF_VAR_STORAGE_PROFILE_ENABLED=true
#export TF_VAR_STORAGE_PROFILE_UNITS="MB"
#export TF_VAR_STORAGE_PROFILE_LIMIT=0
#export TF_VAR_STORAGE_PROFILE_DEFAULT=true

export TF_VAR_RESOURCE_GUARANTEED_MEMORY=0.0
export TF_VAR_RESOURCE_GUARANTEED_CPU=0.0
export TF_VAR_VCPU_IN_MHZ=256

export TF_VAR_IS_THIN_PROVISION=false
export TF_VAR_NETWORK_POOL_NAME=""
export TF_VAR_USES_FAST_PROVISIONING=false
export TF_VAR_OVER_COMMIT_ALLOWED=false
export TF_VAR_VM_DISCOVERY_ENABLED=true

export TF_VAR_IS_ENABLED_TRUE=true
export TF_VAR_IS_ENABLED_FALSE=false

go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/  -v -run TestAccResourceVdc | grep --line-buffered -vE 'DEBUG|TRACE|terraform.|^$'