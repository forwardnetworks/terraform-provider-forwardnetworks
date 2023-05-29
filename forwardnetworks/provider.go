package forwardnetworks

import (
	"context"
	"os"

	"github.com/forwardnetworks/forwardnetworks-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/provider/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &forwardnetworksProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &forwardnetworksProvider{}
}

// forwardnetworksProvider is the provider implementation.
type forwardnetworksProvider struct{}

// Metadata returns the provider type name.
func (p *forwardnetworksProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "forwardnetworks"
}

// Schema defines the provider-level schema for configuration data.
func (p *forwardnetworksProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
		Description: "Interact with Forward Networks API.",
        Attributes: map[string]schema.Attribute{
            "host": schema.StringAttribute{
				Description: "URI for Forward Networks API. May also be provided via FORWARDNETWORKS_HOST environment variable. Defaults to https://fwd.app",
                Optional: true,
            },
            "username": schema.StringAttribute{
				Description: "Username for Forward Networks API. May also be provided via FORWARDNETWORKS_USERNAME environment variable.",
                Optional: true,
            },
            "password": schema.StringAttribute{
				Description: "Password for Forward Networks API. May also be provided via FORWARDNETWORKS_PASSWORD environment variable.",
                Optional:  true,
                Sensitive: true,
            },
			"insecure" : schema.BoolAttribute{
				Description: "Allow for connections to Forward Networks on prem instances without SSL verification.  Defaults to FALSE.",
				Optional:	true,
			},
        },
    }
}

// forwardnetworksProviderModel maps provider schema data to a Go type.
type forwardnetworksProviderModel struct {
    Host     types.String `tfsdk:"host"`
    Username types.String `tfsdk:"username"`
    Password types.String `tfsdk:"password"`
	Insecure types.Bool   `tfsdk:"insecure"`
}



func (p *forwardnetworksProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    // Retrieve provider data from configuration
	tflog.Info(ctx, "Configuring Forward Networks client")
    var config forwardnetworksProviderModel
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // If practitioner provided a configuration value for any of the
    // attributes, it must be a known value.

    if config.Host.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("host"),
            "Unknown Forward Networks API Host",
            "The provider cannot create the Forward Networks API client as there is an unknown configuration value for the Forward Networks API host. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the FORWARDNETWORKS_HOST environment variable.",
        )
    }

    if config.Username.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("username"),
            "Unknown Forward Networks API Username",
            "The provider cannot create the Forward Networks API client as there is an unknown configuration value for the Forward Networks API username. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the FORWARDNETWORKS_USERNAME environment variable.",
        )
    }

    if config.Password.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("password"),
            "Unknown Forward Networks API Password",
            "The provider cannot create the Forward Networks API client as there is an unknown configuration value for the Forward Networks API password. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the FORWARDNETWORKS_PASSWORD environment variable.",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }

    // Default values to environment variables, but override
    // with Terraform configuration value if set.

    host := os.Getenv("FORWARDNETWORKS_HOST")
    username := os.Getenv("FORWARDNETWORKS_USERNAME")
    password := os.Getenv("FORWARDNETWORKS_PASSWORD")
	insecure := false

    if !config.Host.IsNull() {
        host = config.Host.ValueString()
	} else if host == "" {
        host = "https://fwd.app" // Default host
    }

    if !config.Username.IsNull() {
        username = config.Username.ValueString()
    }

    if !config.Password.IsNull() {
        password = config.Password.ValueString()
    }
	
	if !config.Insecure.IsNull() {
        insecure = config.Insecure.ValueBool()
    }

    // If any of the expected configurations are missing, return
    // errors with provider-specific guidance.

    if host == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("host"),
            "Missing Forward Networks API Host",
            "The provider cannot create the Forward Networks API client as there is a missing or empty value for the Forward Networks API host. "+
                "Set the host value in the configuration or use the FWDNET_HOST environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if username == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("username"),
            "Missing Forward Networks API Username",
            "The provider cannot create the Forward Networks API client as there is a missing or empty value for the Forward Networks API username. "+
                "Set the username value in the configuration or use the FWDNET_USERNAME environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if password == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("password"),
            "Missing Forward Networks API Password",
            "The provider cannot create the Forward Networks API client as there is a missing or empty value for the Forward Networks API password. "+
                "Set the password value in the configuration or use the FWDNET_PASSWORD environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }
	
	ctx = tflog.SetField(ctx, "forwardnetworks_host", host)
    ctx = tflog.SetField(ctx, "forwardnetworks_username", username)
    ctx = tflog.SetField(ctx, "forwardnetworks_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "forwardnetworks_password")

    tflog.Debug(ctx, "Creating Forward Networks client")
	
    // Create a new Forward Networks client using the configuration values
    client, err := forwardnetworks.NewClient(&host, &username, &password, insecure)
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Create Forward Networks API Client",
            "An unexpected error occurred when creating the Forward Networks API client. "+
                "If the error is not clear, please contact the provider developers.\n\n"+
                "Forward Networks Client Error: "+err.Error(),
        )
        return
    }

    // Make the Forward Networks client available during DataSource and Resource
    // type Configure methods.
    resp.DataSourceData = client
    resp.ResourceData = client
	
	tflog.Info(ctx, "Configured Forward Networks client", map[string]any{"success": true})
}


// DataSources defines the data sources implemented in the provider.
func (p *forwardnetworksProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource {
		NewVersionDataSource,
        NewExternalIdDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *forwardnetworksProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
    }
}
