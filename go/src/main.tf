#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************
provider "vcloud-director" {
  #value come from ENV VARIALES
}


resource "vcloud-director_vapp_vm" "source_vapp_vm"{
      // vappvm from vapp
      // target_vapp="test2"
      // target_vdc="ACME_PAYG"
      // target_vm_name="pcp_hi_09"
      // source_vm_name="pcp_hi_09"
      // source_vapp = "test1"
      // network = "External-VM-Network"
      // ip_allocation_mode = "dhcp"
      // hostname = "ubuntu"
      // password = "abc"
      // password_auto = false
      // password_reset = false
      // memory = 64
      // cores_per_socket = 2
      // virtual_cpus = 2
      // power_on = true
      // all_eulas_accepted = true

      // vappvm from catalog
      target_vapp="test1"
      target_vdc="ACME_PAYG"
      target_vm_name="pcp_hi_09"
      source_vm_name="Tiny Linux template"
      source_catalog_name="ACME"
      source_template_name="Tiny Linux VM.ova"
      network = "External-VM-Network"
      ip_allocation_mode = "dhcp"
      hostname = "ubuntu"
      password = "abc"
      password_auto = false
      password_reset = false
      memory = 64
      cores_per_socket = 2
      virtual_cpus = 2
      power_on = true
      all_eulas_accepted = true
}
	