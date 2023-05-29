package forwardnetworks

import (
	"context"
	"strconv"
	"time"

	"github.com/forwardnetworks/forwardnetworks-client-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &networkResource{}
	_ resource.ResourceWithConfigure   = &networkResource{}
	_ resource.ResourceWithImportState = &networkResource{}
)

// NewNetworkResource is a helper function to simplify the provider implementation.
func NewNetworkResource() resource.Resource {
	return &networkResource{}
}

// networkResource is the resource implementation.
type networkResource struct {
	client *forwardnetworks.Client
}

// networkResourceModel maps the resource schema data.
type networkResourceModel struct {
	ID        types.String `tfsdk:"id"`
    ParentID  types.String `tfsdk:"parentId"`
    Name      types.String `tfsdk:"name"`
    OrgID     types.String `tfsdk:"orgId"`
    Creator   types.String `tfsdk:"creator"`
    CreatorID types.String `tfsdk:"creatorId"`
    CreatedAt types.Int64  `tfsdk:"createdAt"`
    Note      types.String `tfsdk:"note"`
}

// Metadata returns the data source type name.
func (r *networkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

// Schema defines the schema for the data source.
func (r *networkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an network.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the network.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the network.",
				Required:    true,
			},
			"note": schema.StringAttribute{
				Description: "Note for a network.",
				Optional:    true,
			},
}

// Configure adds the provider configured client to the data source.
func (r *networkResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*forwardnetworks.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *networkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan networkResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var items []forwardnetworks.NetworkItem
	for _, item := range plan.Items {
		items = append(items, forwardnetworks.NetworkItem{
			Coffee: forwardnetworks.Coffee{
				ID: int(item.Coffee.ID.ValueInt64()),
			},
			Quantity: int(item.Quantity.ValueInt64()),
		})
	}

	// Create new network
	network, err := r.client.CreateNetwork(items)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating network",
			"Could not create network, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(strconv.Itoa(network.ID))
	for itemIndex, item := range network.Items {
		plan.Items[itemIndex] = networkItemModel{
			Coffee: networkItemCoffeeModel{
				ID:          types.Int64Value(int64(item.Coffee.ID)),
				Name:        types.StringValue(item.Coffee.Name),
				Teaser:      types.StringValue(item.Coffee.Teaser),
				Description: types.StringValue(item.Coffee.Description),
				Price:       types.Float64Value(item.Coffee.Price),
				Image:       types.StringValue(item.Coffee.Image),
			},
			Quantity: types.Int64Value(int64(item.Quantity)),
		}
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *networkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state networkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed network value from forwardnetworks
	network, err := r.client.GetNetwork(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading forwardnetworks Network",
			"Could not read forwardnetworks network ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.Items = []networkItemModel{}
	for _, item := range network.Items {
		state.Items = append(state.Items, networkItemModel{
			Coffee: networkItemCoffeeModel{
				ID:          types.Int64Value(int64(item.Coffee.ID)),
				Name:        types.StringValue(item.Coffee.Name),
				Teaser:      types.StringValue(item.Coffee.Teaser),
				Description: types.StringValue(item.Coffee.Description),
				Price:       types.Float64Value(item.Coffee.Price),
				Image:       types.StringValue(item.Coffee.Image),
			},
			Quantity: types.Int64Value(int64(item.Quantity)),
		})
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *networkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan networkResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var forwardnetworksItems []forwardnetworks.NetworkItem
	for _, item := range plan.Items {
		forwardnetworksItems = append(forwardnetworksItems, forwardnetworks.NetworkItem{
			Coffee: forwardnetworks.Coffee{
				ID: int(item.Coffee.ID.ValueInt64()),
			},
			Quantity: int(item.Quantity.ValueInt64()),
		})
	}

	// Update existing network
	_, err := r.client.UpdateNetwork(plan.ID.ValueString(), forwardnetworksItems)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating forwardnetworks Network",
			"Could not update network, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetNetwork as UpdateNetwork items are not
	// populated.
	network, err := r.client.GetNetwork(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading forwardnetworks Network",
			"Could not read forwardnetworks network ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update resource state with updated items and timestamp
	plan.Items = []networkItemModel{}
	for _, item := range network.Items {
		plan.Items = append(plan.Items, networkItemModel{
			Coffee: networkItemCoffeeModel{
				ID:          types.Int64Value(int64(item.Coffee.ID)),
				Name:        types.StringValue(item.Coffee.Name),
				Teaser:      types.StringValue(item.Coffee.Teaser),
				Description: types.StringValue(item.Coffee.Description),
				Price:       types.Float64Value(item.Coffee.Price),
				Image:       types.StringValue(item.Coffee.Image),
			},
			Quantity: types.Int64Value(int64(item.Quantity)),
		})
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *networkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state networkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing network
	err := r.client.DeleteNetwork(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting forwardnetworks Network",
			"Could not delete network, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *networkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}