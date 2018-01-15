#!/bin/bash
echo 'install Python-3'
python3 -V
if [ $? -eq 0 ]
then
	echo 'NO ACTION REQ : Python-3 installed'
else 
	cd /opt
	wget https://www.python.org/ftp/python/3.6.3/Python-3.6.3.tgz
	tar -xvf Python-3.6.3.tgz 
	ls
	cd Python-3.6.3
	./configure 
	make
	make install

fi