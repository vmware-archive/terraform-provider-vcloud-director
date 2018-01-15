#!/bin/bash
echo 'install pip3.6 '
pip3.6 -V
if [ $? -eq 0 ]
then
	echo 'NO ACTION REQ : pip3.6 -V installed'
else 
	cd /opt
	curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py -k
	python3 get-pip.py 

fi

