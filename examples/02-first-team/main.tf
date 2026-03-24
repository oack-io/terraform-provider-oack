# Example 02: First Team
#
# Create a team — the organizational unit that owns monitors, channels,
# and status pages.
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

output "team_id" {
  description = "The UUID of the newly created team"
  value       = oack_team.engineering.id
}
