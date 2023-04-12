package main

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

    "terraform-provider-forwardnetworks/forwardnetworks"
)

func resourceCollectors() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceCollectorCreateandUpdate,
        ReadContext:   resourceCollectorRead,
        UpdateContext: resourceCollectorCreateandUpdate,
        DeleteContext: resourceCollectorDelete,

        Schema: map[string]*schema.Schema{
            "network_id": {
                Type:        schema.TypeString,
                Description: "The ID of the network to which the collector belongs.",
                Required:    true,
                ForceNew:    true,
            },
            "collector_name": {
                Type:        schema.TypeString,
                Description: "The name of the collector.",
                Required:    true,
            },
            "collector_username": {
                Type:        schema.TypeString,
                Description: "The username associated with the collector.",
                Computed:    true, // Set to Computed since it's not provided as input
            },
        },
    }
}

func resourceCollectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)

    collectorName := d.Get("collector_name").(string)

    collectors, err := client.GetCollectors()
    if err != nil {
        return diag.FromErr(err)
    }

    var collector *forwardnetworks.Collector
    for _, c := range collectors {
        if c.CollectorName == collectorName {
            collector = &c
            break
        }
    }

    if collector == nil {
        return diag.Errorf("collector with name %q not found", collectorName)
    }

    d.SetId(collectorName)
    d.Set("collector_username", collector.Username)

    return nil
}

func resourceCollectorCreateandUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    // Call resourceCollectorRead to populate the ResourceData with collector_username
    if diags := resourceCollectorRead(ctx, d, m); diags.HasError() {
        return diags
    }

    client := m.(*forwardnetworks.ForwardNetworksClient)
    networkID := d.Get("network_id").(string)
    collectorName := d.Get("collector_name").(string)
    collectorUsername := d.Get("collector_username").(string) // Retrieve the populated collector_username

    err := client.UpdateCollector(networkID, collectorName, collectorUsername)
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId(collectorName)

    return nil
}

func resourceCollectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    d.SetId("")
    fmt.Println("Warning: Deleting collectors is not supported by the Forward Networks API")
    return nil
}
