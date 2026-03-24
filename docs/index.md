---
page_title: "Provider: Oack"
subcategory: ""
description: |-
  Terraform provider for Oack uptime monitoring.
---

# Oack Provider

The Oack provider lets you manage [Oack](https://oack.io) uptime monitoring
infrastructure as code. Create monitors, configure alert channels, build status
pages, and wire everything together with Terraform.

## Authentication

The provider requires an account-level API key and an account ID. Both can be
set in the provider block or via environment variables.

| Provider Attribute | Environment Variable | Required | Description                                   |
|--------------------|----------------------|----------|-----------------------------------------------|
| `api_key`          | `OACK_API_KEY`       | Yes      | Account API key (starts with `oack_acc_...`)  |
| `account_id`       | `OACK_ACCOUNT_ID`    | Yes      | Account UUID                                  |
| `api_url`          | `OACK_API_URL`       | No       | API base URL (default: `https://api.oack.io`) |

## Example Usage

```hcl
terraform {
  required_providers {
    oack = {
      source  = "oack-io/oack"
      version = "~> 0.1"
    }
  }
}

# Configure using environment variables.
provider "oack" {}

# Or configure explicitly.
provider "oack" {
  api_key    = var.oack_api_key
  account_id = var.oack_account_id
}
```

## Argument Reference

- `api_key` - (Optional) Account API key. Falls back to `OACK_API_KEY` env var.
- `account_id` - (Optional) Account ID. Falls back to `OACK_ACCOUNT_ID` env var.
- `api_url` - (Optional) API base URL. Falls back to `OACK_API_URL` env var.
  Default: `https://api.oack.io`.
