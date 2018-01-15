# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause

echo 'set up rum dependencies'
 yum install git -y
 yum group install "Development Tools" -y
 yum install wget -y
 yum install zlib-devel -y
 yum install openssl openssl-devel -y



echo 'install  PIP3.6'

. ./setup_pip3.6.sh

echo 'install  PIP3.6 dependencies'

 pip3.6 install grpcio
 pip3.6 install grpcio-tools
 pip3.6 install grpcio_health_checking
 pip3.6 install vcd_cli
 pip3.6 install pyvcloud


. ./setup_go.sh

 cd /home/
 ls
 export GOPATH=/home/terraform-provider-vcloud-director/go/


cd $GOPATH/src/
./init.sh
./build.sh

. ./setup_protoc.sh



cd /home/terraform-provider-vcloud-director
export PATH=$PATH:$GOPATH/bin
./rebuildproto.sh