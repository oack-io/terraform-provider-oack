---
page_title: "oack_external_link Resource - terraform-provider-oack"
subcategory: ""
description: |-
  Manages an Oack external link.
---

# oack_external_link (Resource)

Manages an Oack external link. External links provide quick access to related
dashboards (e.g. Grafana, Datadog) directly from the Oack UI. The URL template
supports time-window placeholders so links open with the relevant time range.

## Example Usage

```hcl
resource "oack_external_link" "grafana" {
  team_id             = oack_team.engineering.id
  name                = "Grafana Dashboard"
  url_template        = "https://grafana.example.com/d/uptime?from=now-{{.TimeWindow}}&to=now"
  time_window_minutes = 60
}
```

## Argument Reference

- `team_id` - (Required, Forces new resource) Team UUID.
- `name` - (Required) Link display name.
- `url_template` - (Required) URL template for the external link.
- `icon_url` - (Optional) Icon URL for the external link. Default: `""`.
- `time_window_minutes` - (Required) Time window in minutes passed to the URL template.

## Attribute Reference

- `id` - The UUID of the external link.
- `created_at` - Creation timestamp (RFC 3339).
- `updated_at` - Last update timestamp (RFC 3339).

## Import

Import uses the format `team_id/link_id`:

```shell
terraform import oack_external_link.example <team_id>/<link_id>
```
