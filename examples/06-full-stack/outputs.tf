output "team_id" {
  description = "Team UUID"
  value       = oack_team.engineering.id
}

output "api_monitor_id" {
  description = "API monitor UUID"
  value       = oack_monitor.api.id
}

output "website_monitor_id" {
  description = "Website monitor UUID"
  value       = oack_monitor.website.id
}

output "status_page_url" {
  description = "Public URL of the status page"
  value       = "https://status.oack.io/${oack_status_page.public.slug}"
}

output "cicd_api_key" {
  description = "API key for CI/CD deploy events (only visible at creation time)"
  value       = oack_team_api_key.cicd.key
  sensitive   = true
}
