package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/oack-io/terraform-provider-oack/internal/client"
)

var (
	_ resource.Resource                = &AlertChannelResource{}
	_ resource.ResourceWithImportState = &AlertChannelResource{}
)

type AlertChannelResource struct {
	client *client.Client
}

type AlertChannelResourceModel struct {
	ID            types.String `tfsdk:"id"`
	TeamID        types.String `tfsdk:"team_id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	Config        types.Map    `tfsdk:"config"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	EmailVerified types.Bool   `tfsdk:"email_verified"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

func NewAlertChannelResource() resource.Resource {
	return &AlertChannelResource{}
}

func (r *AlertChannelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_channel"
}

func (r *AlertChannelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Oack alert channel (Slack, email, webhook, Telegram, Discord, PagerDuty).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Channel UUID.",
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
				Description: "Channel display name.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Channel type: slack, webhook, email, telegram, discord, pagerduty.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config": schema.MapAttribute{
				Description: "Type-specific config. slack: webhook_url; email: email; webhook: url; telegram: chat_id; discord: webhook_url; pagerduty: routing_key + region.",
				Required:    true,
				Sensitive:   true,
				ElementType: types.StringType,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the channel is active.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"email_verified": schema.BoolAttribute{
				Description: "Whether the email address is verified (email channels only).",
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

func (r *AlertChannelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AlertChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlertChannelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := make(map[string]string)
	resp.Diagnostics.Append(plan.Config.ElementsAs(ctx, &config, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enabled := plan.Enabled.ValueBool()
	ch, err := r.client.CreateAlertChannel(ctx, plan.TeamID.ValueString(),
		&client.CreateAlertChannelRequest{
			Type:    plan.Type.ValueString(),
			Name:    plan.Name.ValueString(),
			Config:  config,
			Enabled: &enabled,
		})
	if err != nil {
		resp.Diagnostics.AddError("Create Alert Channel Failed", err.Error())
		return
	}

	channelToState(ctx, ch, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AlertChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlertChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ch, err := r.client.GetAlertChannel(ctx, state.TeamID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Alert Channel Failed", err.Error())
		return
	}
	if ch == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	channelToState(ctx, ch, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AlertChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AlertChannelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AlertChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := make(map[string]string)
	resp.Diagnostics.Append(plan.Config.ElementsAs(ctx, &config, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enabled := plan.Enabled.ValueBool()
	ch, err := r.client.UpdateAlertChannel(ctx,
		state.TeamID.ValueString(), state.ID.ValueString(),
		&client.CreateAlertChannelRequest{
			Type:    plan.Type.ValueString(),
			Name:    plan.Name.ValueString(),
			Config:  config,
			Enabled: &enabled,
		})
	if err != nil {
		resp.Diagnostics.AddError("Update Alert Channel Failed", err.Error())
		return
	}

	channelToState(ctx, ch, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AlertChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlertChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAlertChannel(ctx,
		state.TeamID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Alert Channel Failed", err.Error())
	}
}

func (r *AlertChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Format: team_id/channel_id
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID",
			"Expected format: team_id/channel_id")
		return
	}
	teamID, channelID := parts[0], parts[1]

	ch, err := r.client.GetAlertChannel(ctx, teamID, channelID)
	if err != nil {
		resp.Diagnostics.AddError("Import Alert Channel Failed", err.Error())
		return
	}
	if ch == nil {
		resp.Diagnostics.AddError("Alert Channel Not Found",
			fmt.Sprintf("Channel %s not found in team %s", channelID, teamID))
		return
	}

	var state AlertChannelResourceModel
	channelToState(ctx, ch, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func channelToState(ctx context.Context, ch *client.AlertChannel, state *AlertChannelResourceModel, diags *diag.Diagnostics) {
	state.ID = types.StringValue(ch.ID)
	state.TeamID = types.StringValue(ch.TeamID)
	state.Name = types.StringValue(ch.Name)
	state.Type = types.StringValue(ch.Type)
	state.Enabled = types.BoolValue(ch.Enabled)
	state.EmailVerified = types.BoolValue(ch.EmailVerified)
	state.CreatedAt = types.StringValue(ch.CreatedAt)
	state.UpdatedAt = types.StringValue(ch.UpdatedAt)

	if ch.Config != nil {
		configMap, d := types.MapValueFrom(ctx, types.StringType, ch.Config)
		diags.Append(d...)
		state.Config = configMap
	} else {
		state.Config = types.MapNull(types.StringType)
	}
}
