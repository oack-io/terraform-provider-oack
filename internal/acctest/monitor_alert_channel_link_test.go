package acctest

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMonitorAlertChannelLink_basic(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	monitorName := fmt.Sprintf("tf-test-monitor-%d", ts)
	channelName := fmt.Sprintf("tf-test-channel-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlertChannelLinkConfig(teamName, monitorName, channelName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_monitor_alert_channel_link.test", "team_id"),
					resource.TestCheckResourceAttrSet("oack_monitor_alert_channel_link.test", "monitor_id"),
					resource.TestCheckResourceAttrSet("oack_monitor_alert_channel_link.test", "channel_id"),
				),
			},
		},
	})
}

func TestAccMonitorAlertChannelLink_import(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	monitorName := fmt.Sprintf("tf-test-monitor-%d", ts)
	channelName := fmt.Sprintf("tf-test-channel-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlertChannelLinkConfig(teamName, monitorName, channelName),
			},
			{
				ResourceName:      "oack_monitor_alert_channel_link.test",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["oack_monitor_alert_channel_link.test"]
					if !ok {
						return "", fmt.Errorf("resource not found: oack_monitor_alert_channel_link.test")
					}
					return fmt.Sprintf("%s/%s/%s",
						rs.Primary.Attributes["team_id"],
						rs.Primary.Attributes["monitor_id"],
						rs.Primary.Attributes["channel_id"],
					), nil
				},
			},
		},
	})
}

func testAccMonitorAlertChannelLinkConfig(teamName, monitorName, channelName string) string {
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

resource "oack_alert_channel" "test" {
  team_id = oack_team.test.id
  name    = %q
  type    = "webhook"
  config = {
    "url" = "https://example.com/webhook"
  }
}

resource "oack_monitor_alert_channel_link" "test" {
  team_id    = oack_team.test.id
  monitor_id = oack_monitor.test.id
  channel_id = oack_alert_channel.test.id
}
`, teamName, monitorName, channelName)
}
