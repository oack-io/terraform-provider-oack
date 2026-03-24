package acctest

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTeam_basic(t *testing.T) {
	name := fmt.Sprintf("tf-test-%d", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_team.test", "id"),
					resource.TestCheckResourceAttr("oack_team.test", "name", name),
					resource.TestCheckResourceAttrSet("oack_team.test", "created_at"),
					resource.TestCheckResourceAttrSet("oack_team.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccTeam_update(t *testing.T) {
	name := fmt.Sprintf("tf-test-%d", time.Now().UnixNano())
	nameUpdated := name + "-updated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("oack_team.test", "name", name),
				),
			},
			{
				Config: testAccTeamConfig(nameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("oack_team.test", "name", nameUpdated),
				),
			},
		},
	})
}

func TestAccTeam_import(t *testing.T) {
	name := fmt.Sprintf("tf-test-%d", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamConfig(name),
			},
			{
				ResourceName:      "oack_team.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTeamConfig(name string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_team" "test" {
  name = %q
}
`, name)
}
