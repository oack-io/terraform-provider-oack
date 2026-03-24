package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	oack "github.com/oack-io/oack-go"
	"github.com/oack-io/terraform-provider-oack/internal/providerdata"
)

var (
	_ resource.Resource                = &PagerDutyIntegrationResource{}
	_ resource.ResourceWithImportState = &PagerDutyIntegrationResource{}
)

type PagerDutyIntegrationResource struct {
	data *providerdata.Data
}

type PagerDutyIntegrationResourceModel struct {
	ID           types.String `tfsdk:"id"`
	APIKey       types.String `tfsdk:"api_key"`
	Region       types.String `tfsdk:"region"`
	ServiceIDs   types.List   `tfsdk:"service_ids"`
	SyncEnabled  types.Bool   `tfsdk:"sync_enabled"`
	SyncError    types.String `tfsdk:"sync_error"`
	LastSyncedAt types.String `tfsdk:"last_synced_at"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

func NewPagerDutyIntegrationResource() resource.Resource {
	return &PagerDutyIntegrationResource{}
}

func (r *PagerDutyIntegrationResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_pagerduty_integration"
}

func (r *PagerDutyIntegrationResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages the Oack PagerDuty integration (singleton per account).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Integration UUID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"api_key": schema.StringAttribute{
				Description: "PagerDuty API key.",
				Required:    true,
				Sensitive:   true,
			},
			"region": schema.StringAttribute{
				Description: "PagerDuty region: us or eu.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("us", "eu"),
				},
			},
			"service_ids": schema.ListAttribute{
				Description: "List of PagerDuty service IDs to sync.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"sync_enabled": schema.BoolAttribute{
				Description: "Whether automatic sync is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"sync_error": schema.StringAttribute{
				Description: "Last sync error message, if any.",
				Computed:    true,
			},
			"last_synced_at": schema.StringAttribute{
				Description: "Last successful sync timestamp (RFC 3339).",
				Computed:    true,
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

func (r *PagerDutyIntegrationResource) Configure(
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

func (r *PagerDutyIntegrationResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan PagerDutyIntegrationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var serviceIDs []string
	resp.Diagnostics.Append(plan.ServiceIDs.ElementsAs(ctx, &serviceIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd, err := r.data.Client.CreatePDIntegration(ctx, r.data.AccountID, &oack.CreatePDIntegrationParams{
		APIKey:     plan.APIKey.ValueString(),
		Region:     plan.Region.ValueString(),
		ServiceIDs: serviceIDs,
	})
	if err != nil {
		resp.Diagnostics.AddError("Create PagerDuty Integration Failed", err.Error())
		return
	}

	pdToState(ctx, pd, &plan, &resp.Diagnostics)
	// Preserve the API key from the plan since the API may not return it.
	plan.APIKey = types.StringValue(plan.APIKey.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PagerDutyIntegrationResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state PagerDutyIntegrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd, err := r.data.Client.GetPDIntegration(ctx, r.data.AccountID)
	if err != nil {
		if oack.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read PagerDuty Integration Failed", err.Error())
		return
	}

	// Preserve the sensitive API key from existing state.
	apiKey := state.APIKey
	pdToState(ctx, pd, &state, &resp.Diagnostics)
	state.APIKey = apiKey
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PagerDutyIntegrationResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan PagerDutyIntegrationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var serviceIDs []string
	resp.Diagnostics.Append(plan.ServiceIDs.ElementsAs(ctx, &serviceIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := plan.APIKey.ValueString()
	region := plan.Region.ValueString()
	syncEnabled := plan.SyncEnabled.ValueBool()
	updateReq := &oack.UpdatePDIntegrationParams{
		APIKey:      &apiKey,
		Region:      &region,
		ServiceIDs:  serviceIDs,
		SyncEnabled: &syncEnabled,
	}

	pd, err := r.data.Client.UpdatePDIntegration(ctx, r.data.AccountID, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Update PagerDuty Integration Failed", err.Error())
		return
	}

	pdToState(ctx, pd, &plan, &resp.Diagnostics)
	// Preserve the API key from the plan.
	plan.APIKey = types.StringValue(plan.APIKey.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PagerDutyIntegrationResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state PagerDutyIntegrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.data.Client.DeletePDIntegration(ctx, r.data.AccountID); err != nil {
		resp.Diagnostics.AddError("Delete PagerDuty Integration Failed", err.Error())
	}
}

func (r *PagerDutyIntegrationResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	// Singleton resource — import ID is ignored, we just read the integration.
	pd, err := r.data.Client.GetPDIntegration(ctx, r.data.AccountID)
	if err != nil {
		resp.Diagnostics.AddError("Import PagerDuty Integration Failed", err.Error())
		return
	}

	var state PagerDutyIntegrationResourceModel
	pdToState(ctx, pd, &state, &resp.Diagnostics)
	// API key is not available on import — set to empty to avoid null.
	state.APIKey = types.StringValue("")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func pdToState(
	ctx context.Context,
	pd *oack.PDIntegration,
	state *PagerDutyIntegrationResourceModel,
	diags *diag.Diagnostics,
) {
	state.ID = types.StringValue(pd.ID)
	state.Region = types.StringValue(pd.Region)
	state.SyncEnabled = types.BoolValue(pd.SyncEnabled)
	state.SyncError = types.StringValue(pd.SyncError)
	state.CreatedAt = types.StringValue(pd.CreatedAt)
	state.UpdatedAt = types.StringValue(pd.UpdatedAt)

	if pd.LastSyncedAt != nil {
		state.LastSyncedAt = types.StringValue(*pd.LastSyncedAt)
	} else {
		state.LastSyncedAt = types.StringNull()
	}

	if len(pd.ServiceIDs) > 0 {
		svcList, d := types.ListValueFrom(ctx, types.StringType, pd.ServiceIDs)
		diags.Append(d...)
		state.ServiceIDs = svcList
	} else {
		state.ServiceIDs = types.ListNull(types.StringType)
	}
}
