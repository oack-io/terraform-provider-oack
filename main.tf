terraform {
  required_providers {
    oack = {
      source  = "oack-io/oack"
      version = "~> 0.1"
    }
  }
}

provider "oack" {
  api_key    = var.oack_api_key    # or set OACK_API_KEY
  account_id = var.oack_account_id # or set OACK_ACCOUNT_ID
}

resource "oack_team" "engineering" {
  name = "Engineering"
}

resource "oack_monitor" "api" {
  team_id           = oack_team.engineering.id
  name              = "API Health"
  url               = "https://httpbin.org/status/200"
  check_interval_ms = 30000
}
