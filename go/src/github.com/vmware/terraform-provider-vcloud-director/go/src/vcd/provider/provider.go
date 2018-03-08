/*****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
******************************************************************/
package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/util/logging"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
)

// .. Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{

			"ip": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VCD_IP", nil),
				Description: "The vcd IP for vcd API operations.",
			},

			"user": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VCD_USER", nil),
				Description: "The user name for vcd API operations.",
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VCD_PASSWORD", nil),
				Description: "The user password for vcd API operations.",
			},

			"org": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VCD_ORG", nil),
				Description: "The vcd org for API operations",
			},

			"use_vcd_cli_profile": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VCD_USE_VCD_CLI_PROFILE", false),
				Description: "If set, VCDClient will use vcd cli profile and token .",
			},

			"allow_unverified_ssl": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VCD_ALLOW_UNVERIFIED_SSL", false),
				Description: "If set, VCDClient will permit unverifiable SSL certificates.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{

			"vcloud-director_catalog":            resourceCatalog(),
			"vcloud-director_catalog_item_media": resourceCatalogItemMedia(),
			"vcloud-director_catalog_item_ova":   resourceCatalogItemOva(),
			"vcloud-director_vapp":               resourceVApp(),
			"vcloud-director_independent_disk":   resourceIndependentDisk(),
			"vcloud-director_org":                resourceOrg(),
			"vcloud-director_disk":               resourceIndependentDisk(),
			"vcloud-director_user":               resourceUser(),
			"vcloud-director_vdc":                resourceVdc(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	logging.Plog("__INIT__func providerConfigure")

	l := proto.LoginCredentials{
		Username: d.Get("user").(string),
		Password: d.Get("password").(string),
		Org:      d.Get("org").(string),

		Ip: d.Get("ip").(string),

		UseVcdCliProfile: d.Get("use_vcd_cli_profile").(bool),

		AllowInsecureFlag: d.Get("allow_unverified_ssl").(bool),
	}

	config := Config{L: l}

	vcdclient, err := config.CreateClient()
	logging.Plog("__DONE__func providerConfigure Successfully configured Provider...\n")

	return vcdclient, err
}
