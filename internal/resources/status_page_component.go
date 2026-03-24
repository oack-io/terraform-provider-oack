package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/oack-io/terraform-provider-oack/internal/client"
)

var (
	_ resource.Resource                = &StatusPageComponentResource{}
	_ resource.ResourceWithImportState = &StatusPageComponentResource{}
)

type StatusPageComponentResource struct {
	client *client.Client
}

type StatusPageComponentResourceModel struct {
	ID            types.String `tfsdk:"id"`
	StatusPageID  types.String `tfsdk:"status_page_id"`
	GroupID       types.String `tfsdk:"group_id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	DisplayUptime types.Bool   `tfsdk:"display_uptime"`
	Position      types.Int64  `tfsdk:"position"`
	Status        types.String `tfsdk:"status"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

func NewStatusPageComponentResource() resource.Resource {
	return &StatusPageComponentResource{}
}

func (r *StatusPageComponentResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_status_page_component"
}

func (r *StatusPageComponentResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a component within an Oack status page.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Component UUID.",
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
			"group_id": schema.StringAttribute{
				Description: "Component group UUID.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Component display name.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Component description.",
				Optional:    true,
			},
			"display_uptime": schema.BoolAttribute{
				Description: "Whether to display uptime for this component.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"position": schema.Int64Attribute{
				Description: "Display order position.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"status": schema.StringAttribute{
				Description: "Current component status.",
				Computed:    true,
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

func (r *StatusPageComponentResource) Configure(
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

func (r *StatusPageComponentResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan StatusPageComponentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := componentToRequest(&plan)
	comp, err := r.client.CreateComponent(ctx, plan.StatusPageID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Create Component Failed", err.Error())
		return
	}

	componentToState(comp, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageComponentResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state StatusPageComponentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	comp, err := r.client.GetComponent(ctx,
		state.StatusPageID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Component Failed", err.Error())
		return
	}
	if comp == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	componentToState(comp, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StatusPageComponentResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan StatusPageComponentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state StatusPageComponentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := componentToRequest(&plan)
	comp, err := r.client.UpdateComponent(ctx,
		state.StatusPageID.ValueString(), state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Update Component Failed", err.Error())
		return
	}

	componentToState(comp, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageComponentResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state StatusPageComponentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteComponent(ctx,
		state.StatusPageID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Component Failed", err.Error())
	}
}

func (r *StatusPageComponentResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	// Format: page_id/component_id
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID",
			"Expected format: page_id/component_id")
		return
	}
	pageID, compID := parts[0], parts[1]

	comp, err := r.client.GetComponent(ctx, pageID, compID)
	if err != nil {
		resp.Diagnostics.AddError("Import Component Failed", err.Error())
		return
	}
	if comp == nil {
		resp.Diagnostics.AddError("Component Not Found",
			fmt.Sprintf("Component %s not found in status page %s", compID, pageID))
		return
	}

	var state StatusPageComponentResourceModel
	componentToState(comp, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func componentToRequest(m *StatusPageComponentResourceModel) *client.CreateComponentRequest {
	r := &client.CreateComponentRequest{
		Name:     m.Name.ValueString(),
		Position: int(m.Position.ValueInt64()),
	}

	if !m.GroupID.IsNull() && !m.GroupID.IsUnknown() {
		r.GroupID = m.GroupID.ValueString()
	}
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		r.Description = m.Description.ValueString()
	}
	if !m.DisplayUptime.IsNull() && !m.DisplayUptime.IsUnknown() {
		v := m.DisplayUptime.ValueBool()
		r.DisplayUptime = &v
	}

	return r
}

func componentToState(comp *client.Component, m *StatusPageComponentResourceModel) {
	m.ID = types.StringValue(comp.ID)
	m.StatusPageID = types.StringValue(comp.StatusPageID)
	m.Name = types.StringValue(comp.Name)
	m.DisplayUptime = types.BoolValue(comp.DisplayUptime)
	m.Position = types.Int64Value(int64(comp.Position))
	m.Status = types.StringValue(comp.Status)
	m.CreatedAt = types.StringValue(comp.CreatedAt)
	m.UpdatedAt = types.StringValue(comp.UpdatedAt)

	if comp.GroupID != "" {
		m.GroupID = types.StringValue(comp.GroupID)
	} else {
		m.GroupID = types.StringNull()
	}

	if comp.Description != "" {
		m.Description = types.StringValue(comp.Description)
	} else {
		m.Description = types.StringNull()
	}
}
