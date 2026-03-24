---
page_title: "oack_team_api_key Resource - terraform-provider-oack"
subcategory: ""
description: |-
  Manages an Oack team API key.
---

# oack_team_api_key (Resource)

Manages an Oack team API key. Team API keys are scoped to a single team and
can be used for CI/CD integrations, deploy event reporting, and programmatic
access. The plaintext key is only available at creation time and is stored in
Terraform state.

## Example Usage

```hcl
resource "oack_team_api_key" "cicd" {
  team_id = oack_team.engineering.id
  name    = "CI/CD Deploy Events"
}
```

### With expiration

```hcl
resource "oack_team_api_key" "temp" {
  team_id    = oack_team.engineering.id
  name       = "Temporary Access"
  expires_at = "2025-12-31T23:59:59Z"
}
```

## Argument Reference

- `team_id` - (Required, Forces new resource) Team UUID.
- `name` - (Required) API key display name.
- `expires_at` - (Optional) Expiration timestamp (RFC 3339). Leave empty for no expiration.

## Attribute Reference

- `id` - The UUID of the API key.
- `key` - (Sensitive) The plaintext API key. Only available at creation time. Stored in state.
- `key_prefix` - Visible prefix of the API key (e.g. `oack_team_abc...`).
- `created_at` - Creation timestamp (RFC 3339).

## Import

Import uses the format `team_id/key_id`:

```shell
terraform import oack_team_api_key.example <team_id>/<key_id>
```

~> **Note:** The `key` attribute will be null after import because the API
does not return plaintext keys after creation.
