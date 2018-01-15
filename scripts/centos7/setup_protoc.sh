#!/bin/bash
echo 'install protoc v3.4.1 '
protoc --version
if [ $? -eq 0 ]
then
	echo 'NO ACTION REQ : protoc is installed'
else 
	cd /opt
	wget https://github.com/google/protobuf/releases/download/v3.4.1/protobuf-cpp-3.4.1.tar.gz
	tar -xvf proto*
	cd proto*
	yum install autoconf automake libtool curl make g++ unzip -y
	./configure
	make
	make install
	ldconfig
	protoc --version


fi