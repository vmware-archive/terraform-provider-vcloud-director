/*****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
******************************************************************/

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

func TestAccResourceCatalogItemOva(t *testing.T) {
	logging.Plog("__INIT__TestAccResourceCatalogItemOva_")
	maintf := fmt.Sprintf("%v \n %v", testAccCatalogItemOva_basic, testAccCatalogItemOva_LocalFile)
	logging.Plogf("maintf \n %v", maintf)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCatalogItemOvaDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: maintf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCatalogItem(os.Getenv("TF_VAR_CATALOG_NAME"), "item1"),
				),
			},
		},
	})

	logging.Plog("__DONE_TestAccResourceCatalogItemOva_")
}

func TestAccResourceCatalogItemCaptureVapp(t *testing.T) {
	logging.Plog("__INIT__TestAccResourceCatalogItemCaptureVapp_")
	maintf := fmt.Sprintf("%v \n %v", testAccCatalogItemOva_basic, testAccCatalogItemOva_CaptureVapp)
	logging.Plogf("%v", maintf)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCatalogItemCaptureDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: maintf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCatalogItem(os.Getenv("TF_VAR_CATALOG_NAME"), "capturedItem"),
				),
			},
		},
	})

	logging.Plog("__DONE_TestAccResourceCatalogItemCaptureVapp_")
}

func testAccCheckCatalogItem(catalogName string, itemName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logging.Plog(fmt.Sprintf("__INIT__testAccCheckCatalogItem_ %v %v", catalogName, itemName))
		m := testAccProvider.Meta()

		vcdClient := m.(*VCDClient)

		provider := vcdClient.getProvider()
		logging.Plog(fmt.Sprintf("%#v", provider))
		//TEST against VCD that the catalog item is uploaded

		catalogIsPresInfo := proto.IsPresentCatalogItemInfo{CatalogName: catalogName, ItemName: itemName}
		isPresResp, isPreErr := provider.IsPresentCatalogItem(catalogIsPresInfo)
		logging.Plog(fmt.Sprintf("%+v %+v", isPresResp, isPreErr))
		//time.Sleep(20000 * time.Millisecond)
		if isPreErr != nil {
			return fmt.Errorf("Error in validating catalog item creation")
		}
		if !isPresResp.Present {
			return fmt.Errorf("Error Catlog Item not created as expected")
		}

		logging.Plog("__DONE_testAccCheckCatalogItemUpload_")
		return nil

	}
}

func checkDestroyItem(s *terraform.State, catalogName string, itemName string) error {

	m := testAccProvider.Meta()

	vcdClient := m.(*VCDClient)

	provider := vcdClient.getProvider()
	logging.Plogf("%#v", provider)
	//TEST against VCD that the catalog item is DELETED

	catalogIsPresInfo := proto.IsPresentCatalogItemInfo{CatalogName: catalogName, ItemName: itemName}
	isPresResp, isPreErr := provider.IsPresentCatalogItem(catalogIsPresInfo)

	if isPreErr != nil {
		return fmt.Errorf("Error in validating catalog item deletion")
	}
	if isPresResp.Present {
		return fmt.Errorf("Error Catlog Item not deleted as expected")
	}
	return nil
}
func testAccCheckCatalogItemCaptureDestroy(s *terraform.State) error {

	return checkDestroyItem(s, os.Getenv("TF_VAR_CATALOG_NAME"), "capturedItem")

}

func testAccCheckCatalogItemOvaDestroy(s *terraform.State) error {

	return checkDestroyItem(s, os.Getenv("TF_VAR_CATALOG_NAME"), "item1")

}

const testAccCatalogItemOva_basic = `

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

const testAccCatalogItemOva_LocalFile = `

variable "OVA_PATH" { 

 type    = "string"
 default = "nullova" 
}


resource "vcloud-director_catalog_item_ova" "item1" {
	item_name = "item1"
	catalog_name= "${vcloud-director_catalog.catalog1.name}"
	source_file_path="${var.OVA_PATH}"
}
`
const testAccCatalogItemOva_CaptureVapp = `

variable "SOURCE_VDC_NAME" { 

 type    = "string"
 default = "notdefined_vdcname" 
}

variable "SOURCE_VAPP_NAME" { 

 type    = "string"
 default = "notdefined_vappname" 
}


resource "vcloud-director_catalog_item_ova" "item2" {
	item_name = "capturedItem"
	catalog_name= "${vcloud-director_catalog.catalog1.name}"
	source_vdc_name="${var.SOURCE_VDC_NAME}"
	source_vapp_name="${var.SOURCE_VAPP_NAME}"
}

`

//variable "OVA_PATH" {'/Users/srinarayana/vmws/tiny.ova'}

//
//
