# Example 09: Browser Monitor with Custom Playwright Script
#
# Run a custom Playwright script on a schedule to verify a multi-step user flow
# (login, navigate, assert content). The script runs in Oack's sandboxed
# Chromium environment — same engine, same browser, same results as running
# locally with `npx playwright test`.
#
# This example monitors a login → dashboard flow. The script:
#   1. Opens the login page
#   2. Fills email + password from encrypted env vars
#   3. Submits the form
#   4. Asserts the dashboard heading and a critical widget
#
# See check.js for the full script. See README.md for how to run it locally.
#
# Usage:
#   export OACK_API_KEY="oack_acc_xxxxxxxxxxxxxxxxxxxx"
#   export OACK_ACCOUNT_ID="your-account-uuid"
#   export TF_VAR_login_email="test@example.com"
#   export TF_VAR_login_password="s3cret"
#   terraform init && terraform apply

terraform {
  required_providers {
    oack = {
      source  = "oack-io/oack"
      version = "~> 0.2"
    }
  }
}

provider "oack" {}

variable "login_email" {
  description = "Email for the login check"
  type        = string
  sensitive   = true
}

variable "login_password" {
  description = "Password for the login check"
  type        = string
  sensitive   = true
}

resource "oack_team" "web" {
  name = "Web Team"
}

# Read the Playwright script from the local file.
# This keeps the script version-controlled alongside the Terraform config.
resource "oack_monitor" "login_flow" {
  team_id           = oack_team.web.id
  name              = "Login Flow — Scripted Browser Check"
  url               = "https://app.example.com/login"
  type              = "browser"
  check_interval_ms = 300000 # 5 minutes
  timeout_ms        = 30000  # 30 seconds

  browser_config_json = jsonencode({
    mode                   = "script"
    script                 = file("${path.module}/check.js")
    screenshot_enabled     = true
    screenshot_full_page   = false
    viewport_width         = 1920
    viewport_height        = 1080
    wait_until             = "load"
    console_error_threshold  = 5
    resource_error_threshold = 10

    # Environment variables injected into the script sandbox.
    # Marked as secret so they're encrypted at rest and masked in logs.
    script_env = [
      { key = "LOGIN_EMAIL",    value = var.login_email,    secret = true },
      { key = "LOGIN_PASSWORD", value = var.login_password, secret = true },
    ]
  })

  failure_threshold = 2

  # SSL and domain expiry still work for browser monitors.
  ssl_expiry_enabled    = true
  ssl_expiry_thresholds = [30, 14, 7]
}

# Wire up a Slack alert channel.
resource "oack_alert_channel" "slack" {
  team_id = oack_team.web.id
  type    = "slack"
  name    = "Web Alerts"
  config = {
    webhook_url = "https://hooks.slack.com/services/T00/B00/xxxx"
  }
}

resource "oack_monitor_alert_channel_link" "login_slack" {
  team_id    = oack_team.web.id
  monitor_id = oack_monitor.login_flow.id
  channel_id = oack_alert_channel.slack.id
}

output "monitor_id" {
  value = oack_monitor.login_flow.id
}
