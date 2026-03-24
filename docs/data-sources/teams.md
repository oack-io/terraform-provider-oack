---
page_title: "oack_teams Data Source - terraform-provider-oack"
subcategory: ""
description: |-
  List all teams in the account.
---

# oack_teams (Data Source)

Lists all teams in the account. Useful for discovering existing teams or for
referencing team IDs in other resources without hardcoding UUIDs.

## Example Usage

```hcl
data "oack_teams" "all" {}

output "team_names" {
  value = [for t in data.oack_teams.all.teams : t.name]
}
```

### Look up a team by name

```hcl
data "oack_teams" "all" {}

locals {
  engineering_team = [
    for t in data.oack_teams.all.teams : t if t.name == "Engineering"
  ][0]
}

resource "oack_monitor" "api" {
  team_id           = local.engineering_team.id
  name              = "API"
  url               = "https://api.example.com/healthz"
  check_interval_ms = 60000
}
```

## Argument Reference

This data source takes no arguments.

## Attribute Reference

- `teams` - List of teams. Each team has the following attributes:
  - `id` - Team UUID.
  - `name` - Team display name.
  - `created_at` - Creation timestamp (RFC 3339).
  - `updated_at` - Last update timestamp (RFC 3339).
