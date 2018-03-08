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

func TestAccResourceVappVm(t *testing.T) {
	logging.Plog("__INIT__TestAccResourceVappVm")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVappVmDestroy,
		Steps: []resource.TestStep{

			resource.TestStep{
				Config: testAccVappVm_basic + "\n" + testAccVappVm_enable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCreateVappVm(),
				),
			},

			resource.TestStep{
				Config: testAccVappVm_basic + "\n" + testAccVappVm_disable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdateVappVm(),
				),
			},
		},
	})

	logging.Plog("__DONE__TestAccResourceVappVm_")
}

func testAccCheckCreateVappVm() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logging.Plog("__INIT_testAccCheckCreateVappVm__")

		provider := providerGlobalRefPointer.vappVmProvider

		name := os.Getenv("TF_VAR_VappVm_NAME")
		roleName := os.Getenv("TF_VAR_ROLE_NAME")
		fullName := os.Getenv("TF_VAR_FULL_NAME")
		email := os.Getenv("TF_VAR_EMAIL")
		telephone := os.Getenv("Telephone")
		im := os.Getenv("TF_VAR_IM")

		storedVmQuota, _ := strconv.ParseInt(os.Getenv("TF_VAR_STORED_VM_QUOTA"))
		storedVmQuota32 := int32(storedVmQuota)

		deployedVmQuota, _ := strconv.ParseInt(os.Getenv("TF_VAR_DEPLOYED_VN_QUOTA"))
		deployedVmQuota32 := int32(deployedVmQuota)

		isGroupRole, _ := strconv.ParseBool(os.Getenv("TF_VAR_IS_GROUP_ROLE"))
		isExternal, _ := strconv.ParseBool(os.Getenv("TF_VAR_IS_EXTERNAL"))
		isEnabled, _ := strconv.ParseBool(os.Getenv("TF_VAR_IS_ENABLED_TRUE"))

		readResp, readErrp := provider.Read(proto.ReadVappVmInfo{Name: name})

		if readErrp != nil {
			return fmt.Errorf("__ERROR__.... in reading VappVm  creation", readErrp)
		}
		if readResp.Present {
			if !(name == readResp.Name) {
				return fmt.Errorf("__ERROR__.... name [expected %v, found %v]  do not match", name, readResp.Name)
			}

			if !(roleName == readResp.RoleName) {
				return fmt.Errorf("__ERROR__.... roleName [expected %v, found %v]  do not match", roleName, readResp.RoleName)
			}

			if !(fullName == readResp.FullName) {
				return fmt.Errorf("__ERROR__.... fullName [expected %v, found %v]  do not match", fullName, readResp.FullName)
			}

			if !(email == readResp.Email) {
				return fmt.Errorf("__ERROR__.... email [expected %v, found %v]  do not match", email, readResp.Email)
			}

			if !(telephone == readResp.Telephone) {
				return fmt.Errorf("__ERROR__.... VappVm telephone [expected %v, found %v]  do not match", telephone, readResp.Telephone)
			}

			if !(im == readResp.Im) {
				return fmt.Errorf("__ERROR__.... im [expected %v, found %v]  do not match", im, readResp.Im)
			}

			if !(storedVmQuota32 == readResp.StoredVmQuota) {
				return fmt.Errorf("__ERROR__.... storedVmQuota [expected %v, found %v]  do not match", storedVmQuota, readResp.StoredVmQuota)
			}

			if !(deployedVmQuota32 == readResp.DeployedVmQuota) {
				return fmt.Errorf("__ERROR__.... deployedVmQuota [expected %v, found %v]  do not match", deployedVmQuota, readResp.DeployedVmQuota)
			}

			if !(isGroupRole == readResp.IsGroupRole) {
				return fmt.Errorf("__ERROR__.... isGroupRole [expected %v, found %v]  do not match", isGroupRole, readResp.IsGroupRole)
			}

			if !(isExternal == readResp.IsExternal) {
				return fmt.Errorf("__ERROR__.... isExternal [expected %v, found %v]  do not match", isExternal, readResp.IsExternal)
			}

			if !(isEnabled == readResp.IsEnabled) {
				return fmt.Errorf("__ERROR__.... isEnabled [expected %v, found %v]  do not match", isEnabled, readResp.IsEnabled)
			}

			logging.Plog("VappVm creation varified")
		} else {
			return fmt.Errorf("__ERROR__.... VappVm[%v]  not found", name)
		}
		logging.Plog("__DONE_testAccCheckCreateVappVm__")
		return nil
	}
}

func testAccCheckUpdateVappVm() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logging.Plog("__INIT_testAccCheckUpdateVappVm__")

		provider := providerGlobalRefPointer.vappVmProvider

		name := os.Getenv("TF_VAR_VappVm_NAME")

		readResp, readErrp := provider.Read(proto.ReadVappVmInfo{Name: name})
		if readErrp != nil {
			return fmt.Errorf("__ERROR__.... in reading VappVm ", readErrp)
		}
		if readResp.Present && !readResp.IsEnabled {
			logging.Plog("VappVm update varified")
		} else {
			return fmt.Errorf("__ERROR__.... updating VappVm[%v], VappVm read response[%v]", name, readResp)
		}
		logging.Plog("__DONE_testAccCheckUpdateVappVm__")
		return nil
	}
}

func testAccCheckVappVmDestroy(s *terraform.State) error {

	logging.Plog("__INIT__testAccCheckVappVmDestroy_")

	provider := providerGlobalRefPointer.vappVmProvider

	name := os.Getenv("TF_VAR_VappVm_NAME")

	readResp, readErrp := provider.Read(proto.ReadVappVmInfo{Name: name})

	if readResp.Present {
		return fmt.Errorf("__ERROR__.... VappVm[%v] found", name)
	}

	logging.Plog("__DONE__testAccCheckVappVmDestroy_")
	return nil

}

const testAccVappVm_basic = `
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
`

const testAccVappVm_enable = `resource "vcloud-director_VappVm" "source_VappVm"{
        name = "${var.VappVm_NAME}"
		password = "${var.VappVm_PASSWORD}"
		role_href = "${var.ROLE_NAME}"
		full_name = "${var.FULL_NAME}"
		description = "${var.DESCRIPTION}"
		email = "${var.EMAIL}"
		telephone = "${var.TELEPHONE}"
		im = "${var.IM}"
		alert_email = "${var.ALERT_EMAIL}"
		alert_email_prefix = "${var.ALERT_EMAIL_PREFIX}"
		stored_vm_quota = "${var.STORED_VM_QUOTA}"
		deployed_vm_quota = "${var.DEPLOYED_VN_QUOTA}"
		is_group_role = "${var.IS_GROUP_ROLE}"
		is_default_cached = "${var.IS_DEFAULT_CACHED}"
		is_external = "${var.IS_EXTERNAL}"
		is_alert_enabled = "${var.IS_ALERT_ENABLED}"
		is_enabled = "${var.IS_ENABLED_TRUE}"
}
`

const testAccVappVm_disable = `resource "vcloud-director_VappVm" "source_VappVm"{
        name = "${var.VappVm_NAME}"
		password = "${var.VappVm_PASSWORD}"
		role_name = "${var.ROLE_NAME}"
		full_name = "${var.FULL_NAME}"
		description = "${var.DESCRIPTION}"
		email = "${var.EMAIL}"
		telephone = "${var.TELEPHONE}"
		im = "${var.IM}"
		alert_email = "${var.ALERT_EMAIL}"
		alert_email_prefix = "${var.ALERT_EMAIL_PREFIX}"
		stored_vm_quota = "${var.STORED_VM_QUOTA}"
		deployed_vm_quota = "${var.DEPLOYED_VN_QUOTA}"
		is_group_role = "${var.IS_GROUP_ROLE}"
		is_default_cached = "${var.IS_DEFAULT_CACHED}"
		is_external = "${var.IS_EXTERNAL}"
		is_alert_enabled = "${var.IS_ALERT_ENABLED}"
		is_enabled = "${var.IS_ENABLED_FALSE}"
}
`
