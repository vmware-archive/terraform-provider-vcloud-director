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

func resourceVApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceVAppCreate,
		Read:   resourceVAppRead,
		Update: resourceVAppUpdate,
		Delete: resourceVAppDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"template_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"catalog_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vdc": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			//optional parameters
			"network": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"ip_allocation_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"memory": &schema.Schema{
				Type:     schema.TypeString, //Keeping it String to accomodate the large values and keep the logic simple
				Optional: true,
				ForceNew: false,
			},
			"cpu": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},

			//optional parameters
			"storage_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"power_on": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func getVAppId(d *schema.ResourceData) string {

	vappName := d.Get("name").(string)

	return vappName
}

func getVAppInfo(d *schema.ResourceData) (string, string, string, string, string, string, string, int32, string) {
	vAppName := d.Get("name").(string)
	templateName := d.Get("template_name").(string)
	catalogName := d.Get("catalog_name").(string)
	vdc := d.Get("vdc").(string)
	network := d.Get("network").(string)
	ipAllocationMode := d.Get("ip_allocation_mode").(string)
	memory := d.Get("memory").(string)
	cpu := int32(d.Get("cpu").(int))
	storageProfile := d.Get("storage_profile").(string)
	return vAppName, templateName, catalogName, vdc, network, ipAllocationMode, memory, cpu, storageProfile
}

func resourceVAppCreate(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT__resourceVAppCreate_")

	provider := getProvider(m)

	vAppName, templateName, catalogName, vdc, network, ipAllocationMode, memory, cpu, storageProfile := getVAppInfo(d)

	readvAppInfo := proto.ReadVAppInfo{

		Name: vAppName,
		Vdc:  vdc,
	}

	readresp, readerr := provider.ReadVApp(readvAppInfo)
	if readerr != nil {
		return fmt.Errorf("__ERROR__.... in reading vapp  creation", readerr)
	}

	if readresp.Present {

		d.SetId(getVAppId(d))
		logging.Plog("__DONE__ resourceVAppRead_EARLY EXIT")
		return nil
	}

	vAppInfo := proto.CreateVAppInfo{

		Name:             vAppName,
		TemplateName:     templateName,
		CatalogName:      catalogName,
		Vdc:              vdc,
		Network:          network,
		IpAllocationMode: ipAllocationMode,
		Memory:           memory,
		Cpu:              cpu,
		StorageProfile:   storageProfile,
	}

	logging.Plog(fmt.Sprintf("__LOG__vAppInfo %#v", vAppInfo))
	resp, errp := provider.CreateVApp(vAppInfo)

	if errp != nil {

		return fmt.Errorf("__ERROR__ Creating VApp Failed: %v", errp)
	}

	logging.Plog(fmt.Sprintf("__LOG__resp %v %#v", resp.Created, *resp.InVappInfo))

	if resp.Created {

		d.SetId(getVAppId(d))
		return nil
	}

	logging.Plog(" __DONE__resourceVAppCreate_")
	return nil
}

func resourceVAppRead(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT__resourceVAppRead_ ")
	vAppName, _, _, vdc, _, _, _, _, _ := getVAppInfo(d)

	vAppInfo := proto.ReadVAppInfo{

		Name: vAppName,
		Vdc:  vdc,
	}

	provider := getProvider(m)
	resp, errp := provider.ReadVApp(vAppInfo)
	if errp != nil {
		return fmt.Errorf("__ERROR__.... in reading vapp  creation", errp)
	}

	if resp.Present {
		logging.Plogf("__LOG__ setting id %v", getVAppId(d))
		d.SetId(getVAppId(d))
	} else {
		d.SetId("")
	}
	logging.Plog("__DONE__ resourceVAppRead_")
	return nil

}

func resourceVAppUpdate(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceVAppUpdate_")

	vAppName := d.Get("name").(string)
	vdc := d.Get("vdc").(string)

	powerOnOldRaw, powerOnNewRaw := d.GetChange("power_on")
	powerOnOld := powerOnOldRaw.(bool)
	powerOnNew := powerOnNewRaw.(bool)

	logging.Plog(fmt.Sprintf("[powerOnOld %v ]  [powerOnNew %v ]", powerOnOld, powerOnNew))

	if !(powerOnOld == powerOnNew) {

		updateVAppInfo := proto.UpdateVAppInfo{
			Name:    vAppName,
			Vdc:     vdc,
			PowerOn: powerOnNew,
		}

		provider := getProvider(m)
		resp, errp := provider.UpdateVApp(updateVAppInfo)
		if errp != nil {
			return fmt.Errorf("Error Updating VApp :[%+v] %#v", updateVAppInfo, errp)
		}
		if resp.Updated {
			logging.Plog(fmt.Sprintf("VApp [%+v]  updated  ", resp))
		}
	} else {
		return fmt.Errorf("__ERROR__ Can not update field other then power_on")
	}

	return nil
}

func resourceVAppDelete(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceVAppDelete_")
	vAppName, _, _, vdc, _, _, _, _, _ := getVAppInfo(d)

	provider := getProvider(m)

	vappInfo := proto.DeleteVAppInfo{Name: vAppName, Vdc: vdc}

	_, err := provider.DeleteVApp(vappInfo)

	logging.Plog("__DONE__resourceVAppDelete_")
	return err

}
