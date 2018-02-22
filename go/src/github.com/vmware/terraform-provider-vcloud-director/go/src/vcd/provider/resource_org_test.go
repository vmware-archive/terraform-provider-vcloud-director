// terraform-provider-vcloud-director
// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause
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

func TestAccResourceOrg(t *testing.T) {
	logging.Plog("__INIT__TestAccResourceOrg")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrgDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOrg_basic + "\n" + testAccOrg_resource_enable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCreateOrg(),
				),
			},
			resource.TestStep{
				Config: testAccOrg_basic + "\n" + testAccOrg_resource_disable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdateOrg(false),
				),
			},
			resource.TestStep{
				Config: testAccOrg_basic + "\n" + testAccOrg_resource_enable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdateOrg(true),
				),
			},
		},
	})

	logging.Plog("__DONE__TestAccResourceOrg_")
}

func testAccCheckCreateOrg() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logging.Plog("__INIT_testAccCheckCreateOrg__")

		provider := providerGlobalRefPointer.orgProvider

		orgName := os.Getenv("TF_VAR_ORG_NAME")
		orgFullName := os.Getenv("TF_VAR_FULL_NAME")

		isEnabled, _ := strconv.ParseBool(os.Getenv("TF_VAR_IS_ENABLED"))

		//Read org created
		readResp, _ := provider.Read(proto.ReadOrgInfo{Name: orgName})
		if readResp.Present {
			if !(orgName == readResp.Name) {
				return fmt.Errorf("__ERROR__.... Name [expected : %v, found : %v]  do not match", orgName, readResp.Name)
			}

			if !(orgFullName == readResp.OrgFullName) {
				return fmt.Errorf("__ERROR__.... OrgFullName [expected : %v, found : %v]  do not match", orgFullName, readResp.OrgFullName)
			}

			if !(isEnabled == readResp.IsEnabled) {
				return fmt.Errorf("__ERROR__.... IsEnabled [expected %v, found %v]  do not match", isEnabled, readResp.IsEnabled)
			}
		} else {
			return fmt.Errorf("__ERROR__.... org[ %v ] do not exist", orgName)
		}

		logging.Plog("__DONE_testAccCheckCreateOrg__")
		return nil
	}
}

func testAccCheckUpdateOrg(isEnabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logging.Plog("__INIT_testAccCheckUpdateOrg__")

		provider := providerGlobalRefPointer.orgProvider

		name := os.Getenv("TF_VAR_ORG_NAME")

		readResp, readErrp := provider.Read(proto.ReadOrgInfo{Name: name})
		if readErrp != nil {
			return fmt.Errorf("__ERROR__.... in reading org ", readErrp)
		}
		if readResp.Present && (isEnabled == readResp.IsEnabled) {
			logging.Plog("Org update varified")
		} else {
			return fmt.Errorf("__ERROR__.... updating org [%v], org read response[%v]", name, readResp)
		}
		logging.Plog("__DONE_testAccCheckUpdateOrg__")
		return nil
	}
}

func testAccCheckOrgDestroy(s *terraform.State) error {

	logging.Plog("__INIT__testAccCheckOrgDestroy_")

	provider := providerGlobalRefPointer.orgProvider
	orgName := os.Getenv("TF_VAR_ORG_NAME")

	resp, _ := provider.Read(proto.ReadOrgInfo{Name: orgName})

	if resp.Present {
		return fmt.Errorf("__ERROR__.... user[%v] found", orgName)
	}

	logging.Plog("__DONE__testAccCheckOrgDestroy_")
	return nil

}

const testAccOrg_basic = `
provider "vcloud-director" {
  #value come from ENV VARIALES
}

variable "ORG_NAME" {
 type    = "string"
 default = "NOT DEFINED"
}

variable "FULL_NAME" {

 type    = "string"
 default = ""
}

variable "IS_ENABLED" {
 #type    = bool
 default = true
}


variable "IS_DISABLED" {
 #type    = bool
 default = false
}

variable "FORCE" {

 #type    = "bool"
 default = true
}

variable "RECURSIVE" {

 #type    = "bool"
 default = true
}


`

const testAccOrg_resource_enable = `
resource "vcloud-director_org" "o1" {
        name    = "${var.ORG_NAME}"
        full_name = "${var.FULL_NAME}"
        #description = "${var.DESCRIPTION}"
        is_enabled = "${var.IS_ENABLED}"
        force = "${var.FORCE}"
        recursive = "${var.RECURSIVE}"
}

`
const testAccOrg_resource_disable = `

resource "vcloud-director_org" "o1" {
        name    = "${var.ORG_NAME}"
        full_name = "${var.FULL_NAME}"
        #description = "${var.DESCRIPTION}"
        is_enabled = "${var.IS_DISABLED}"
        force = "${var.FORCE}"
        recursive = "${var.RECURSIVE}"
}

`
