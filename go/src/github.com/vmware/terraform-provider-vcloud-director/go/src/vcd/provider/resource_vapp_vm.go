/* ****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
***************************************************************** */
package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/util/logging"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
)

func resourceVappVm() *schema.Resource {
	return &schema.Resource{
		Create: resourceVappVmCreate,
		Read:   resourceVappVmRead,
		Update: resourceVappVmUpdate,
		Delete: resourceVappVmDelete,

		Schema: map[string]*schema.Schema{
			"target_vm_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_vapp": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"target_vdc": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"source_vapp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"source_vm_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"source_catalog_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"source_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"password_auto": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"password_reset": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"cust_script": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"network": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"storage_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"power_on": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"all_eulas_accepted": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"ip_allocation_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceVappVmCreate(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceVappVmCreate_")

	targetVmName := d.Get("target_vm_name").(string)
	targetVapp := d.Get("target_vapp").(string)
	targetVdc := d.Get("target_vdc").(string)
	sourceVapp := d.Get("source_vapp").(string)
	sourceVmName := d.Get("source_vm_name").(string)
	hostname := d.Get("hostname").(string)
	password := d.Get("password").(string)

	passwordAuto := d.Get("password_auto").(bool)
	passwordReset := d.Get("password_reset").(bool)

	custScript := d.Get("cust_script").(string)
	network := d.Get("network").(string)
	storageProfile := d.Get("storage_profile").(string)

	powerOn := d.Get("power_on").(bool)
	allEulasAccepted := d.Get("all_eulas_accepted").(bool)
	sourceCatalogName := d.Get("source_catalog_name").(string)
	sourceTemplateName := d.Get("source_template_name").(string)

	ipAllocationMode := d.Get("ip_allocation_mode").(string)

	provider := providerGlobalRefPointer.vappVmProvider

	createVappVmInfo := proto.CreateVappVmInfo{
		TargetVmName:       targetVmName,
		TargetVapp:         targetVapp,
		TargetVdc:          targetVdc,
		SourceVapp:         sourceVapp,
		SourceVmName:       sourceVmName,
		SourceCatalogName:  sourceCatalogName,
		SourceTemplateName: sourceTemplateName,
		Hostname:           hostname,
		Password:           password,
		PasswordAuto:       passwordAuto,
		PasswordReset:      passwordReset,
		CustScript:         custScript,
		Network:            network,
		StorageProfile:     storageProfile,
		PowerOn:            powerOn,
		AllEulasAccepted:   allEulasAccepted,
		IpAllocationMode:   ipAllocationMode,
	}

	res, err := provider.Create(createVappVmInfo)

	if err != nil {
		return fmt.Errorf("Error Creating VappVm :[%+v] %#v", createVappVmInfo, err)
	}

	if res.Created {
		logging.Plog(fmt.Sprintf("VappVm [%+v]  created  ", targetVmName))
		d.SetId(targetVmName)
	}

	logging.Plog("__DONE__resourceVappVmCreate_")
	return nil
}

func resourceVappVmDelete(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceVappVmDelete_")

	targetVmName := d.Get("target_vm_name").(string)
	targetVapp := d.Get("target_vapp").(string)
	targetVdc := d.Get("target_vdc").(string)

	provider := providerGlobalRefPointer.vappVmProvider

	deleteVappVmInfo := proto.DeleteVappVmInfo{
		TargetVmName: targetVmName,
		TargetVapp:   targetVapp,
		TargetVdc:    targetVdc,
	}
	res, err := provider.Delete(deleteVappVmInfo)

	if err != nil {
		return fmt.Errorf("Error Deleting VappVm :[%+v] %#v", deleteVappVmInfo, err)
	}

	if res.Deleted {
		logging.Plog(fmt.Sprintf("VappVm [%+v]  deleted  ", res))
	}

	logging.Plog("__DONE__resourceVappVmDelete_")
	return nil
}

func resourceVappVmUpdate(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceVappVmUpdate_")

	name := d.Get("target_vm_name").(string)

	oldIsEnabledRaw, newIsEnabledRaw := d.GetChange("is_enabled")
	oldIsEnabled := oldIsEnabledRaw.(bool)
	newIsEnabled := newIsEnabledRaw.(bool)

	provider := providerGlobalRefPointer.vappVmProvider

	if !(oldIsEnabled == newIsEnabled) {
		updateVappVmInfo := proto.UpdateVappVmInfo{Name: name, IsEnabled: newIsEnabled}
		res, err := provider.Update(updateVappVmInfo)
		if err != nil {
			return fmt.Errorf("Error updating VappVm :[%+v] %#v", updateVappVmInfo, err)
		}

		if res.Updated {
			logging.Plog(fmt.Sprintf("VappVm [%+v]  updated  ", res))
			d.SetId(name)
		}
	} else {
		return fmt.Errorf("Error updating VappVm :[%+v]. "+
			"Can not update the given fields ", name)
	}

	logging.Plog("__DONE__resourceVappVmUpdate_")
	return nil
}

func resourceVappVmRead(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceVappVmRead_")

	targetVmName := d.Get("target_vm_name").(string)
	targetVapp := d.Get("target_vapp").(string)
	targetVdc := d.Get("target_vdc").(string)

	provider := providerGlobalRefPointer.vappVmProvider

	readVappVmInfo := proto.ReadVappVmInfo{
		TargetVmName: targetVmName,
		TargetVapp:   targetVapp,
		TargetVdc:    targetVdc,
	}
	res, err := provider.Read(readVappVmInfo)

	if err != nil {
		logging.Plog(fmt.Sprintf("Error Reading VappVm :[%+v] %#v", readVappVmInfo, err))
		//return fmt.Errorf("Error Reading VappVm :[%+v] %#v", readVappVmInfo, err)
	}

	if res.Present {
		d.SetId(targetVmName)
		logging.Plog(fmt.Sprintf("__DONE__resourceVappVmRead_ +setting id %v", targetVmName))
		return nil
	} else {
		d.SetId("")
		logging.Plog(fmt.Sprintf("__DONE__resourceVappVmRead_ +unsetting id,resource got deleted %v", targetVmName))
	}

	logging.Plog("__DONE__resourceVappVmRead_")
	return nil
}
