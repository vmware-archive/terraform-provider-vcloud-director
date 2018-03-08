#*****************************************************************
# terraform-provider-vcloud-director
# Copyright (c) 2017 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: BSD-2-Clause
#*****************************************************************


provider "vcloud-director" {
  #value come from ENV VARIALES
}

variable "STORAGE_PROFILE_1" {
   type = "string"
   default = "{ \"name\" : \"Performance\",\"enabled\"  : true, \"units\" : \"MB\", \"limit\" : 0, \"default\" : true }"
   #default = "{ 'name' = 'Performance','is_enabled'  = true, 'units' = 'MB', 'limit' = 0, 'default' = true }"
}


resource "vcloud-director_vdc" "source_vdc"{
        name = "vdc_100"
		provider_vdc = "PVDC1"
		description = "Desc VDC"
		allocation_model = "AllocationVApp"
		
		cpu_units = "MHz"
		cpu_allocated = 0
		cpu_limit = 0
		
		mem_units = "MB"
		mem_allocated = 0
		mem_limit = 0
		
		nic_quota = 0
		network_quota = 0
		vm_quota = 0
		
		#storage_profiles = "${var.STORAGE_PROFILE_1}"
		storage_profiles = "[${var.STORAGE_PROFILE_1}]"
		
		resource_guaranteed_memory = 0.0
		resource_guaranteed_cpu = 0.0
		vcpu_in_mhz = 256
		
		is_thin_provision = false
		network_pool_name = ""
		uses_fast_provisioning = false
		over_commit_allowed = false
		vm_discovery_enabled = true
		is_enabled = true
		
		#storage_profiles_name = "Performance"
		#storage_profiles_is_enabled  = true
		#storage_profiles_units = "MB"
		#storage_profiles_limit = 0
		#storage_profiles_default = false
}		
