---
page_title: "oack_status_page_component_group Resource - terraform-provider-oack"
subcategory: ""
description: |-
  Manages a component group within an Oack status page.
---

# oack_status_page_component_group (Resource)

Manages a component group within an Oack status page. Component groups let you
organize related components under a collapsible section on the status page.

## Example Usage

```hcl
resource "oack_status_page_component_group" "backend" {
  status_page_id = oack_status_page.public.id
  name           = "Backend Services"
  description    = "Core API and backend infrastructure"
  position       = 0
  collapsed      = false
}
```

## Argument Reference

- `status_page_id` - (Required, Forces new resource) Status page UUID.
- `name` - (Required) Group display name.
- `description` - (Optional) Group description.
- `position` - (Optional) Display order position. Default: `0`.
- `collapsed` - (Optional) Whether the group is collapsed by default. Default: `false`.

## Attribute Reference

- `id` - The UUID of the component group.
- `created_at` - Creation timestamp (RFC 3339).
- `updated_at` - Last update timestamp (RFC 3339).

## Import

Import uses the format `page_id/group_id`:

```shell
terraform import oack_status_page_component_group.example <page_id>/<group_id>
```
