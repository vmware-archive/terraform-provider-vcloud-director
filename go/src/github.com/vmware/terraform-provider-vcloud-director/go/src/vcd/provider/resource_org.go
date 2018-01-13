/*****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
******************************************************************/

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/util/logging"
	//"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
)

func resourceOrg() *schema.Resource {
	return &schema.Resource{
		Create: resourceOrgCreate,
		Read:   resourceOrgRead,
		Update: resourceOrgUpdate,
		Delete: resourceOrgDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"full_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func getOrgId(d *schema.ResourceData) string {

	OrgName := d.Get("name").(string)

	return OrgName
}

func getOrgInfo(d *schema.ResourceData) (string, string) {
	name := d.Get("name").(string)
	full_name := d.Get("full_name").(string)

	return name, full_name
}

func resourceOrgCreate(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT__resourceOrgCreate_")

	logging.Plog(" __DONE__resourceOrgCreate_")
	return nil
}

func resourceOrgRead(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT__resourceOrgRead_ ")

	logging.Plog("__DONE__ resourceOrgRead_")
	return nil

}

func resourceOrgUpdate(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__NOT IMPL_resourceOrgUpdate_")

	return fmt.Errorf("__ERROR__Not Updating NAME , NOT IMPLEMENTED !!!!")
}

func resourceOrgDelete(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceOrgDelete_")

	logging.Plog("__DONE__resourceOrgDelete_")
	return nil

}
