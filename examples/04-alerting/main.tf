# Example 04: Alerting
#
# Create a monitor with Slack and email alert channels, then link both
# channels to the monitor so alerts fire when the monitor goes down.
#
# Usage:
#   export OACK_API_KEY="oack_acc_xxxxxxxxxxxxxxxxxxxx"
#   export OACK_ACCOUNT_ID="your-account-uuid"
#   export TF_VAR_slack_webhook_url="https://hooks.slack.com/services/XXX/YYY/ZZZ"
#   terraform init && terraform apply

terraform {
  required_providers {
    oack = {
      source  = "oack-io/oack"
      version = "~> 0.1"
    }
  }
}

provider "oack" {}

variable "slack_webhook_url" {
  description = "Slack incoming webhook URL for alert notifications"
  type        = string
  sensitive   = true
}

# --- Team ---

resource "oack_team" "engineering" {
  name = "Engineering"
}

# --- Monitor ---

resource "oack_monitor" "api" {
  team_id           = oack_team.engineering.id
  name              = "API Health Check"
  url               = "https://httpbin.org/status/200"
  check_interval_ms = 60000
  failure_threshold = 3
}

# --- Alert Channels ---

resource "oack_alert_channel" "slack" {
  team_id = oack_team.engineering.id
  name    = "Engineering Slack"
  type    = "slack"

  config = {
    webhook_url = var.slack_webhook_url
  }
}

resource "oack_alert_channel" "email" {
  team_id = oack_team.engineering.id
  name    = "On-Call Email"
  type    = "email"

  config = {
    email = "oncall@example.com"
  }
}

# --- Link channels to the monitor ---

resource "oack_monitor_alert_channel_link" "api_slack" {
  team_id    = oack_team.engineering.id
  monitor_id = oack_monitor.api.id
  channel_id = oack_alert_channel.slack.id
}

resource "oack_monitor_alert_channel_link" "api_email" {
  team_id    = oack_team.engineering.id
  monitor_id = oack_monitor.api.id
  channel_id = oack_alert_channel.email.id
}
