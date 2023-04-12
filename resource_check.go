package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-forwardnetworks/forwardnetworks"
)

func resourceCheck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCheckCreate,
		ReadContext:   resourceCheckRead,
		UpdateContext: resourceCheckUpdate,
		DeleteContext: resourceCheckDelete,

		Schema: map[string]*schema.Schema{
			"snapshot_id": {
				Type:        schema.TypeString,
				Description: "The ID of the snapshot to which the check belongs.",
				Required:    true,
				ForceNew:    true,
			},
			"check_id": {
				Type:        schema.TypeString,
				Description: "The ID of the check.",
				Computed:    true,
			},
			"check_type": {
				Type:        schema.TypeString,
				Description: "The type of the check.",
				Required:    true,
			},
			"query_id": {
				Type:        schema.TypeString,
				Description: "The query ID of the check.",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the check is enabled or not.",
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceCheckCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*forwardnetworks.ForwardNetworksClient)
	snapshotID := d.Get("snapshot_id").(string)
	checkType := d.Get("check_type").(string)
	queryID := d.Get("query_id").(string)
	enabled := d.Get("enabled").(bool)

	check, err := client.ActivateCheck(snapshotID, checkType, queryID, enabled)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(check.ID)

	return resourceCheckRead(ctx, d, m)
}

func resourceCheckRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*forwardnetworks.ForwardNetworksClient)
	snapshotID := d.Get("snapshot_id").(string)
	checkID := d.Id()

	check, err := client.GetCheck(snapshotID, checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("check_type", check.Definition.CheckType)
	d.Set("query_id", check.Definition.QueryID)
	d.Set("enabled", check.Enabled)

	return nil
}

func resourceCheckUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// For this example, we only support enabling/disabling the check.
	if d.HasChange("enabled") {
		client := m.(*forwardnetworks.ForwardNetworksClient)
		snapshotID := d.Get("snapshot_id").(string)
		checkID := d.Id()
		enabled := d.Get("enabled").(bool)

		// If the check is being disabled, use DeactivateCheck to disable it.
		if !enabled {
			err := client.DeactivateCheck(snapshotID, checkID)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			// If the check is being enabled, recreate it with ActivateCheck.
			checkType := d.Get("check_type").(string)
			queryID := d.Get("query_id").(string)

			_, err := client.ActivateCheck(snapshotID, checkType, queryID, enabled)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceCheckRead(ctx, d, m)
}

func resourceCheckDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*forwardnetworks.ForwardNetworksClient)
	snapshotID := d.Get("snapshot_id").(string)
	checkID := d.Id()

	err := client.DeactivateCheck(snapshotID, checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

