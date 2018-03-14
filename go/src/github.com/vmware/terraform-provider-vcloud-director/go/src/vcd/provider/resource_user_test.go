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

func TestAccResourceUser(t *testing.T) {
	logging.Plog("__INIT__TestAccResourceUser")

	conf := testAccUser_basic + "\n" + testAccUser_enable
	logging.Plog("conf : \n" + conf)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{

			resource.TestStep{
				Config: testAccUser_basic + "\n" + testAccUser_enable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCreateUser(),
				),
			},

			resource.TestStep{
				Config: testAccUser_basic + "\n" + testAccUser_disable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdateUser(),
				),
			},
		},
	})

	logging.Plog("__DONE__TestAccResourceUser_")
}

func testAccCheckCreateUser() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logging.Plog("__INIT_testAccCheckCreateUser__")

		provider := providerGlobalRefPointer.userProvider

		name := os.Getenv("TF_VAR_USER_NAME")
		roleName := os.Getenv("TF_VAR_ROLE_NAME")
		fullName := os.Getenv("TF_VAR_FULL_NAME")
		email := os.Getenv("TF_VAR_EMAIL")
		telephone := os.Getenv("TF_VAR_TELEPHONE")
		im := os.Getenv("TF_VAR_IM")
		storedVmQuota, _ := strconv.Atoi(os.Getenv("TF_VAR_STORED_VM_QUOTA"))
		storedVmQuota32 := int32(storedVmQuota)

		deployedVmQuota, _ := strconv.Atoi(os.Getenv("TF_VAR_DEPLOYED_VN_QUOTA"))
		deployedVmQuota32 := int32(deployedVmQuota)

		isGroupRole, _ := strconv.ParseBool(os.Getenv("TF_VAR_IS_GROUP_ROLE"))
		isExternal, _ := strconv.ParseBool(os.Getenv("TF_VAR_IS_EXTERNAL"))
		isEnabled, _ := strconv.ParseBool(os.Getenv("TF_VAR_IS_ENABLED_TRUE"))

		readResp, readErrp := provider.Read(proto.ReadUserInfo{Name: name})

		if readErrp != nil {
			return fmt.Errorf("__ERROR__.... in reading user  creation", readErrp)
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
				return fmt.Errorf("__ERROR__.... user telephone [expected %v, found %v]  do not match", telephone, readResp.Telephone)
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

			logging.Plog("User creation varified")
		} else {
			return fmt.Errorf("__ERROR__.... user[%v]  not found", name)
		}
		logging.Plog("__DONE_testAccCheckCreateUser__")
		return nil
	}
}

func testAccCheckUpdateUser() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logging.Plog("__INIT_testAccCheckUpdateUser__")

		provider := providerGlobalRefPointer.userProvider

		name := os.Getenv("TF_VAR_USER_NAME")

		readResp, readErrp := provider.Read(proto.ReadUserInfo{Name: name})
		if readErrp != nil {
			return fmt.Errorf("__ERROR__.... in reading user ", readErrp)
		}
		if readResp.Present && !readResp.IsEnabled {
			logging.Plog("User update varified")
		} else {
			return fmt.Errorf("__ERROR__.... updating user[%v], user read response[%v]", name, readResp)
		}
		logging.Plog("__DONE_testAccCheckUpdateUser__")
		return nil
	}
}

func testAccCheckUserDestroy(s *terraform.State) error {

	logging.Plog("__INIT__testAccCheckUSerDestroy_")

	provider := providerGlobalRefPointer.userProvider

	name := os.Getenv("TF_VAR_USER_NAME")

	readResp, _ := provider.Read(proto.ReadUserInfo{Name: name})

	if readResp.Present {
		return fmt.Errorf("__ERROR__.... user[%v] found", name)
	}

	logging.Plog("__DONE__testAccCheckUserDestroy_")
	return nil

}

const testAccUser_basic = `
	provider "vcloud-director" {
		  #value come from ENV VARIALES
	}
	variable "USER_NAME" {
		 type    = "string"
		 default = "NOT DEFINED"
	}
	variable "USER_PASSWORD" {
		 type    = "string"
		 default = ""
	}
	variable "ROLE_NAME" {
		 type    = "string"
		 default = ""
	}
	variable "FULL_NAME" {
		 type    = "string"
		 default = ""
	}
	variable "DESCRIPTION" {
		 type    = "string"
		 default = ""
	}
	variable "EMAIL" {
		 type    = "string"
		 default = ""
	}
	variable "TELEPHONE" {
		 type    = "string"
		 default = ""
	}
	variable "IM" {
		 type    = "string"
		 default = ""
	}
	variable "ALERT_EMAIL" {
		 type    = "string"
		 default = ""
	}
	variable "ALERT_EMAIL_PREFIX" {
		 type    = "string"
		 default = ""
	}
	variable "STORED_VM_QUOTA" {
		 #type    = "string"
		 default = 0
	}
	variable "DEPLOYED_VN_QUOTA" {
		 #type    = "string"
		 default = 0
	}
	variable "IS_GROUP_ROLE" {
		 default = false
	}
	variable "IS_EXTERNAL" {
		 default = false
	}
	variable "IS_ALERT_ENABLED" {
		 default = false
	}
	variable "IS_ENABLED_TRUE" {
		 default = true
	}
	variable "IS_DEFAULT_CACHED" {
		 default = false
	}
`

const testAccUser_enable = `
resource "vcloud-director_user" "source_user"{
        name = "${var.USER_NAME}"
		password = "${var.USER_PASSWORD}"
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
		is_enabled = "${var.IS_ENABLED_TRUE}"
}
`

const testAccUser_disable = `
resource "vcloud-director_user" "source_user"{
        name = "${var.USER_NAME}"
		password = "${var.USER_PASSWORD}"
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
		#is_enabled = "${var.IS_ENABLED_FALSE}"
		is_enabled = false
}

`
