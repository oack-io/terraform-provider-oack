package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/oack-io/terraform-provider-oack/internal/client"
)

var (
	_ resource.Resource                = &TeamAPIKeyResource{}
	_ resource.ResourceWithImportState = &TeamAPIKeyResource{}
)

type TeamAPIKeyResource struct {
	client *client.Client
}

type TeamAPIKeyResourceModel struct {
	ID        types.String `tfsdk:"id"`
	TeamID    types.String `tfsdk:"team_id"`
	Name      types.String `tfsdk:"name"`
	ExpiresAt types.String `tfsdk:"expires_at"`
	Key       types.String `tfsdk:"key"`
	KeyPrefix types.String `tfsdk:"key_prefix"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func NewTeamAPIKeyResource() resource.Resource {
	return &TeamAPIKeyResource{}
}

func (r *TeamAPIKeyResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_team_api_key"
}

func (r *TeamAPIKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Oack team API key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "API key UUID.",
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
				Description: "API key display name.",
				Required:    true,
			},
			"expires_at": schema.StringAttribute{
				Description: "Expiration timestamp (RFC 3339). Leave empty for no expiration.",
				Optional:    true,
			},
			"key": schema.StringAttribute{
				Description: "The plaintext API key. Only available at creation time.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key_prefix": schema.StringAttribute{
				Description: "Visible prefix of the API key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Creation timestamp (RFC 3339).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *TeamAPIKeyResource) Configure(
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

func (r *TeamAPIKeyResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan TeamAPIKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateTeamAPIKeyRequest{
		Name: plan.Name.ValueString(),
	}
	if !plan.ExpiresAt.IsNull() && !plan.ExpiresAt.IsUnknown() {
		v := plan.ExpiresAt.ValueString()
		createReq.ExpiresAt = &v
	}

	result, err := r.client.CreateTeamAPIKey(ctx, plan.TeamID.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Create Team API Key Failed", err.Error())
		return
	}

	plan.ID = types.StringValue(result.APIKey.ID)
	plan.Key = types.StringValue(result.Key)
	plan.KeyPrefix = types.StringValue(result.APIKey.KeyPrefix)
	plan.CreatedAt = types.StringValue(result.APIKey.CreatedAt)
	if result.APIKey.ExpiresAt != nil {
		plan.ExpiresAt = types.StringValue(*result.APIKey.ExpiresAt)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *TeamAPIKeyResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state TeamAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey, err := r.client.GetTeamAPIKey(ctx, state.TeamID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Team API Key Failed", err.Error())
		return
	}
	if apiKey == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(apiKey.Name)
	state.KeyPrefix = types.StringValue(apiKey.KeyPrefix)
	state.CreatedAt = types.StringValue(apiKey.CreatedAt)
	if apiKey.ExpiresAt != nil {
		state.ExpiresAt = types.StringValue(*apiKey.ExpiresAt)
	} else {
		state.ExpiresAt = types.StringNull()
	}
	// key is not returned by the API on Read — preserve existing state value.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TeamAPIKeyResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	// API keys cannot be updated — all mutable fields trigger replacement.
	// This method is required by the Resource interface but should not be called.
	resp.Diagnostics.AddError("Update Not Supported",
		"Team API keys cannot be updated in place. Change triggers replacement.")
}

func (r *TeamAPIKeyResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state TeamAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteTeamAPIKey(ctx,
		state.TeamID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Team API Key Failed", err.Error())
	}
}

func (r *TeamAPIKeyResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	// Format: team_id/key_id
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID",
			"Expected format: team_id/key_id")
		return
	}
	teamID, keyID := parts[0], parts[1]

	apiKey, err := r.client.GetTeamAPIKey(ctx, teamID, keyID)
	if err != nil {
		resp.Diagnostics.AddError("Import Team API Key Failed", err.Error())
		return
	}
	if apiKey == nil {
		resp.Diagnostics.AddError("Team API Key Not Found",
			fmt.Sprintf("API key %s not found in team %s", keyID, teamID))
		return
	}

	state := TeamAPIKeyResourceModel{
		ID:        types.StringValue(apiKey.ID),
		TeamID:    types.StringValue(apiKey.TeamID),
		Name:      types.StringValue(apiKey.Name),
		KeyPrefix: types.StringValue(apiKey.KeyPrefix),
		Key:       types.StringNull(), // Plaintext key is not available after import.
		CreatedAt: types.StringValue(apiKey.CreatedAt),
	}
	if apiKey.ExpiresAt != nil {
		state.ExpiresAt = types.StringValue(*apiKey.ExpiresAt)
	} else {
		state.ExpiresAt = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
