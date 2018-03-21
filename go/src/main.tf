#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************
provider "vcloud-director" {
  #value come from ENV VARIALES
}

resource "vcloud-director_catalog" "catalog1" {
  name = "test_catalog"
  description = "Test Catalog"
  shared  = false
}

# resource "vcloud-director_catalog_item_ova" "centos" {
#   item_name = "${vcloud-director_vapp.vapp1.template_name}"
#   catalog_name = "${vcloud-director_catalog.catalog1.name}"
#   source_vdc_name ="${vcloud-director_vapp.vapp1.vdc}"
#   source_vapp_name = "${vcloud-director_vapp.vapp1.name}"
#   customize_on_instantiate = "true"
# }

resource "vcloud-director_catalog_item_ova" "tiny_linux_vm" {
  item_name = "Tiny Linux VM"
  catalog_name= "${vcloud-director_catalog.catalog1.name}"
  source_file_path="/Users/mtaneja/Downloads/Tiny Linux VM.ovf"
}

resource "vcloud-director_catalog_item_media" "tiny_core_current" {
  item_name = "TinyCore-current.iso"
  catalog_name= "${vcloud-director_catalog.catalog1.name}"
  source_file_path="/Users/mtaneja/Downloads/TinyCore-current.iso"
}

resource "vcloud-director_vapp" "vapp1" {
  name = "acvapp"
  template_name = "CentOS7"
  catalog_name = "Components"
  vdc = "Terraform_VDC"
  network = "192.168.10.0/247"
  ip_allocation_mode = "dhcp"
  cpu = "1"
  memory = "64"
  power_on = false
  accept_all_eulas = true
}

resource "vcloud-director_vapp" "vapp2" {
  name = "test_vapp"
  template_name = "CentOS7"
  catalog_name = "Components"
  vdc = "Terraform_VDC"
  network = "192.168.10.0/247"
  ip_allocation_mode = "dhcp"
  cpu = "1"
  memory = "64"
  power_on = false
  accept_all_eulas = true
}

resource "vcloud-director_vapp_vm" "vapp_vm_catalog"{
  target_vapp = "${vcloud-director_vapp.vapp2.name}"
  target_vdc = "${vcloud-director_vapp.vapp2.vdc}"
  target_vm_name = "Tiny Linux template"
  source_vm_name = "Tiny Linux template"
  source_catalog_name = "${vcloud-director_catalog_item_ova.tiny_linux_vm.catalog_name}"
  source_template_name = "${vcloud-director_catalog_item_ova.tiny_linux_vm.item_name}"
  network = "192.168.10.0/247"
  ip_allocation_mode = "dhcp"
  hostname = "ubuntu"
  password = "abc"
  password_auto = false
  password_reset = false
  memory = 64
  cores_per_socket = 2
  virtual_cpus = 2
  power_on = false
  all_eulas_accepted = true
}

resource "vcloud-director_vapp_vm" "vapp_vm_vapp"{
  target_vapp = "${vcloud-director_vapp.vapp2.name}"
  target_vdc = "${vcloud-director_vapp.vapp2.vdc}"
  target_vm_name = "centos3"
  source_vm_name = "${vcloud-director_vapp.vapp1.template_name}"
  source_vapp = "${vcloud-director_vapp.vapp1.name}"
  network = "192.168.10.0/247"
  ip_allocation_mode = "dhcp"
  hostname = "ubuntu"
  password = "abc"
  password_auto = false
  password_reset = false
  memory = 64
  cores_per_socket = 2
  virtual_cpus = 2
  power_on = false
  all_eulas_accepted = true
}