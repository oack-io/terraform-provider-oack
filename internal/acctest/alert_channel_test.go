package acctest

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAlertChannel_email(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	channelName := fmt.Sprintf("tf-test-email-ch-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelEmailConfig(teamName, channelName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_alert_channel.test", "id"),
					resource.TestCheckResourceAttrSet("oack_alert_channel.test", "team_id"),
					resource.TestCheckResourceAttr("oack_alert_channel.test", "name", channelName),
					resource.TestCheckResourceAttr("oack_alert_channel.test", "type", "email"),
					resource.TestCheckResourceAttr("oack_alert_channel.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("oack_alert_channel.test", "created_at"),
					resource.TestCheckResourceAttrSet("oack_alert_channel.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccAlertChannel_webhook(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	channelName := fmt.Sprintf("tf-test-webhook-ch-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelWebhookConfig(teamName, channelName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_alert_channel.test", "id"),
					resource.TestCheckResourceAttr("oack_alert_channel.test", "name", channelName),
					resource.TestCheckResourceAttr("oack_alert_channel.test", "type", "webhook"),
					resource.TestCheckResourceAttr("oack_alert_channel.test", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccAlertChannel_update(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	channelName := fmt.Sprintf("tf-test-ch-%d", ts)
	channelNameUpdated := channelName + "-updated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelWebhookConfig(teamName, channelName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("oack_alert_channel.test", "name", channelName),
				),
			},
			{
				Config: testAccAlertChannelWebhookConfig(teamName, channelNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("oack_alert_channel.test", "name", channelNameUpdated),
				),
			},
		},
	})
}

func TestAccAlertChannel_import(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	channelName := fmt.Sprintf("tf-test-ch-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelWebhookConfig(teamName, channelName),
			},
			{
				ResourceName:            "oack_alert_channel.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config"},
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["oack_alert_channel.test"]
					if !ok {
						return "", fmt.Errorf("resource not found: oack_alert_channel.test")
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

func testAccAlertChannelEmailConfig(teamName, channelName string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_team" "test" {
  name = %q
}

resource "oack_alert_channel" "test" {
  team_id = oack_team.test.id
  name    = %q
  type    = "email"
  config = {
    "email" = "tf-test@example.com"
  }
}
`, teamName, channelName)
}

func testAccAlertChannelWebhookConfig(teamName, channelName string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_team" "test" {
  name = %q
}

resource "oack_alert_channel" "test" {
  team_id = oack_team.test.id
  name    = %q
  type    = "webhook"
  config = {
    "url" = "https://example.com/webhook"
  }
}
`, teamName, channelName)
}
