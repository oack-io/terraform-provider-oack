---
page_title: "oack_monitor_alert_channel_link Resource - terraform-provider-oack"
subcategory: ""
description: |-
  Links an alert channel to a monitor.
---

# oack_monitor_alert_channel_link (Resource)

Links an alert channel to a monitor. When the monitor goes down, the linked
channel receives notifications. When the monitor recovers, the channel receives
a recovery notification.

All attributes force resource replacement -- this resource cannot be updated in
place.

## Example Usage

```hcl
resource "oack_monitor_alert_channel_link" "api_slack" {
  team_id    = oack_team.engineering.id
  monitor_id = oack_monitor.api.id
  channel_id = oack_alert_channel.slack.id
}
```

## Argument Reference

- `team_id` - (Required, Forces new resource) Team UUID.
- `monitor_id` - (Required, Forces new resource) Monitor UUID.
- `channel_id` - (Required, Forces new resource) Alert channel UUID.

## Attribute Reference

This resource has no additional computed attributes beyond the arguments.

## Import

Import uses the format `team_id/monitor_id/channel_id`:

```shell
terraform import oack_monitor_alert_channel_link.example <team_id>/<monitor_id>/<channel_id>
```
