---
page_title: "oack_monitor Resource - terraform-provider-oack"
subcategory: ""
description: |-
  Manages an Oack uptime monitor.
---

# oack_monitor (Resource)

Manages an Oack HTTP/HTTPS uptime monitor. Monitors periodically check a URL
from distributed checker nodes and report health status.

## Example Usage

```hcl
resource "oack_monitor" "api" {
  team_id           = oack_team.engineering.id
  name              = "API Health Check"
  url               = "https://httpbin.org/status/200"
  check_interval_ms = 30000
  timeout_ms        = 10000
  http_method       = "GET"
  failure_threshold = 3

  follow_redirects     = true
  allowed_status_codes = ["2xx"]

  ssl_expiry_enabled       = true
  ssl_expiry_thresholds    = [30, 14, 7]
  domain_expiry_enabled    = true
  domain_expiry_thresholds = [30, 14]

  uptime_threshold_good     = 99.9
  uptime_threshold_degraded = 99.0
  uptime_threshold_critical = 95.0
}
```

## Argument Reference

- `team_id` - (Required, Forces new resource) Team UUID.
- `name` - (Required) Monitor display name.
- `url` - (Required) URL to monitor.
- `status` - (Optional) Monitor status: `active` or `paused`. Default: `active`.
- `check_interval_ms` - (Optional) Check interval in milliseconds (minimum 30000). Default: `60000`.
- `timeout_ms` - (Optional) Request timeout in milliseconds. Default: `10000`.
- `http_method` - (Optional) HTTP method: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`. Default: `GET`.
- `http_version` - (Optional) HTTP version: empty string (auto), `1.1`, or `2`. Default: `""` (auto).
- `headers` - (Optional) Map of custom request headers.
- `follow_redirects` - (Optional) Whether to follow HTTP redirects. Default: `true`.
- `allowed_status_codes` - (Optional) List of allowed status codes (e.g. `["2xx", "200", "404"]`).
- `failure_threshold` - (Optional) Number of consecutive failures before marking the monitor as down. Default: `3`.
- `latency_threshold_ms` - (Optional) Latency threshold in milliseconds. `0` disables the check. Default: `0`.
- `ssl_expiry_enabled` - (Optional) Whether to monitor SSL certificate expiry. Default: `true`.
- `ssl_expiry_thresholds` - (Optional) List of day counts before SSL expiry to trigger alerts (e.g. `[30, 14, 7]`).
- `domain_expiry_enabled` - (Optional) Whether to monitor domain expiry. Default: `true`.
- `domain_expiry_thresholds` - (Optional) List of day counts before domain expiry to trigger alerts (e.g. `[30, 14]`).
- `uptime_threshold_good` - (Optional) Uptime percentage for "good" status. Default: `99.9`.
- `uptime_threshold_degraded` - (Optional) Uptime percentage for "degraded" status. Default: `99.0`.
- `uptime_threshold_critical` - (Optional) Uptime percentage for "critical" status. Default: `95.0`.
- `checker_region` - (Optional) Preferred checker region. Default: `""` (any).
- `checker_country` - (Optional) Preferred checker country. Default: `""` (any).
- `resolve_override_ip` - (Optional) Override DNS resolution with a specific IPv4/IPv6 address. Default: `""`.

## Attribute Reference

- `id` - The UUID of the monitor.
- `health_status` - Current health status (`up` or `down`). Read-only.
- `created_at` - Creation timestamp (RFC 3339).
- `updated_at` - Last update timestamp (RFC 3339).

## Import

Import uses the format `team_id/monitor_id`:

```shell
terraform import oack_monitor.example <team_id>/<monitor_id>
```
