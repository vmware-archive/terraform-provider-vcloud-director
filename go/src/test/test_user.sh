#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

export TF_ACC=1
export TF_LOG=TRACE



export TF_ACC=1
export TF_LOG=TRACE

export VCD_ALLOW_UNVERIFIED_SSL=true
export VCD_IP="255.255.255.255"

export VCD_USER="acmeadmin"
export VCD_PASSWORD="*******"
export VCD_ORG="Acme"

export TF_VAR_USER_NAME="pcp_pcp_google_4"
export TF_VAR_USER_PASSWORD="123456"
export TF_VAR_ROLE_NAME="Organization Administrator"
export TF_VAR_FULL_NAME="Prakash Pandey"
export TF_VAR_DESCRIPTION="desc"
export TF_VAR_EMAIL="pcp@mail.com"
export TF_VAR_TELEPHONE="12345678"
export TF_VAR_IM="i_m_val"
export TF_VAR_ALERT_EMAIL="pcp_alert@mail.com"
export TF_VAR_ALERT_EMAIL_PREFIX="mail_alert_prefix"
export TF_VAR_STORED_VM_QUOTA=0
export TF_VAR_DEPLOYED_VN_QUOTA=0
export TF_VAR_IS_GROUP_ROLE=false
export TF_VAR_IS_DEFAULT_CACHED=false
export TF_VAR_IS_EXTERNAL=false
export TF_VAR_IS_ALERT_ENABLED=true
export TF_VAR_IS_ENABLED_TRUE=true
export TF_VAR_IS_ENABLED_FALSE=false

go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/  -v -run TestAccResourceUser | grep --line-buffered -vE 'DEBUG|TRACE|terraform.|^$'