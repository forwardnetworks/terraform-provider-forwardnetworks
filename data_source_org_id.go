package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-forwardnetworks/forwardnetworks"
)

func dataSourceOrgId() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrgIdRead,

		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOrgIdRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardnetworks.ForwardNetworksClient)

	orgId, orgName, err := client.GetOrgId()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("forwardnetworks_org_id")
	d.Set("org_id", orgId)
	d.Set("org_name", orgName)

	return nil
}
