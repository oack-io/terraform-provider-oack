variable "slack_webhook_url" {
  description = "Slack incoming webhook URL for alert notifications"
  type        = string
  sensitive   = true
}

variable "api_monitor_url" {
  description = "URL to monitor for the API health check"
  type        = string
  default     = "https://httpbin.org/status/200"
}

variable "web_monitor_url" {
  description = "URL to monitor for the website"
  type        = string
  default     = "https://httpbin.org/html"
}

variable "grafana_base_url" {
  description = "Base URL for the Grafana instance (no trailing slash)"
  type        = string
  default     = "https://grafana.example.com"
}
