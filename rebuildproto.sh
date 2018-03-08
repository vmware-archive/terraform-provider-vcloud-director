#!/bin/bash
#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************


echo 'rebuild go proto'


DIR="go/src/github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/"

echo $DIR

protoc -I $DIR  -I $DIR/proto proto/pyvcloudprovider.proto  proto/vapp.proto proto/catalog_item.proto proto/disk.proto proto/org.proto proto/user.proto proto/vdc.proto --go_out=plugins=grpc:$DIR


DIR="go/src/github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/"
echo 'rebuild python'


python3 -m grpc_tools.protoc -I ./$DIR -I ./$DIR/proto   --python_out=./plugin-python/ --grpc_python_out=./plugin-python/  proto/vapp.proto proto/pyvcloudprovider.proto proto/catalog_item.proto proto/disk.proto proto/org.proto proto/user.proto proto/vdc.proto        
