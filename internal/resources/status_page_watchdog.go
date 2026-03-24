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

	"github.com/oack-io/terraform-provider-oack/internal/client"
)

var (
	_ resource.Resource                = &StatusPageWatchdogResource{}
	_ resource.ResourceWithImportState = &StatusPageWatchdogResource{}
)

type StatusPageWatchdogResource struct {
	client *client.Client
}

type StatusPageWatchdogResourceModel struct {
	ID                types.String `tfsdk:"id"`
	StatusPageID      types.String `tfsdk:"status_page_id"`
	ComponentID       types.String `tfsdk:"component_id"`
	MonitorID         types.String `tfsdk:"monitor_id"`
	Severity          types.String `tfsdk:"severity"`
	AutoCreate        types.Bool   `tfsdk:"auto_create"`
	AutoResolve       types.Bool   `tfsdk:"auto_resolve"`
	NotifySubscribers types.Bool   `tfsdk:"notify_subscribers"`
	TemplateID        types.String `tfsdk:"template_id"`
	CreatedAt         types.String `tfsdk:"created_at"`
}

func NewStatusPageWatchdogResource() resource.Resource {
	return &StatusPageWatchdogResource{}
}

func (r *StatusPageWatchdogResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_status_page_watchdog"
}

func (r *StatusPageWatchdogResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a watchdog linking a monitor to a status page component.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Watchdog UUID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status_page_id": schema.StringAttribute{
				Description: "Status page UUID (required for API routing).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"component_id": schema.StringAttribute{
				Description: "Component UUID.",
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
			"severity": schema.StringAttribute{
				Description: "Incident severity: minor, medium, major, or critical.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auto_create": schema.BoolAttribute{
				Description: "Automatically create incidents when the monitor goes down.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"auto_resolve": schema.BoolAttribute{
				Description: "Automatically resolve incidents when the monitor recovers.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"notify_subscribers": schema.BoolAttribute{
				Description: "Notify status page subscribers on incident changes.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"template_id": schema.StringAttribute{
				Description: "Incident template UUID.",
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Creation timestamp (RFC 3339).",
				Computed:    true,
			},
		},
	}
}

func (r *StatusPageWatchdogResource) Configure(
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

func (r *StatusPageWatchdogResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan StatusPageWatchdogResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := watchdogToRequest(&plan)
	w, err := r.client.CreateWatchdog(ctx,
		plan.StatusPageID.ValueString(), plan.ComponentID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Create Watchdog Failed", err.Error())
		return
	}

	watchdogToState(w, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageWatchdogResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state StatusPageWatchdogResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	w, err := r.client.GetWatchdog(ctx,
		state.StatusPageID.ValueString(),
		state.ComponentID.ValueString(),
		state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Watchdog Failed", err.Error())
		return
	}
	if w == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	watchdogToState(w, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StatusPageWatchdogResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	// Watchdogs do not support update — all mutable fields use RequiresReplace,
	// so Terraform will delete+recreate. This method should never be called,
	// but is required by the resource.Resource interface.
	resp.Diagnostics.AddError("Update Not Supported",
		"Watchdogs cannot be updated in place. Terraform should delete and recreate.")
}

func (r *StatusPageWatchdogResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state StatusPageWatchdogResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteWatchdog(ctx,
		state.StatusPageID.ValueString(),
		state.ComponentID.ValueString(),
		state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Watchdog Failed", err.Error())
	}
}

func (r *StatusPageWatchdogResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	// Format: page_id/component_id/watchdog_id
	parts := strings.SplitN(req.ID, "/", 3)
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Invalid Import ID",
			"Expected format: page_id/component_id/watchdog_id")
		return
	}
	pageID, compID, watchdogID := parts[0], parts[1], parts[2]

	w, err := r.client.GetWatchdog(ctx, pageID, compID, watchdogID)
	if err != nil {
		resp.Diagnostics.AddError("Import Watchdog Failed", err.Error())
		return
	}
	if w == nil {
		resp.Diagnostics.AddError("Watchdog Not Found",
			fmt.Sprintf("Watchdog %s not found for component %s in status page %s",
				watchdogID, compID, pageID))
		return
	}

	var state StatusPageWatchdogResourceModel
	state.StatusPageID = types.StringValue(pageID)
	watchdogToState(w, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func watchdogToRequest(m *StatusPageWatchdogResourceModel) *client.CreateWatchdogRequest {
	r := &client.CreateWatchdogRequest{
		MonitorID: m.MonitorID.ValueString(),
		Severity:  m.Severity.ValueString(),
	}

	if !m.AutoCreate.IsNull() && !m.AutoCreate.IsUnknown() {
		v := m.AutoCreate.ValueBool()
		r.AutoCreate = &v
	}
	if !m.AutoResolve.IsNull() && !m.AutoResolve.IsUnknown() {
		v := m.AutoResolve.ValueBool()
		r.AutoResolve = &v
	}
	if !m.NotifySubscribers.IsNull() && !m.NotifySubscribers.IsUnknown() {
		v := m.NotifySubscribers.ValueBool()
		r.NotifySubscribers = &v
	}
	if !m.TemplateID.IsNull() && !m.TemplateID.IsUnknown() {
		r.TemplateID = m.TemplateID.ValueString()
	}

	return r
}

func watchdogToState(w *client.Watchdog, m *StatusPageWatchdogResourceModel) {
	m.ID = types.StringValue(w.ID)
	m.ComponentID = types.StringValue(w.ComponentID)
	m.MonitorID = types.StringValue(w.MonitorID)
	m.Severity = types.StringValue(w.Severity)
	m.AutoCreate = types.BoolValue(w.AutoCreate)
	m.AutoResolve = types.BoolValue(w.AutoResolve)
	m.NotifySubscribers = types.BoolValue(w.NotifySubscribers)
	m.CreatedAt = types.StringValue(w.CreatedAt)

	if w.TemplateID != "" {
		m.TemplateID = types.StringValue(w.TemplateID)
	} else {
		m.TemplateID = types.StringNull()
	}
}
