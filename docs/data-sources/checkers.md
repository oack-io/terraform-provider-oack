---
page_title: "oack_checkers Data Source - terraform-provider-oack"
subcategory: ""
description: |-
  List available Oack checker nodes.
---

# oack_checkers (Data Source)

Lists all available Oack checker nodes. Checker nodes are the distributed
probes that perform uptime checks from different regions around the world.

This data source is useful for verifying your API credentials work and for
discovering which regions are available for monitor configuration.

## Example Usage

```hcl
data "oack_checkers" "all" {}

output "checker_count" {
  value = length(data.oack_checkers.all.checkers)
}

output "regions" {
  value = distinct([for c in data.oack_checkers.all.checkers : c.region])
}
```

## Argument Reference

This data source takes no arguments.

## Attribute Reference

- `checkers` - List of checker nodes. Each checker has the following attributes:
  - `id` - Checker UUID.
  - `region` - Checker region (e.g. `us-east`, `eu-west`, `ap-southeast`).
  - `country` - Country code where the checker is located.
  - `ip` - IP address of the checker node.
  - `asn` - Autonomous System Number of the checker.
  - `mode` - Checker mode.
  - `status` - Current checker status.
