package main

import (
    "context"
    "fmt"
    "time"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "terraform-provider-forwardnetworks/forwardnetworks"
)

func dataSourceForwardNetworksNQEQuery() *schema.Resource {
    return &schema.Resource{
        ReadContext: dataSourceForwardNetworksNQEQueryRead,

        Schema: map[string]*schema.Schema{
            "path": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "query_id": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "entries": {
                Type:     schema.TypeList,
                Computed: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "query_id": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "path": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "intent": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "repository": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                    },
                },
            },
        },
    }
}

func dataSourceForwardNetworksNQEQueryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)

    path := d.Get("path").(string)
    queryID, hasQueryID := d.GetOk("query_id")

    nqeQueries, err := client.GetNqeQueries(path)
    if err != nil {
        return diag.FromErr(err)
    }

    var entries []map[string]interface{}
    if !hasQueryID {
        for _, query := range nqeQueries {
            entry := map[string]interface{}{
                "query_id":    query.QueryID,
                "path":        query.Path,
                "intent":      query.Intent,
                "repository":  query.Repository,
            }
            entries = append(entries, entry)
        }
        d.SetId(fmt.Sprintf("%d", time.Now().Unix()))
    } else {
        queryIDStr := queryID.(string)
        for _, query := range nqeQueries {
            if query.QueryID == queryIDStr {
                entry := map[string]interface{}{
                    "query_id":    query.QueryID,
                    "path":        query.Path,
                    "intent":      query.Intent,
                    "repository":  query.Repository,
                }
                entries = append(entries, entry)
                break
            }
        }
        d.SetId(queryIDStr)
    }

    d.Set("entries", entries)

    return nil
}

func dataSourceForwardNetworksNQEQueryExecution() *schema.Resource {
    return &schema.Resource{
        ReadContext: dataSourceForwardNetworksNQEQueryExecutionRead,

        Schema: map[string]*schema.Schema{
            "network_id": {
                Type:     schema.TypeString,
                Required: true,
            },
            "query_id": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "query": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "offset": {
                Type:     schema.TypeInt,
                Optional: true,
            },
            "limit": {
                Type:     schema.TypeInt,
                Optional: true,
            },
            "sort_column_name": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "sort_order": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "filter_column_name": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "filter_value": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "parameters": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "result": {
                Type:     schema.TypeString,
                Computed: true,
            },
        },
    }
}

func dataSourceForwardNetworksNQEQueryExecutionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*forwardnetworks.ForwardNetworksClient)

    networkID := d.Get("network_id").(string)
    queryID := d.Get("query_id").(string)
    query := d.Get("query").(string)
    offset := d.Get("offset").(int)
    limit := d.Get("limit").(int)
    sortColumnName := d.Get("sort_column_name").(string)
    sortOrder := d.Get("sort_order").(string)
    filterColumnName := d.Get("filter_column_name").(string)
    filterValue := d.Get("filter_value").(string)
    parameters := d.Get("parameters").(string)

    queryOptions := make(map[string]interface{})

    if offset > 0 {
        queryOptions["offset"] = offset
    }

    if limit > 0 {
        queryOptions["limit"] = limit
    }

    if sortColumnName != "" && sortOrder != "" {
        queryOptions["sortBy"] = map[string]string{
            "columnName": sortColumnName,
            "order":      sortOrder,
        }
    }

    if filterColumnName != "" && filterValue != "" {
        queryOptions["columnFilters"] = []map[string]string{{
            "columnName": filterColumnName,
            "value":      filterValue,
        }}
    }

    var paramMap map[string]interface{}
    if parameters != "" {
        err := json.Unmarshal([]byte(parameters), &paramMap)
        if err != nil {
            return diag.FromErr(err)
        }
    }

    result, err := client.ExecuteNQEQuery(networkID, "", query, queryOptions, queryID)
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId(fmt.Sprintf("query-%s-%d", networkID, time.Now().Unix()))
    d.Set("result", string(result))

    return nil
}