# Example 06: Full Stack
#
# Production-ready configuration with everything wired together:
# - Team
# - Two monitors (API + website)
# - Slack and email alert channels linked to both monitors
# - Status page with component groups, components, and watchdogs
# - External Grafana dashboard link
# - Team API key for CI/CD deploy events
#
# Usage:
#   export OACK_API_KEY="oack_acc_xxxxxxxxxxxxxxxxxxxx"
#   export OACK_ACCOUNT_ID="your-account-uuid"
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

# ─── Team ─────────────────────────────────────────────────────────────────────

resource "oack_team" "engineering" {
  name = "Engineering"
}

# ─── Monitors ─────────────────────────────────────────────────────────────────

resource "oack_monitor" "api" {
  team_id           = oack_team.engineering.id
  name              = "API"
  url               = var.api_monitor_url
  check_interval_ms = 30000
  timeout_ms        = 10000
  http_method       = "GET"
  failure_threshold = 3

  ssl_expiry_enabled       = true
  ssl_expiry_thresholds    = [30, 14, 7]
  domain_expiry_enabled    = true
  domain_expiry_thresholds = [30, 14]

  uptime_threshold_good     = 99.9
  uptime_threshold_degraded = 99.0
  uptime_threshold_critical = 95.0
}

resource "oack_monitor" "website" {
  team_id           = oack_team.engineering.id
  name              = "Website"
  url               = var.web_monitor_url
  check_interval_ms = 60000
  timeout_ms        = 15000
  http_method       = "GET"
  failure_threshold = 3

  follow_redirects     = true
  allowed_status_codes = ["2xx", "3xx"]

  ssl_expiry_enabled       = true
  ssl_expiry_thresholds    = [30, 14, 7]
  domain_expiry_enabled    = true
  domain_expiry_thresholds = [30, 14]
}

# ─── Alert Channels ──────────────────────────────────────────────────────────

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

# ─── Monitor <-> Channel Links ───────────────────────────────────────────────

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

resource "oack_monitor_alert_channel_link" "website_slack" {
  team_id    = oack_team.engineering.id
  monitor_id = oack_monitor.website.id
  channel_id = oack_alert_channel.slack.id
}

resource "oack_monitor_alert_channel_link" "website_email" {
  team_id    = oack_team.engineering.id
  monitor_id = oack_monitor.website.id
  channel_id = oack_alert_channel.email.id
}

# ─── Status Page ──────────────────────────────────────────────────────────────

resource "oack_status_page" "public" {
  name                   = "Example Status"
  slug                   = "example-status"
  description            = "Current operational status of Example services."
  show_historical_uptime = true
}

# ─── Component Groups ────────────────────────────────────────────────────────

resource "oack_status_page_component_group" "backend" {
  status_page_id = oack_status_page.public.id
  name           = "Backend Services"
  description    = "Core API and backend infrastructure"
  position       = 0
}

resource "oack_status_page_component_group" "frontend" {
  status_page_id = oack_status_page.public.id
  name           = "Frontend"
  description    = "Web application and marketing site"
  position       = 1
}

# ─── Components ───────────────────────────────────────────────────────────────

resource "oack_status_page_component" "api" {
  status_page_id = oack_status_page.public.id
  group_id       = oack_status_page_component_group.backend.id
  name           = "REST API"
  description    = "Public-facing REST API"
  display_uptime = true
  position       = 0
}

resource "oack_status_page_component" "website" {
  status_page_id = oack_status_page.public.id
  group_id       = oack_status_page_component_group.frontend.id
  name           = "Website"
  description    = "Marketing website and documentation"
  display_uptime = true
  position       = 0
}

# ─── Watchdogs ────────────────────────────────────────────────────────────────

resource "oack_status_page_watchdog" "api" {
  status_page_id     = oack_status_page.public.id
  component_id       = oack_status_page_component.api.id
  monitor_id         = oack_monitor.api.id
  severity           = "major"
  auto_create        = true
  auto_resolve       = true
  notify_subscribers = true
}

resource "oack_status_page_watchdog" "website" {
  status_page_id     = oack_status_page.public.id
  component_id       = oack_status_page_component.website.id
  monitor_id         = oack_monitor.website.id
  severity           = "minor"
  auto_create        = true
  auto_resolve       = true
  notify_subscribers = true
}

# ─── External Link (Grafana) ─────────────────────────────────────────────────

resource "oack_external_link" "grafana" {
  team_id            = oack_team.engineering.id
  name               = "Grafana Dashboard"
  url_template       = "${var.grafana_base_url}/d/uptime?from=now-{{.TimeWindow}}&to=now"
  time_window_minutes = 60
}

# ─── Team API Key (for CI/CD) ────────────────────────────────────────────────

resource "oack_team_api_key" "cicd" {
  team_id = oack_team.engineering.id
  name    = "CI/CD Deploy Events"
}
