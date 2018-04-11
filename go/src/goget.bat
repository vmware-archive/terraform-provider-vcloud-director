rem !/bin/bash
rem *****************************************************************
rem terraform-provider-vcloud-director
rem Copyright (c) 2017 VMware, Inc. All Rights Reserved.
rem SPDX-License-Identifier: BSD-2-Clause
rem *****************************************************************

go get github.com/hashicorp/terraform/
 
go get github.com/golang/protobuf/proto
 
go get github.com/hashicorp/go-plugin

rem go get google.golang.org/grpc 

go get github.com/hashicorp/logutils

go get -u github.com/golang/protobuf/protoc-gen-go
