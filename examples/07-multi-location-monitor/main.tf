# Example 07: Multi-Location Monitor
#
# Monitor a URL from multiple geographic locations simultaneously.
# The monitor is marked "down" only when at least 2 out of 3 locations fail,
# which eliminates false positives from regional network issues.
#
# Usage:
#   export OACK_API_KEY="oack_acc_xxxxxxxxxxxxxxxxxxxx"
#   export OACK_ACCOUNT_ID="your-account-uuid"
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

# Look up available checker locations.
data "oack_checkers" "all" {}

resource "oack_team" "production" {
  name = "Production"
}

# HTTP monitor checked from Dallas, Amsterdam, and Singapore.
# Marked down only when 2+ locations report failure.
resource "oack_monitor" "api_multi_region" {
  team_id           = oack_team.production.id
  name              = "API — Multi-Region"
  url               = "https://httpbin.org/status/200"
  check_interval_ms = 60000  # 60 seconds
  timeout_ms        = 15000  # 15 seconds
  http_method       = "GET"
  failure_threshold = 2

  # Require at least 2 locations to fail before marking the monitor down.
  # Options: "any", "majority", "all", "at_least_n".
  # aggregate_failure_mode  = "at_least_n"
  # aggregate_failure_count = 2

  ssl_expiry_enabled    = true
  ssl_expiry_thresholds = [30, 14, 7]
}

output "monitor_id" {
  value = oack_monitor.api_multi_region.id
}

output "available_checker_regions" {
  description = "Regions you can use in locations[].checker_region"
  value       = distinct([for c in data.oack_checkers.all.checkers : c.region])
}
