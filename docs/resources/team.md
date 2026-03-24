---
page_title: "oack_team Resource - terraform-provider-oack"
subcategory: ""
description: |-
  Manages an Oack team.
---

# oack_team (Resource)

Manages an Oack team. Teams are the organizational unit that owns monitors,
alert channels, external links, and API keys.

## Example Usage

```hcl
resource "oack_team" "engineering" {
  name = "Engineering"
}
```

## Argument Reference

- `name` - (Required) Team display name.

## Attribute Reference

- `id` - The UUID of the team.
- `created_at` - Creation timestamp (RFC 3339).
- `updated_at` - Last update timestamp (RFC 3339).

## Import

```shell
terraform import oack_team.example <team_id>
```
