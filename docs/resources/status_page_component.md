---
page_title: "oack_status_page_component Resource - terraform-provider-oack"
subcategory: ""
description: |-
  Manages a component within an Oack status page.
---

# oack_status_page_component (Resource)

Manages a component within an Oack status page. Components represent individual
services or systems displayed on the status page. They can optionally be placed
inside a component group.

## Example Usage

```hcl
resource "oack_status_page_component" "api" {
  status_page_id = oack_status_page.public.id
  group_id       = oack_status_page_component_group.backend.id
  name           = "REST API"
  description    = "Public-facing REST API"
  display_uptime = true
  position       = 0
}
```

### Without a group

```hcl
resource "oack_status_page_component" "docs" {
  status_page_id = oack_status_page.public.id
  name           = "Documentation"
  display_uptime = false
  position       = 10
}
```

## Argument Reference

- `status_page_id` - (Required, Forces new resource) Status page UUID.
- `group_id` - (Optional) Component group UUID. If omitted, the component appears at the top level.
- `name` - (Required) Component display name.
- `description` - (Optional) Component description.
- `display_uptime` - (Optional) Whether to display uptime for this component. Default: `true`.
- `position` - (Optional) Display order position. Default: `0`.

## Attribute Reference

- `id` - The UUID of the component.
- `status` - Current component status.
- `created_at` - Creation timestamp (RFC 3339).
- `updated_at` - Last update timestamp (RFC 3339).

## Import

Import uses the format `page_id/component_id`:

```shell
terraform import oack_status_page_component.example <page_id>/<component_id>
```
