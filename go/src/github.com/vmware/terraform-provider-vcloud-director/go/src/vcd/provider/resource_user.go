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

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"role_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"full_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"telephone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"im": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"alert_email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"alert_email_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"stored_vm_quota": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"deployed_vm_quota": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"is_group_role": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"is_default_cached": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"is_external": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"is_alert_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
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

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceUserCreate_")

	name := d.Get("name").(string)
	password := d.Get("password").(string)
	roleName := d.Get("role_name").(string)
	fullName := d.Get("full_name").(string)
	description := d.Get("description").(string)
	email := d.Get("email").(string)
	telephone := d.Get("telephone").(string)
	im := d.Get("im").(string)
	alertEmail := d.Get("alert_email").(string)
	alertEmailPrefix := d.Get("alert_email_prefix").(string)

	storedVmQuota := int32(d.Get("stored_vm_quota").(int))
	deployedVmQuota := int32(d.Get("deployed_vm_quota").(int))

	isGroupRole := d.Get("is_group_role").(bool)
	isDefaultCached := d.Get("is_default_cached").(bool)
	isExternal := d.Get("is_external").(bool)
	isAlertEnabled := d.Get("is_alert_enabled").(bool)
	isEnabled := d.Get("is_enabled").(bool)

	provider := providerGlobalRefPointer.userProvider

	createUserInfo := proto.CreateUserInfo{Name: name,
		Password:         password,
		RoleName:         roleName,
		FullName:         fullName,
		Description:      description,
		Email:            email,
		Telephone:        telephone,
		Im:               im,
		AlertEmail:       alertEmail,
		AlertEmailPrefix: alertEmailPrefix,
		StoredVmQuota:    storedVmQuota,
		DeployedVmQuota:  deployedVmQuota,
		IsGroupRole:      isGroupRole,
		IsDefaultCached:  isDefaultCached,
		IsExternal:       isExternal,
		IsAlertEnabled:   isAlertEnabled,
		IsEnabled:        isEnabled}

	res, err := provider.Create(createUserInfo)

	if err != nil {
		return fmt.Errorf("Error Creating User :[%+v] %#v", createUserInfo, err)
	}

	if res.Created {
		logging.Plog(fmt.Sprintf("User [%+v]  created  ", res))
		d.SetId(name)
	}

	logging.Plog("__DONE__resourceUserCreate_")
	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceUserDelete_")

	name := d.Get("name").(string)

	provider := providerGlobalRefPointer.userProvider

	deleteUserInfo := proto.DeleteUserInfo{Name: name}
	res, err := provider.Delete(deleteUserInfo)

	if err != nil {
		return fmt.Errorf("Error Deleting User :[%+v] %#v", deleteUserInfo, err)
	}

	if res.Deleted {
		logging.Plog(fmt.Sprintf("User [%+v]  deleted  ", res))
	}

	logging.Plog("__DONE__resourceUserDelete_")
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceUserUpdate_")

	name := d.Get("name").(string)

	oldIsEnabledRaw, newIsEnabledRaw := d.GetChange("is_enabled")
	oldIsEnabled := oldIsEnabledRaw.(bool)
	newIsEnabled := newIsEnabledRaw.(bool)

	provider := providerGlobalRefPointer.userProvider

	if !(oldIsEnabled == newIsEnabled) {
		updateUserInfo := proto.UpdateUserInfo{Name: name, IsEnabled: newIsEnabled}
		res, err := provider.Update(updateUserInfo)
		if err != nil {
			return fmt.Errorf("Error updating User :[%+v] %#v", updateUserInfo, err)
		}

		if res.Updated {
			logging.Plog(fmt.Sprintf("User [%+v]  updated  ", res))
			d.SetId(name)
		}
	} else {
		return fmt.Errorf("Error updating User :[%+v]. "+
			"Can not update the given fields ", name)
	}

	logging.Plog("__DONE__resourceUserUpdate_")
	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	logging.Plog("__INIT__resourceUserRead_")

	name := d.Get("name").(string)

	provider := providerGlobalRefPointer.userProvider

	readUserInfo := proto.ReadUserInfo{Name: name}
	res, err := provider.Read(readUserInfo)

	if err != nil {
		logging.Plog(fmt.Sprintf("Error Reading User :[%+v] %#v", readUserInfo, err))
		//return fmt.Errorf("Error Reading User :[%+v] %#v", readUserInfo, err)
	}

	if res.Present {
		d.SetId(name)
		d.Set("is_enabled", res.IsEnabled)
		logging.Plog(fmt.Sprintf("__DONE__resourceUserRead_ +setting id %v", name))
		return nil
	} else {
		d.SetId("")
		logging.Plog(fmt.Sprintf("__DONE__resourceUserRead_ +unsetting id,resource got deleted %v", name))
	}

	logging.Plog("__DONE__resourceUserRead_")
	return nil
}
