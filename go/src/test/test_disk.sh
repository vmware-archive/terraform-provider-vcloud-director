#!/bin/bash
#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************
#go test disk_test.go -v run  TestDiskInterface

export TF_ACC=1
export TF_LOG=DEBUG

export VCD_ALLOW_UNVERIFIED_SSL=true
export VCD_IP="10.112.83.27"


export VCD_USER="user1"
export VCD_PASSWORD="Admin!23"
export VCD_ORG="O1"

export TF_VAR_DISK_NAME="test_acc_disk1" #DESTINATION CATALOG





go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/ -v -run TestAccResourceIndependentDiskBasic 

#| grep --line-buffered -vE 'TRACE|terraform|^$'
