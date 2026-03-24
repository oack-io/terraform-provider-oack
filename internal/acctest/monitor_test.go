package acctest

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMonitor_basic(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	monitorName := fmt.Sprintf("tf-test-monitor-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorBasicConfig(teamName, monitorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_monitor.test", "id"),
					resource.TestCheckResourceAttrSet("oack_monitor.test", "team_id"),
					resource.TestCheckResourceAttr("oack_monitor.test", "name", monitorName),
					resource.TestCheckResourceAttr("oack_monitor.test", "url", "https://example.com"),
					resource.TestCheckResourceAttr("oack_monitor.test", "status", "active"),
					resource.TestCheckResourceAttr("oack_monitor.test", "check_interval_ms", "60000"),
					resource.TestCheckResourceAttr("oack_monitor.test", "timeout_ms", "10000"),
					resource.TestCheckResourceAttr("oack_monitor.test", "http_method", "GET"),
					resource.TestCheckResourceAttr("oack_monitor.test", "follow_redirects", "true"),
					resource.TestCheckResourceAttr("oack_monitor.test", "failure_threshold", "3"),
					resource.TestCheckResourceAttr("oack_monitor.test", "ssl_expiry_enabled", "true"),
					resource.TestCheckResourceAttr("oack_monitor.test", "domain_expiry_enabled", "true"),
					resource.TestCheckResourceAttrSet("oack_monitor.test", "created_at"),
					resource.TestCheckResourceAttrSet("oack_monitor.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccMonitor_full(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	monitorName := fmt.Sprintf("tf-test-monitor-full-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorFullConfig(teamName, monitorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_monitor.test", "id"),
					resource.TestCheckResourceAttr("oack_monitor.test", "name", monitorName),
					resource.TestCheckResourceAttr("oack_monitor.test", "url", "https://example.com/health"),
					resource.TestCheckResourceAttr("oack_monitor.test", "status", "active"),
					resource.TestCheckResourceAttr("oack_monitor.test", "check_interval_ms", "120000"),
					resource.TestCheckResourceAttr("oack_monitor.test", "timeout_ms", "5000"),
					resource.TestCheckResourceAttr("oack_monitor.test", "http_method", "POST"),
					resource.TestCheckResourceAttr("oack_monitor.test", "follow_redirects", "false"),
					resource.TestCheckResourceAttr("oack_monitor.test", "failure_threshold", "5"),
					resource.TestCheckResourceAttr("oack_monitor.test", "latency_threshold_ms", "2000"),
					resource.TestCheckResourceAttr("oack_monitor.test", "ssl_expiry_enabled", "true"),
					resource.TestCheckResourceAttr("oack_monitor.test", "domain_expiry_enabled", "false"),
					resource.TestCheckResourceAttr("oack_monitor.test", "uptime_threshold_good", "99.9"),
					resource.TestCheckResourceAttr("oack_monitor.test", "uptime_threshold_degraded", "99"),
					resource.TestCheckResourceAttr("oack_monitor.test", "uptime_threshold_critical", "95"),
					resource.TestCheckResourceAttr("oack_monitor.test", "headers.X-Custom-Header", "test-value"),
					resource.TestCheckResourceAttr("oack_monitor.test", "allowed_status_codes.0", "2xx"),
					resource.TestCheckResourceAttr("oack_monitor.test", "allowed_status_codes.1", "301"),
					resource.TestCheckResourceAttr("oack_monitor.test", "ssl_expiry_thresholds.0", "30"),
					resource.TestCheckResourceAttr("oack_monitor.test", "ssl_expiry_thresholds.1", "14"),
					resource.TestCheckResourceAttr("oack_monitor.test", "ssl_expiry_thresholds.2", "7"),
				),
			},
		},
	})
}

func TestAccMonitor_update(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	monitorName := fmt.Sprintf("tf-test-monitor-%d", ts)
	monitorNameUpdated := monitorName + "-updated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorBasicConfig(teamName, monitorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("oack_monitor.test", "name", monitorName),
					resource.TestCheckResourceAttr("oack_monitor.test", "check_interval_ms", "60000"),
				),
			},
			{
				Config: testAccMonitorUpdateConfig(teamName, monitorNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("oack_monitor.test", "name", monitorNameUpdated),
					resource.TestCheckResourceAttr("oack_monitor.test", "check_interval_ms", "120000"),
				),
			},
		},
	})
}

func TestAccMonitor_import(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	monitorName := fmt.Sprintf("tf-test-monitor-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorBasicConfig(teamName, monitorName),
			},
			{
				ResourceName:      "oack_monitor.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["oack_monitor.test"]
					if !ok {
						return "", fmt.Errorf("resource not found: oack_monitor.test")
					}
					return fmt.Sprintf("%s/%s",
						rs.Primary.Attributes["team_id"],
						rs.Primary.Attributes["id"],
					), nil
				},
			},
		},
	})
}

func testAccMonitorBasicConfig(teamName, monitorName string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_team" "test" {
  name = %q
}

resource "oack_monitor" "test" {
  team_id = oack_team.test.id
  name    = %q
  url     = "https://example.com"
}
`, teamName, monitorName)
}

func testAccMonitorFullConfig(teamName, monitorName string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_team" "test" {
  name = %q
}

resource "oack_monitor" "test" {
  team_id          = oack_team.test.id
  name             = %q
  url              = "https://example.com/health"
  status           = "active"
  check_interval_ms = 120000
  timeout_ms       = 5000
  http_method      = "POST"
  follow_redirects = false
  failure_threshold = 5
  latency_threshold_ms = 2000

  ssl_expiry_enabled     = true
  ssl_expiry_thresholds  = [30, 14, 7]
  domain_expiry_enabled  = false

  headers = {
    "X-Custom-Header" = "test-value"
  }

  allowed_status_codes = ["2xx", "301"]
}
`, teamName, monitorName)
}

func testAccMonitorUpdateConfig(teamName, monitorName string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_team" "test" {
  name = %q
}

resource "oack_monitor" "test" {
  team_id           = oack_team.test.id
  name              = %q
  url               = "https://example.com"
  check_interval_ms = 120000
}
`, teamName, monitorName)
}
