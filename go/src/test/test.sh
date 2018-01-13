#!/bin/bash
#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************


scripts[0]=test_media.sh
scripts[1]=test_ova.sh
#scripts[0]=test_disk.sh
#scripts[0]=test_media.sh
#scripts[0]=test_ova.sh
#scripts[0]=test_token_login.sh
#scripts[0]=test_vapp.sh

for script in "${scripts[@]}"
do
  echo "${script}"
  

	sh ./$script
	status=$?
	echo $status
	if [ $status -ne 0 ]
		then
		echo "FAILED ${script}"
		exit $status
	fi
  # do something on $var
done





# catalog item test follows







#./test_capture_vapp.sh


#vapp with params
#./test_vapp.sh


#independent disk
#./test_disk.sh