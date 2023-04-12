package main

import (
    "context"
    "encoding/base64"
    "strings"
    "fmt"
    "strconv"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "terraform-provider-forwardnetworks/forwardnetworks"
)

func dataSourceSnapshots() *schema.Resource {
    return &schema.Resource{
        ReadContext: dataSourceSnapshotsRead,
        Schema: map[string]*schema.Schema{
            "network_id": {
                Type:        schema.TypeString,
                Description: "The ID of the network.",
                Required:    true,
            },
            "latest_processed": {
                Type:        schema.TypeBool,
                Description: "Whether to return only the latest processed snapshot.",
                Optional:    true,
                Default:     false,
            },
            "metrics": {
                Type:        schema.TypeBool,
                Description: "Whether to include metrics for each snapshot.",
                Optional:    true,
                Default:     false,
            },
            "params": {
                Type:        schema.TypeMap,
                Description: "Optional parameters for exporting the snapshot.",
                Optional:    true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "exported_snapshot": {
                Type:        schema.TypeString,
                Description: "The base64-encoded content of the exported snapshot zip file.",
                Computed:    true,
            },
            "snapshots": {
                Type:        schema.TypeList,
                Computed:    true,
                Description: "The list of snapshots.",
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "id": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "processing_trigger": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "total_devices": {
                            Type:     schema.TypeInt,
                            Computed: true,
                        },
                        "creation_date_millis": {
                            Type:     schema.TypeInt,
                            Computed: true,
                        },
                        "processed_at_millis": {
                            Type:     schema.TypeInt,
                            Computed: true,
                        },
                        "oldest_collection_millis": {
                            Type:     schema.TypeInt,
                            Computed: true,
                        },
                        "latest_collection_millis": {
                            Type:     schema.TypeInt,
                            Computed: true,
                        },
                        "is_draft": {
                            Type:     schema.TypeBool,
                            Computed: true,
                        },
                        "state": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "metrics": {
                            Type:     schema.TypeMap,
                            Optional: true,
                            Elem: &schema.Schema{
                                Type: schema.TypeString,
                            },
                        },
                    },
                },
            },
        },
    }
}

func dataSourceSnapshotsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)
    networkID := d.Get("network_id").(string)
    latestProcessed := d.Get("latest_processed").(bool)
    metrics := d.Get("metrics").(bool)

    snapshots, err := client.ListSnapshots(networkID, latestProcessed, metrics)
    if err != nil {
        return diag.FromErr(err)
    }

    err = d.Set("snapshots", flattenSnapshots(snapshots))
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId(networkID)

    return nil
}

func flattenSnapshots(snapshots []forwardnetworks.SnapshotWithMetrics) []map[string]interface{} {
    result := make([]map[string]interface{}, len(snapshots))

    for i, snapshot := range snapshots {
        metrics := make(map[string]interface{})
        if snapshot.Metrics != nil {
            metrics["collection_concurrency"] = strconv.Itoa(snapshot.Metrics.CollectionConcurrency)
            metrics["collection_duration"] = strconv.Itoa(snapshot.Metrics.CollectionDuration)
            metrics["collection_failures"] = fmt.Sprintf("%v", snapshot.Metrics.CollectionFailures) // updated
            metrics["creation_date_millis"] = strconv.FormatInt(snapshot.Metrics.CreationDateMillis, 10) // updated
            metrics["host_computation_status"] = snapshot.Metrics.HostComputationStatus
            metrics["ip_location_indexing_status"] = snapshot.Metrics.IpLocationIndexingStatus
            metrics["jump_server_collection_concurrency"] = strconv.Itoa(snapshot.Metrics.JumpServerCollectionConcurrency)
            metrics["l2_indexing_status"] = snapshot.Metrics.L2IndexingStatus
            metrics["needs_reprocessing"] = strconv.FormatBool(snapshot.Metrics.NeedsReprocessing)
            metrics["num_collection_failure_devices"] = strconv.Itoa(snapshot.Metrics.NumCollectionFailureDevices)
            metrics["num_parsing_failure_devices"] = strconv.Itoa(snapshot.Metrics.NumParsingFailureDevices)
            metrics["num_successful_devices"] = strconv.Itoa(snapshot.Metrics.NumSuccessfulDevices)
            metrics["parsing_failures"] = fmt.Sprintf("%v", snapshot.Metrics.ParsingFailures) // updated
            metrics["path_search_indexing_status"] = snapshot.Metrics.PathSearchIndexingStatus
            metrics["processing_duration"] = strconv.Itoa(snapshot.Metrics.ProcessingDuration)
            metrics["search_indexing_status"] = snapshot.Metrics.SearchIndexingStatus
            metrics["snapshot_id"] = snapshot.Metrics.SnapshotID
        }

        result[i] = map[string]interface{}{
            "id":                     snapshot.ID,
            "processing_trigger":     snapshot.ProcessingTrigger,
            "total_devices":          snapshot.TotalDevices,
            "creation_date_millis":   snapshot.CreationDateMillis,
            "processed_at_millis":    snapshot.ProcessedAtMillis,
            "oldest_collection_millis": snapshot.OldestCollectionMillis,
            "latest_collection_millis": snapshot.LatestCollectionMillis,
            "is_draft":               snapshot.IsDraft,
            "state":                  snapshot.State,
            "metrics":                metrics,
        }
    }

    return result
}

func dataSourceExportSnapshotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)
    snapshotID := d.Get("snapshot_id").(string)
    paramsRaw := d.Get("params").(map[string]interface{})

    var data []byte
    var err error

    exportParams := forwardnetworks.ExportParams{
        IncludeDevices: convertCSVStringToStringSlice(paramsRaw["includeDevices"].(string)),
        ExcludeDevices: convertCSVStringToStringSlice(paramsRaw["excludeDevices"].(string)),
        ObfuscationKey: paramsRaw["obfuscationKey"].(string),
        ObfuscateNames: parseBool(paramsRaw["obfuscateNames"].(string)),
    }
    data, err = client.ExportSnapshotWithParams(snapshotID, &exportParams)

    if err != nil {
        return diag.FromErr(err)
    }

    encodedData := base64.StdEncoding.EncodeToString(data)

    d.Set("exported_snapshot", encodedData)
    d.SetId(snapshotID)

    return nil
}

func convertCSVStringToStringSlice(input string) []string {
    if input == "" {
        return nil
    }
    return strings.Split(input, ",")
}

func parseBool(input string) bool {
    return strings.ToLower(input) == "true"
}