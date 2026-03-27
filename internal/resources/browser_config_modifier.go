package resources

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	oack "github.com/oack-io/oack-go"
)

// browserConfigSemanticEqual normalizes browser_config_json through the Go
// BrowserConfig struct so that both plan and state use the same canonical
// JSON serialization (Go struct field order with all zero-value fields).
// This prevents spurious diffs when the user's jsonencode() uses alphabetical
// order or omits zero-value fields like user_agent.
type browserConfigSemanticEqual struct{}

func (m browserConfigSemanticEqual) Description(_ context.Context) string {
	return "Normalize browser_config_json through the BrowserConfig struct."
}

func (m browserConfigSemanticEqual) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m browserConfigSemanticEqual) PlanModifyString(
	_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse,
) {
	// Nothing to normalize if config is null or unknown.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Always normalize the config value through the Go struct so the plan
	// value matches what monitorToState will produce after Create/Read.
	var bc oack.BrowserConfig
	if err := json.Unmarshal([]byte(req.ConfigValue.ValueString()), &bc); err != nil {
		return // let validation catch it later
	}
	normalized, _ := json.Marshal(bc)
	resp.PlanValue = types.StringValue(string(normalized))
}

func BrowserConfigSemanticEqual() planmodifier.String {
	return browserConfigSemanticEqual{}
}
