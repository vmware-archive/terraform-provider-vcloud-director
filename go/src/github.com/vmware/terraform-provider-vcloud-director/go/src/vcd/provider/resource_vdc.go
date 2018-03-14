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

func resourceVdc() *schema.Resource {
	return &schema.Resource{
		Create: resourceVdcCreate,
		Read:   resourceVdcRead,
		Update: resourceVdcUpdate,
		Delete: resourceVdcDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"provider_vdc": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"allocation_model": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"cpu_units": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"cpu_allocated": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"cpu_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"mem_units": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"mem_allocated": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"mem_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"nic_quota": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"network_quota": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"vm_quota": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"storage_profiles": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"resource_guaranteed_memory": &schema.Schema{
				Type:     schema.TypeFloat,
				Required: true,
				ForceNew: false,
			},
			"resource_guaranteed_cpu": &schema.Schema{
				Type:     schema.TypeFloat,
				Required: true,
				ForceNew: false,
			},
			"vcpu_in_mhz": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"is_thin_provision": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
			"network_pool_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"uses_fast_provisioning": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
			"over_commit_allowed": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
			"vm_discovery_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
			"is_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceVdcCreate(d *schema.ResourceData, m interface{}) error {

	logging.Plog("__INIT__resourceVdcCreate_")

	name := d.Get("name").(string)
	providerVdc := d.Get("provider_vdc").(string)
	description := d.Get("description").(string)
	allocationModel := d.Get("allocation_model").(string)

	cpuUnits := d.Get("cpu_units").(string)
	cpuAllocated := int32(d.Get("cpu_allocated").(int))
	cpuLimit := int32(d.Get("cpu_limit").(int))

	memUnits := d.Get("mem_units").(string)
	memAllocated := int32(d.Get("mem_allocated").(int))
	memLimit := int32(d.Get("mem_limit").(int))

	nicQuota := int32(d.Get("nic_quota").(int))
	networkQuota := int32(d.Get("network_quota").(int))
	vmQuota := int32(d.Get("vm_quota").(int))

	storageProfiles := d.Get("storage_profiles").(string)
	logging.Plog(fmt.Sprintf("StorageProfiles: [%v]", storageProfiles))

	resourceGuaranteedMemory := float32(d.Get("resource_guaranteed_memory").(float64))
	resourceGuaranteedCpu := float32(d.Get("resource_guaranteed_cpu").(float64))

	vcpuInMhz := int32(d.Get("vcpu_in_mhz").(int))

	isThinProvision := d.Get("is_thin_provision").(bool)

	networkPoolName := d.Get("network_pool_name").(string)

	usesFastProvisioning := d.Get("uses_fast_provisioning").(bool)
	overCommitAllowed := d.Get("over_commit_allowed").(bool)
	vmDiscoveryEnabled := d.Get("vm_discovery_enabled").(bool)
	isEnabled := d.Get("is_enabled").(bool)

	provider := providerGlobalRefPointer.vdcProvider

	createVdcInfo := proto.CreateVdcInfo{
		Name:                     name,
		ProviderVdc:              providerVdc,
		Description:              description,
		AllocationModel:          allocationModel,
		CpuUnits:                 cpuUnits,
		CpuAllocated:             cpuAllocated,
		CpuLimit:                 cpuLimit,
		MemUnits:                 memUnits,
		MemAllocated:             memAllocated,
		MemLimit:                 memLimit,
		NicQuota:                 nicQuota,
		NetworkQuota:             networkQuota,
		VmQuota:                  vmQuota,
		StorageProfiles:          storageProfiles,
		ResourceGuaranteedMemory: resourceGuaranteedMemory,
		ResourceGuaranteedCpu:    resourceGuaranteedCpu,
		VcpuInMhz:                vcpuInMhz,
		IsThinProvision:          isThinProvision,
		NetworkPoolName:          networkPoolName,
		UsesFastProvisioning:     usesFastProvisioning,
		OverCommitAllowed:        overCommitAllowed,
		VmDiscoveryEnabled:       vmDiscoveryEnabled,
		IsEnabled:                isEnabled,
	}

	res, err := provider.Create(createVdcInfo)

	if err != nil {
		return fmt.Errorf("Error Creating vdc :[%+v] %#v", createVdcInfo, err)
	}

	if res.Created {
		logging.Plog(fmt.Sprintf("Vdc [%+v]  created  ", res))
		d.SetId(name)
	}

	logging.Plog("__DONE__resourceVdcCreate_")
	return nil
}

//Delete vdc
func resourceVdcDelete(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceVdcDelete_")
	vdcName := d.Get("name").(string)
	provider := providerGlobalRefPointer.vdcProvider

	deleteVdcInfo := proto.DeleteVdcInfo{Name: vdcName}
	res, err := provider.Delete(deleteVdcInfo)

	if err != nil {
		return fmt.Errorf("Error Deleting Vdc :[%+v] %#v", deleteVdcInfo, err)
	}

	if res.Deleted {
		logging.Plog(fmt.Sprintf("Vdc [%+v]  deleted  ", res))
	}

	logging.Plog("__DONE__resourceVdcDelete_")
	return nil
}

func resourceVdcUpdate(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceVdcUpdate_")

	name := d.Get("name").(string)

	oldIsEnabledRaw, newIsEnabledRaw := d.GetChange("is_enabled")
	oldIsEnabled := oldIsEnabledRaw.(bool)
	newIsEnabled := newIsEnabledRaw.(bool)

	provider := providerGlobalRefPointer.vdcProvider

	if !(oldIsEnabled == newIsEnabled) {
		updateVdcInfo := proto.UpdateVdcInfo{Name: name, IsEnabled: newIsEnabled}
		res, err := provider.Update(updateVdcInfo)
		if err != nil {
			return fmt.Errorf("Error updating Vdc :[%+v] %#v", updateVdcInfo, err)
		}

		if res.Updated {
			logging.Plog(fmt.Sprintf("Vdc [%+v]  updated  ", res))
			d.SetId(name)
		}
	} else {
		return fmt.Errorf("Error updating Vdc :[%+v]. "+
			"Can not update the given fields ", name)
	}

	logging.Plog("__DONE__resourceVdcUpdate_")
	return nil
}

func resourceVdcRead(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceVdcRead_")

	name := d.Get("name").(string)

	provider := providerGlobalRefPointer.vdcProvider

	readVdcInfo := proto.ReadVdcInfo{Name: name}
	res, err := provider.Read(readVdcInfo)

	if err != nil {
		logging.Plog(fmt.Sprintf("Error Reading Vdc :[%+v] %#v", readVdcInfo, err))
	}

	if res.Present {
		d.SetId(name)
		d.Set("is_enabled", res.IsEnabled)
		logging.Plog(fmt.Sprintf("setting id %v", name))
		return nil
	} else {
		d.SetId("")
		logging.Plog(fmt.Sprintf("unsetting id,resource got deleted %v", name))
	}

	logging.Plog("__DONE__resourceVdcRead_")
	return nil
}
