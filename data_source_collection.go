package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-forwardnetworks/forwardnetworks"
)

func dataCollection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCollectionRead,
		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeString,
				Description: "The ID of the network to start collection.",
				Required:    true,
			},
		},
	}
}

func dataCollectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*forwardnetworks.ForwardNetworksClient)
	networkID := d.Get("network_id").(string)

	err := client.StartCollection(networkID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(networkID)

	return nil
}