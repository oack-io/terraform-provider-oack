# Example 01: Hello World
#
# Verify your API key works by reading the list of available checker nodes.
# This makes no changes to your account — it only reads data.
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

provider "oack" {
  # Credentials are read from OACK_API_KEY and OACK_ACCOUNT_ID env vars.
  # You can also set them explicitly:
  # api_key    = "oack_acc_xxxxxxxxxxxxxxxxxxxx"
  # account_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
}

# Read all available checker nodes.
data "oack_checkers" "all" {}

output "checker_count" {
  description = "Number of available checker nodes"
  value       = length(data.oack_checkers.all.checkers)
}

output "checker_regions" {
  description = "Regions where checkers are deployed"
  value       = distinct([for c in data.oack_checkers.all.checkers : c.region])
}
