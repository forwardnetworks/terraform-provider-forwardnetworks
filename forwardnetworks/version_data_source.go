package forwardnetworks

import (
	"context"
	"github.com/forwardnetworks/forwardnetworks-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &versionDataSource{}
	_ datasource.DataSourceWithConfigure = &versionDataSource{}
)

// NewVersionDataSource is a helper function to simplify the provider implementation.
func NewVersionDataSource() datasource.DataSource {
	return &versionDataSource{}
}

// versionDataSource is the data source implementation.
type versionDataSource struct {
	client *forwardnetworks.Client
}

// versionDataSourceModel maps the data source schema data.
type versionDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Version      types.String `tfsdk:"version"`
}

// Metadata returns the data source type name.
func (d *versionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_version"
}

// versionDataSourceModel maps the data source schema data.


// Schema defines the schema for the data source.
func (d *versionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the Forward Networks API version.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute.",
                Computed: true,
            },
            "version": schema.StringAttribute{
				Description: "The version of the Forward Networks API.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *versionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*forwardnetworks.Client)
}

// Read refreshes the Terraform state with the latest data.
func (d *versionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state versionDataSourceModel

	version, err := d.client.GetVersion()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Forward Networks Version",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(version.Version)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
