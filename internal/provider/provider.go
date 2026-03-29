package provider

import (
	"context"
	"os"

	oack "github.com/oack-io/oack-go"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/oack-io/terraform-provider-oack/internal/datasources"
	"github.com/oack-io/terraform-provider-oack/internal/providerdata"
	"github.com/oack-io/terraform-provider-oack/internal/resources"
)

var _ provider.Provider = &OackProvider{}

type OackProvider struct {
	version string
}

type OackProviderModel struct {
	APIKey    types.String `tfsdk:"api_key"`
	AccountID types.String `tfsdk:"account_id"`
	APIURL    types.String `tfsdk:"api_url"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OackProvider{version: version}
	}
}

func (p *OackProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "oack"
	resp.Version = p.version
}

func (p *OackProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for Oack uptime monitoring.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "Account API key (oack_acc_...). Can also be set via OACK_API_KEY env var.",
				Optional:    true,
				Sensitive:   true,
			},
			"account_id": schema.StringAttribute{
				Description: "Account ID. Can also be set via OACK_ACCOUNT_ID env var.",
				Optional:    true,
			},
			"api_url": schema.StringAttribute{
				Description: "API base URL. Defaults to https://api.oack.io. Can also be set via OACK_API_URL env var.",
				Optional:    true,
			},
		},
	}
}

func (p *OackProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config OackProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := envOrValue(config.APIKey, "OACK_API_KEY")
	accountID := envOrValue(config.AccountID, "OACK_ACCOUNT_ID")
	apiURL := envOrValue(config.APIURL, "OACK_API_URL")
	if apiURL == "" {
		apiURL = "https://api.oack.io"
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"Set api_key in the provider block or OACK_API_KEY environment variable.",
		)
		return
	}
	if accountID == "" {
		resp.Diagnostics.AddError(
			"Missing Account ID",
			"Set account_id in the provider block or OACK_ACCOUNT_ID environment variable.",
		)
		return
	}

	c := oack.New(oack.BearerToken(apiKey), oack.WithBaseURL(apiURL))
	data := &providerdata.Data{Client: c, AccountID: accountID}
	resp.DataSourceData = data
	resp.ResourceData = data
}

func (p *OackProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewTeamResource,
		resources.NewMonitorResource,
		resources.NewAlertChannelResource,
		resources.NewMonitorAlertChannelLinkResource,
		resources.NewStatusPageResource,
		resources.NewStatusPageComponentGroupResource,
		resources.NewStatusPageComponentResource,
		resources.NewStatusPageWatchdogResource,
		resources.NewExternalLinkResource,
		resources.NewPagerDutyIntegrationResource,
		resources.NewTeamAPIKeyResource,
		resources.NewAccountAPIKeyResource,
		resources.NewEnvVarResource,
	}
}

func (p *OackProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewCheckersDataSource,
		datasources.NewTeamsDataSource,
	}
}

func envOrValue(val types.String, envKey string) string {
	if !val.IsNull() && !val.IsUnknown() {
		return val.ValueString()
	}
	return os.Getenv(envKey)
}
