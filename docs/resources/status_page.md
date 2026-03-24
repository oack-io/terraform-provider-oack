---
page_title: "oack_status_page Resource - terraform-provider-oack"
subcategory: ""
description: |-
  Manages an Oack status page.
---

# oack_status_page (Resource)

Manages an Oack status page. Status pages provide a public or private view of
your service health. They support custom domains, branding, password protection,
and iframe embedding.

## Example Usage

```hcl
resource "oack_status_page" "public" {
  name                   = "Acme Status"
  slug                   = "acme-status"
  description            = "Current operational status of Acme services."
  show_historical_uptime = true
}
```

### With custom branding and password

```hcl
resource "oack_status_page" "internal" {
  name                   = "Internal Status"
  slug                   = "internal"
  description            = "Internal service health dashboard."
  custom_domain          = "status.internal.example.com"
  password               = var.status_page_password
  show_historical_uptime = true
  allow_iframe           = true
  branding_logo_url      = "https://cdn.example.com/logo.svg"
  branding_favicon_url   = "https://cdn.example.com/favicon.ico"
  branding_primary_color = "#4F46E5"
}
```

## Argument Reference

- `name` - (Required) Status page display name.
- `slug` - (Required, Forces new resource) URL slug for the status page.
- `description` - (Optional) Status page description.
- `custom_domain` - (Optional) Custom domain for the status page.
- `password` - (Optional, Sensitive) Password to protect the status page. Write-only; not returned by the API on read.
- `allow_iframe` - (Optional) Whether to allow embedding in iframes. Default: `false`.
- `show_historical_uptime` - (Optional) Whether to show historical uptime data. Default: `true`.
- `branding_logo_url` - (Optional) URL for a custom logo.
- `branding_favicon_url` - (Optional) URL for a custom favicon.
- `branding_primary_color` - (Optional) Primary color hex code for branding (e.g. `#4F46E5`).

## Attribute Reference

- `id` - The UUID of the status page.
- `has_password` - Whether the status page is password-protected.
- `created_at` - Creation timestamp (RFC 3339).
- `updated_at` - Last update timestamp (RFC 3339).

## Import

```shell
terraform import oack_status_page.example <status_page_id>
```
