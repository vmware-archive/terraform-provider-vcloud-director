#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

OS="`uname`"
case $OS in
  'Linux')
    go build -o ../../builds/linux/terraform-provider-vcloud-director
    ;;
  'WindowsNT')
    OS='Windows'
    ;;
  'Darwin') 
    go build -o ../../builds/mac/terraform-provider-vcloud-director
    ;;
esac

