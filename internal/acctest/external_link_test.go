package acctest

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccExternalLink_basic(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	linkName := fmt.Sprintf("tf-test-link-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExternalLinkBasicConfig(teamName, linkName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_external_link.test", "id"),
					resource.TestCheckResourceAttrSet("oack_external_link.test", "team_id"),
					resource.TestCheckResourceAttr("oack_external_link.test", "name", linkName),
					resource.TestCheckResourceAttr(
						"oack_external_link.test",
						"url_template",
						"https://grafana.example.com/d/abc?from={{start}}&to={{end}}",
					),
					resource.TestCheckResourceAttr(
						"oack_external_link.test", "time_window_minutes", "60",
					),
					resource.TestCheckResourceAttrSet("oack_external_link.test", "created_at"),
					resource.TestCheckResourceAttrSet("oack_external_link.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccExternalLink_update(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	linkName := fmt.Sprintf("tf-test-link-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExternalLinkBasicConfig(teamName, linkName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"oack_external_link.test",
						"url_template",
						"https://grafana.example.com/d/abc?from={{start}}&to={{end}}",
					),
				),
			},
			{
				Config: testAccExternalLinkUpdatedConfig(teamName, linkName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"oack_external_link.test",
						"url_template",
						"https://kibana.example.com/app/discover?from={{start}}&to={{end}}",
					),
					resource.TestCheckResourceAttr(
						"oack_external_link.test", "time_window_minutes", "120",
					),
				),
			},
		},
	})
}

func TestAccExternalLink_import(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	linkName := fmt.Sprintf("tf-test-link-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExternalLinkBasicConfig(teamName, linkName),
			},
			{
				ResourceName:      "oack_external_link.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["oack_external_link.test"]
					if !ok {
						return "", fmt.Errorf("resource not found: oack_external_link.test")
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

func testAccExternalLinkBasicConfig(teamName, linkName string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_team" "test" {
  name = %q
}

resource "oack_external_link" "test" {
  team_id             = oack_team.test.id
  name                = %q
  url_template        = "https://grafana.example.com/d/abc?from={{start}}&to={{end}}"
  time_window_minutes = 60
}
`, teamName, linkName)
}

func testAccExternalLinkUpdatedConfig(teamName, linkName string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_team" "test" {
  name = %q
}

resource "oack_external_link" "test" {
  team_id             = oack_team.test.id
  name                = %q
  url_template        = "https://kibana.example.com/app/discover?from={{start}}&to={{end}}"
  time_window_minutes = 120
}
`, teamName, linkName)
}
