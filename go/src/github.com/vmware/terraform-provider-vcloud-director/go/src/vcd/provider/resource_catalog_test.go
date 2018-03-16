/*****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
******************************************************************/

package provider

import (
	"fmt"
	"testing"
	//"time"
	"os"
	"strings"
	//"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	//"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/util/logging"
	"strconv"
)

func TestAccResourceCatalogBasic(t *testing.T) {
	logging.Plog("__INIT_TestAccResourceCatalog_")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCatalogDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCatalog_basic + testAccCatalog_create,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCatalogCreate(),
				),
			},
			resource.TestStep{
				Config: testAccCatalog_basic + testAccCatalog_unshared,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckShareCatalog(false),
				),
			},
			resource.TestStep{
				Config: testAccCatalog_basic + testAccCatalog_shared,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckShareCatalog(true),
				),
			},
			resource.TestStep{
				Config: testAccCatalog_basic + testAccCatalog_update_all_fileds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdateAllFileds(),
				),
			},
		},
	})

	logging.Plog("__DONE__TestAccResourceCatalog_")
}

func testAccCheckCatalogCreate() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logging.Plog("__INIT__testAccCheckCatalogCreate_")
		m := testAccProvider.Meta()

		vcdClient := m.(*VCDClient)

		provider := vcdClient.getProvider()

		catalog, isPreErr := provider.ReadCatalog(os.Getenv("TF_VAR_CATALOG_NAME_OLD"))
		logging.Plog(fmt.Sprintf("__LOG__Read catalog  [%#v]", catalog))

		if isPreErr != nil {
			return fmt.Errorf("__ERROR__.... in validating catalog  creation")
		}
		if !catalog.Present {
			return fmt.Errorf("__ERROR__.... Catalog  NOT created as expected, catalog[%#v]", catalog)
		}

		name := os.Getenv("TF_VAR_CATALOG_NAME_OLD")
		if !(name == catalog.Name) {
			return fmt.Errorf("__ERROR__.... name do not match [expected: %v, found: %v]", name, catalog.Name)
		}

		desc := os.Getenv("TF_VAR_CATALOG_DESCRIPTION_1")

		if strings.Compare(catalog.Description, desc) != 0 {
			return fmt.Errorf("ERROR.... Catalog  Description  NOT as expected")
		}

		shared, _ := strconv.ParseBool(os.Getenv("TF_VAR_CATALOG_SHARED"))
		if !(shared == catalog.Shared) {
			return fmt.Errorf("__ERROR__.... shared do not match [expected: %v, found: %v]", shared, catalog.Shared)
		}

		logging.Plog("__DONE__testAccCheckCatalogCreate_")
		return nil

	}
}

func testAccCheckUpdateAllFileds() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logging.Plog("__INIT_testAccCheckUpdateAllFileds__")

		name := os.Getenv("TF_VAR_CATALOG_NAME_NEW")
		description := os.Getenv("TF_VAR_CATALOG_DESCRIPTION_2")
		shared, _ := strconv.ParseBool(os.Getenv("TF_VAR_CATALOG_SHARED"))
		shared = !shared

		m := testAccProvider.Meta()
		vcdClient := m.(*VCDClient)
		provider := vcdClient.getProvider()
		readResp, readErrp := provider.ReadCatalog(name)

		if readErrp != nil {
			return fmt.Errorf("__ERROR__.... in reading Catalog ", readErrp)
		}
		if readResp.Present {
			if !(name == readResp.Name) {
				return fmt.Errorf("__ERROR__.... name do not match [expected: %v, found: %v]", name, readResp.Name)
			}

			if !(description == readResp.Description) {
				return fmt.Errorf("__ERROR__.... description do not match [expected: %v, found: %v]", description, readResp.Description)
			}

			if !(shared == readResp.Shared) {
				return fmt.Errorf("__ERROR__.... shared do not match [expected: %v, found: %v]", shared, readResp.Shared)
			}
			logging.Plog("Catalog update varified")
		} else {
			return fmt.Errorf("__ERROR__.... updating Catalog[%v]. Catalog not found, Catalog read response[%v]", name, readResp)
		}
		logging.Plog("__DONE_testAccCheckUpdateAllFileds__")
		return nil
	}
}

func testAccCheckShareCatalog(shared bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logging.Plog("__INIT_testAccCheckShareCatalog__")
		logging.Plog(fmt.Sprintf("shared[%v] ", shared))

		name := os.Getenv("TF_VAR_CATALOG_NAME_OLD")

		m := testAccProvider.Meta()
		vcdClient := m.(*VCDClient)
		provider := vcdClient.getProvider()
		readResp, readErrp := provider.ReadCatalog(name)

		if readErrp != nil {
			return fmt.Errorf("__ERROR__.... in reading Catalog ", readErrp)
		}
		if readResp.Present {
			if !(shared == readResp.Shared) {
				return fmt.Errorf("__ERROR__.... shared do not match [expected: %v, found: %v]", shared, readResp.Shared)
			}
			logging.Plog("Catalog update varified")
		} else {
			return fmt.Errorf("__ERROR__.... updating Catalog[%v]. Catalog not found, Catalog read response[%v]", name, readResp)
		}
		logging.Plog("__DONE_testAccCheckShareCatalog__")
		return nil
	}
}

func testAccCheckCatalogDestroy(s *terraform.State) error {

	logging.Plog("__INIT__testAccCheckCatalogDestroy_")

	m := testAccProvider.Meta()

	vcdClient := m.(*VCDClient)

	provider := vcdClient.getProvider()

	isPresResp, isPreErr := provider.ReadCatalog(os.Getenv("TF_VAR_CATALOG_NAME_NEW"))

	if isPreErr != nil {
		return fmt.Errorf("__ERROR__ in validating catalog  deletion")
	}
	if isPresResp.Present {
		return fmt.Errorf("__ERROR__ Catlog  not deleted as expected")
	}

	logging.Plog("__DONE__testAccCheckCatalogDestroy_")
	return nil

}

const testAccCatalog_basic = `
variable "CATALOG_DESCRIPTION_1" { 

 type    = "string"
 default = "NOT DEFINED" 
}

variable "CATALOG_DESCRIPTION_2" { 

 type    = "string"
 default = "NOT DEFINED" 
}

variable "CATALOG_NAME_OLD" { 

 type    = "string"
 default = "NOT DEFINED" 
}

variable "CATALOG_NAME_NEW" { 

 type    = "string"
 default = "NOT DEFINED" 
}

variable "CATALOG_SHARED" { 

 default = false
}

provider "vcloud-director" {
  #value come from ENV VARIALES
}

`

const testAccCatalog_create = `
resource "vcloud-director_catalog" "catalog1" {
        name    ="${var.CATALOG_NAME_OLD}"
        description = "${var.CATALOG_DESCRIPTION_1}"
        shared  = true
}

`
const testAccCatalog_shared = `
resource "vcloud-director_catalog" "catalog1" {
        name    ="${var.CATALOG_NAME_OLD}"
        description = "${var.CATALOG_DESCRIPTION_1}"
        shared  = true
}

`

const testAccCatalog_unshared = `
resource "vcloud-director_catalog" "catalog1" {
        name    ="${var.CATALOG_NAME_OLD}"
        description = "${var.CATALOG_DESCRIPTION_1}"
        shared  = false
}

`

const testAccCatalog_update_all_fileds = `
resource "vcloud-director_catalog" "catalog1" {
        name    ="${var.CATALOG_NAME_NEW}"
        description = "${var.CATALOG_DESCRIPTION_2}"
        shared  = false
}

`
