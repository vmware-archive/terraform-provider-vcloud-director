/*****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
******************************************************************/

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/util/logging"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
	"os"
	"strconv"
	"testing"
)

func TestAccResourceVdc(t *testing.T) {
	logging.Plog("__INIT__TestAccResourceVdc__")

	conf := testAccVdc_basic + "\n" + testAccVdc_enable
	logging.Plog("conf : \n" + conf)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVdcDestroy,
		Steps: []resource.TestStep{

			resource.TestStep{
				Config: testAccVdc_basic + "\n" + testAccVdc_enable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCreateVdc(),
				),
			},

			resource.TestStep{
				Config: testAccVdc_basic + "\n" + testAccVdc_disable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdateVdc(),
				),
			},
		},
	})

	logging.Plog("__DONE__TestAccResourceVdc__")
}

func testAccCheckCreateVdc() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logging.Plog("__INIT_testAccCheckCreateVdc__")

		provider := providerGlobalRefPointer.vdcProvider

		name := os.Getenv("TF_VAR_VDC_NAME")
		//providerVdc := os.Getenv("TF_VAR_PROVIDER_VDC")
		description := os.Getenv("TF_VAR_DESCRIPTION")
		allocationModel := os.Getenv("TF_VAR_ALLOCATION_MODEL")

		cpuUnits := os.Getenv("TF_VAR_CPU_UNITS")

		cpuAllocated, _ := strconv.Atoi(os.Getenv("TF_VAR_CPU_ALLOCATED"))
		cpuAllocated32 := int32(cpuAllocated)

		cpuLimit, _ := strconv.Atoi(os.Getenv("TF_VAR_CPU_LIMIT"))
		cpuLimit32 := int32(cpuLimit)

		memUnits := os.Getenv("TF_VAR_MEM_UNITS")

		memAllocated, _ := strconv.Atoi(os.Getenv("TF_VAR_MEM_ALLOCATED"))
		memAllocated32 := int32(memAllocated)

		memLimit, _ := strconv.Atoi(os.Getenv("TF_VAR_MEM_LIMIT"))
		memLimit32 := int32(memLimit)

		nicQuota, _ := strconv.Atoi(os.Getenv("TF_VAR_NIC_QUOTA"))
		nicQuota32 := int32(nicQuota)

		networkQuota, _ := strconv.Atoi(os.Getenv("TF_VAR_NETWORK_QUOTA"))
		networkQuota32 := int32(networkQuota)

		vmQuota, _ := strconv.Atoi(os.Getenv("TF_VAR_VM_QUOTA"))
		vmQuota32 := int32(vmQuota)

		//storageProfiles := os.Getenv("storage_profiles")

		//resourceGuaranteedMemory := os.Getenv("TF_VAR_RESOURCE_GUARANTEED_MEMORY").(float32)
		//resourceGuaranteedCpu := os.Getenv("TF_VAR_RESOURCE_GUARANTEED_CPU").(float32)
		vcpuInMhz, _ := strconv.Atoi(os.Getenv("TF_VAR_VCPU_IN_MHZ"))
		vcpuInMhz32 := int32(vcpuInMhz)

		//isThinProvision := os.Getenv("TF_VAR_IS_THIN_PROVISION").(bool)
		//networkPoolName := os.Getenv("TF_VAR_NETWORK_POOL_NAME")
		//usesFastProvisioning := os.Getenv("TF_VAR_USES_FAST_PROVISIONING").(float32)
		//overCommitAllowed := os.Getenv("TF_VAR_OVER_COMMIT_ALLOWED").(bool)
		//vmDiscoveryEnabled := os.Getenv("TF_VAR_VM_DISCOVERY_ENABLED").(bool)
		isEnabled, _ := strconv.ParseBool(os.Getenv("TF_VAR_IS_ENABLED_TRUE"))

		readResp, readErrp := provider.Read(proto.ReadVdcInfo{Name: name})

		if readErrp != nil {
			return fmt.Errorf("__ERROR__.... in reading vdc  creation", readErrp)
		}
		if readResp.Present {
			if !(name == readResp.Name) {
				return fmt.Errorf("__ERROR__.... name  do not match [expected: %v, found: %v]", name, readResp.Name)
			}

			if !(description == readResp.Description) {
				return fmt.Errorf("__ERROR__.... description do not match [expected: %v, found: %v]", description, readResp.Description)
			}

			if !(allocationModel == readResp.AllocationModel) {
				return fmt.Errorf("__ERROR__.... allocationModel do not match [expected: %v, found: %v]", allocationModel, readResp.AllocationModel)
			}

			if !(cpuUnits == readResp.CpuUnits) {
				return fmt.Errorf("__ERROR__.... cpuUnits do not match [expected: %v, found: %v]", cpuUnits, readResp.CpuUnits)
			}

			if !(cpuAllocated32 == readResp.CpuAllocated) {
				return fmt.Errorf("__ERROR__.... vdc cpuAllocated do not match [expected: %v, found: %v]", cpuAllocated32, readResp.CpuAllocated)
			}

			if !(cpuLimit32 == readResp.CpuLimit) {
				return fmt.Errorf("__ERROR__.... cpuLimit do not match [expected: %v, found: %v]", cpuLimit32, readResp.CpuLimit)
			}

			if !(memUnits == readResp.MemUnits) {
				return fmt.Errorf("__ERROR__.... memUnits do not match [expected: %v, found: %v]", memUnits, readResp.MemUnits)
			}

			if !(memAllocated32 == readResp.MemAllocated) {
				return fmt.Errorf("__ERROR__.... memAllocated do not match [expected: %v, found: %v]", memAllocated32, readResp.MemAllocated)
			}

			if !(memLimit32 == readResp.MemLimit) {
				return fmt.Errorf("__ERROR__.... memLimit do not match [expected: %v, found: %v]", memLimit32, readResp.MemLimit)
			}

			if !(nicQuota32 == readResp.NicQuota) {
				return fmt.Errorf("__ERROR__.... nicQuota do not match [expected: %v, found: %v]", nicQuota32, readResp.NicQuota)
			}

			if !(networkQuota32 == readResp.NetworkQuota) {
				return fmt.Errorf("__ERROR__.... networkQuota do not match [expected: %v, found: %v]", networkQuota32, readResp.NetworkQuota)
			}

			if !(vmQuota32 == readResp.VmQuota) {
				return fmt.Errorf("__ERROR__.... vmQuota do not match [expected: %v, found: %v]", vmQuota32, readResp.VmQuota)
			}

			if !(vcpuInMhz32 == readResp.VcpuInMhz) {
				return fmt.Errorf("__ERROR__.... vcpuInMhz do not match [expected: %v, found: %v]", vcpuInMhz32, readResp.VcpuInMhz)
			}

			if !(isEnabled == readResp.IsEnabled) {
				return fmt.Errorf("__ERROR__.... isEnabled do not match [expected: %v, found: %v]", isEnabled, readResp.IsEnabled)
			}

			logging.Plog("Vdc creation varified")
		} else {
			return fmt.Errorf("__ERROR__.... vdc[%v]  not found:", name)
		}
		logging.Plog("__DONE_testAccCheckCreatevdc__")
		return nil
	}
}

func testAccCheckUpdateVdc() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logging.Plog("__INIT_testAccCheckUpdateVdc__")

		provider := providerGlobalRefPointer.vdcProvider

		name := os.Getenv("TF_VAR_VDC_NAME")

		readResp, readErrp := provider.Read(proto.ReadVdcInfo{Name: name})
		if readErrp != nil {
			return fmt.Errorf("__ERROR__.... in reading Vdc ", readErrp)
		}
		if readResp.Present && !readResp.IsEnabled {
			logging.Plog("Vdc update varified")
		} else {
			return fmt.Errorf("__ERROR__.... updating Vdc[%v], Vdc read response[%v]", name, readResp)
		}
		logging.Plog("__DONE_testAccCheckUpdateVdc__")
		return nil
	}
}

func testAccCheckVdcDestroy(s *terraform.State) error {

	logging.Plog("__INIT__testAccCheckVdcDestroy_")

	provider := providerGlobalRefPointer.vdcProvider

	name := os.Getenv("TF_VAR_VDC_NAME")

	readResp, _ := provider.Read(proto.ReadVdcInfo{Name: name})

	if readResp.Present {
		return fmt.Errorf("__ERROR__.... Vdc[%v] found", name)
	}

	logging.Plog("__DONE__testAccCheckVdcDestroy_")
	return nil

}

const testAccVdc_basic = `
	provider "vcloud-director" {
		  #value come from ENV VARIALES
	}
	variable "VDC_NAME" {
		 type    = "string"
		 default = "NOT DEFINED"
	}
	variable "PROVIDER_VDC" {
		 type    = "string"
		 default = ""
	}
	variable "DESCRIPTION" {
		 type    = "string"
		 default = ""
	}
	variable "ALLOCATION_MODEL" {
		 type    = "string"
		 default = ""
	}
	variable "CPU_UNITS" {
		 type    = "string"
		 default = ""
	}
	variable "CPU_ALLOCATED" {
		 default = 0
	}
	variable "CPU_LIMIT" {
		 default = 0
	}
	variable "MEM_UNITS" {
		 type    = "string"
		 default = ""
	}
	variable "MEM_ALLOCATED" {
		 default = 0
	}
	variable "MEM_LIMIT" {
		 default = 0
	}
	variable "NIC_QUOTA" {
		 default = 0
	}
	variable "NETWORK_QUOTA" {
		 default = 0
	}
	variable "VM_QUOTA" {
		 default = 0
	}
	
	variable "STORAGE_PROFILE_NAME" {
		 type    = "string"
		 default = ""
	}
	variable "STORAGE_PROFILE_ENABLED" {
		 type    = "string"
		 default = ""
	}
	variable "STORAGE_PROFILE_LIMIT" {
		 default = 0
	}
	variable "STORAGE_PROFILE_DEFAULT" {
		 default = true
	}
	
	
	variable "RESOURCE_GUARANTEED_MEMORY" {
		 default = 0.0
	}
	
	variable "RESOURCE_GUARANTEED_CPU" {
		 default = 0.0
	}
	
	variable "VCPU_IN_MHZ" {
		 default = 0
	}
	
	variable "IS_THIN_PROVISION" {
		 default = false
	}
	
	variable "NETWORK_POOL_NAME" {
		 type    = "string"
		 default = ""
	}
	
	variable "USES_FAST_PROVISIONING" {
		 default = false
	}
	
	variable "OVER_COMMIT_ALLOWED" {
		 default = false
	}
	
	variable "VM_DISCOVERY_ENABLED" {
		 default = true
	}
	
	variable "IS_ENABLED_TRUE" {
		 default = true
	}

`

const testAccVdc_enable = `
variable "STORAGE_PROFILE_1" {
   type = "string"
   default = "{ \"name\" : \"Performance\",\"enabled\"  : true, \"units\" : \"MB\", \"limit\" : 0, \"default\" : true }"
}

resource "vcloud-director_vdc" "source_vdc"{
        name = "${var.VDC_NAME}"
		provider_vdc = "${var.PROVIDER_VDC}"
		description = "${var.DESCRIPTION}"
		allocation_model = "${var.ALLOCATION_MODEL}"
		
		cpu_units = "${var.CPU_UNITS}"
		cpu_allocated = "${var.CPU_ALLOCATED}"
		cpu_limit = "${var.CPU_LIMIT}"
		
		mem_units = "${var.MEM_UNITS}"
		mem_allocated = "${var.MEM_ALLOCATED}"
		mem_limit = "${var.MEM_LIMIT}"
		nic_quota = "${var.NIC_QUOTA}"
		network_quota = "${var.NETWORK_QUOTA}"
		vm_quota = "${var.VM_QUOTA}"
		
		storage_profiles = "[${var.STORAGE_PROFILE_1}]"
		
		resource_guaranteed_memory = "${var.RESOURCE_GUARANTEED_MEMORY}"
		resource_guaranteed_cpu = "${var.RESOURCE_GUARANTEED_CPU}"
		vcpu_in_mhz = "${var.VCPU_IN_MHZ}"
		
		is_thin_provision = "${var.IS_THIN_PROVISION}"
		network_pool_name = "${var.NETWORK_POOL_NAME}"
		uses_fast_provisioning = "${var.USES_FAST_PROVISIONING}"
		over_commit_allowed = "${var.OVER_COMMIT_ALLOWED}"
		vm_discovery_enabled = "${var.VM_DISCOVERY_ENABLED}"
		is_enabled = true
}
`

const testAccVdc_disable = `

variable "STORAGE_PROFILE_1" {
   type = "string"
   default = "{ \"name\" : \"Performance\",\"enabled\"  : true, \"units\" : \"MB\", \"limit\" : 0, \"default\" : true }"
}

resource "vcloud-director_vdc" "source_vdc"{
        name = "${var.VDC_NAME}"
		provider_vdc = "${var.PROVIDER_VDC}"
		description = "${var.DESCRIPTION}"
		allocation_model = "${var.ALLOCATION_MODEL}"
		
		cpu_units = "${var.CPU_UNITS}"
		cpu_allocated = "${var.CPU_ALLOCATED}"
		cpu_limit = "${var.CPU_LIMIT}"
		
		mem_units = "${var.MEM_UNITS}"
		mem_allocated = "${var.MEM_ALLOCATED}"
		mem_limit = "${var.MEM_LIMIT}"
		nic_quota = "${var.NIC_QUOTA}"
		network_quota = "${var.NETWORK_QUOTA}"
		vm_quota = "${var.VM_QUOTA}"
		
		storage_profiles = "[${var.STORAGE_PROFILE_1}]"
		
		resource_guaranteed_memory = "${var.RESOURCE_GUARANTEED_MEMORY}"
		resource_guaranteed_cpu = "${var.RESOURCE_GUARANTEED_CPU}"
		vcpu_in_mhz = "${var.VCPU_IN_MHZ}"
		
		is_thin_provision = "${var.IS_THIN_PROVISION}"
		network_pool_name = "${var.NETWORK_POOL_NAME}"
		uses_fast_provisioning = "${var.USES_FAST_PROVISIONING}"
		over_commit_allowed = "${var.OVER_COMMIT_ALLOWED}"
		vm_discovery_enabled = "${var.VM_DISCOVERY_ENABLED}"
		is_enabled = false
}
`
