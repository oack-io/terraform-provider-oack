---
page_title: "oack_alert_channel Resource - terraform-provider-oack"
subcategory: ""
description: |-
  Manages an Oack alert channel (Slack, email, webhook, Telegram, Discord, PagerDuty).
---

# oack_alert_channel (Resource)

Manages an Oack alert channel. Alert channels define where notifications are
sent when a linked monitor goes down or recovers. Supported types: Slack,
email, webhook, Telegram, Discord, and PagerDuty.

## Example Usage

### Slack

```hcl
resource "oack_alert_channel" "slack" {
  team_id = oack_team.engineering.id
  name    = "Engineering Slack"
  type    = "slack"

  config = {
    webhook_url = "https://hooks.slack.com/services/XXX/YYY/ZZZ"
  }
}
```

### Email

```hcl
resource "oack_alert_channel" "email" {
  team_id = oack_team.engineering.id
  name    = "On-Call Email"
  type    = "email"

  config = {
    email = "oncall@example.com"
  }
}
```

### Webhook

```hcl
resource "oack_alert_channel" "webhook" {
  team_id = oack_team.engineering.id
  name    = "Incident Webhook"
  type    = "webhook"

  config = {
    url = "https://example.com/hooks/oack"
  }
}
```

### Telegram

```hcl
resource "oack_alert_channel" "telegram" {
  team_id = oack_team.engineering.id
  name    = "Ops Telegram"
  type    = "telegram"

  config = {
    chat_id = "-1001234567890"
  }
}
```

### Discord

```hcl
resource "oack_alert_channel" "discord" {
  team_id = oack_team.engineering.id
  name    = "Alerts Discord"
  type    = "discord"

  config = {
    webhook_url = "https://discord.com/api/webhooks/1234567890/abcdef"
  }
}
```

### PagerDuty

```hcl
resource "oack_alert_channel" "pagerduty" {
  team_id = oack_team.engineering.id
  name    = "PagerDuty Routing"
  type    = "pagerduty"

  config = {
    routing_key = "your-pagerduty-routing-key"
    region      = "us"
  }
}
```

## Argument Reference

- `team_id` - (Required, Forces new resource) Team UUID.
- `name` - (Required) Channel display name.
- `type` - (Required, Forces new resource) Channel type. One of: `slack`, `email`, `webhook`, `telegram`, `discord`, `pagerduty`.
- `config` - (Required, Sensitive) Type-specific configuration map. See the table below.
- `enabled` - (Optional) Whether the channel is active. Default: `true`.

### Config Keys by Type

| Type       | Required Keys          | Description                       |
|------------|------------------------|-----------------------------------|
| `slack`    | `webhook_url`          | Slack incoming webhook URL        |
| `email`    | `email`                | Recipient email address           |
| `webhook`  | `url`                  | Webhook endpoint URL              |
| `telegram` | `chat_id`              | Telegram chat ID                  |
| `discord`  | `webhook_url`          | Discord webhook URL               |
| `pagerduty`| `routing_key`, `region`| PagerDuty routing key and region  |

## Attribute Reference

- `id` - The UUID of the alert channel.
- `email_verified` - Whether the email address is verified (email channels only).
- `created_at` - Creation timestamp (RFC 3339).
- `updated_at` - Last update timestamp (RFC 3339).

## Import

Import uses the format `team_id/channel_id`:

```shell
terraform import oack_alert_channel.example <team_id>/<channel_id>
```
