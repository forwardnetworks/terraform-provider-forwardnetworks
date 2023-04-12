package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-forwardnetworks/forwardnetworks"
)

func dataSourceCloudAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudAccountsRead,

		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cloud_accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"regions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"assume_role_infos": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"accountid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"accountname": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rolearn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"subscriptions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subscriptionid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"clientid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tenant": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proxy_server_id": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
					},
				},
			},
		},
	}
}

func dataSourceCloudAccountsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardnetworks.ForwardNetworksClient)

	networkId := d.Get("network_id").(string)
	accountName := d.Get("account_name").(string)

	cloudAccounts, err := client.GetCloudAccounts(networkId, accountName)
	if err != nil {
		return diag.FromErr(err)
	}

	cloudAccountsList := make([]map[string]interface{}, 0, len(cloudAccounts))

	for _, account := range cloudAccounts {
		assumeRoleInfos := make([]map[string]interface{}, len(account.AssumeRoleInfos))
		for j, info := range account.AssumeRoleInfos {
			assumeRoleInfos[j] = map[string]interface{}{
				"accountid":   info.AccountId,
				"accountname": info.AccountName,
				"rolearn":     info.RoleArn,
				"enabled":     info.Enabled,
			}
		}

		subscriptions := make([]map[string]interface{}, len(account.Subscriptions))
		for j, sub := range account.Subscriptions {
			subscriptions[j] = map[string]interface{}{
				"subscriptionid": sub.SubscriptionId,
				"clientid":       sub.ClientId,
				"tenant":         sub.Tenant,
				"enabled":        sub.Enabled,
			}
		}

		cloudAccountMap := map[string]interface{}{
			"name":             account.Name,
			"type":             account.Type,
			"collect":          account.Collect,
			"regions":          account.Regions,
			"assume_role_infos": assumeRoleInfos,
			"subscriptions":    subscriptions,
			"username":     account.AWSUsername,
			"proxy_server_id":  account.ProxyServerId,
		}
		cloudAccountsList = append(cloudAccountsList, cloudAccountMap)
	}

	d.SetId("forwardnetworks_cloud_accounts")
	if err := d.Set("cloud_accounts", cloudAccountsList); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
