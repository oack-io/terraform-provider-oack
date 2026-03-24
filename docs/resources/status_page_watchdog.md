---
page_title: "oack_status_page_watchdog Resource - terraform-provider-oack"
subcategory: ""
description: |-
  Manages a watchdog linking a monitor to a status page component.
---

# oack_status_page_watchdog (Resource)

Manages a watchdog that links a monitor to a status page component. When the
monitor goes down, the watchdog can automatically create an incident on the
status page and notify subscribers. When the monitor recovers, the incident
can be automatically resolved.

This resource is immutable -- any attribute change triggers a destroy and
recreate.

## Example Usage

```hcl
resource "oack_status_page_watchdog" "api" {
  status_page_id     = oack_status_page.public.id
  component_id       = oack_status_page_component.api.id
  monitor_id         = oack_monitor.api.id
  severity           = "major"
  auto_create        = true
  auto_resolve       = true
  notify_subscribers = true
}
```

## Argument Reference

- `status_page_id` - (Required, Forces new resource) Status page UUID.
- `component_id` - (Required, Forces new resource) Component UUID.
- `monitor_id` - (Required, Forces new resource) Monitor UUID.
- `severity` - (Required, Forces new resource) Incident severity. One of: `minor`, `medium`, `major`, `critical`.
- `auto_create` - (Optional) Automatically create incidents when the monitor goes down. Default: `true`.
- `auto_resolve` - (Optional) Automatically resolve incidents when the monitor recovers. Default: `true`.
- `notify_subscribers` - (Optional) Notify status page subscribers on incident changes. Default: `true`.
- `template_id` - (Optional) Incident template UUID to use for auto-created incidents.

## Attribute Reference

- `id` - The UUID of the watchdog.
- `created_at` - Creation timestamp (RFC 3339).

## Import

Import uses the format `page_id/component_id/watchdog_id`:

```shell
terraform import oack_status_page_watchdog.example <page_id>/<component_id>/<watchdog_id>
```
