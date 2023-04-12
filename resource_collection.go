package main

import (
    "context"
    "errors"
    "time"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "terraform-provider-forwardnetworks/forwardnetworks"
)

func resourceCollection() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceCollectionCreate,
        ReadContext:   resourceCollectionRead,
        UpdateContext: resourceCollectionUpdate,
        DeleteContext: resourceCollectionDelete,
        Schema: map[string]*schema.Schema{
            "network_id": {
                Type:        schema.TypeString,
                Description: "The ID of the network to start collection.",
                Required:    true,
            },
            "force_refresh": {
                Type:        schema.TypeBool,
                Description: "Force a new snapshot collection.",
                Optional:    true,
                Default:     false,
                ForceNew:    true,
            },
        },
    }
}

func resourceCollectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)
    networkID := d.Get("network_id").(string)

    snapshotIDBefore, err := getLatestProcessedSnapshotID(client, networkID)
    if err != nil {
        return diag.FromErr(err)
    }

    // Start a new collection only if there are no existing snapshots
    if snapshotIDBefore == "" {
        err = client.StartCollection(networkID)
        if err != nil {
            return diag.FromErr(err)
        }
    }

    snapshotIDAfter, err := getLatestProcessedSnapshotIDWithRetry(client, networkID, snapshotIDBefore)
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId(snapshotIDAfter)

    return resourceCollectionRead(ctx, d, m)
}

func resourceCollectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)
    networkID := d.Get("network_id").(string)

    // Get the latest processed snapshot ID
    snapshotID, err := getLatestProcessedSnapshotID(client, networkID)
    if err != nil {
        return diag.FromErr(err)
    }

    if snapshotID != "" {
        // Check if the resource state is not set
        if d.Id() == "" {
            // Update the resource state with the new snapshot ID
            d.SetId(snapshotID)
        }
    } else {
        // If there are no processed snapshots, clear the resource state
        d.SetId("")
    }

    return nil
}


func resourceCollectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)
    networkID := d.Get("network_id").(string)

    forceRefresh := d.Get("force_refresh").(bool)

    if forceRefresh {
        // Start a new collection
        err := client.StartCollection(networkID)
        if err != nil {
            return diag.FromErr(err)
        }

        // Wait for the StartCollection to finish and get the new snapshot ID
        currentStateSnapshotID := d.Id()
        newSnapshotID, err := getLatestProcessedSnapshotIDWithRetry(client, networkID, currentStateSnapshotID)
        if err != nil {
            return diag.FromErr(err)
        }

        d.SetId(newSnapshotID)
    }

    return resourceCollectionRead(ctx, d, m)
}


func resourceCollectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    // Do nothing
    return nil
}

func getLatestProcessedSnapshotID(client *forwardnetworks.ForwardNetworksClient, networkID string) (string, error) {
    snapshots, err := client.ListSnapshots(networkID, true, false)
    if err != nil {
        if apiErr, ok := err.(forwardnetworks.ApiError); ok {
            if apiErr.Message == "The network has no Snapshots." || apiErr.Message == "None of the Snapshots in the network are processed. Processing of the latest Snapshot has begun." {
                // Return empty string if there are no snapshots or no processed snapshots
                return "", nil
            } else {
                return "", err
            }
        } else {
            return "", err
        }
    }

    if len(snapshots) == 0 {
        // No processed snapshot found, return empty string
        return "", nil
    }

    return snapshots[0].ID, nil
}
func getLatestProcessedSnapshotIDWithRetry(client *forwardnetworks.ForwardNetworksClient, networkID string, previousSnapshotID string) (string, error) {
    maxRetries := 10
    sleepInterval := 5 * time.Second

    for i := 0; i < maxRetries; i++ {
        snapshotID, err := getLatestProcessedSnapshotID(client, networkID)
        if err != nil {
            return "", err
        }

        if snapshotID != "" {
            return snapshotID, nil
        }

        time.Sleep(sleepInterval)
    }

    return "", errors.New("Could not get the new snapshot ID after retries")
}