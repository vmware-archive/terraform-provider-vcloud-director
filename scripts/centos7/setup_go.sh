#!/bin/bash
echo 'install go 1.9 '
go version
if [ $? -eq 0 ]
then
	echo 'NO ACTION REQ : go is installed'
else 

 	cd /opt
	wget https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz --no-check-certificate
	tar -xvf go1.9.linux-amd64.tar.gz 
	export PATH=/opt/go/bin:$PATH
	export  GOROOT=/opt/go

fi