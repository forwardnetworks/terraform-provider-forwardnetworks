package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-forwardnetworks/forwardnetworks"
	"log"
)

// This resource manages the proxy configuration for a Forward Networks network.
//
// Example Usage
//
// ```hcl
// resource "forwardnetworks_proxy" "example" {
//   network_id           = "your_network_id"
//   protocol             = "https"
//   host                 = "proxy.example.com"
//   port                 = 8080
//   username             = "proxyuser"
//   password             = "proxypassword"
//   disable_cert_checking = false
// }
// ```

func resourceProxy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProxyCreateOrUpdate,
		ReadContext:   resourceProxyRead,
		UpdateContext: resourceProxyCreateOrUpdate,
		DeleteContext: resourceProxyDelete,

		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network ID used to manage the proxy.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The protocol used by the proxy.",
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The host of the proxy.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The port of the proxy.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username for the proxy.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The password for the proxy.",
				Sensitive:   true,
			},
			"disable_cert_checking": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether certificate checking is disabled for the proxy.",
			},
		},
	}
}

func resourceProxyCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardnetworks.ForwardNetworksClient)

	networkId := d.Get("network_id").(string)
	proxy := &forwardnetworks.Proxy{
		Protocol:            d.Get("protocol").(string),
		Host:                d.Get("host").(string),
		Port:                d.Get("port").(int),
		Username:            d.Get("username").(string),
		Password:            d.Get("password").(string),
		DisableCertChecking: d.Get("disable_cert_checking").(bool),
	}

	err := client.CreateOrUpdateProxy(networkId, proxy)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(networkId)
	return resourceProxyRead(ctx, d, meta)
}

func resourceProxyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardnetworks.ForwardNetworksClient)

	networkId := d.Get("network_id").(string)

	proxy, err := client.GetProxy(networkId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(proxy.Id)
	d.Set("protocol", proxy.Protocol)
	d.Set("host", proxy.Host)
	d.Set("port", proxy.Port)
	d.Set("username", proxy.Username)
	d.Set("disable_cert_checking", proxy.DisableCertChecking)

	return nil
}

func resourceProxyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// There is no delete operation for the proxy.
	log.Printf("[WARN] Delete operation called for Forward Networks proxy, but delete is not supported. You can only update the proxy.")
	return nil
}
