terraform {
  required_providers {
    oack = {
      source = "oack-io/oack"
    }
  }
}

provider "oack" {
  # Set via environment variables:
  # OACK_API_KEY    = "oack_acc_..."
  # OACK_ACCOUNT_ID = "..."
  # OACK_API_URL    = "http://localhost:8080"  (optional, for local dev)
}

# List available checkers
data "oack_checkers" "all" {}

# Create a team
resource "oack_team" "prod" {
  name = "Production"
}

# Create a monitor
resource "oack_monitor" "api" {
  team_id           = oack_team.prod.id
  name              = "API Health"
  url               = "https://api.example.com/health"
  check_interval_ms = 60000
  timeout_ms        = 10000
  failure_threshold = 3

  allowed_status_codes = ["2xx"]

  ssl_expiry_enabled    = true
  ssl_expiry_thresholds = [30, 14, 7, 1]
}

output "team_id" {
  value = oack_team.prod.id
}

output "monitor_id" {
  value = oack_monitor.api.id
}

output "checker_count" {
  value = length(data.oack_checkers.all.checkers)
}
