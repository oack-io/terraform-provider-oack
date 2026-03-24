package acctest

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyIntegration_basic(t *testing.T) {
	pdAPIKey := os.Getenv("OACK_PD_API_KEY")
	if pdAPIKey == "" {
		t.Skip("OACK_PD_API_KEY not set, skipping PagerDuty integration test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyIntegrationConfig(pdAPIKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_pagerduty_integration.test", "id"),
					resource.TestCheckResourceAttr(
						"oack_pagerduty_integration.test", "region", "us",
					),
					resource.TestCheckResourceAttr(
						"oack_pagerduty_integration.test", "sync_enabled", "true",
					),
					resource.TestCheckResourceAttrSet(
						"oack_pagerduty_integration.test", "created_at",
					),
					resource.TestCheckResourceAttrSet(
						"oack_pagerduty_integration.test", "updated_at",
					),
				),
			},
		},
	})
}

func testAccPagerDutyIntegrationConfig(apiKey string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_pagerduty_integration" "test" {
  api_key = %q
  region  = "us"
}
`, apiKey)
}
