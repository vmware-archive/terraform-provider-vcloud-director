#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************


echo "go get running ...."
./goget.sh


echo "go get complete ...."

./rmtrace.sh


echo "now builidng ......"
./build.sh 
