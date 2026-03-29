package acctest

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccEnvVar_basic(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-env-team-%d", ts)
	key := fmt.Sprintf("TF_TEST_KEY_%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvVarBasicConfig(teamName, key, "test-value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_env_var.test", "id"),
					resource.TestCheckResourceAttrSet("oack_env_var.test", "team_id"),
					resource.TestCheckResourceAttr("oack_env_var.test", "key", key),
					resource.TestCheckResourceAttr("oack_env_var.test", "value", "test-value"),
					resource.TestCheckResourceAttr("oack_env_var.test", "is_secret", "false"),
					resource.TestCheckResourceAttrSet("oack_env_var.test", "created_at"),
					resource.TestCheckResourceAttrSet("oack_env_var.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccEnvVar_update(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-env-team-%d", ts)
	key := fmt.Sprintf("TF_TEST_KEY_%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvVarBasicConfig(teamName, key, "original"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("oack_env_var.test", "value", "original"),
				),
			},
			{
				Config: testAccEnvVarBasicConfig(teamName, key, "updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("oack_env_var.test", "value", "updated"),
				),
			},
		},
	})
}

func TestAccEnvVar_import(t *testing.T) {
	ts := time.Now().UnixNano()
	teamName := fmt.Sprintf("tf-test-env-team-%d", ts)
	key := fmt.Sprintf("TF_TEST_KEY_%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvVarBasicConfig(teamName, key, "import-value"),
			},
			{
				ResourceName: "oack_env_var.test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["oack_env_var.test"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}
					return rs.Primary.Attributes["team_id"] + "/" + rs.Primary.Attributes["key"], nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
		},
	})
}

func testAccEnvVarBasicConfig(teamName, key, value string) string {
	return fmt.Sprintf(`
resource "oack_team" "test" {
  name = %q
}

resource "oack_env_var" "test" {
  team_id   = oack_team.test.id
  key       = %q
  value     = %q
  is_secret = false
}
`, teamName, key, value)
}
