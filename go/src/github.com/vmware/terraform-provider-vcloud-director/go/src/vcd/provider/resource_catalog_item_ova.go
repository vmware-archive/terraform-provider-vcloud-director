/*****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
******************************************************************/

package provider

import (
	"fmt"

	//"time"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/util/logging"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
)

func resourceCatalogItemOva() *schema.Resource {
	return &schema.Resource{
		Create: resourceCatalogItemOvaCreate,
		Read:   resourceCatalogItemOvaRead,
		Update: resourceCatalogItemOvaUpdate,
		Delete: resourceCatalogItemOvaDelete,

		Schema: map[string]*schema.Schema{
			"catalog_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"item_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"source_file_path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"source_vapp_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_vdc_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"customize_on_instantiate": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

//def capture_vapp(client,vapp_name,vdc_name,catalog_name, item_name,desc,customize_on_instantiate=False):
func getOvaId(d *schema.ResourceData) string {

	itemName := d.Get("item_name").(string)
	catalogName := d.Get("catalog_name").(string)
	return catalogName + "_" + itemName
}

func getCatalogUploadOvaInfo(d *schema.ResourceData) proto.CatalogUploadOvaInfo {
	itemName := d.Get("item_name").(string)
	filePath := d.Get("source_file_path").(string)
	catalogName := d.Get("catalog_name").(string)

	p := proto.CatalogUploadOvaInfo{CatalogName: catalogName,
		FilePath: filePath,
		ItemName: itemName}

	return p
}

func resourceCatalogItemOvaCreate(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT__resourceCatalogItemOvaCreate_ ")
	catalogItem := buildCatalogItem(d, m)

	return catalogItem.create()
}

func resourceCatalogItemOvaRead(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT__resourceCatalogItemOvaRead_  ")

	catalogName, itemName, _ := getCatalogItemInfo(d)
	provider := getProvider(m)

	isPrInfo := proto.IsPresentCatalogItemInfo{CatalogName: catalogName,
		ItemName: itemName}

	res, err := provider.IsPresentCatalogItem(isPrInfo)

	if err != nil {
		return fmt.Errorf("__ERROR__Error checking Catalog Item  [%#v]", err)
	}
	if res.Present {
		logging.Plog(fmt.Sprintf("__LOG__catalog item [%v] is present / setting state ", itemName))
		d.SetId(getOvaId(d))
	} else {
		// resource catalog item not present / clear id for recreate
		d.SetId("")
	}
	logging.Plog("__DONE__resourceCatalogItemOvaRead")
	return nil

}

func resourceCatalogItemOvaUpdate(d *schema.ResourceData, m interface{}) error {

	oraw, nraw := d.GetChange("item_name")
	o := oraw.(string)
	n := nraw.(string)

	msg := o + " , " + n
	logging.Plog(msg)
	return nil
}

func resourceCatalogItemOvaDelete(d *schema.ResourceData, m interface{}) error {

	catalogName, itemName, _ := getCatalogItemInfo(d)
	provider := getProvider(m)

	cInfo := proto.DeleteCatalogItemInfo{CatalogName: catalogName,
		ItemName: itemName}

	_, err := provider.DeleteCatalogItem(cInfo)
	return err
}
