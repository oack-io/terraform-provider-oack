package acctest

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccTeamAPIKey_basic(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	keyName := fmt.Sprintf("tf-test-key-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamAPIKeyConfig(teamName, keyName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_team_api_key.test", "id"),
					resource.TestCheckResourceAttrSet("oack_team_api_key.test", "team_id"),
					resource.TestCheckResourceAttr("oack_team_api_key.test", "name", keyName),
					resource.TestCheckResourceAttrSet("oack_team_api_key.test", "key"),
					resource.TestCheckResourceAttrSet("oack_team_api_key.test", "key_prefix"),
					resource.TestCheckResourceAttrSet("oack_team_api_key.test", "created_at"),
				),
			},
		},
	})
}

func TestAccTeamAPIKey_import(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-team-%d", ts)
	keyName := fmt.Sprintf("tf-test-key-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamAPIKeyConfig(teamName, keyName),
			},
			{
				ResourceName:            "oack_team_api_key.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key"},
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["oack_team_api_key.test"]
					if !ok {
						return "", fmt.Errorf("resource not found: oack_team_api_key.test")
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

func testAccTeamAPIKeyConfig(teamName, keyName string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_team" "test" {
  name = %q
}

resource "oack_team_api_key" "test" {
  team_id = oack_team.test.id
  name    = %q
}
`, teamName, keyName)
}
