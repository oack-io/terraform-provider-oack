package resources

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/oack-io/terraform-provider-oack/internal/client"
)

var (
	regexpSlug     = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,98}[a-z0-9])?$`)
	regexpHexColor = regexp.MustCompile(`^#([0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
)

var (
	_ resource.Resource                = &StatusPageResource{}
	_ resource.ResourceWithImportState = &StatusPageResource{}
)

type StatusPageResource struct {
	client *client.Client
}

type StatusPageResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Slug                 types.String `tfsdk:"slug"`
	Description          types.String `tfsdk:"description"`
	CustomDomain         types.String `tfsdk:"custom_domain"`
	Password             types.String `tfsdk:"password"`
	HasPassword          types.Bool   `tfsdk:"has_password"`
	AllowIframe          types.Bool   `tfsdk:"allow_iframe"`
	ShowHistoricalUptime types.Bool   `tfsdk:"show_historical_uptime"`
	BrandingLogoURL      types.String `tfsdk:"branding_logo_url"`
	BrandingFaviconURL   types.String `tfsdk:"branding_favicon_url"`
	BrandingPrimaryColor types.String `tfsdk:"branding_primary_color"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

func NewStatusPageResource() resource.Resource {
	return &StatusPageResource{}
}

func (r *StatusPageResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_status_page"
}

func (r *StatusPageResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages an Oack status page.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Status page UUID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Status page display name.",
				Required:    true,
			},
			"slug": schema.StringAttribute{
				Description: "URL slug. Lowercase letters, digits, and hyphens (1-100 chars).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexpSlug,
						"must contain only lowercase letters, digits, and hyphens (1-100 chars), start and end with alphanumeric",
					),
				},
			},
			"description": schema.StringAttribute{
				Description: "Status page description.",
				Optional:    true,
			},
			"custom_domain": schema.StringAttribute{
				Description: "Custom domain for the status page.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password to protect the status page (write-only).",
				Optional:    true,
				Sensitive:   true,
			},
			"has_password": schema.BoolAttribute{
				Description: "Whether the status page is password-protected.",
				Computed:    true,
			},
			"allow_iframe": schema.BoolAttribute{
				Description: "Whether to allow embedding in iframes.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"show_historical_uptime": schema.BoolAttribute{
				Description: "Whether to show historical uptime data.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"branding_logo_url": schema.StringAttribute{
				Description: "URL for a custom logo.",
				Optional:    true,
			},
			"branding_favicon_url": schema.StringAttribute{
				Description: "URL for a custom favicon.",
				Optional:    true,
			},
			"branding_primary_color": schema.StringAttribute{
				Description: "Primary color hex code for branding (e.g. #FF5733).",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexpHexColor,
						"must be a valid hex color code (e.g. #RRGGBB or #RRGGBBAA)",
					),
				},
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

func (r *StatusPageResource) Configure(
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

func (r *StatusPageResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan StatusPageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := statusPageToRequest(&plan)
	sp, err := r.client.CreateStatusPage(ctx, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Create Status Page Failed", err.Error())
		return
	}

	statusPageToState(sp, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state StatusPageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sp, err := r.client.GetStatusPage(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Status Page Failed", err.Error())
		return
	}
	if sp == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	statusPageToState(sp, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StatusPageResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan StatusPageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state StatusPageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := statusPageToRequest(&plan)
	sp, err := r.client.UpdateStatusPage(ctx, state.ID.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Update Status Page Failed", err.Error())
		return
	}

	statusPageToState(sp, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state StatusPageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteStatusPage(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Status Page Failed", err.Error())
	}
}

func (r *StatusPageResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	sp, err := r.client.GetStatusPage(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Import Status Page Failed", err.Error())
		return
	}
	if sp == nil {
		resp.Diagnostics.AddError("Status Page Not Found",
			fmt.Sprintf("Status page %s not found", req.ID))
		return
	}

	var state StatusPageResourceModel
	statusPageToState(sp, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func statusPageToRequest(m *StatusPageResourceModel) *client.CreateStatusPageRequest {
	r := &client.CreateStatusPageRequest{
		Name: m.Name.ValueString(),
		Slug: m.Slug.ValueString(),
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		r.Description = m.Description.ValueString()
	}
	if !m.CustomDomain.IsNull() && !m.CustomDomain.IsUnknown() {
		v := m.CustomDomain.ValueString()
		r.CustomDomain = &v
	}
	if !m.Password.IsNull() && !m.Password.IsUnknown() {
		v := m.Password.ValueString()
		r.Password = &v
	}
	if !m.AllowIframe.IsNull() && !m.AllowIframe.IsUnknown() {
		v := m.AllowIframe.ValueBool()
		r.AllowIframe = &v
	}
	if !m.ShowHistoricalUptime.IsNull() && !m.ShowHistoricalUptime.IsUnknown() {
		v := m.ShowHistoricalUptime.ValueBool()
		r.ShowHistoricalUptime = &v
	}
	if !m.BrandingLogoURL.IsNull() && !m.BrandingLogoURL.IsUnknown() {
		v := m.BrandingLogoURL.ValueString()
		r.BrandingLogoURL = &v
	}
	if !m.BrandingFaviconURL.IsNull() && !m.BrandingFaviconURL.IsUnknown() {
		v := m.BrandingFaviconURL.ValueString()
		r.BrandingFaviconURL = &v
	}
	if !m.BrandingPrimaryColor.IsNull() && !m.BrandingPrimaryColor.IsUnknown() {
		v := m.BrandingPrimaryColor.ValueString()
		r.BrandingPrimaryColor = &v
	}

	return r
}

func statusPageToState(sp *client.StatusPage, m *StatusPageResourceModel) {
	m.ID = types.StringValue(sp.ID)
	m.Name = types.StringValue(sp.Name)
	m.Slug = types.StringValue(sp.Slug)
	m.HasPassword = types.BoolValue(sp.HasPassword)
	m.AllowIframe = types.BoolValue(sp.AllowIframe)
	m.ShowHistoricalUptime = types.BoolValue(sp.ShowHistoricalUptime)
	m.CreatedAt = types.StringValue(sp.CreatedAt)
	m.UpdatedAt = types.StringValue(sp.UpdatedAt)

	if sp.Description != "" {
		m.Description = types.StringValue(sp.Description)
	} else {
		m.Description = types.StringNull()
	}

	if sp.CustomDomain != nil {
		m.CustomDomain = types.StringValue(*sp.CustomDomain)
	} else {
		m.CustomDomain = types.StringNull()
	}

	if sp.BrandingLogoURL != nil {
		m.BrandingLogoURL = types.StringValue(*sp.BrandingLogoURL)
	} else {
		m.BrandingLogoURL = types.StringNull()
	}

	if sp.BrandingFaviconURL != nil {
		m.BrandingFaviconURL = types.StringValue(*sp.BrandingFaviconURL)
	} else {
		m.BrandingFaviconURL = types.StringNull()
	}

	if sp.BrandingPrimaryColor != nil {
		m.BrandingPrimaryColor = types.StringValue(*sp.BrandingPrimaryColor)
	} else {
		m.BrandingPrimaryColor = types.StringNull()
	}
}
