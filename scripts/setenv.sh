# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause

echo ' Please execute this file as $ . ./setenv.sh in the current shell. '
echo ' Default project location -/home/terraform-provider-vcloud-director , edit as necessary'
export PATH=/opt/Python-3.6.3/:$PATH
export PATH=/opt/go/bin:$PATH
export GOROOT=/opt/go
export GOPATH=/home/terraform-provider-vcloud-director/go/
export PATH=$PATH:$GOPATH/bin

export PY_PLUGIN='python3 /home/terraform-provider-vcloud-director/plugin-python/plugin.py'