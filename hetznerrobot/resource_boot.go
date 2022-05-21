package hetznerrobot

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func resourceBoot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBootCreate,
		ReadContext:   resourceBootRead,
		UpdateContext: resourceBootUpdate,
		DeleteContext: resourceBootDelete,

		Importer: &schema.ResourceImporter{
			State: resourceBootImportState,
		},

		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Server ID",
			},
			// optional
			"active_profile": {
				Type:        schema.TypeString, // Enum should be better (linux/rescue/...)
				Optional:    true,
				Description: "Active boot profile",
			},
			"architecture": {
				Type:        schema.TypeString, // Enum should be better (amd64/...)
				Optional:    true,
				Description: "Active Architecture",
			},
			"language": {
				Type:        schema.TypeString, // Enum should be better (amd64/...)
				Optional:    true,
				Description: "Language",
			},
			"operating_system": {
				Type:        schema.TypeString, // Enum should be better (ubuntu_20.04/...)
				Optional:    true,
				Description: "Active Operating System / Distribution",
			},
			"authorized_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "One or more SSH key fingerprints",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// read-only / computed
			"ipv4_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server main IPv4 address",
			},
			"ipv6_network": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server main IPv6 net address",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current Rescue System root password / Linux installation password or null",
				Sensitive:   true,
			},
		},
	}
}

func resourceBootImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(HetznerRobotClient)

	serverID, _ := strconv.Atoi(d.Id())

	boot, err := c.getBoot(serverID)
	if err != nil {
		return nil, err
	}

	d.Set("active_profile", boot.ActiveProfile)
	d.Set("architecture", boot.Architecture)
	d.Set("ipv4_address", boot.ServerIPv4)
	d.Set("ipv6_network", boot.ServerIPv6)
	d.Set("language", boot.Language)
	d.Set("operating_system", boot.OperatingSystem)
	d.Set("password", boot.Password)
	d.Set("server_id", serverID)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceBootCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	serverID := d.Get("server_id").(int)
	activeBootProfile := d.Get("active_profile").(string)
	arch := d.Get("architecture").(string)
	os := d.Get("operating_system").(string)
	lang := d.Get("language").(string)
	authorizedKeys := make([]string, 0)
	if input := d.Get("authorized_keys"); input != nil {
		for _, key := range input.([]interface{}) {
			authorizedKeys = append(authorizedKeys, key.(string))
		}
	}

	bootProfile, err := c.setBootProfile(serverID, activeBootProfile, arch, os, lang, authorizedKeys)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("ipv4_address", bootProfile.ServerIPv4)
	d.Set("ipv6_network", bootProfile.ServerIPv6)
	d.Set("password", bootProfile.Password)
	d.SetId(strconv.Itoa(serverID))

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceBootRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	serverID, _ := strconv.Atoi(d.Id())
	boot, err := c.getBoot(serverID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("active_profile", boot.ActiveProfile)
	d.Set("architecture", boot.Architecture)
	d.Set("ipv4_address", boot.ServerIPv4)
	d.Set("ipv6_network", boot.ServerIPv6)
	d.Set("language", boot.Language)
	d.Set("operating_system", boot.OperatingSystem)
	d.Set("password", boot.Password)
	d.Set("server_id", serverID)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceBootUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	serverID, _ := strconv.Atoi(d.Id())
	activeBootProfile := d.Get("active_profile").(string)
	arch := d.Get("architecture").(string)
	os := d.Get("operating_system").(string)
	lang := d.Get("language").(string)
	authorizedKeys := make([]string, 0)
	if input := d.Get("authorized_keys"); input != nil {
		for _, key := range input.([]interface{}) {
			authorizedKeys = append(authorizedKeys, key.(string))
		}
	}

	bootProfile, err := c.setBootProfile(serverID, activeBootProfile, arch, os, lang, authorizedKeys)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("ipv4_address", bootProfile.ServerIPv4)
	d.Set("ipv6_network", bootProfile.ServerIPv6)
	d.Set("password", bootProfile.Password)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceBootDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
