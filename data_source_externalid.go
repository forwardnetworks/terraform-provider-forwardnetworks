package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-forwardnetworks/forwardnetworks"
)

func dataSourceExternalId() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExternalIdRead,

		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network ID used to fetch the external ID.",
			},
			"external_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The external ID associated with the given network ID.",
			},
		},
	}
}

func dataSourceExternalIdRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardnetworks.ForwardNetworksClient)

	networkId := d.Get("network_id").(string)

	externalId, err := client.GetExternalId(networkId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("forwardnetworks_external_id")
	d.Set("external_id", externalId)

	return nil
}