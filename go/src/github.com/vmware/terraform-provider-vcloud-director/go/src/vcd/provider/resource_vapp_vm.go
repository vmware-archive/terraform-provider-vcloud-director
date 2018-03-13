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
			"virtual_cpus": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"cores_per_socket": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"memory": &schema.Schema{
				Type:     schema.TypeInt,
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

	if len(sourceCatalogName) > 0 {
		createVappVmInfo := proto.CreateVappVmInfo{
			TargetVmName:       targetVmName,
			TargetVapp:         targetVapp,
			TargetVdc:          targetVdc,
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

		res, err := provider.CreateFromCatalog(createVappVmInfo)

		if err != nil {
			return fmt.Errorf("Error Creating VappVm :[%+v] %#v", createVappVmInfo, err)
		}

		if res.Created {
			logging.Plog(fmt.Sprintf("VappVm [%+v]  created  ", targetVmName))
			d.SetId(targetVmName)
		}

	} else {
		createVappVmInfo := proto.CreateVappVmInfo{
			TargetVmName:     targetVmName,
			TargetVapp:       targetVapp,
			TargetVdc:        targetVdc,
			SourceVapp:       sourceVapp,
			SourceVmName:     sourceVmName,
			Hostname:         hostname,
			Password:         password,
			PasswordAuto:     passwordAuto,
			PasswordReset:    passwordReset,
			CustScript:       custScript,
			Network:          network,
			StorageProfile:   storageProfile,
			PowerOn:          powerOn,
			AllEulasAccepted: allEulasAccepted,
			IpAllocationMode: ipAllocationMode,
		}

		res, err := provider.CreateFromVapp(createVappVmInfo)

		if err != nil {
			return fmt.Errorf("Error Creating VappVm :[%+v] %#v", createVappVmInfo, err)
		}

		if res.Created {
			logging.Plog(fmt.Sprintf("VappVm [%+v]  created  ", targetVmName))
			d.SetId(targetVmName)
		}
	}

	memCPUUpdateErr := updateMemCPUAfterCreate(d)
	if memCPUUpdateErr != nil {
		return fmt.Errorf("Error Updating VappVm memory and cpu after create :[%+v] %#v", memCPUUpdateErr)
	}

	logging.Plog("__DONE__resourceVappVmCreate_")
	return nil
}

/**
+ ** update Memory and CPU if the variables are present.
+ ** As per current PyVcloud API, it do not support cpu and memory info
+ ** during vapp_vm creation.
+**/
func updateMemCPUAfterCreate(d *schema.ResourceData) error {
	logging.Plog("__INIT__updateMemCPUAfterCreate_")

	targetVmName := d.Get("target_vm_name").(string)
	targetVapp := d.Get("target_vapp").(string)
	targetVdc := d.Get("target_vdc").(string)

	memory := int32(d.Get("memory").(int))
	coresPerSocket := int32(d.Get("cores_per_socket").(int))
	virtualCPUs := int32(d.Get("virtual_cpus").(int))

	logging.Plog(fmt.Sprintf("Memory, cpu values received [memory %d ], [coresPerSocket %d ], [virtualCPUs, %d ]", memory, coresPerSocket, virtualCPUs))

	if memory > 0 {
		modifyMemoryRes, modifyMemoryErr := modifyMemory(targetVmName, targetVapp, targetVdc, memory)
		if modifyMemoryErr != nil {
			return fmt.Errorf("Error updating memory VappVm :[%+v] %#v", targetVmName, modifyMemoryErr)
		} else {
			if modifyMemoryRes.Modified {
				logging.Plog(fmt.Sprintf("VappVm [%+v]  memory modified ", modifyMemoryRes))
			}
		}
	}

	if coresPerSocket > 0 && virtualCPUs > 0 {
		modifyCPURes, modifyCPUErr := modifyCPU(targetVmName, targetVapp, targetVdc, coresPerSocket, virtualCPUs)

		if modifyCPUErr != nil {
			return fmt.Errorf("Error updating cpus VappVm :[%+v] %#v", targetVmName, modifyCPUErr)
		} else {
			if modifyCPURes.Modified {
				logging.Plog(fmt.Sprintf("VappVm [%+v]  cpus modified ", modifyCPURes))
			}
		}
	}

	logging.Plog("__DONE__updateMemCPUAfterCreate_")

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

	targetVmName := d.Get("target_vm_name").(string)
	targetVapp := d.Get("target_vapp").(string)
	targetVdc := d.Get("target_vdc").(string)

	oldMemoryRaw, newMemoryRaw := d.GetChange("memory")
	oldMemory := int32(oldMemoryRaw.(int))
	newMemory := int32(newMemoryRaw.(int))

	oldCoresPerSocketRaw, newCoresPerSocketRaw := d.GetChange("cores_per_socket")
	oldCoresPerSocket := int32(oldCoresPerSocketRaw.(int))
	newCoresPerSocket := int32(newCoresPerSocketRaw.(int))

	oldVirtualCPUsRaw, newVirtualCPUsRaw := d.GetChange("virtual_cpus")
	oldVirtualCPUs := int32(oldVirtualCPUsRaw.(int))
	newVirtualCPUs := int32(newVirtualCPUsRaw.(int))

	oldPoweredOnRaw, newPoweredOnRaw := d.GetChange("power_on")
	oldPoweredOn := oldPoweredOnRaw.(bool)
	newPoweredOn := newPoweredOnRaw.(bool)

	if !(oldPoweredOn == newPoweredOn) || !(oldMemory == newMemory) || !(newCoresPerSocket == oldCoresPerSocket) || !(newVirtualCPUs == oldVirtualCPUs) {
		if !(oldPoweredOn == newPoweredOn) {
			if newPoweredOn {
				//call power_on
				res, err := powerOnVM(targetVmName, targetVapp, targetVdc)
				if err != nil {
					return fmt.Errorf("Error updating PowerOn VappVm :[%+v] %#v", targetVmName, err)
				}
				if res.PoweredOn {
					logging.Plog(fmt.Sprintf("VappVm [%+v]  powered On  ", res))
				}

			} else {
				//call power_off
				res, err := powerOffVM(targetVmName, targetVapp, targetVdc)
				if err != nil {
					return fmt.Errorf("Error updating PowerOff VappVm :[%+v] %#v", targetVmName, err)
				}
				if res.PoweredOff {
					logging.Plog(fmt.Sprintf("VappVm [%+v]  powered Off  ", res))
				}
			}
		}

		if !(oldMemory == newMemory) && newMemory > 0 {
			res, err := modifyMemory(targetVmName, targetVapp, targetVdc, newMemory)
			if err != nil {
				return fmt.Errorf("Error updating memory VappVm :[%+v] %#v", targetVmName, err)
			}
			if res.Modified {
				logging.Plog(fmt.Sprintf("VappVm [%+v]  memory modified ", res))
			}
		}

		if !(newCoresPerSocket == oldCoresPerSocket) || !(newVirtualCPUs == oldVirtualCPUs) {

			if newCoresPerSocket > 0 && newVirtualCPUs > 0 {
				res, err := modifyCPU(targetVmName, targetVapp, targetVdc, newCoresPerSocket, newVirtualCPUs)
				if err != nil {
					return fmt.Errorf("Error updating cpus VappVm :[%+v] %#v", targetVmName, err)
				}
				if res.Modified {
					logging.Plog(fmt.Sprintf("VappVm [%+v]  cpus modified ", res))
				}
			}
		}
	} else {
		return fmt.Errorf("Error modifying VappVm :[%+v]. "+
			"Can not update the given fields ", targetVmName)
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

func powerOffVM(targetVmName, targetVapp, targetVdc string) (*proto.PowerOffVappVmResult, error) {
	logging.Plog("__INIT__powerOffVM_")

	provider := providerGlobalRefPointer.vappVmProvider

	powerOffVappVmInfo := proto.PowerOffVappVmInfo{
		TargetVmName: targetVmName,
		TargetVapp:   targetVapp,
		TargetVdc:    targetVdc,
	}
	res, err := provider.PowerOff(powerOffVappVmInfo)
	logging.Plog("__DONE__powerOffVM_")
	return res, err
}

func powerOnVM(targetVmName, targetVapp, targetVdc string) (*proto.PowerOnVappVmResult, error) {
	logging.Plog("__INIT__powerOnVM_")

	provider := providerGlobalRefPointer.vappVmProvider

	powerOnVappVmInfo := proto.PowerOnVappVmInfo{
		TargetVmName: targetVmName,
		TargetVapp:   targetVapp,
		TargetVdc:    targetVdc,
	}
	res, err := provider.PowerOn(powerOnVappVmInfo)
	logging.Plog("__DONE__powerOnVM_")
	return res, err

}

func modifyMemory(targetVmName string, targetVapp string, targetVdc string, memory int32) (*proto.ModifyVappVmMemoryResult, error) {

	logging.Plog("__INIT__modifyMemory_")

	provider := providerGlobalRefPointer.vappVmProvider

	modifyVappVmMemoryInfo := proto.ModifyVappVmMemoryInfo{
		TargetVmName: targetVmName,
		TargetVapp:   targetVapp,
		TargetVdc:    targetVdc,
		Memory:       memory,
	}
	res, err := provider.ModifyMemory(modifyVappVmMemoryInfo)
	logging.Plog("__DONE__modifyMemory_")
	return res, err

}

func modifyCPU(targetVmName string, targetVapp string, targetVdc string, coresPerSocket int32, virtualCpus int32) (*proto.ModifyVappVmCPUResult, error) {

	logging.Plog("__INIT__modifyCPU_")

	provider := providerGlobalRefPointer.vappVmProvider

	modifyVappVmCPUInfo := proto.ModifyVappVmCPUInfo{
		TargetVmName:   targetVmName,
		TargetVapp:     targetVapp,
		TargetVdc:      targetVdc,
		VirtualCpus:    virtualCpus,
		CoresPerSocket: coresPerSocket,
	}
	res, err := provider.ModifyCPU(modifyVappVmCPUInfo)
	logging.Plog("__DONE__modifyCPU_")
	return res, err
}
