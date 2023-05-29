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
	_ datasource.DataSource              = &externalIdDataSource{}
	_ datasource.DataSourceWithConfigure = &externalIdDataSource{}
)

// NewExternalIdDataSource is a helper function to simplify the provider implementation.
func NewExternalIdDataSource() datasource.DataSource {
	return &externalIdDataSource{}
}

// externalIdDataSource is the data source implementation.
type externalIdDataSource struct {
	client *forwardnetworks.Client
}

// externalIdDataSourceModel maps the data source schema data.
type externalIdDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	NetworkID  types.String `tfsdk:"network_id"`
	ExternalID types.String `tfsdk:"external_id"`
}

// Metadata returns the data source type name.
func (d *externalIdDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_external_id"
}

// Schema defines the schema for the data source.
func (d *externalIdDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the Forward Networks external ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute.",
				Computed:    true,
			},
			"network_id": schema.StringAttribute{
				Description: "The network ID used to fetch the external ID.",
				Required:    true,			
			},
			"external_id": schema.StringAttribute{
				Description: "The external ID associated with the network ID.",
				Optional:    true,			
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *externalIdDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*forwardnetworks.Client)
}

// Read refreshes the Terraform state with the latest data.
func (d *externalIdDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state externalIdDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	externalId, err := d.client.GetExternalId(state.NetworkID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable Reading External ID",
			"Could not read Network ID "+state.NetworkID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.ID = types.StringValue(externalId.ExternalId)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
