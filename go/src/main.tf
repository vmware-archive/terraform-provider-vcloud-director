#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************


provider "vcloud-director" {
  #value come from ENV VARIALES
}




variable "VAPP_VDC" { 

 type    = "string"
 default = "NOT DEFINED" 
}

variable "VAPP_NETWORK" { 

 type    = "string"
 default = "" 
}

variable "VAPP_IP_ALLOCATION_MODE" { 

 type    = "string"
 default = "dhcp" 
}

variable "VAPP_CPU" { 

 type    = "string"
 default = "-1" 
}

variable "VAPP_MEMORY" { 

 type    = "string"
 default = "-1" 
}



variable "OVA_PATH" { 

 type    = "string"
 default = "nullova" 
}


resource "vcloud-director_catalog" "source_catalog" {
        name    = "source_catalog"
        description = "source_catalog"
        shared  = "true"

}

resource "vcloud-director_catalog_item_ova" "source_tinyova" {
    item_name = "source_tinyova"
    catalog_name= "${vcloud-director_catalog.source_catalog.name}"
    source_file_path="${var.OVA_PATH}"
}


resource "vcloud-director_vapp" "tinyVapp" {
        name                    = "tinyVapp"
        template_name           = "${vcloud-director_catalog_item_ova.source_tinyova.item_name}"
        catalog_name            = "${vcloud-director_catalog.source_catalog.name}"
        vdc                     = "${var.VAPP_VDC}"
        network                 = "${var.VAPP_NETWORK}"
        ip_allocation_mode      = "${var.VAPP_IP_ALLOCATION_MODE}"
        cpu                     = "${var.VAPP_CPU}"
        memory                  = "${var.VAPP_MEMORY}"
        storage_profile         = "DS2_173"
       

}


resource "vcloud-director_catalog" "dest_catalog" {
        name    = "dest_catalog"
        description = "dest_catalog"
        shared  = "true"

}


resource "vcloud-director_catalog_item_ova" "tinyVapp_capture" {
    item_name = "tinyVapp_capture"
    catalog_name= "${vcloud-director_catalog.dest_catalog.name}"
    source_vdc_name="${vcloud-director_vapp.tinyVapp.vdc}"
    source_vapp_name="${vcloud-director_vapp.tinyVapp.name}"
    customize_on_instantiate=true
}

