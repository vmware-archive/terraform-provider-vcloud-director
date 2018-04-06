#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

OS="`uname`"
mkdir ../bin
case $OS in
  'Linux')
    go build -o ../bin/terraform-provider-vcloud-director
    ;;
  'WindowsNT')
    OS='Windows'
    ;;
  'Darwin') 
    go build -o ../bin/terraform-provider-vcloud-director
    ;;
esac

