package acctest

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccWatchdog_basic(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	pageName := fmt.Sprintf("tf-test-%d-status", ts)
	pageSlug := fmt.Sprintf("tf-test-wd-%d-status", ts)
	compName := fmt.Sprintf("tf-test-comp-%d", ts)
	monitorName := fmt.Sprintf("tf-test-monitor-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWatchdogBasicConfig(
					teamName, pageName, pageSlug, compName, monitorName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_status_page_watchdog.test", "id"),
					resource.TestCheckResourceAttrSet(
						"oack_status_page_watchdog.test", "status_page_id",
					),
					resource.TestCheckResourceAttrSet(
						"oack_status_page_watchdog.test", "component_id",
					),
					resource.TestCheckResourceAttrSet(
						"oack_status_page_watchdog.test", "monitor_id",
					),
					resource.TestCheckResourceAttr(
						"oack_status_page_watchdog.test", "severity", "major",
					),
					resource.TestCheckResourceAttr(
						"oack_status_page_watchdog.test", "auto_create", "true",
					),
					resource.TestCheckResourceAttr(
						"oack_status_page_watchdog.test", "auto_resolve", "true",
					),
					resource.TestCheckResourceAttr(
						"oack_status_page_watchdog.test", "notify_subscribers", "true",
					),
					resource.TestCheckResourceAttrSet(
						"oack_status_page_watchdog.test", "created_at",
					),
				),
			},
		},
	})
}

func TestAccWatchdog_import(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	pageName := fmt.Sprintf("tf-test-%d-status", ts)
	pageSlug := fmt.Sprintf("tf-test-wdi-%d-status", ts)
	compName := fmt.Sprintf("tf-test-comp-%d", ts)
	monitorName := fmt.Sprintf("tf-test-monitor-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWatchdogBasicConfig(
					teamName, pageName, pageSlug, compName, monitorName,
				),
			},
			{
				ResourceName:      "oack_status_page_watchdog.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["oack_status_page_watchdog.test"]
					if !ok {
						return "", fmt.Errorf(
							"resource not found: oack_status_page_watchdog.test",
						)
					}
					return fmt.Sprintf("%s/%s/%s",
						rs.Primary.Attributes["status_page_id"],
						rs.Primary.Attributes["component_id"],
						rs.Primary.Attributes["id"],
					), nil
				},
			},
		},
	})
}

func testAccWatchdogBasicConfig(
	teamName, pageName, pageSlug, compName, monitorName string,
) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_team" "test" {
  name = %q
}

resource "oack_status_page" "test" {
  name = %q
  slug = %q
}

resource "oack_status_page_component" "test" {
  status_page_id = oack_status_page.test.id
  name           = %q
}

resource "oack_monitor" "test" {
  team_id = oack_team.test.id
  name    = %q
  url     = "https://example.com"
}

resource "oack_status_page_watchdog" "test" {
  status_page_id = oack_status_page.test.id
  component_id   = oack_status_page_component.test.id
  monitor_id     = oack_monitor.test.id
  severity       = "major"
}
`, teamName, pageName, pageSlug, compName, monitorName)
}
