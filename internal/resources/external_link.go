package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/oack-io/terraform-provider-oack/internal/client"
)

var (
	_ resource.Resource                = &ExternalLinkResource{}
	_ resource.ResourceWithImportState = &ExternalLinkResource{}
)

type ExternalLinkResource struct {
	client *client.Client
}

type ExternalLinkResourceModel struct {
	ID                types.String `tfsdk:"id"`
	TeamID            types.String `tfsdk:"team_id"`
	Name              types.String `tfsdk:"name"`
	URLTemplate       types.String `tfsdk:"url_template"`
	IconURL           types.String `tfsdk:"icon_url"`
	TimeWindowMinutes types.Int64  `tfsdk:"time_window_minutes"`
	CreatedAt         types.String `tfsdk:"created_at"`
	UpdatedAt         types.String `tfsdk:"updated_at"`
}

func NewExternalLinkResource() resource.Resource {
	return &ExternalLinkResource{}
}

func (r *ExternalLinkResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_external_link"
}

func (r *ExternalLinkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Oack external link.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "External link UUID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"team_id": schema.StringAttribute{
				Description: "Team UUID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Link display name (max 255 characters).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"url_template": schema.StringAttribute{
				Description: "URL template for the external link.",
				Required:    true,
			},
			"icon_url": schema.StringAttribute{
				Description: "Icon URL for the external link.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"time_window_minutes": schema.Int64Attribute{
				Description: "Time window in minutes (must be positive).",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Creation timestamp (RFC 3339).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "Last update timestamp (RFC 3339).",
				Computed:    true,
			},
		},
	}
}

func (r *ExternalLinkResource) Configure(
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

func (r *ExternalLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ExternalLinkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.client.CreateExternalLink(ctx, plan.TeamID.ValueString(),
		&client.CreateExternalLinkRequest{
			Name:              plan.Name.ValueString(),
			URLTemplate:       plan.URLTemplate.ValueString(),
			IconURL:           plan.IconURL.ValueString(),
			TimeWindowMinutes: int(plan.TimeWindowMinutes.ValueInt64()),
		})
	if err != nil {
		resp.Diagnostics.AddError("Create External Link Failed", err.Error())
		return
	}

	linkToState(link, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ExternalLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ExternalLinkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.client.GetExternalLink(ctx, state.TeamID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read External Link Failed", err.Error())
		return
	}
	if link == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	linkToState(link, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ExternalLinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ExternalLinkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ExternalLinkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.client.UpdateExternalLink(ctx,
		state.TeamID.ValueString(), state.ID.ValueString(),
		&client.CreateExternalLinkRequest{
			Name:              plan.Name.ValueString(),
			URLTemplate:       plan.URLTemplate.ValueString(),
			IconURL:           plan.IconURL.ValueString(),
			TimeWindowMinutes: int(plan.TimeWindowMinutes.ValueInt64()),
		})
	if err != nil {
		resp.Diagnostics.AddError("Update External Link Failed", err.Error())
		return
	}

	linkToState(link, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ExternalLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ExternalLinkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteExternalLink(ctx,
		state.TeamID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete External Link Failed", err.Error())
	}
}

func (r *ExternalLinkResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	// Format: team_id/link_id
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID",
			"Expected format: team_id/link_id")
		return
	}
	teamID, linkID := parts[0], parts[1]

	link, err := r.client.GetExternalLink(ctx, teamID, linkID)
	if err != nil {
		resp.Diagnostics.AddError("Import External Link Failed", err.Error())
		return
	}
	if link == nil {
		resp.Diagnostics.AddError("External Link Not Found",
			fmt.Sprintf("Link %s not found in team %s", linkID, teamID))
		return
	}

	var state ExternalLinkResourceModel
	linkToState(link, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func linkToState(link *client.ExternalLink, state *ExternalLinkResourceModel) {
	state.ID = types.StringValue(link.ID)
	state.TeamID = types.StringValue(link.TeamID)
	state.Name = types.StringValue(link.Name)
	state.URLTemplate = types.StringValue(link.URLTemplate)
	state.IconURL = types.StringValue(link.IconURL)
	state.TimeWindowMinutes = types.Int64Value(int64(link.TimeWindowMinutes))
	state.CreatedAt = types.StringValue(link.CreatedAt)
	state.UpdatedAt = types.StringValue(link.UpdatedAt)
}
