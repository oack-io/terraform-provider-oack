package acctest

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccComponent_basic(t *testing.T) {
	ts := time.Now().UnixNano()
	pageName := fmt.Sprintf("tf-test-page-%d", ts)
	pageSlug := fmt.Sprintf("tf-test-comp-%d", ts)
	compName := fmt.Sprintf("tf-test-component-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComponentBasicConfig(pageName, pageSlug, compName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_status_page_component.test", "id"),
					resource.TestCheckResourceAttrSet(
						"oack_status_page_component.test", "status_page_id",
					),
					resource.TestCheckResourceAttr(
						"oack_status_page_component.test", "name", compName,
					),
					resource.TestCheckResourceAttr(
						"oack_status_page_component.test", "display_uptime", "true",
					),
					resource.TestCheckResourceAttr(
						"oack_status_page_component.test", "position", "0",
					),
					resource.TestCheckResourceAttrSet(
						"oack_status_page_component.test", "status",
					),
					resource.TestCheckResourceAttrSet(
						"oack_status_page_component.test", "created_at",
					),
					resource.TestCheckResourceAttrSet(
						"oack_status_page_component.test", "updated_at",
					),
				),
			},
		},
	})
}

func TestAccComponent_withGroup(t *testing.T) {
	ts := time.Now().UnixNano()
	pageName := fmt.Sprintf("tf-test-page-%d", ts)
	pageSlug := fmt.Sprintf("tf-test-compg-%d", ts)
	groupName := fmt.Sprintf("tf-test-group-%d", ts)
	compName := fmt.Sprintf("tf-test-component-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComponentWithGroupConfig(pageName, pageSlug, groupName, compName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_status_page_component.test", "id"),
					resource.TestCheckResourceAttr(
						"oack_status_page_component.test", "name", compName,
					),
					resource.TestCheckResourceAttrSet(
						"oack_status_page_component.test", "group_id",
					),
				),
			},
		},
	})
}

func TestAccComponent_import(t *testing.T) {
	ts := time.Now().UnixNano()
	pageName := fmt.Sprintf("tf-test-page-%d", ts)
	pageSlug := fmt.Sprintf("tf-test-compi-%d", ts)
	compName := fmt.Sprintf("tf-test-component-%d", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComponentBasicConfig(pageName, pageSlug, compName),
			},
			{
				ResourceName:      "oack_status_page_component.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["oack_status_page_component.test"]
					if !ok {
						return "", fmt.Errorf(
							"resource not found: oack_status_page_component.test",
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

func testAccComponentBasicConfig(pageName, pageSlug, compName string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_status_page" "test" {
  name = %q
  slug = %q
}

resource "oack_status_page_component" "test" {
  status_page_id = oack_status_page.test.id
  name           = %q
}
`, pageName, pageSlug, compName)
}

func testAccComponentWithGroupConfig(pageName, pageSlug, groupName, compName string) string {
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

resource "oack_status_page_component" "test" {
  status_page_id = oack_status_page.test.id
  group_id       = oack_status_page_component_group.test.id
  name           = %q
}
`, pageName, pageSlug, groupName, compName)
}
