package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-forwardnetworks/forwardnetworks"
)

func dataSourceChecks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceChecksRead,
		Schema: map[string]*schema.Schema{
			"snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"check_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"checks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"predefined_check_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date_millis": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"definition_date_millis": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"execution_date_millis": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"execution_duration_millis": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceChecksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*forwardnetworks.ForwardNetworksClient)

	snapshotID := d.Get("snapshot_id").(string)
	checkType := d.Get("type").(string)
	priority := d.Get("priority").(string)
	status := d.Get("status").(string)
	checkID, hasCheckID := d.GetOk("check_id")

	checks, err := client.GetChecks(snapshotID, checkType, priority, status)
	if err != nil {
		return diag.FromErr(err)
	}

	if hasCheckID {
		checkIDStr := checkID.(string)
		var foundCheck *forwardnetworks.Check
		for _, check := range checks {
			if check.ID == checkIDStr {
				foundCheck = &check
				break
			}
		}

		if foundCheck == nil {
			return diag.Errorf("No check found with ID: %s", checkIDStr)
		}

		checks = []forwardnetworks.Check{*foundCheck}
	}

	checkList := make([]map[string]interface{}, len(checks))
	for i, check := range checks {
		checkList[i] = map[string]interface{}{
			"id":                       check.ID,
			"predefined_check_type":    check.Definition.PredefinedCheckType,
			"check_type":               check.Definition.CheckType,
			"enabled":                  check.Enabled,
			"priority":                 check.Priority,
			"name":                     check.Name,
			"creation_date_millis":     check.CreationDateMillis,
			"creator_id":               check.CreatorId,
			"definition_date_millis":   check.DefinitionDateMillis,
			"description":              check.Description,
			"status":                   check.Status,
			"execution_date_millis":    check.ExecutionDateMillis,
			"execution_duration_millis": check.ExecutionDurationMillis,
		}
	}

	if err := d.Set("checks", checkList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(snapshotID)

	return nil
}
