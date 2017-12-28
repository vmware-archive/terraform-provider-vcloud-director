#!/bin/bash
#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************


echo 'rebuild go proto'


DIR="go/src/github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"

echo $DIR
                                                                                                                                                            
protoc -I $DIR  $DIR/pyvcloudprovider.proto --go_out=plugins=grpc:$DIR


echo 'rebuild python'


python3 -m grpc_tools.protoc -I ./$DIR --python_out=./plugin-python/ --grpc_python_out=./plugin-python/ ./$DIR/pyvcloudprovider.proto

