package main

import (
    "context"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "terraform-provider-forwardnetworks/forwardnetworks"
)

func resourceNetworks() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceNetworkCreate,
        ReadContext:   resourceNetworkRead,
        UpdateContext: resourceNetworkUpdate,
        DeleteContext: resourceNetworkDelete,

        Schema: map[string]*schema.Schema{
            "name": {
                Type:     schema.TypeString,
                Required: true,
            },
            "note": {
                Type:     schema.TypeString,
                Optional: true,
            },
        },
    }
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)

    name := d.Get("name").(string)
    note := d.Get("note").(string)

    newNetwork := forwardnetworks.Network{
        Name: name,
        Note: note,
    }

    network, err := client.CreateNetwork(&newNetwork)
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId(network.ID)

    // Force an update to add the note field
    updateDiags := resourceNetworkUpdate(ctx, d, m)
    if updateDiags != nil {
        return updateDiags
    }

    return resourceNetworkRead(ctx, d, m)
}


func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)

    networkId := d.Id()

    network, err := client.GetNetwork(networkId)
    if err != nil {
        return diag.FromErr(err)
    }

    if network == nil {
        d.SetId("")
        return nil
    }

    d.Set("name", network.Name)
    d.Set("note", network.Note)

    return nil
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)

    networkId := d.Id()

    update := forwardnetworks.NetworkUpdate{
        Name: d.Get("name").(string),
        Note: d.Get("note").(string),
    }

    _, err := client.UpdateNetwork(networkId, update)
    if err != nil {
        return diag.FromErr(err)
    }

    return resourceNetworkRead(ctx, d, m)
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)

    networkId := d.Id()

    err := client.DeleteNetwork(networkId)
    if err != nil {
        return diag.Errorf("Error deleting network: %s", err)
    }

    return nil
}
