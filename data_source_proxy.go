package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-forwardnetworks/forwardnetworks"
)

func dataSourceProxy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProxyRead,

		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network ID used to fetch the proxy information.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the proxy.",
			},
			"host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host of the proxy.",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port of the proxy.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username for the proxy.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol used by the proxy.",
			},
			"disable_cert_checking": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether certificate checking is disabled for the proxy.",
			},
		},
	}
}

func dataSourceProxyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardnetworks.ForwardNetworksClient)

	networkId := d.Get("network_id").(string)

	proxy, err := client.GetProxy(networkId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(proxy.Id)
	d.Set("id", proxy.Id)
	d.Set("host", proxy.Host)
	d.Set("port", proxy.Port)
	d.Set("username", proxy.Username)
	d.Set("protocol", proxy.Protocol)
	d.Set("disable_cert_checking", proxy.DisableCertChecking)

	return nil
}
