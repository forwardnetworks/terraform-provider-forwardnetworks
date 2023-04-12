package main

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-forwardnetworks/forwardnetworks"
)

func resourceCloudAccount() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceCloudAccountCreate,
        ReadContext:   resourceCloudAccountRead,
        UpdateContext: resourceCloudAccountUpdate,
        DeleteContext: resourceCloudAccountDelete,

        Schema: map[string]*schema.Schema{
            "network_id": {
                Type:     schema.TypeString,
                Required: true,
            },
            "type": {
                Type:     schema.TypeString,
                Required: true,
            },
            "name": {
                Type:     schema.TypeString,
                Required: true,
            },
            "collect": {
                Type:     schema.TypeBool,
                Optional: true,
                Default:  true,
            },
            "account_id": {
                Type:     schema.TypeList,
                Optional: true,
                Elem:     &schema.Schema{Type: schema.TypeString},
            },
            "rolearn": {
                Type:     schema.TypeList,
                Optional: true,
                Elem:     &schema.Schema{Type: schema.TypeString},
            },
            "account_name": {
                Type:     schema.TypeList,
                Optional: true,
                Elem:     &schema.Schema{Type: schema.TypeString},
            },
            "external_id": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "username": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "secret": {
                Type:       schema.TypeString,
                Optional:   true,
                Sensitive:  true,
            },
            "regions": {
                Type:     schema.TypeList,
                Optional: true,
                Elem:     &schema.Schema{Type: schema.TypeString},
            },
            "client_id": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "client_email": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "tenant": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "private_key_id": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "private_key": {
                Type:     schema.TypeString,
                Optional: true,
                Sensitive: true,
            },
            "proxy_server_id": {
                Type:     schema.TypeString,
                Optional: true,
                Sensitive: true,
            },
            "testinstant": {
                Type:     schema.TypeInt,
                Optional: true,
            },
            "subscription_id": {
                Type:     schema.TypeList,
                Optional: true,
                Elem:     &schema.Schema{Type: schema.TypeString},
            },
            "environment": {
                Type:     schema.TypeString,
                Optional: true,
                Default:  "AZURE",
            },
        },
    }
}


func resourceCloudAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*forwardnetworks.ForwardNetworksClient)

    networkId := d.Get("network_id").(string)
    accountType := d.Get("type").(string)

    var account *forwardnetworks.CloudAccount
    if accountType == "AWS" {
        username := d.Get("username").(string)
        if username != "" {
            account = createAWSIAMUserAccount(d)
        } else {
            account = createAWSAssumeRoleAccount(d)
        }
    } else if accountType == "AZURE" {
        account = createAzureCloudAccount(d)
    } else if accountType == "GCP" {
        account = createGCPCloudAccount(d)
    } else {
        return diag.Errorf("Invalid account type: %s", accountType)
    }

    err := client.CreateCloudAccount(networkId, account)
    if err != nil {
        return diag.FromErr(err)
    }

    // Set the ID of the created cloud account.
    d.SetId(account.Name)

    return resourceCloudAccountRead(ctx, d, meta)
}


func createAWSAssumeRoleAccount(d *schema.ResourceData) *forwardnetworks.CloudAccount {
    name := d.Get("name").(string)
    collect := d.Get("collect").(bool)

    roleArnsRaw := d.Get("rolearn").([]interface{})
    accountIdsRaw := d.Get("account_id").([]interface{})
    accountNamesRaw := d.Get("account_name").([]interface{})
    externalId := d.Get("external_id").(string)

    assumeRoleInfos := make([]forwardnetworks.AssumeRoleInfo, len(roleArnsRaw))
    for i := range roleArnsRaw {
        assumeRoleInfos[i] = forwardnetworks.AssumeRoleInfo{
            RoleArn:    roleArnsRaw[i].(string),
            AccountId:  accountIdsRaw[i].(string),
            AccountName: accountNamesRaw[i].(string),
            ExternalId:  externalId,
            Enabled:     true,
        }
    }

    regionsRaw := d.Get("regions").([]interface{})

    regions := make(map[string]interface{})
    for _, region := range regionsRaw {
        regionName := region.(string)
        regions[regionName] = time.Now().Unix() * 1000
    }

    account := forwardnetworks.CloudAccount{
        Name:            name,
        Type:            "AWS",
        Collect:         collect,
        AssumeRoleInfos: assumeRoleInfos,
        Regions:         regions,
    }

proxyServerIdRaw := d.Get("proxy_server_id")

if proxyServerIdRaw != nil {
    proxyServerId := proxyServerIdRaw.(string)

    // Only include ProxyServerId if it exists
    if proxyServerId != "" {
        account.ProxyServerId = proxyServerId
    }
}

return &account
}

func createAWSIAMUserAccount(d *schema.ResourceData) *forwardnetworks.CloudAccount {
    name := d.Get("name").(string)
    collect := d.Get("collect").(bool)

    roleArnsRaw := d.Get("rolearn").([]interface{})
    accountIdsRaw := d.Get("account_id").([]interface{})
    accountNamesRaw := d.Get("account_name").([]interface{})

    assumeRoleInfos := make([]forwardnetworks.AssumeRoleInfo, len(roleArnsRaw))
    for i := range roleArnsRaw {
        assumeRoleInfos[i] = forwardnetworks.AssumeRoleInfo{
            RoleArn:     roleArnsRaw[i].(string),
            AccountId:   accountIdsRaw[i].(string),
            AccountName: accountNamesRaw[i].(string),
            Enabled:     true,
        }
    }

    regionsRaw := d.Get("regions").([]interface{})

    regions := make(map[string]interface{})
    for _, region := range regionsRaw {
        regionName := region.(string)
        regions[regionName] = time.Now().Unix() * 1000
    }

    account := forwardnetworks.CloudAccount{
        Name:            name,
        Type:            "AWS",
        Collect:         collect,
        AssumeRoleInfos: assumeRoleInfos,
        Regions:         regions,
        AWSUsername:     d.Get("username").(string),
        Secret:     d.Get("secret").(string),
    }

    proxyServerId := d.Get("proxy_server_id").(string)

    // Only include ProxyServerId if it exists
    if proxyServerId != "" {
        account.ProxyServerId = proxyServerId
    }

    return &account
}


func createAzureCloudAccount(d *schema.ResourceData) *forwardnetworks.CloudAccount {
    name := d.Get("name").(string)
    collect := d.Get("collect").(bool)

    subscriptionIdsRaw := d.Get("subscription_id").([]interface{})
    clientId := d.Get("client_id").(string)
    tenant := d.Get("tenant").(string)
    secret := d.Get("secret").(string)
    environment := d.Get("environment").(string)

    subscriptions := make([]forwardnetworks.Subscription, len(subscriptionIdsRaw))
    for i := range subscriptionIdsRaw {
        subscriptions[i] = forwardnetworks.Subscription{
            SubscriptionId: subscriptionIdsRaw[i].(string),
            ClientId:       clientId,
            Tenant:         tenant,
            Secret:         secret,
            TestInstant:    time.Now().Unix() * 1000,
            Environment:    environment,
            Enabled:        true,
        }
    }

    account := forwardnetworks.CloudAccount{
        Name:          name,
        Type:          "AZURE",
        Collect:       collect,
        Subscriptions: subscriptions,
    }

    proxyServerId := d.Get("proxy_server_id").(string)

    // Only include ProxyServerId if it exists
    if proxyServerId != "" {
        account.ProxyServerId = proxyServerId
    }

    return &account

}

func createGCPCloudAccount(d *schema.ResourceData) *forwardnetworks.CloudAccount {
    regionsRaw := d.Get("regions").([]interface{})

    regions := make(map[string]interface{})
    for _, region := range regionsRaw {
        regionName := region.(string)
        regions[regionName] = time.Now().Unix() * 1000
    }

    account := forwardnetworks.CloudAccount{
        Name: d.Get("name").(string),
        Type: d.Get("type").(string),
        Collect: d.Get("collect").(bool),
        ClientID: d.Get("client_id").(string),
        ClientEmail: d.Get("client_email").(string),
        PrivateKeyID: d.Get("private_key_id").(string),
        PrivateKey: d.Get("private_key").(string),
        Regions: regions,
    }
    proxyServerId := d.Get("proxy_server_id").(string)

    // Only include ProxyServerId if it exists
    if proxyServerId != "" {
        account.ProxyServerId = proxyServerId
    }

    return &account
}

func resourceCloudAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*forwardnetworks.ForwardNetworksClient)

    networkId := d.Get("network_id").(string)
    accountName := d.Get("name").(string)

    cloudAccounts, err := client.GetCloudAccounts(networkId, accountName)
    if err != nil {
        return diag.FromErr(err)
    }

    if len(cloudAccounts) != 1 {
        return diag.Errorf("Unable to find a single cloud account with name %s", accountName)
    }

    account, ok := cloudAccounts[accountName]
if !ok {
    return diag.Errorf("Unable to find a cloud account with name %s", accountName)
}

    // Set common attributes for all account types
    d.Set("name", account.Name)
    d.Set("type", account.Type)
    d.Set("collect", account.Collect)

    switch account.Type {
    case "AWS":
        if account.AWSUsername != "" {
            d.Set("username", account.AWSUsername)
            d.Set("secret", account.Secret)
        } else {
            roleArns := []string{}
            accountIds := []string{}
            accountNames := []string{}
            for _, assumeRoleInfo := range account.AssumeRoleInfos {
                roleArns = append(roleArns, assumeRoleInfo.RoleArn)
                accountIds = append(accountIds, assumeRoleInfo.AccountId)
                accountNames = append(accountNames, assumeRoleInfo.AccountName)
            }
            d.Set("rolearn", roleArns)
            d.Set("account_id", accountIds)
            d.Set("account_name", accountNames)
            d.Set("external_id", account.AssumeRoleInfos[0].ExternalId)
        }
        d.Set("proxy_server_id", account.ProxyServerId)
        regions := []string{}
        for region := range account.Regions {
            regions = append(regions, region)
        }
        d.Set("regions", regions)
    case "AZURE":
        subscriptionIds := []string{}
        for _, subscription := range account.Subscriptions {
            subscriptionIds = append(subscriptionIds, subscription.SubscriptionId)
        }
        d.Set("subscription_id", subscriptionIds)
        d.Set("client_id", account.Subscriptions[0].ClientId)
        d.Set("tenant", account.Subscriptions[0].Tenant)
        d.Set("secret", account.Subscriptions[0].Secret)
        d.Set("environment", account.Subscriptions[0].Environment)
        d.Set("proxy_server_id", account.ProxyServerId)
    case "GCP":
        d.Set("client_id", account.ClientID)
        d.Set("client_email", account.ClientEmail)
        d.Set("private_key_id", account.PrivateKeyID)
        d.Set("private_key", account.PrivateKey)
        d.Set("proxy_server_id", account.ProxyServerId)
        regions := []string{}
        for region := range account.Regions {
            regions = append(regions, region)
        }
        d.Set("regions", regions)
    }

    return nil
}


func deleteCloudAccount(client *forwardnetworks.ForwardNetworksClient, networkId string, accountName string) error {
    return client.DeleteCloudAccount(networkId, accountName)
}

func resourceCloudAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*forwardnetworks.ForwardNetworksClient)

    networkId := d.Get("network_id").(string)
    accountName := d.Get("name").(string)
    accountType := d.Get("type").(string)

    if accountType == "AZURE" {
        err := deleteCloudAccount(client, networkId, accountName)
        if err != nil {
            return diag.FromErr(err)
        }

        account := createAzureCloudAccount(d)
        err = client.CreateCloudAccount(networkId, account)
        if err != nil {
            return diag.FromErr(err)
        }

        return resourceCloudAccountRead(ctx, d, meta)
    } else if accountType == "AWS" {
        username := d.Get("username").(string)
        var account *forwardnetworks.CloudAccount
        if username != "" {
            account = createAWSIAMUserAccount(d)
        } else {
            account = createAWSAssumeRoleAccount(d)
        }

        err := client.UpdateCloudAccount(networkId, accountName, *account)
        if err != nil {
            return diag.FromErr(err)
        }

        return resourceCloudAccountRead(ctx, d, meta)
    } else if accountType == "GCP" {
        err := deleteCloudAccount(client, networkId, accountName)
        if err != nil {
            return diag.FromErr(err)
        }

        account := createGCPCloudAccount(d)
        err = client.CreateCloudAccount(networkId, account)
        if err != nil {
            return diag.FromErr(err)
        }

        return resourceCloudAccountRead(ctx, d, meta)
    }

    return diag.Errorf("Invalid account type: %s", accountType)
}


func resourceCloudAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardnetworks.ForwardNetworksClient)

	networkId := d.Get("network_id").(string)
	accountName := d.Get("name").(string)

	err := client.DeleteCloudAccount(networkId, accountName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
