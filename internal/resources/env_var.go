package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	oack "github.com/oack-io/oack-go"
	"github.com/oack-io/terraform-provider-oack/internal/providerdata"
)

var (
	_ resource.Resource                = &EnvVarResource{}
	_ resource.ResourceWithImportState = &EnvVarResource{}
)

type EnvVarResource struct {
	data *providerdata.Data
}

type EnvVarResourceModel struct {
	ID        types.String `tfsdk:"id"`
	TeamID    types.String `tfsdk:"team_id"`
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	IsSecret  types.Bool   `tfsdk:"is_secret"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func NewEnvVarResource() resource.Resource {
	return &EnvVarResource{}
}

func (r *EnvVarResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_env_var"
}

func (r *EnvVarResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Oack team environment variable or secret.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Environment variable UUID.",
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
			"key": schema.StringAttribute{
				Description: "Variable name (unique per team).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				Description: "Variable value.",
				Required:    true,
				Sensitive:   true,
			},
			"is_secret": schema.BoolAttribute{
				Description: "Whether the variable is a secret (encrypted, masked in responses).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
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

func (r *EnvVarResource) Configure(
	_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
) {
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

func (r *EnvVarResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EnvVarResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ev, err := r.data.Client.CreateEnvVar(ctx, plan.TeamID.ValueString(),
		&oack.CreateEnvVarParams{
			Key:      plan.Key.ValueString(),
			Value:    plan.Value.ValueString(),
			IsSecret: plan.IsSecret.ValueBool(),
		})
	if err != nil {
		resp.Diagnostics.AddError("Create Env Var Failed", err.Error())
		return
	}

	envVarToState(ev, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EnvVarResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnvVarResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vars, err := r.data.Client.ListEnvVars(ctx, state.TeamID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Env Var Failed", err.Error())
		return
	}

	var found *oack.EnvVar
	for i := range vars {
		if vars[i].ID == state.ID.ValueString() {
			found = &vars[i]
			break
		}
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	envVarToState(found, &state)
	// Preserve the planned value — secret values are masked in list responses.
	if found.IsSecret {
		// Keep the existing state value since the API returns a masked value.
	} else {
		state.Value = types.StringValue(found.Value)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnvVarResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnvVarResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ev, err := r.data.Client.UpdateEnvVar(ctx,
		plan.TeamID.ValueString(), plan.Key.ValueString(),
		&oack.UpdateEnvVarParams{
			Value:    plan.Value.ValueString(),
			IsSecret: plan.IsSecret.ValueBool(),
		})
	if err != nil {
		resp.Diagnostics.AddError("Update Env Var Failed", err.Error())
		return
	}

	envVarToState(ev, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EnvVarResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnvVarResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.data.Client.DeleteEnvVar(ctx,
		state.TeamID.ValueString(), state.Key.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Env Var Failed", err.Error())
	}
}

func (r *EnvVarResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	// Format: team_id/key
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID",
			"Expected format: team_id/key")
		return
	}
	teamID, key := parts[0], parts[1]

	vars, err := r.data.Client.ListEnvVars(ctx, teamID)
	if err != nil {
		resp.Diagnostics.AddError("Import Env Var Failed", err.Error())
		return
	}

	var found *oack.EnvVar
	for i := range vars {
		if vars[i].Key == key {
			found = &vars[i]
			break
		}
	}
	if found == nil {
		resp.Diagnostics.AddError("Import Env Var Failed",
			fmt.Sprintf("env var with key %q not found in team %s", key, teamID))
		return
	}

	var state EnvVarResourceModel
	envVarToState(found, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func envVarToState(ev *oack.EnvVar, state *EnvVarResourceModel) {
	state.ID = types.StringValue(ev.ID)
	state.TeamID = types.StringValue(ev.TeamID)
	state.Key = types.StringValue(ev.Key)
	state.Value = types.StringValue(ev.Value)
	state.IsSecret = types.BoolValue(ev.IsSecret)
	state.CreatedAt = types.StringValue(ev.CreatedAt)
	state.UpdatedAt = types.StringValue(ev.UpdatedAt)
}
