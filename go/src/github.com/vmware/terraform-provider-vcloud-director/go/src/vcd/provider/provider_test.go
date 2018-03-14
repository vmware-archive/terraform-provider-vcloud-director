/*****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
******************************************************************/
package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/vmware/terraform-provider-vcloud-director/go/src/util/logging"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]terraform.ResourceProvider

func init() {

	logging.Plog("__INIT__init")
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"vcloud-director": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	logging.Plog("__INIT__TestProvider_impl_")

	//var _ terraform.ResourceProvider = Provider()

}

func testAccPreCheck(t *testing.T) {
	logging.Init()
	logging.Plog("__INIT__testAccPreCheck_")
	logging.Plog("__DONE__testAccPreCheck_")
}
