# Example 08: Simple Browser Monitor
#
# Monitor a web page using a real Chromium browser. The checker loads the page,
# captures Web Vitals (LCP, FCP, CLS, TTFB), takes a screenshot, and records a
# HAR file. The monitor is marked "down" if the page returns an error status,
# exceeds the timeout, or has too many console/resource errors.
#
# Usage:
#   export OACK_API_KEY="oack_acc_xxxxxxxxxxxxxxxxxxxx"
#   export OACK_ACCOUNT_ID="your-account-uuid"
#   terraform init && terraform apply

terraform {
  required_providers {
    oack = {
      source  = "oack-io/oack"
      version = "~> 0.2"
    }
  }
}

provider "oack" {}

resource "oack_team" "web" {
  name = "Web Team"
}

# Browser monitor — loads the page in headless Chromium.
resource "oack_monitor" "homepage" {
  team_id           = oack_team.web.id
  name              = "Homepage — Browser Check"
  url               = "https://example.com"
  type              = "browser"
  check_interval_ms = 300000 # 5 minutes (browser checks are heavier)
  timeout_ms        = 30000  # 30 seconds

  # Browser-specific configuration.
  # Mode "simple" just loads the URL — no scripting needed.
  browser_config_json = jsonencode({
    mode                   = "simple"
    screenshot_enabled     = true
    screenshot_full_page   = false
    viewport_width         = 1920
    viewport_height        = 1080
    wait_until             = "load" # "load", "domcontentloaded", or "networkidle"
    extra_wait_ms          = 0
    console_error_threshold  = 0 # 0 = ignore console errors
    resource_error_threshold = 5 # fail if > 5 resource errors (broken images, etc.)
  })

  failure_threshold = 2
}

output "monitor_id" {
  value = oack_monitor.homepage.id
}

output "monitor_type" {
  value = oack_monitor.homepage.type
}
