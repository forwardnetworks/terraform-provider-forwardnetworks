package main

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-forwardnetworks/forwardnetworks"
)

func dataSourceCloud() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudRead,

		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"setup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cloud_accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"collect": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"num_virtualized_devices": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"regions": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assume_role_infos": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardnetworks.ForwardNetworksClient)

	networkID := d.Get("network_id").(string)
	setupID := d.Get("setup_id").(string)

	var cloudAccounts []*forwardnetworks.CloudAccount
	var err error

	if setupID != "" {
		account, err := client.GetCloudAccount(networkID, setupID)
		if err != nil {
			return diag.FromErr(err)
		}
		cloudAccounts = []*forwardnetworks.CloudAccount{account}
	} else {
		cloudAccounts, err = client.ListCloudAccounts(networkID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	accounts := make([]map[string]interface{}, len(cloudAccounts))

	for i, account := range cloudAccounts {
		accountMap := make(map[string]interface{})
		accountMap["name"] = account.Name
		accountMap["type"] = account.Type
		accountMap["collect"] = account.Collect
		accountMap["num_virtualized_devices"] = account.NumVirtualizedDevices

		regionsJSON, _ := json.Marshal(account.Regions)
		accountMap["regions"] = string(regionsJSON)

		if account.AssumeRoleInfos != nil {
			assumeRoleInfosJSON, _ := json.Marshal(account.AssumeRoleInfos)
			accountMap["assume_role_infos"] = string(assumeRoleInfosJSON)
		}

		accounts[i] = accountMap
	}

	d.SetId("forwardnetworks_cloud_" + networkID)
	d.Set("cloud_accounts", accounts)

	return nil
}
