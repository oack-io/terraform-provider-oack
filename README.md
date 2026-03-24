# Terraform Provider for Oack

The Oack provider enables [Terraform](https://www.terraform.io) to manage
[Oack](https://oack.io) uptime monitoring infrastructure as code: monitors,
alert channels, status pages, and more.

## Quick Start

```hcl
terraform {
  required_providers {
    oack = {
      source  = "oack-io/oack"
      version = "~> 0.1"
    }
  }
}

provider "oack" {
  api_key    = var.oack_api_key   # or set OACK_API_KEY
  account_id = var.oack_account_id # or set OACK_ACCOUNT_ID
}

resource "oack_team" "engineering" {
  name = "Engineering"
}

resource "oack_monitor" "api" {
  team_id          = oack_team.engineering.id
  name             = "API Health"
  url              = "https://api.example.com/healthz"
  check_interval_ms = 30000
}
```

```shell
terraform init
terraform plan
terraform apply
```

## Authentication

The provider requires an **account-level API key** and an **account ID**.
Values can be supplied in the provider block or via environment variables:

| Provider Attribute | Environment Variable | Description                                    |
|--------------------|----------------------|------------------------------------------------|
| `api_key`          | `OACK_API_KEY`       | Account API key (starts with `oack_acc_...`)   |
| `account_id`       | `OACK_ACCOUNT_ID`    | Account UUID                                   |
| `api_url`          | `OACK_API_URL`       | API base URL (default: `https://api.oack.io`)  |

Environment variables take effect when the corresponding provider attribute is
not set. This lets you keep secrets out of your HCL files:

```shell
export OACK_API_KEY="oack_acc_xxxxxxxxxxxxxxxxxxxx"
export OACK_ACCOUNT_ID="aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
terraform plan
```

## Resources and Data Sources

### Resources (11)

| Resource                                | Description                                           |
|-----------------------------------------|-------------------------------------------------------|
| `oack_team`                             | Team (organizational unit for monitors and channels)  |
| `oack_monitor`                          | HTTP/HTTPS uptime monitor                             |
| `oack_alert_channel`                    | Alert channel (Slack, email, webhook, Telegram, etc.) |
| `oack_monitor_alert_channel_link`       | Link a monitor to an alert channel                    |
| `oack_status_page`                      | Public or private status page                         |
| `oack_status_page_component_group`      | Logical group of components on a status page          |
| `oack_status_page_component`            | Individual component on a status page                 |
| `oack_status_page_watchdog`             | Auto-incident creation from monitor health            |
| `oack_external_link`                    | External link (e.g. Grafana dashboard)                |
| `oack_pagerduty_integration`            | PagerDuty integration (singleton per account)         |
| `oack_team_api_key`                     | Team-scoped API key                                   |

### Data Sources (2)

| Data Source      | Description                          |
|------------------|--------------------------------------|
| `oack_checkers`  | List available checker nodes         |
| `oack_teams`     | List all teams in the account        |

## Examples

The [`examples/`](examples/) directory contains progressive configurations,
each building on the previous one:

| Directory               | What it demonstrates                                    |
|-------------------------|---------------------------------------------------------|
| `01-hello-world`        | Verify credentials with a data source read              |
| `02-first-team`         | Create your first team                                  |
| `03-first-monitor`      | Add an HTTP monitor to the team                         |
| `04-alerting`           | Set up Slack and email alert channels, link to monitor  |
| `05-status-page`        | Full status page with component groups and watchdogs    |
| `06-full-stack`         | Production-ready setup with everything wired together   |

## Development

### Build

```shell
make build
```

### Install Locally

```shell
make install
```

This places the binary into
`~/.terraform.d/plugins/registry.terraform.io/oack-io/oack/0.1.0/darwin_arm64/`.

### Run Acceptance Tests

```shell
export OACK_API_KEY="oack_acc_test_key"
export OACK_ACCOUNT_ID="test-account-id"
make testacc
```

### Project Structure

```
.
├── main.go                        # Provider entry point
├── internal/
│   ├── provider/provider.go       # Provider config and registration
│   ├── client/client.go           # HTTP API client
│   ├── resources/                 # All 11 resources
│   └── datasources/               # Both data sources
├── examples/                      # Progressive HCL examples
├── docs/                          # Terraform Registry documentation
└── GNUmakefile                    # Build targets
```

## License

MPL-2.0. See [LICENSE](LICENSE) for details.
