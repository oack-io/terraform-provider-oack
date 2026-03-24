package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/oack-io/terraform-provider-oack/internal/client"
)

var (
	_ resource.Resource                = &StatusPageComponentGroupResource{}
	_ resource.ResourceWithImportState = &StatusPageComponentGroupResource{}
)

type StatusPageComponentGroupResource struct {
	client *client.Client
}

type StatusPageComponentGroupResourceModel struct {
	ID           types.String `tfsdk:"id"`
	StatusPageID types.String `tfsdk:"status_page_id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Position     types.Int64  `tfsdk:"position"`
	Collapsed    types.Bool   `tfsdk:"collapsed"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

func NewStatusPageComponentGroupResource() resource.Resource {
	return &StatusPageComponentGroupResource{}
}

func (r *StatusPageComponentGroupResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_status_page_component_group"
}

func (r *StatusPageComponentGroupResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a component group within an Oack status page.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Component group UUID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status_page_id": schema.StringAttribute{
				Description: "Status page UUID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Group display name.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Group description.",
				Optional:    true,
			},
			"position": schema.Int64Attribute{
				Description: "Display order position.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"collapsed": schema.BoolAttribute{
				Description: "Whether the group is collapsed by default.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"created_at": schema.StringAttribute{
				Description: "Creation timestamp (RFC 3339).",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Last update timestamp (RFC 3339).",
				Computed:    true,
			},
		},
	}
}

func (r *StatusPageComponentGroupResource) Configure(
	_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *StatusPageComponentGroupResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan StatusPageComponentGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := componentGroupToRequest(&plan)
	g, err := r.client.CreateComponentGroup(ctx, plan.StatusPageID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Create Component Group Failed", err.Error())
		return
	}

	componentGroupToState(g, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageComponentGroupResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state StatusPageComponentGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	g, err := r.client.GetComponentGroup(ctx,
		state.StatusPageID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Component Group Failed", err.Error())
		return
	}
	if g == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	componentGroupToState(g, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StatusPageComponentGroupResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan StatusPageComponentGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state StatusPageComponentGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := componentGroupToRequest(&plan)
	g, err := r.client.UpdateComponentGroup(ctx,
		state.StatusPageID.ValueString(), state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Update Component Group Failed", err.Error())
		return
	}

	componentGroupToState(g, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageComponentGroupResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state StatusPageComponentGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteComponentGroup(ctx,
		state.StatusPageID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Component Group Failed", err.Error())
	}
}

func (r *StatusPageComponentGroupResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	// Format: page_id/group_id
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID",
			"Expected format: page_id/group_id")
		return
	}
	pageID, groupID := parts[0], parts[1]

	g, err := r.client.GetComponentGroup(ctx, pageID, groupID)
	if err != nil {
		resp.Diagnostics.AddError("Import Component Group Failed", err.Error())
		return
	}
	if g == nil {
		resp.Diagnostics.AddError("Component Group Not Found",
			fmt.Sprintf("Group %s not found in status page %s", groupID, pageID))
		return
	}

	var state StatusPageComponentGroupResourceModel
	componentGroupToState(g, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func componentGroupToRequest(m *StatusPageComponentGroupResourceModel) *client.CreateComponentGroupRequest {
	r := &client.CreateComponentGroupRequest{
		Name:     m.Name.ValueString(),
		Position: int(m.Position.ValueInt64()),
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		r.Description = m.Description.ValueString()
	}
	if !m.Collapsed.IsNull() && !m.Collapsed.IsUnknown() {
		v := m.Collapsed.ValueBool()
		r.Collapsed = &v
	}

	return r
}

func componentGroupToState(g *client.ComponentGroup, m *StatusPageComponentGroupResourceModel) {
	m.ID = types.StringValue(g.ID)
	m.StatusPageID = types.StringValue(g.StatusPageID)
	m.Name = types.StringValue(g.Name)
	m.Position = types.Int64Value(int64(g.Position))
	m.Collapsed = types.BoolValue(g.Collapsed)
	m.CreatedAt = types.StringValue(g.CreatedAt)
	m.UpdatedAt = types.StringValue(g.UpdatedAt)

	if g.Description != "" {
		m.Description = types.StringValue(g.Description)
	} else {
		m.Description = types.StringNull()
	}
}
