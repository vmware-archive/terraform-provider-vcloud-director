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

func resourceIndependentDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceIndependentDiskCreate,
		Read:   resourceIndependentDiskRead,
		Update: resourceIndependentDiskUpdate,
		Delete: resourceIndependentDiskDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"vdc": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"storage_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"disk_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceIndependentDiskCreate(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT__resourceIndependentDiskCreate_")
	cname := d.Get("name").(string)
	description := d.Get("description").(string)
	size := d.Get("size").(string)
	vdc := d.Get("vdc").(string)
	storage_profile := d.Get("storage_profile").(string)
	disk_id := d.Get("disk_id").(string)

	logging.Plog(fmt.Sprintf("%v %v %v %v %v ", cname, size, vdc, storage_profile, description))

	provider := providerGlobalRefPointer.independentDiskProvider

	//TODO THIS SECTIONS
	readInfo := proto.ReadDiskInfo{Name: cname, Vdc: vdc, DiskId: disk_id}
	resp, errp := provider.Read(readInfo)

	if errp != nil {
		return fmt.Errorf("__ERROR__ Creating IndependentDisk :[%v] %#v", cname, errp)
	}

	if resp.Present {
		logging.Plog(fmt.Sprintf("IndependentDisk %v is present/ setting internal state / this is to allow adding to IndependentDisk already present ", cname))
		d.SetId(resp.DiskId)
		return nil
	}

	creatDiskInfo := proto.CreateDiskInfo{Name: cname,
		Size:           size,
		StorageProfile: storage_profile,
		Description:    description,
		Vdc:            vdc}

	logging.Plog(fmt.Sprintf("__LOG__calling create IndependentDisk  %+v ", creatDiskInfo))
	res, err := provider.Create(creatDiskInfo)

	if err != nil {
		return fmt.Errorf("Error Creating IndependentDisk :[%+v] %#v", creatDiskInfo, err)
	}
	if res.Created {
		logging.Plog(fmt.Sprintf("IndependentDisk [%+v]  created  ", res))
		d.SetId(res.DiskId)
	}
	logging.Plog("__DONE__resourceIndependentDiskCreate ")
	return nil
}

func resourceIndependentDiskRead(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT__resourceIndependentDiskRead_")

	disk_name := d.Get("name").(string)
	disk_id := d.Get("disk_id").(string)
	vdc := d.Get("vdc").(string)
	logging.Plog(fmt.Sprintf("read validate =%v %v %v", disk_name, disk_id, vdc))
	provider := providerGlobalRefPointer.independentDiskProvider
	//TODO THIS SECTIONS
	readInfo := proto.ReadDiskInfo{Name: disk_name, Vdc: vdc, DiskId: disk_id}
	resp, errp := provider.Read(readInfo)

	if errp != nil {
		return fmt.Errorf("__ERROR__ Creating IndependentDisk :[%v] %#v", disk_name, errp)
	}

	if resp.Present {

		d.SetId(resp.DiskId)
		logging.Plog(fmt.Sprintf("__DONE__resourceIndependentDiskRead_ +setting id %v", resp.DiskId))
		return nil
	} else {
		d.SetId("")
		logging.Plog(fmt.Sprintf("__DONE__resourceIndependentDiskRead_ +unsetting id,resource got deleted %v", disk_name))
	}

	logging.Plog("__DONE__resourceIndependentDiskRead_ NOT SET ID ")
	return nil

}

func resourceIndependentDiskUpdate(d *schema.ResourceData, m interface{}) error {

	logging.Plog(fmt.Sprintf("__INIT__DONE__ NO IMPl !!! resourceIndependentDiskUpdate "))
	return nil
}

func resourceIndependentDiskDelete(d *schema.ResourceData, m interface{}) error {

	logging.Plog(fmt.Sprintf("__INIT__resourceIndependentDiskDelete_"))
	cname := d.Get("name").(string)
	vdc := d.Get("vdc").(string)
	disk_id := d.Get("disk_id").(string)

	diskProvider := providerGlobalRefPointer.independentDiskProvider
	diskInfo := proto.DeleteDiskInfo{Name: cname, Vdc: vdc, DiskId: disk_id}

	_, err := diskProvider.Delete(diskInfo)

	if err != nil {
		return fmt.Errorf("Error Deleting IndependentDisk :%v %#v", cname, err)
	}
	logging.Plog(fmt.Sprintf("__DONE__resourceIndependentDiskDelete_"))
	return nil
}
