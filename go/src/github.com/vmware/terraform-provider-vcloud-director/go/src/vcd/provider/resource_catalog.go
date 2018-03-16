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
)

func resourceCatalog() *schema.Resource {
	return &schema.Resource{
		Create: resourceCatalogCreate,
		Read:   resourceCatalogRead,
		Update: resourceCatalogUpdate,
		Delete: resourceCatalogDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"shared": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceCatalogCreate(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT__resourceCatalogCreate_")
	cname := d.Get("name").(string)
	desc := d.Get("description").(string)
	shared := d.Get("shared").(bool)

	//vcdClient := m.(*VCDClient)

	//provider := vcdClient.getProvider()
	//TEST
	provider := providerGlobalRefPointer.pyVcloudProvider

	resp, errp := provider.ReadCatalog(cname)

	if errp != nil {
		return fmt.Errorf("__ERROR__ Creating catalog :[%v] %#v", cname, errp)
	}

	if resp.Present {
		logging.Plog(fmt.Sprintf("catalog %v is present/ setting internal state / this is to allow adding to catalog already present ", cname))
		d.SetId(cname)
		return nil
	}

	catalog := proto.Catalog{Name: cname, Description: desc, Shared: shared}

	logging.Plog(fmt.Sprintf("__LOG__calling create catalog  name=[%v]  description=[%v]  ", cname, desc))
	res, err := provider.CreateCatalog(catalog)

	if err != nil {
		return fmt.Errorf("Error Creating catalog :[%v] %#v", cname, err)
	}
	if res.Created {
		logging.Plog(fmt.Sprintf("catalog [%v]  created  ", cname))
		d.SetId(cname)
	}
	logging.Plog("__DONE__resourceCatalogCreate ")
	return nil
}

func resourceCatalogRead(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT__resourceCatalogRead_")

	cname := d.Get("name").(string)

	vcdClient := m.(*VCDClient)

	provider := vcdClient.getProvider()

	res, err := provider.ReadCatalog(cname)
	if err != nil {
		return fmt.Errorf("Error checking resourceCatalogRead  %#v", err)
	}
	if res.Present {
		logging.Plog(fmt.Sprintf("__LOG__ catalog %v is present / setting state ", cname))
		d.SetId(cname)
	} else {
		// resource catalog not present / clear id for recreate
		d.SetId("")
	}

	logging.Plog("__DONE__resourceCatalogRead_  ")
	return nil

}

func resourceCatalogUpdate(d *schema.ResourceData, m interface{}) error {

	logging.Plog(fmt.Sprintf("__INIT__resourceCatalogUpdate__ "))

	cNameOldRaw, cNameNewRaw := d.GetChange("name")
	cNameOld := cNameOldRaw.(string)
	cNameNew := cNameNewRaw.(string)

	desc := d.Get("description").(string)
	shared := d.Get("shared").(bool)

	vcdClient := m.(*VCDClient)
	provider := vcdClient.getProvider()

	updateCatalogInfo := proto.UpdateCatalogInfo{
		Name:        cNameNew,
		OldName:     cNameOld,
		Description: desc,
		Shared:      shared,
	}
	res, err := provider.UpdateCatalog(updateCatalogInfo)

	if err != nil {
		return fmt.Errorf("Error updating Catalog :[%+v] %#v", updateCatalogInfo, err)
	}
	if res.Updated {
		logging.Plog(fmt.Sprintf("Catalog [%+v]  updated", res))
		d.SetId(cNameNew)
	}

	logging.Plog(fmt.Sprintf("__DONE__resourceCatalogUpdate__ "))
	return nil
}

func resourceCatalogDelete(d *schema.ResourceData, m interface{}) error {

	logging.Plog(fmt.Sprintf("__INIT__resourceCatalogDelete_"))
	cname := d.Get("name").(string)

	vcdClient := m.(*VCDClient)

	provider := vcdClient.getProvider()

	_, err := provider.DeleteCatalog(cname)

	if err != nil {
		return fmt.Errorf("Error Creating catalog :%v %#v", cname, err)
	}
	logging.Plog(fmt.Sprintf("__DONE__resourceCatalogDelete_"))
	return nil
}
