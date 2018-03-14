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
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
	"strings"
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
			"is_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
			"force": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
			"recursive": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func getOrgId(d *schema.ResourceData) string {
	//We are using orgName as id
	orgName := d.Get("name").(string)
	return orgName
}

func resourceOrgCreate(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT__resourceOrgCreate_")

	orgName := d.Get("name").(string)
	orgFullName := d.Get("full_name").(string)
	isEnabled := d.Get("is_enabled").(bool)
	provider := providerGlobalRefPointer.orgProvider

	createOrgInfo := proto.CreateOrgInfo{Name: orgName, OrgFullName: orgFullName, IsEnabled: isEnabled}
	res, err := provider.Create(createOrgInfo)

	if err != nil {
		return fmt.Errorf("Error Creating Org :[%+v] %#v", createOrgInfo, err)
	}

	if res.Created {
		logging.Plog(fmt.Sprintf("Org [%+v]  created  ", res))
		d.SetId(orgName)
	}

	logging.Plog("__DONE__resourceOrgCreate_")
	return nil
}

//Delete org
func resourceOrgDelete(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceOrgDelete_")
	orgName := d.Get("name").(string)
	force := d.Get("force").(bool)
	recursive := d.Get("recursive").(bool)
	provider := providerGlobalRefPointer.orgProvider

	deleteOrgInfo := proto.DeleteOrgInfo{Name: orgName, Force: force, Recursive: recursive}
	res, err := provider.Delete(deleteOrgInfo)

	if err != nil {
		return fmt.Errorf("Error Deleting Org :[%+v] %#v", deleteOrgInfo, err)
	}

	if res.Deleted {
		logging.Plog(fmt.Sprintf("Org [%+v]  deleted  ", res))
		//d.SetId(orgName)
	}

	logging.Plog("__DONE__resourceOrgDelete_")
	return nil
}

//Currently API's not implemented in pyVCloud
func resourceOrgUpdate(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceOrgUpdate_")

	orgName := d.Get("name").(string)
	oldOrgFullNameRaw, newOrgFullNameRaw := d.GetChange("full_name")
	oldOrgFullName := oldOrgFullNameRaw.(string)
	newOrgFullName := newOrgFullNameRaw.(string)
	isEnabled := d.Get("is_enabled").(bool)

	if !strings.EqualFold(oldOrgFullName, newOrgFullName) {
		return fmt.Errorf("__ERROR__ Not Updating org_full_name , API NOT IMPLEMENTED !!!!")
	}

	provider := providerGlobalRefPointer.orgProvider
	updateOrgInfo := proto.UpdateOrgInfo{Name: orgName, OrgFullName: newOrgFullName, IsEnabled: isEnabled}
	res, err := provider.Update(updateOrgInfo)

	if err != nil {
		return fmt.Errorf("Error Updating Org :[%+v] %#v", updateOrgInfo, err)
	}

	if res.Updated {
		logging.Plog(fmt.Sprintf("Org [%+v]  updated  ", res))
		//Set updated Id, discuss with Sri
	}
	logging.Plog("__DONE__resourceOrgUpdate_")
	return nil
}

func resourceOrgRead(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceOrgRead_ ")

	orgName := d.Get("name").(string)

	provider := providerGlobalRefPointer.orgProvider
	readOrgInfo := proto.ReadOrgInfo{Name: orgName}
	res, _ := provider.Read(readOrgInfo)

	//if err != nil {
	//return fmt.Errorf("Error while getting Org :[%+v] %#v", readOrgInfo, err)
	//}

	if res.Present {
		d.SetId(res.Name)
		d.Set("is_enabled", res.IsEnabled)
		logging.Plog(fmt.Sprintf("__DONE__resourceOrgRead_ +setting id %v", orgName))
		return nil
	} else {
		d.SetId("")
		logging.Plog(fmt.Sprintf("__DONE__resourceOrgRead_ +unsetting id,resource got deleted %v", orgName))
	}
	logging.Plog("__DONE__ resourceOrgRead_")
	return nil

}
