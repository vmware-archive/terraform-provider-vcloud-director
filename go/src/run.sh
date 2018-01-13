#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************

export TF_LOG=1
terraform init
echo call PLAN 
terraform plan -out o1 > f.txt 2>&1

echo call APPLY

terraform apply o1 >> f.txt 2>&1
cat f.txt | grep PLUGINLOG
