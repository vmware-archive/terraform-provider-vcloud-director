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

func resourceCatalogItemMedia() *schema.Resource {
	return &schema.Resource{
		Create: resourceCatalogItemMediaCreate,
		Read:   resourceCatalogItemMediaRead,
		Update: resourceCatalogItemMediaUpdate,
		Delete: resourceCatalogItemMediaDelete,

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
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func getMediaId(d *schema.ResourceData) string {

	itemName := d.Get("item_name").(string)
	filePath := d.Get("source_file_path").(string)
	catalogName := d.Get("catalog_name").(string)
	return catalogName + "_" + itemName + "_" + filePath
}

func getCatalogUploadMediaInfo(d *schema.ResourceData) proto.CatalogUploadMediaInfo {
	itemName := d.Get("item_name").(string)
	filePath := d.Get("source_file_path").(string)
	catalogName := d.Get("catalog_name").(string)

	p := proto.CatalogUploadMediaInfo{CatalogName: catalogName,
		FilePath: filePath,
		ItemName: itemName}

	return p
}

func getCatalogItemInfo(d *schema.ResourceData) (string, string, string) {
	itemName := d.Get("item_name").(string)
	filePath := d.Get("source_file_path").(string)
	catalogName := d.Get("catalog_name").(string)

	return catalogName, itemName, filePath
}

func resourceCatalogItemMediaCreate(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT__resourceCatalogItemMediaCreate_")

	provider := getProvider(m)

	catalogUploadInfo := getCatalogUploadMediaInfo(d)

	catalogName, itemName, _ := getCatalogItemInfo(d)

	catalogIsPresInfo := proto.IsPresentCatalogItemInfo{CatalogName: catalogName, ItemName: itemName}

	isPresResp, isPreErr := provider.IsPresentCatalogItem(catalogIsPresInfo)

	if isPreErr != nil {
		return fmt.Errorf("Error Creating Item :[%v] %#v", catalogName, isPreErr)
	}
	if isPresResp.Present {
		logging.PlogWarn(fmt.Sprintf("__LOG__ catalog item [%v] is already present /setting state as created  ", catalogName))
		d.SetId(getMediaId(d))
		return nil
	}

	resp, errp := provider.CatalogUploadMedia(catalogUploadInfo)

	if errp != nil {
		logging.Plog("__ERROR__ creating catalog ITEM ")
		return fmt.Errorf("Error Creating catalog Item: %#v", errp)
	}

	logging.Plog(fmt.Sprintf("__LOG__resp.Created [%v]", resp.Created))
	if resp.Created {

		d.SetId(getMediaId(d))
		return nil
	}

	logging.Plog("__DONE__resourceCatalogItemMediaCreate_")
	return nil
}

func resourceCatalogItemMediaRead(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT_resourceCatalogItemMediaRead_")

	catalogName, itemName, _ := getCatalogItemInfo(d)
	provider := getProvider(m)

	isPrInfo := proto.IsPresentCatalogItemInfo{CatalogName: catalogName,
		ItemName: itemName}

	res, err := provider.IsPresentCatalogItem(isPrInfo)

	if err != nil {
		return fmt.Errorf("Error checking Catalog Item  %#v", err)
	}
	if res.Present {

		d.SetId(getMediaId(d))
	} else {
		// resource catalog item not present / clear id for recreate
		d.SetId("")
	}
	logging.Plog("__DONE__resourceCatalogItemMediaRead_")
	return nil

}

func resourceCatalogItemMediaUpdate(d *schema.ResourceData, m interface{}) error {

	oraw, nraw := d.GetChange("item_name")
	o := oraw.(string)
	n := nraw.(string)

	msg := o + " , " + n
	logging.Plog(msg)
	return nil
}

func resourceCatalogItemMediaDelete(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT__resourceCatalogItemMediaDelete_")
	catalogName, itemName, _ := getCatalogItemInfo(d)
	provider := getProvider(m)

	cInfo := proto.DeleteCatalogItemInfo{CatalogName: catalogName,
		ItemName: itemName}

	_, err := provider.DeleteCatalogItem(cInfo)
	logging.Plog("__DONE__resourceCatalogItemMediaDelete_")
	return err
}
