package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	oack "github.com/oack-io/oack-go"
	"github.com/oack-io/terraform-provider-oack/internal/providerdata"
)

var (
	_ resource.Resource                = &TeamResource{}
	_ resource.ResourceWithImportState = &TeamResource{}
)

type TeamResource struct {
	data *providerdata.Data
}

type TeamResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func NewTeamResource() resource.Resource {
	return &TeamResource{}
}

func (r *TeamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (r *TeamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Oack team.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Team UUID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Team name.",
				Required:    true,
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

func (r *TeamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*providerdata.Data)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providerdata.Data, got: %T", req.ProviderData))
		return
	}
	r.data = c
}

func (r *TeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, err := r.data.Client.CreateTeam(ctx, r.data.AccountID, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Create Team Failed", err.Error())
		return
	}

	plan.ID = types.StringValue(team.ID)
	plan.CreatedAt = types.StringValue(team.CreatedAt)
	plan.UpdatedAt = types.StringValue(team.UpdatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, err := r.data.Client.GetTeam(ctx, state.ID.ValueString())
	if err != nil {
		if oack.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Team Failed", err.Error())
		return
	}

	state.Name = types.StringValue(team.Name)
	state.CreatedAt = types.StringValue(team.CreatedAt)
	state.UpdatedAt = types.StringValue(team.UpdatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TeamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, err := r.data.Client.UpdateTeam(ctx, state.ID.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Update Team Failed", err.Error())
		return
	}

	plan.ID = state.ID
	plan.CreatedAt = types.StringValue(team.CreatedAt)
	plan.UpdatedAt = types.StringValue(team.UpdatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.data.Client.DeleteTeam(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Team Failed", err.Error())
		return
	}
}

func (r *TeamResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	team, err := r.data.Client.GetTeam(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Import Team Failed", err.Error())
		return
	}

	state := TeamResourceModel{
		ID:        types.StringValue(team.ID),
		Name:      types.StringValue(team.Name),
		CreatedAt: types.StringValue(team.CreatedAt),
		UpdatedAt: types.StringValue(team.UpdatedAt),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
