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
	_ resource.Resource                = &MonitorAlertChannelLinkResource{}
	_ resource.ResourceWithImportState = &MonitorAlertChannelLinkResource{}
)

type MonitorAlertChannelLinkResource struct {
	client *client.Client
}

type MonitorAlertChannelLinkModel struct {
	TeamID    types.String `tfsdk:"team_id"`
	MonitorID types.String `tfsdk:"monitor_id"`
	ChannelID types.String `tfsdk:"channel_id"`
}

func NewMonitorAlertChannelLinkResource() resource.Resource {
	return &MonitorAlertChannelLinkResource{}
}

func (r *MonitorAlertChannelLinkResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_monitor_alert_channel_link"
}

func (r *MonitorAlertChannelLinkResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Links an alert channel to a monitor. " +
			"When the monitor goes down, the linked channel receives notifications.",
		Attributes: map[string]schema.Attribute{
			"team_id": schema.StringAttribute{
				Description: "Team UUID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"monitor_id": schema.StringAttribute{
				Description: "Monitor UUID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"channel_id": schema.StringAttribute{
				Description: "Alert channel UUID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *MonitorAlertChannelLinkResource) Configure(
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

func (r *MonitorAlertChannelLinkResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan MonitorAlertChannelLinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.LinkMonitorChannel(ctx,
		plan.TeamID.ValueString(),
		plan.MonitorID.ValueString(),
		plan.ChannelID.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Link Channel Failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MonitorAlertChannelLinkResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state MonitorAlertChannelLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	channelIDs, err := r.client.ListMonitorChannels(ctx,
		state.TeamID.ValueString(), state.MonitorID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Monitor Channels Failed", err.Error())
		return
	}

	// Check if the link still exists.
	found := false
	for _, id := range channelIDs {
		if id == state.ChannelID.ValueString() {
			found = true
			break
		}
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MonitorAlertChannelLinkResource) Update(
	_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	// All attributes require replace — Update should never be called.
	resp.Diagnostics.AddError("Update Not Supported",
		"All attributes of monitor_alert_channel_link require replacement")
}

func (r *MonitorAlertChannelLinkResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state MonitorAlertChannelLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UnlinkMonitorChannel(ctx,
		state.TeamID.ValueString(),
		state.MonitorID.ValueString(),
		state.ChannelID.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Unlink Channel Failed", err.Error())
	}
}

func (r *MonitorAlertChannelLinkResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	// Format: team_id/monitor_id/channel_id
	parts := strings.SplitN(req.ID, "/", 3)
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Invalid Import ID",
			"Expected format: team_id/monitor_id/channel_id")
		return
	}

	state := MonitorAlertChannelLinkModel{
		TeamID:    types.StringValue(parts[0]),
		MonitorID: types.StringValue(parts[1]),
		ChannelID: types.StringValue(parts[2]),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
