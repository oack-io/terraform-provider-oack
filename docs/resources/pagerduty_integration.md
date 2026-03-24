---
page_title: "oack_pagerduty_integration Resource - terraform-provider-oack"
subcategory: ""
description: |-
  Manages the Oack PagerDuty integration (singleton per account).
---

# oack_pagerduty_integration (Resource)

Manages the Oack PagerDuty integration. This is a singleton resource -- there
can only be one PagerDuty integration per account. It syncs PagerDuty services
into Oack so you can route alerts through PagerDuty escalation policies.

## Example Usage

```hcl
resource "oack_pagerduty_integration" "main" {
  api_key      = var.pagerduty_api_key
  region       = "us"
  service_ids  = ["P1234AB", "P5678CD"]
  sync_enabled = true
}
```

## Argument Reference

- `api_key` - (Required, Sensitive) PagerDuty API key.
- `region` - (Required) PagerDuty region: `us` or `eu`.
- `service_ids` - (Optional) List of PagerDuty service IDs to sync.
- `sync_enabled` - (Optional) Whether automatic sync is enabled. Default: `true`.

## Attribute Reference

- `id` - The UUID of the integration.
- `sync_error` - Last sync error message, if any.
- `last_synced_at` - Last successful sync timestamp (RFC 3339).
- `created_at` - Creation timestamp (RFC 3339).
- `updated_at` - Last update timestamp (RFC 3339).

## Import

The PagerDuty integration is a singleton. The import ID is ignored -- the
provider reads the account's integration directly.

```shell
terraform import oack_pagerduty_integration.main _
```

~> **Note:** The `api_key` attribute will be empty after import. You must
supply it in configuration and run `terraform apply` to reconcile state.
