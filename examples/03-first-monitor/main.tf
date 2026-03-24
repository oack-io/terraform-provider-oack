# Example 03: First Monitor
#
# Create a team and an HTTP monitor that checks a URL every 60 seconds.
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

resource "oack_team" "engineering" {
  name = "Engineering"
}

resource "oack_monitor" "api" {
  team_id           = oack_team.engineering.id
  name              = "API Health Check"
  url               = "https://httpbin.org/status/200"
  check_interval_ms = 60000 # 60 seconds
  timeout_ms        = 10000 # 10 seconds
  http_method       = "GET"
  failure_threshold = 3

  # Alert when SSL certificate is expiring soon.
  ssl_expiry_enabled    = true
  ssl_expiry_thresholds = [30, 14, 7]

  # Alert when domain is expiring soon.
  domain_expiry_enabled    = true
  domain_expiry_thresholds = [30, 14, 7]
}

output "monitor_id" {
  description = "The UUID of the monitor"
  value       = oack_monitor.api.id
}

output "health_status" {
  description = "Current health status of the monitor"
  value       = oack_monitor.api.health_status
}
