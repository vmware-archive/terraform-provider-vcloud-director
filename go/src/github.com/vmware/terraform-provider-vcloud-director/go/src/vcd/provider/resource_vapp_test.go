// terraform-provider-vcloud-director
// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause
package provider

import (
	"fmt"
	"os"
	"testing"
	//"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/util/logging"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
)

func TestAccResourceVApp(t *testing.T) {
	logging.Plog("__INIT__TestAccResourceVApp")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVAppDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVApp_basic + testAccVdc_create,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCreateVApp(),
				),
			},
			resource.TestStep{
				Config: testAccVApp_basic + testAccVdc_power_off,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdateVApp(),
				),
			},
			resource.TestStep{
				Config: testAccVApp_basic + testAccVdc_power_on,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdateVApp(),
				),
			},
		},
	})

	logging.Plog("__DONE__TestAccResourceVApp_")
}

func testAccCheckCreateVApp() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logging.Plog("__INIT_testAccCheckCreateVAPP__")
		m := testAccProvider.Meta()

		vcdClient := m.(*VCDClient)
		provider := vcdClient.getProvider()
		name := os.Getenv("TF_VAR_VAPP_NAME")
		vdc := os.Getenv("TF_VAR_VAPP_VDC")

		logging.Plogf("__LOG__ vapp name = [%v]", name)
		resp, errp := provider.ReadVApp(proto.ReadVAppInfo{Name: name, Vdc: vdc})
		if errp != nil {
			return fmt.Errorf("__ERROR__.... in validating vapp  creation", errp)
		}
		logging.Plogf("THIS SHOULD BE RESP %#v", *resp)
		logging.Plog("__INIT_testAccCheckCreateVAPP__")
		return nil
	}
}

func testAccCheckUpdateVApp() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logging.Plog("__INIT_testAccCheckUpdateVApp__")

		//Update when read power on, off API available

		logging.Plog("__DONE_testAccCheckUpdateVApp__")
		return nil
	}
}

func testAccCheckVAppDestroy(s *terraform.State) error {

	logging.Plog("__INIT__testAccCheckVAppDestroy_")

	m := testAccProvider.Meta()

	vcdClient := m.(*VCDClient)
	vcdClient.getProvider()

	//TODO check with read read
	logging.Plog("__DONE__testAccCheckVAppDestroy_")
	return nil

}

const testAccVApp_basic = `
variable "VAPP_NAME" { 

 type    = "string"
 default = "NOT DEFINED" 
}

variable "VAPP_TEMPLATE_NAME" { 

 type    = "string"
 default = "NOT DEFINED" 
}

variable "VAPP_CATALOG_NAME" { 

 type    = "string"
 default = "NOT DEFINED" 
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

provider "vcloud-director" {
 
  allow_unverified_ssl = "true"
}

`

const testAccVdc_create = `
resource "vcloud-director_vapp" "vapp1" {
        name  					= "${var.VAPP_NAME}"
        template_name 			= "${var.VAPP_TEMPLATE_NAME}"
        catalog_name  			= "${var.VAPP_CATALOG_NAME}"
        vdc 					= "${var.VAPP_VDC}"
        network 				= "${var.VAPP_NETWORK}"
        ip_allocation_mode 		= "${var.VAPP_IP_ALLOCATION_MODE}"
        cpu 					= "${var.VAPP_CPU}"
        memory 					= "${var.VAPP_MEMORY}"
        power_on				= true
}

`

const testAccVdc_power_on = `
resource "vcloud-director_vapp" "vapp1" {
        name  					= "${var.VAPP_NAME}"
        template_name 			= "${var.VAPP_TEMPLATE_NAME}"
        catalog_name  			= "${var.VAPP_CATALOG_NAME}"
        vdc 					= "${var.VAPP_VDC}"
        network 				= "${var.VAPP_NETWORK}"
        ip_allocation_mode 		= "${var.VAPP_IP_ALLOCATION_MODE}"
        cpu 					= "${var.VAPP_CPU}"
        memory 					= "${var.VAPP_MEMORY}"
        power_on				= true
}

`

const testAccVdc_power_off = `
resource "vcloud-director_vapp" "vapp1" {
        name  					= "${var.VAPP_NAME}"
        template_name 			= "${var.VAPP_TEMPLATE_NAME}"
        catalog_name  			= "${var.VAPP_CATALOG_NAME}"
        vdc 					= "${var.VAPP_VDC}"
        network 				= "${var.VAPP_NETWORK}"
        ip_allocation_mode 		= "${var.VAPP_IP_ALLOCATION_MODE}"
        cpu 					= "${var.VAPP_CPU}"
        memory 					= "${var.VAPP_MEMORY}"
	    power_on				= false

}

`
