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
	"testing"
)

func TestAccResourceVappVm(t *testing.T) {
	logging.Plog("__INIT__TestAccResourceVappVm")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVappVmDestroy,
		Steps: []resource.TestStep{

			resource.TestStep{
				Config: testAccVappVm_basic + "\n" + testAccVappVm,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCreateVappVm(),
				),
			},
			resource.TestStep{
				Config: testAccVappVm_basic + "\n" + testAccVappVm,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdateVappVm(),
				),
			},
		},
	})

	logging.Plog("__DONE__TestAccResourceVappVm_")
}

// func TestAccResourceVappVmFromVapp(t *testing.T) {
// 	logging.Plog("__INIT__TestAccResourceVappVm")
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckVappVmDestroy,
// 		Steps: []resource.TestStep{

// 			resource.TestStep{
// 				Config: testAccVappVm_basic + "\n" + testAccVappVmFromVapp,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckCreateVappVm(),
// 				),
// 			},
// 			resource.TestStep{
// 				Config: testAccVappVm_basic + "\n" + testAccVappVmFromVapp,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckUpdateVappVm(),
// 				),
// 			},
// 		},
// 	})

// 	logging.Plog("__DONE__TestAccResourceVappVm_")
// }

func testAccCheckCreateVappVm() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logging.Plog("__INIT_testAccCheckCreateVappVm__")

		targetVmName := os.Getenv("TF_VAR_TARGET_VM_NAME")
		targetVapp := os.Getenv("TF_VAR_TARGET_VAPP_NAME")
		targetVdc := os.Getenv("TF_VAR_TARGET_VAPP_VDC")

		provider := providerGlobalRefPointer.vappVmProvider

		readVappVmInfo := proto.ReadVappVmInfo{
			TargetVmName: targetVmName,
			TargetVapp:   targetVapp,
			TargetVdc:    targetVdc,
		}

		readResp, readErrp := provider.Read(readVappVmInfo)

		if readErrp != nil {
			return fmt.Errorf("__ERROR__.... in reading VappVm  creation", readErrp)
		}

		if readResp.Present {
			logging.Plog("VappVm create varified")
		} else {
			return fmt.Errorf("__ERROR__.... VappVm[%v]  not found:", targetVmName)
		}

		logging.Plog("__DONE_testAccCheckCreateVappVm__")
		return nil
	}
}

func testAccCheckUpdateVappVm() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logging.Plog("__INIT_testAccCheckUpdateVappVm__")

		//targetVmName := os.Getenv("TF_VAR_TARGET_VM_NAME")
		//targetVapp := os.Getenv("TF_VAR_TARGET_VAPP_NAME")
		//targetVdc := os.Getenv("TF_VAR_TARGET_VAPP_VDC")

		logging.Plog("__DONE_testAccCheckUpdateVappVm__")
		return nil
	}
}

func testAccCheckVappVmDestroy(s *terraform.State) error {

	logging.Plog("__INIT__testAccCheckVappVmDestroy_")

	targetVmName := os.Getenv("TF_VAR_TARGET_VM_NAME")
	targetVapp := os.Getenv("TF_VAR_TARGET_VAPP_NAME")
	targetVdc := os.Getenv("TF_VAR_TARGET_VAPP_VDC")

	provider := providerGlobalRefPointer.vappVmProvider

	readVappVmInfo := proto.ReadVappVmInfo{
		TargetVmName: targetVmName,
		TargetVapp:   targetVapp,
		TargetVdc:    targetVdc,
	}

	readResp, readErrp := provider.Read(readVappVmInfo)

	if readErrp != nil {
		return fmt.Errorf("__ERROR__.... in reading VappVm  creation", readErrp)
	}

	if readResp.Present {
		return fmt.Errorf("__ERROR__.... VappVm[%v] found:", targetVmName)

	} else {
		logging.Plog("VappVm delete varified")
	}

	logging.Plog("__DONE__testAccCheckVappVmDestroy_")
	return nil

}

const testAccVappVm_basic = `

provider "vcloud-director" {
  #value come from ENV VARIALES
}

variable "TARGET_VAPP_NAME" {
 type    = "string"
 default = "NOT DEFINED"
}

variable "TARGET_VAPP_VDC" {
 type    = "string"
 default = "NOT DEFINED"
}

variable "TARGET_VM_NAME" {
 type    = "string"
 default = "NOT DEFINED"
}


variable "SOURCE_VM_NAME" {
 type    = "string"
 default = "NOT DEFINED"
}

variable "SOURCE_CATALOG_NAME" {
 type    = "string"
 default = "NOT DEFINED"
}

variable "SOURCE_VAPP" {
 type    = "string"
 default = "NOT DEFINED"
}

variable "TEMPLATE_NAME" {
 type    = "string"
 default = "NOT DEFINED"
}

variable "NETWORK" {
 type    = "string"
 default = "NOT DEFINED"
}

variable "VAPP_IP_ALLOCATION_MODE" {
 type    = "string"
 default = "NOT DEFINED"
}

variable "HOST_NAME" {
 type    = "string"
 default = "NOT DEFINED"
}

variable "PASSWORD" {
 type    = "string"
 default = "NOT DEFINED"
}

variable "PASSWORD_AUTO" {
 type    = "string"
 default = false
}

variable "PASSWORD_RESET" {
 type    = "string"
 default = false
}

variable "MEMORY" {
 type    = "string"
 default = 64
}

variable "VIRTUAL_CPUS" {
 type    = "string"
 default = 2
}

variable "CORES_PER_SOCKET" {
 type    = "string"
 default = 2
}

variable "POWER_ON" {
 type    = "string"
 default = true
}

variable "ALL_EULAS_ACCEPTED" {
 type    = "string"
 default = true
}

`

const testAccVappVm = `
resource "vcloud-director_vapp_vm" "source_vapp_vm"{
    target_vapp="${var.TARGET_VAPP_NAME}"
    target_vdc="${var.TARGET_VAPP_VDC}"
    target_vm_name="${var.TARGET_VM_NAME}"
    source_vm_name="${var.SOURCE_VM_NAME}"
    source_catalog_name="${var.SOURCE_CATALOG_NAME}"
    source_template_name="${var.TEMPLATE_NAME}"
    network = "${var.NETWORK}"
    ip_allocation_mode = "${var.VAPP_IP_ALLOCATION_MODE}"
    hostname = "${var.HOST_NAME}"
    password = "${var.PASSWORD}"
	password_auto = "${var.PASSWORD_AUTO}"
	password_reset = "${var.PASSWORD_RESET}"
	memory = "${var.MEMORY}"
	cores_per_socket = "${var.CORES_PER_SOCKET}"
	virtual_cpus = "${var.VIRTUAL_CPUS}"
	power_on = "${var.POWER_ON}"
	all_eulas_accepted = "${var.ALL_EULAS_ACCEPTED}"
}

`

// const testAccVappVmFromVapp = `
// resource "vcloud-director_vapp_vm" "source_vapp_vm"{
// 	target_vapp="${var.TARGET_VAPP_NAME}"
// 	target_vdc="${var.TARGET_VAPP_VDC}"
// 	target_vm_name="${var.TARGET_VM_NAME}"
// 	source_vm_name="${var.SOURCE_VM_NAME}"
// 	source_vapp="${var.SOURCE_VAPP}"
// 	network = "${var.NETWORK}"
// 	ip_allocation_mode = "${var.VAPP_IP_ALLOCATION_MODE}"
// 	hostname = "${var.HOST_NAME}"
//     password = "${var.PASSWORD}"
// 	password_auto = "${var.PASSWORD_AUTO}"
// 	password_reset = "${var.PASSWORD_RESET}"
// 	memory = "${var.MEMORY}"
// 	cores_per_socket = "${var.CORES_PER_SOCKET}"
// 	virtual_cpus = "${var.VIRTUAL_CPUS}"
// 	power_on = "${var.POWER_ON}"
// 	all_eulas_accepted = "${var.ALL_EULAS_ACCEPTED}"
// }

// `

// const testAccVappVmUpdate = `resource "vcloud-director_vapp_vm" "source_vapp_vm"{
//     target_vapp="${var.TARGET_VAPP_NAME}"
//     target_vdc="${var.TARGET_VAPP_VDC}"
//     target_vm_name="${var.TARGET_VM_NAME}"
//     source_vm_name="${var.SOURCE_VM_NAME}"
//     source_catalog_name="${var.SOURCE_CATALOG_NAME}"
//     source_template_name="${var.TEMPLATE_NAME}"
//     network = "${var.NETWORK}"
//     ip_allocation_mode = "${var.VAPP_IP_ALLOCATION_MODE}"
//     hostname = "${var.HOST_NAME}"
//     password = "${var.PASSWORD}"
// 	password_auto = "${var.PASSWORD_AUTO}"
// 	password_reset = "${var.PASSWORD_RESET}"
// 	memory = "${var.MEMORY}"
// 	cores_per_socket = "${var.CORES_PER_SOCKET}"
// 	virtual_cpus = "${var.VIRTUAL_CPUS}"
// 	power_on = "${var.POWER_ON}"
// 	all_eulas_accepted = "${var.ALL_EULAS_ACCEPTED}"
// }

// `
