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
)

func TestAccResourceCatalogBasic(t *testing.T) {
	logging.Plog("__INIT_TestAccResourceCatalog_")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCatalogDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCatalog_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCatalogCreate(),
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

		catalog, isPreErr := provider.ReadCatalog(os.Getenv("TF_VAR_CATALOG_NAME"))

		if isPreErr != nil {
			return fmt.Errorf("__ERROR__.... in validating catalog  creation")
		}
		if !catalog.Present {
			return fmt.Errorf("__ERROR__.... Catalog  NOT created as expected")
		}
		logging.Plog(fmt.Sprintf("__LOG__Read catalog  [%#v]", catalog))

		desc := os.Getenv("TF_VAR_CATALOG_DESCRIPTION")

		if strings.Compare(catalog.Description, desc) != 0 {
			return fmt.Errorf("ERROR.... Catalog  Description  NOT as expected")
		}

		logging.Plog("__DONE__testAccCheckCatalogCreate_")
		return nil

	}
}

func testAccCheckCatalogDestroy(s *terraform.State) error {

	logging.Plog("__INIT__testAccCheckCatalogDestroy_")

	m := testAccProvider.Meta()

	vcdClient := m.(*VCDClient)

	provider := vcdClient.getProvider()

	isPresResp, isPreErr := provider.ReadCatalog("testcata_acc1")

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



variable "CATALOG_DESCRIPTION" { 

 type    = "string"
 default = "NOT DEFINED" 
}


variable "CATALOG_NAME" { 

 type    = "string"
 default = "NOT DEFINED" 
}



provider "vcloud-director" {
  #value come from ENV VARIALES
}




resource "vcloud-director_catalog" "catalog1" {
        name    ="${var.CATALOG_NAME}"
        description = "${var.CATALOG_DESCRIPTION}"
        shared  = "true"

}


`
