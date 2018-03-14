#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************
provider "vcloud-director" {
  #value come from ENV VARIALES
}


resource "vcloud-director_vapp" "vapp1" {
        name  					= "pcp_vapp_4"
        template_name 			= "tinyova"
        catalog_name  			= "ACME"
        vdc 					= "ACME_PAYG"
        network 				= "172.16.1.0"
        ip_allocation_mode 		= "dhcp"
        cpu 					= 2
        memory 					= 64
        power_on 				= false
       

}
	