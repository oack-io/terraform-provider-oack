package acctest

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccComponentGroup_basic(t *testing.T) {
	ts := time.Now().UnixNano()
	pageName := fmt.Sprintf("tf-test-%d-status", ts)
	pageSlug := fmt.Sprintf("tf-test-cg-%d-status", ts)
	groupName := fmt.Sprintf("tf-test-group-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComponentGroupBasicConfig(pageName, pageSlug, groupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_status_page_component_group.test", "id"),
					resource.TestCheckResourceAttrSet(
						"oack_status_page_component_group.test", "status_page_id",
					),
					resource.TestCheckResourceAttr(
						"oack_status_page_component_group.test", "name", groupName,
					),
					resource.TestCheckResourceAttr(
						"oack_status_page_component_group.test", "position", "0",
					),
					resource.TestCheckResourceAttr(
						"oack_status_page_component_group.test", "collapsed", "false",
					),
					resource.TestCheckResourceAttrSet(
						"oack_status_page_component_group.test", "created_at",
					),
					resource.TestCheckResourceAttrSet(
						"oack_status_page_component_group.test", "updated_at",
					),
				),
			},
		},
	})
}

func TestAccComponentGroup_import(t *testing.T) {
	ts := time.Now().UnixNano()
	pageName := fmt.Sprintf("tf-test-%d-status", ts)
	pageSlug := fmt.Sprintf("tf-test-cgi-%d-status", ts)
	groupName := fmt.Sprintf("tf-test-group-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComponentGroupBasicConfig(pageName, pageSlug, groupName),
			},
			{
				ResourceName:      "oack_status_page_component_group.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["oack_status_page_component_group.test"]
					if !ok {
						return "", fmt.Errorf(
							"resource not found: oack_status_page_component_group.test",
						)
					}
					return fmt.Sprintf("%s/%s",
						rs.Primary.Attributes["status_page_id"],
						rs.Primary.Attributes["id"],
					), nil
				},
			},
		},
	})
}

func testAccComponentGroupBasicConfig(pageName, pageSlug, groupName string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_status_page" "test" {
  name = %q
  slug = %q
}

resource "oack_status_page_component_group" "test" {
  status_page_id = oack_status_page.test.id
  name           = %q
}
`, pageName, pageSlug, groupName)
}
