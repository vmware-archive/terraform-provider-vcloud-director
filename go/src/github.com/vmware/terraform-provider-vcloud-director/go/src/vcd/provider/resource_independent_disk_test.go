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

func TestAccResourceIndependentDiskBasic(t *testing.T) {
	logging.Plog("__INIT_TestAccResourceIndependentDisk_")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIndependentDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccIndependentDisk_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndependentDiskCreate(),
				),
			},
		},
	})

	logging.Plog("__DONE__TestAccResourceIndependentDisk_")
}

func testAccCheckIndependentDiskCreate() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logging.Plog("__INIT__testAccCheckIndependentDiskCreate_")
		provider := providerGlobalRefPointer.independentDiskProvider
		disk := proto.ReadDiskInfo{Name: os.Getenv("TF_VAR_DISK_NAME")}
		independentDisk, isPreErr := provider.Read(disk)

		if isPreErr != nil {
			return fmt.Errorf("__ERROR__.... in validating IndependentDisk  creation")
		}
		if !independentDisk.Present {
			return fmt.Errorf("__ERROR__.... IndependentDisk  NOT created as expected")
		}
		logging.Plog(fmt.Sprintf("__LOG__Read IndependentDisk  [%#v]", independentDisk))

		/*
			desc := os.Getenv("TF_VAR_IndependentDisk_DESCRIPTION")

			if strings.Compare(independentDisk.Description, desc) != 0 {
				return fmt.Errorf("ERROR.... IndependentDisk  Description  NOT as expected")
			}
		*/
		logging.Plog("__DONE__testAccCheckIndependentDiskCreate_")
		return nil

	}
}

func testAccCheckIndependentDiskDestroy(s *terraform.State) error {

	logging.Plog("__INIT__testAccCheckIndependentDiskDestroy_")

	logging.Plog("__DONE__testAccCheckIndependentDiskDestroy_")
	return nil

}

const testAccIndependentDisk_basic = `





variable "DISK_NAME" { 

 type    = "string"
 default = "NOT DEFINED" 
}



provider "vcloud-director" {
  #value come from ENV VARIALES
}




resource "vcloud-director_independent_disk" "IndependentDisk1" {
        name    ="${var.DISK_NAME}"
        size 	= "100"
        VDC="OVD4"

}

`
