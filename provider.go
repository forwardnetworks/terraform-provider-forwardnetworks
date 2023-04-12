package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-forwardnetworks/forwardnetworks"
)

// Provider returns the Forward Networks provider schema.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// The username for the Forward Networks API.
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FORWARDNETWORKS_USERNAME", nil),
				Description: "The username for the Forward Networks API.",
			},
			// The password for the Forward Networks API.
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("FORWARDNETWORKS_PASSWORD", nil),
				Description: "The password for the Forward Networks API.",
			},
			// The base URL for the Forward Networks API.
			"apphost": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FORWARDNETWORKS_APPHOST", nil),
				Description: "The base URL for the Forward Networks API.",
				Default:     "fwd.app",
			},
			// Set this to true to allow insecure connections to the Forward Networks API.
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set this to true to allow insecure connections to the Forward Networks API.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			// Resource to set the cloud proxy server
			"forwardnetworks_proxy": resourceProxy(),
			// Resource to create annd destroy cloud accounts
			"forwardnetworks_cloud": resourceCloudAccount(),
			// Resource to create and destroy networks
			"forwardnetworks_network": resourceNetworks(),
			// Forward Networks collector resource.
			"forwardnetworks_collector": resourceCollectors(),
			// Forward Networks collection schedule resource.
			"forwardnetworks_collection_schedule": resourceCollectionSchedule(),
			// Forward Networks collection resource.
			"forwardnetworks_collection": resourceCollection(),
			// Forward Networks check resource.
			"forwardnetworks_check": resourceCheck(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			// Forward Networks version data source.
			"forwardnetworks_version": dataSourceVersion(),
			// Forward Networks external ID data source.
			"forwardnetworks_externalid": dataSourceExternalId(),
			// Forward Networks cloud accounts data source.
			"forwardnetworks_cloud": dataSourceCloudAccounts(),
			// Forward Networks proxy data source.
			"forwardnetworks_proxy": dataSourceProxy(),
			// Forward Networks collection data source.
			"forwardnetworks_collection": dataCollection(),
			// Forward Networks snapshot data source.
			"forwardnetworks_snapshot": dataSourceSnapshots(),
			// Forward Networks NQE query data source.
			"forwardnetworks_nqe": dataSourceForwardNetworksNQEQuery(),
			// Forward Networks NQE query execution data source.
			"forwardnetworks_nqe_execute": dataSourceForwardNetworksNQEQueryExecution(),
			// Forward Networks checks data source.
			"forwardnetworks_checks": dataSourceChecks(),
			// Forward Networks AWS policy data source.
			"forwardnetworks_aws_policy": dataSourceAWSPolicy(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure configures the Forward Networks provider.
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	baseURL := d.Get("apphost").(string)
	insecure := d.Get("insecure").(bool)

	client := forwardnetworks.NewForwardNetworksClient("https://"+baseURL, username, password, insecure)

	var diags diag.Diagnostics

	if insecure {
		warning := "You have enabled insecure mode. This may expose your connection to security risks. Use with caution."
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Insecure mode enabled",
			Detail:   warning,
		})
	}

	return client, diags
}
