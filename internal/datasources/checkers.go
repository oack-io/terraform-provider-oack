package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/oack-io/terraform-provider-oack/internal/client"
)

var _ datasource.DataSource = &CheckersDataSource{}

type CheckersDataSource struct {
	client *client.Client
}

type CheckersDataSourceModel struct {
	Checkers []CheckerModel `tfsdk:"checkers"`
}

type CheckerModel struct {
	ID      types.String `tfsdk:"id"`
	Region  types.String `tfsdk:"region"`
	Country types.String `tfsdk:"country"`
	IP      types.String `tfsdk:"ip"`
	ASN     types.String `tfsdk:"asn"`
	Mode    types.String `tfsdk:"mode"`
	Status  types.String `tfsdk:"status"`
}

func NewCheckersDataSource() datasource.DataSource {
	return &CheckersDataSource{}
}

func (d *CheckersDataSource) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_checkers"
}

func (d *CheckersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List available Oack checker nodes.",
		Attributes: map[string]schema.Attribute{
			"checkers": schema.ListNestedAttribute{
				Description: "List of checker nodes.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":      schema.StringAttribute{Computed: true},
						"region":  schema.StringAttribute{Computed: true},
						"country": schema.StringAttribute{Computed: true},
						"ip":      schema.StringAttribute{Computed: true},
						"asn":     schema.StringAttribute{Computed: true},
						"mode":    schema.StringAttribute{Computed: true},
						"status":  schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *CheckersDataSource) Configure(
	_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = c
}

func (d *CheckersDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	checkers, err := d.client.ListCheckers(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Read Checkers Failed", err.Error())
		return
	}

	var state CheckersDataSourceModel
	for _, c := range checkers {
		asn := fmt.Sprintf("%v", c.ASN)
		state.Checkers = append(state.Checkers, CheckerModel{
			ID:      types.StringValue(c.ID),
			Region:  types.StringValue(c.Region),
			Country: types.StringValue(c.Country),
			IP:      types.StringValue(c.IP),
			ASN:     types.StringValue(asn),
			Mode:    types.StringValue(c.Mode),
			Status:  types.StringValue(c.Status),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
