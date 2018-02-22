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


resource "vcloud-director_org" "source_org" {
        name    = "pcp_5555"
        full_name = "oo"
        #description = "oo"
        is_enabled = true
        force = true
        recursive = true
}











