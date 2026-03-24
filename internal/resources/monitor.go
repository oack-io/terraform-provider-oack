package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	oack "github.com/oack-io/oack-go"
	"github.com/oack-io/terraform-provider-oack/internal/providerdata"
)

var (
	_ resource.Resource                = &MonitorResource{}
	_ resource.ResourceWithImportState = &MonitorResource{}
)

type MonitorResource struct {
	data *providerdata.Data
}

type MonitorResourceModel struct {
	ID                      types.String  `tfsdk:"id"`
	TeamID                  types.String  `tfsdk:"team_id"`
	Name                    types.String  `tfsdk:"name"`
	URL                     types.String  `tfsdk:"url"`
	Status                  types.String  `tfsdk:"status"`
	CheckIntervalMs         types.Int64   `tfsdk:"check_interval_ms"`
	TimeoutMs               types.Int64   `tfsdk:"timeout_ms"`
	HTTPMethod              types.String  `tfsdk:"http_method"`
	HTTPVersion             types.String  `tfsdk:"http_version"`
	Headers                 types.Map     `tfsdk:"headers"`
	FollowRedirects         types.Bool    `tfsdk:"follow_redirects"`
	AllowedStatusCodes      types.List    `tfsdk:"allowed_status_codes"`
	FailureThreshold        types.Int64   `tfsdk:"failure_threshold"`
	LatencyThresholdMs      types.Int64   `tfsdk:"latency_threshold_ms"`
	SSLExpiryEnabled        types.Bool    `tfsdk:"ssl_expiry_enabled"`
	SSLExpiryThresholds     types.List    `tfsdk:"ssl_expiry_thresholds"`
	DomainExpiryEnabled     types.Bool    `tfsdk:"domain_expiry_enabled"`
	DomainExpiryThresholds  types.List    `tfsdk:"domain_expiry_thresholds"`
	UptimeThresholdGood     types.Float64 `tfsdk:"uptime_threshold_good"`
	UptimeThresholdDegraded types.Float64 `tfsdk:"uptime_threshold_degraded"`
	UptimeThresholdCritical types.Float64 `tfsdk:"uptime_threshold_critical"`
	CheckerRegion           types.String  `tfsdk:"checker_region"`
	CheckerCountry          types.String  `tfsdk:"checker_country"`
	ResolveOverrideIP       types.String  `tfsdk:"resolve_override_ip"`
	HealthStatus            types.String  `tfsdk:"health_status"`
	CreatedAt               types.String  `tfsdk:"created_at"`
	UpdatedAt               types.String  `tfsdk:"updated_at"`
}

func NewMonitorResource() resource.Resource {
	return &MonitorResource{}
}

func (r *MonitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (r *MonitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Oack uptime monitor.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Monitor UUID.",
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
				Description: "Monitor display name.",
				Required:    true,
			},
			"url": schema.StringAttribute{
				Description: "URL to monitor.",
				Required:    true,
			},
			"status": schema.StringAttribute{
				Description: "Monitor status: active or paused.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("active"),
				Validators: []validator.String{
					stringvalidator.OneOf("active", "paused"),
				},
			},
			"check_interval_ms": schema.Int64Attribute{
				Description: "Check interval in milliseconds (min 30000).",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(60000),
				Validators: []validator.Int64{
					int64validator.AtLeast(30000),
				},
			},
			"timeout_ms": schema.Int64Attribute{
				Description: "Request timeout in milliseconds.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(10000),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"http_method": schema.StringAttribute{
				Description: "HTTP method: GET, POST, PUT, PATCH, DELETE, HEAD.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("GET"),
				Validators: []validator.String{
					stringvalidator.OneOf("GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"),
				},
			},
			"http_version": schema.StringAttribute{
				Description: "HTTP version: empty (auto), 1.1, or 2.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf("", "1.1", "2"),
				},
			},
			"headers": schema.MapAttribute{
				Description: "Custom request headers.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"follow_redirects": schema.BoolAttribute{
				Description: "Follow HTTP redirects.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"allowed_status_codes": schema.ListAttribute{
				Description: "Allowed status codes (e.g. 2xx, 200, 404).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"failure_threshold": schema.Int64Attribute{
				Description: "Consecutive failures before marking down.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(3),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"latency_threshold_ms": schema.Int64Attribute{
				Description: "Latency threshold in ms (0 = disabled).",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"ssl_expiry_enabled": schema.BoolAttribute{
				Description: "Monitor SSL certificate expiry.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"ssl_expiry_thresholds": schema.ListAttribute{
				Description: "Days before SSL expiry to alert.",
				Optional:    true,
				Computed:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_expiry_enabled": schema.BoolAttribute{
				Description: "Monitor domain expiry.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"domain_expiry_thresholds": schema.ListAttribute{
				Description: "Days before domain expiry to alert.",
				Optional:    true,
				Computed:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"uptime_threshold_good": schema.Float64Attribute{
				Description: "Uptime percentage for good status (0-100).",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(99.9),
				Validators: []validator.Float64{
					float64validator.Between(0, 100),
				},
			},
			"uptime_threshold_degraded": schema.Float64Attribute{
				Description: "Uptime percentage for degraded status (0-100).",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(99.0),
				Validators: []validator.Float64{
					float64validator.Between(0, 100),
				},
			},
			"uptime_threshold_critical": schema.Float64Attribute{
				Description: "Uptime percentage for critical status (0-100).",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(95.0),
				Validators: []validator.Float64{
					float64validator.Between(0, 100),
				},
			},
			"checker_region": schema.StringAttribute{
				Description: "Preferred checker region.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"checker_country": schema.StringAttribute{
				Description: "Preferred checker country.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"resolve_override_ip": schema.StringAttribute{
				Description: "Override DNS resolution (IPv4/IPv6).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"health_status": schema.StringAttribute{
				Description: "Current health status (up/down). Read-only.",
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
			"updated_at": schema.StringAttribute{
				Description: "Last update timestamp (RFC 3339).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *MonitorResource) Configure(
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

func (r *MonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MonitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := planToCreateRequest(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := r.data.Client.CreateMonitor(ctx, plan.TeamID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Create Monitor Failed", err.Error())
		return
	}

	monitorToState(ctx, monitor, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := r.data.Client.GetMonitor(ctx, state.TeamID.ValueString(), state.ID.ValueString())
	if err != nil {
		if oack.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Monitor Failed", err.Error())
		return
	}

	monitorToState(ctx, monitor, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MonitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state MonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := planToCreateRequest(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := r.data.Client.UpdateMonitor(ctx, state.TeamID.ValueString(), state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Update Monitor Failed", err.Error())
		return
	}

	monitorToState(ctx, monitor, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.data.Client.DeleteMonitor(ctx, state.TeamID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Monitor Failed", err.Error())
		return
	}
}

func (r *MonitorResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	// Import ID format: team_id/monitor_id
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID",
			"Expected format: team_id/monitor_id")
		return
	}
	teamID, monitorID := parts[0], parts[1]

	monitor, err := r.data.Client.GetMonitor(ctx, teamID, monitorID)
	if err != nil {
		resp.Diagnostics.AddError("Import Monitor Failed", err.Error())
		return
	}

	var state MonitorResourceModel
	monitorToState(ctx, monitor, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// ── helpers ──────────────────────────────────────────────────────────────────

func planToCreateRequest(
	ctx context.Context, plan *MonitorResourceModel, diags *diag.Diagnostics,
) *oack.CreateMonitorParams {
	req := &oack.CreateMonitorParams{
		Name:              plan.Name.ValueString(),
		URL:               plan.URL.ValueString(),
		CheckIntervalMs:   plan.CheckIntervalMs.ValueInt64(),
		TimeoutMs:         plan.TimeoutMs.ValueInt64(),
		HTTPMethod:        plan.HTTPMethod.ValueString(),
		HTTPVersion:       plan.HTTPVersion.ValueString(),
		Status:            plan.Status.ValueString(),
		CheckerRegion:     plan.CheckerRegion.ValueString(),
		CheckerCountry:    plan.CheckerCountry.ValueString(),
		ResolveOverrideIP: plan.ResolveOverrideIP.ValueString(),
	}

	fr := plan.FollowRedirects.ValueBool()
	req.FollowRedirects = &fr

	ft := int(plan.FailureThreshold.ValueInt64())
	req.FailureThreshold = ft

	lt := int(plan.LatencyThresholdMs.ValueInt64())
	req.LatencyThresholdMs = lt

	ssl := plan.SSLExpiryEnabled.ValueBool()
	req.SSLExpiryEnabled = &ssl

	dom := plan.DomainExpiryEnabled.ValueBool()
	req.DomainExpiryEnabled = &dom

	utg := plan.UptimeThresholdGood.ValueFloat64()
	req.UptimeThresholdGood = &utg

	utd := plan.UptimeThresholdDegraded.ValueFloat64()
	req.UptimeThresholdDegraded = &utd

	utc := plan.UptimeThresholdCritical.ValueFloat64()
	req.UptimeThresholdCritical = &utc

	if !plan.Headers.IsNull() && !plan.Headers.IsUnknown() {
		headers := make(map[string]string)
		diags.Append(plan.Headers.ElementsAs(ctx, &headers, false)...)
		req.Headers = headers
	}

	if !plan.AllowedStatusCodes.IsNull() && !plan.AllowedStatusCodes.IsUnknown() {
		var codes []string
		diags.Append(plan.AllowedStatusCodes.ElementsAs(ctx, &codes, false)...)
		req.AllowedStatusCodes = codes
	}

	if !plan.SSLExpiryThresholds.IsNull() && !plan.SSLExpiryThresholds.IsUnknown() {
		var thresholds []int
		var int64s []int64
		diags.Append(plan.SSLExpiryThresholds.ElementsAs(ctx, &int64s, false)...)
		for _, v := range int64s {
			thresholds = append(thresholds, int(v))
		}
		req.SSLExpiryThresholds = thresholds
	}

	if !plan.DomainExpiryThresholds.IsNull() && !plan.DomainExpiryThresholds.IsUnknown() {
		var int64s []int64
		diags.Append(plan.DomainExpiryThresholds.ElementsAs(ctx, &int64s, false)...)
		thresholds := make([]int, len(int64s))
		for i, v := range int64s {
			thresholds[i] = int(v)
		}
		req.DomainExpiryThresholds = thresholds
	}

	return req
}

func monitorToState(ctx context.Context, m *oack.Monitor, state *MonitorResourceModel, diags *diag.Diagnostics) {
	state.ID = types.StringValue(m.ID)
	state.TeamID = types.StringValue(m.TeamID)
	state.Name = types.StringValue(m.Name)
	state.URL = types.StringValue(m.URL)
	state.Status = types.StringValue(m.Status)
	state.CheckIntervalMs = types.Int64Value(m.CheckIntervalMs)
	state.TimeoutMs = types.Int64Value(m.TimeoutMs)
	state.HTTPMethod = types.StringValue(m.HTTPMethod)
	state.HTTPVersion = types.StringValue(m.HTTPVersion)
	state.FollowRedirects = types.BoolValue(m.FollowRedirects)
	state.FailureThreshold = types.Int64Value(int64(m.FailureThreshold))
	state.LatencyThresholdMs = types.Int64Value(int64(m.LatencyThresholdMs))
	state.SSLExpiryEnabled = types.BoolValue(m.SSLExpiryEnabled)
	state.DomainExpiryEnabled = types.BoolValue(m.DomainExpiryEnabled)
	state.UptimeThresholdGood = types.Float64Value(m.UptimeThresholdGood)
	state.UptimeThresholdDegraded = types.Float64Value(m.UptimeThresholdDegraded)
	state.UptimeThresholdCritical = types.Float64Value(m.UptimeThresholdCritical)
	state.CheckerRegion = types.StringValue(m.CheckerRegion)
	state.CheckerCountry = types.StringValue(m.CheckerCountry)
	state.ResolveOverrideIP = types.StringValue(m.ResolveOverrideIP)
	state.HealthStatus = types.StringValue(m.HealthStatus)
	state.CreatedAt = types.StringValue(m.CreatedAt)
	state.UpdatedAt = types.StringValue(m.UpdatedAt)

	// Always set list/map fields from API response to avoid null→non-null
	// inconsistency when Terraform plan has null but API returns defaults.
	if m.Headers == nil {
		m.Headers = map[string]string{}
	}
	headers, d := types.MapValueFrom(ctx, types.StringType, m.Headers)
	diags.Append(d...)
	state.Headers = headers

	if m.AllowedStatusCodes == nil {
		m.AllowedStatusCodes = []string{}
	}
	codes, d2 := types.ListValueFrom(ctx, types.StringType, m.AllowedStatusCodes)
	diags.Append(d2...)
	state.AllowedStatusCodes = codes

	sslInts := make([]int64, len(m.SSLExpiryThresholds))
	for i, v := range m.SSLExpiryThresholds {
		sslInts[i] = int64(v)
	}
	sslThresholds, d3 := types.ListValueFrom(ctx, types.Int64Type, sslInts)
	diags.Append(d3...)
	state.SSLExpiryThresholds = sslThresholds

	domInts := make([]int64, len(m.DomainExpiryThresholds))
	for i, v := range m.DomainExpiryThresholds {
		domInts[i] = int64(v)
	}
	domThresholds, d4 := types.ListValueFrom(ctx, types.Int64Type, domInts)
	diags.Append(d4...)
	state.DomainExpiryThresholds = domThresholds
}
