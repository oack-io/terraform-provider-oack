package acctest

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStatusPage_basic(t *testing.T) {
	ts := time.Now().UnixNano()
	name := fmt.Sprintf("tf-test-page-%d", ts)
	slug := fmt.Sprintf("tf-test-%d-status", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStatusPageBasicConfig(name, slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_status_page.test", "id"),
					resource.TestCheckResourceAttr("oack_status_page.test", "name", name),
					resource.TestCheckResourceAttr("oack_status_page.test", "slug", slug),
					resource.TestCheckResourceAttr("oack_status_page.test", "allow_iframe", "false"),
					resource.TestCheckResourceAttr("oack_status_page.test", "show_historical_uptime", "true"),
					resource.TestCheckResourceAttr("oack_status_page.test", "has_password", "false"),
					resource.TestCheckResourceAttrSet("oack_status_page.test", "created_at"),
					resource.TestCheckResourceAttrSet("oack_status_page.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccStatusPage_full(t *testing.T) {
	ts := time.Now().UnixNano()
	name := fmt.Sprintf("tf-test-page-full-%d", ts)
	slug := fmt.Sprintf("tf-test-full-%d-status", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStatusPageFullConfig(name, slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("oack_status_page.test", "id"),
					resource.TestCheckResourceAttr("oack_status_page.test", "name", name),
					resource.TestCheckResourceAttr("oack_status_page.test", "slug", slug),
					resource.TestCheckResourceAttr("oack_status_page.test", "description", "Test status page description"),
					resource.TestCheckResourceAttr("oack_status_page.test", "allow_iframe", "false"),
					resource.TestCheckResourceAttr("oack_status_page.test", "show_historical_uptime", "true"),
				),
			},
		},
	})
}

func TestAccStatusPage_update(t *testing.T) {
	ts := time.Now().UnixNano()
	name := fmt.Sprintf("tf-test-page-%d", ts)
	slug := fmt.Sprintf("tf-test-%d-status", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStatusPageBasicConfig(name, slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("oack_status_page.test", "name", name),
				),
			},
			{
				Config: testAccStatusPageWithDescriptionConfig(name+"-updated", slug, "Updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("oack_status_page.test", "name", name+"-updated"),
					resource.TestCheckResourceAttr("oack_status_page.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccStatusPage_import(t *testing.T) {
	ts := time.Now().UnixNano()
	name := fmt.Sprintf("tf-test-page-%d", ts)
	slug := fmt.Sprintf("tf-test-%d-status", ts)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStatusPageBasicConfig(name, slug),
			},
			{
				ResourceName:            "oack_status_page.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccStatusPageBasicConfig(name, slug string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_status_page" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccStatusPageFullConfig(name, slug string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_status_page" "test" {
  name        = %q
  slug        = %q
  description = "Test status page description"
}
`, name, slug)
}

func testAccStatusPageWithDescriptionConfig(name, slug, description string) string {
	return fmt.Sprintf(`
provider "oack" {}

resource "oack_status_page" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}
