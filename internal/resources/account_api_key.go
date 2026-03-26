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

	oack "github.com/oack-io/oack-go"
	"github.com/oack-io/terraform-provider-oack/internal/providerdata"
)

var (
	_ resource.Resource                = &AccountAPIKeyResource{}
	_ resource.ResourceWithImportState = &AccountAPIKeyResource{}
)

type AccountAPIKeyResource struct {
	data *providerdata.Data
}

type AccountAPIKeyResourceModel struct {
	ID        types.String `tfsdk:"id"`
	AccountID types.String `tfsdk:"account_id"`
	Name      types.String `tfsdk:"name"`
	ExpiresAt types.String `tfsdk:"expires_at"`
	Key       types.String `tfsdk:"key"`
	KeyPrefix types.String `tfsdk:"key_prefix"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func NewAccountAPIKeyResource() resource.Resource {
	return &AccountAPIKeyResource{}
}

func (r *AccountAPIKeyResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_account_api_key"
}

func (r *AccountAPIKeyResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages an Oack account API key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "API key UUID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_id": schema.StringAttribute{
				Description: "Account UUID.",
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

func (r *AccountAPIKeyResource) Configure(
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

func (r *AccountAPIKeyResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan AccountAPIKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := oack.CreateAccountAPIKeyParams{
		Name: plan.Name.ValueString(),
	}
	if !plan.ExpiresAt.IsNull() && !plan.ExpiresAt.IsUnknown() {
		v := plan.ExpiresAt.ValueString()
		createReq.ExpiresAt = &v
	}

	result, err := r.data.Client.CreateAccountAPIKey(ctx, plan.AccountID.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Create Account API Key Failed", err.Error())
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

func (r *AccountAPIKeyResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state AccountAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey, err := getAccountAPIKey(ctx, r.data.Client,
		state.AccountID.ValueString(), state.ID.ValueString())
	if err != nil {
		if oack.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Account API Key Failed", err.Error())
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
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AccountAPIKeyResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	resp.Diagnostics.AddError("Update Not Supported",
		"Account API keys cannot be updated in place. Change triggers replacement.")
}

func (r *AccountAPIKeyResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state AccountAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.data.Client.DeleteAccountAPIKey(ctx,
		state.AccountID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Account API Key Failed", err.Error())
	}
}

func (r *AccountAPIKeyResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID",
			"Expected format: account_id/key_id")
		return
	}
	accountID, keyID := parts[0], parts[1]

	apiKey, err := getAccountAPIKey(ctx, r.data.Client, accountID, keyID)
	if err != nil {
		resp.Diagnostics.AddError("Import Account API Key Failed", err.Error())
		return
	}

	state := AccountAPIKeyResourceModel{
		ID:        types.StringValue(apiKey.ID),
		AccountID: types.StringValue(apiKey.AccountID),
		Name:      types.StringValue(apiKey.Name),
		KeyPrefix: types.StringValue(apiKey.KeyPrefix),
		Key:       types.StringNull(),
		CreatedAt: types.StringValue(apiKey.CreatedAt),
	}
	if apiKey.ExpiresAt != nil {
		state.ExpiresAt = types.StringValue(*apiKey.ExpiresAt)
	} else {
		state.ExpiresAt = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getAccountAPIKey(
	ctx context.Context, c *oack.Client, accountID, keyID string,
) (*oack.AccountAPIKey, error) {
	keys, err := c.ListAccountAPIKeys(ctx, accountID)
	if err != nil {
		return nil, err
	}
	for _, k := range keys {
		if k.ID == keyID {
			return &k, nil
		}
	}
	return nil, &oack.APIError{StatusCode: 404, Message: "account API key not found"}
}
