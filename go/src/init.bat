rem *****************************************************************
rem terraform-provider-vcloud-director
rem Copyright (c) 2017 VMware, Inc. All Rights Reserved.
rem SPDX-License-Identifier: BSD-2-Clause
rem *****************************************************************


echo "go get running ...."
./goget.bat


echo "go get complete ...."
./rmtrace.bat


echo "now builidng ......"
./build.bat
