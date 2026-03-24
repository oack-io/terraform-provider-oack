# Example 05: Status Page
#
# Create a public status page with a component group, a component, and a
# watchdog that automatically creates incidents when the monitor goes down.
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

# --- Team ---

resource "oack_team" "engineering" {
  name = "Engineering"
}

# --- Monitor ---

resource "oack_monitor" "api" {
  team_id           = oack_team.engineering.id
  name              = "API"
  url               = "https://httpbin.org/status/200"
  check_interval_ms = 30000
  failure_threshold = 3
}

# --- Status Page ---

resource "oack_status_page" "public" {
  name                   = "Example Status"
  slug                   = "example-status"
  description            = "Current operational status of Example services."
  show_historical_uptime = true
}

# --- Component Group ---

resource "oack_status_page_component_group" "backend" {
  status_page_id = oack_status_page.public.id
  name           = "Backend Services"
  description    = "Core API and backend infrastructure"
  position       = 0
}

# --- Component ---

resource "oack_status_page_component" "api" {
  status_page_id = oack_status_page.public.id
  group_id       = oack_status_page_component_group.backend.id
  name           = "REST API"
  description    = "Public-facing REST API"
  display_uptime = true
  position       = 0
}

# --- Watchdog ---
# Automatically creates and resolves incidents based on monitor health.

resource "oack_status_page_watchdog" "api" {
  status_page_id     = oack_status_page.public.id
  component_id       = oack_status_page_component.api.id
  monitor_id         = oack_monitor.api.id
  severity           = "major"
  auto_create        = true
  auto_resolve       = true
  notify_subscribers = true
}

output "status_page_url" {
  description = "Public URL of the status page"
  value       = "https://status.oack.io/${oack_status_page.public.slug}"
}
