package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the Oack API HTTP client.
type Client struct {
	BaseURL    string
	APIKey     string
	AccountID  string
	HTTPClient *http.Client
}

// NewClient creates a new Oack API client.
func NewClient(baseURL, apiKey, accountID string) *Client {
	return &Client{
		BaseURL:   baseURL,
		APIKey:    apiKey,
		AccountID: accountID,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body any) ([]byte, int, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

// APIError represents an error response from the Oack API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("oack API error (%d): %s", e.StatusCode, e.Message)
}

func parseError(statusCode int, body []byte) error {
	var errResp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
		return &APIError{StatusCode: statusCode, Message: errResp.Error}
	}
	return &APIError{StatusCode: statusCode, Message: string(body)}
}

// ── Teams ────────────────────────────────────────────────────────────────────

type Team struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (c *Client) CreateTeam(ctx context.Context, name string) (*Team, error) {
	body, status, err := c.do(ctx, http.MethodPost,
		"/api/v1/accounts/"+c.AccountID+"/teams",
		map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var team Team
	if err := json.Unmarshal(body, &team); err != nil {
		return nil, fmt.Errorf("unmarshal team: %w", err)
	}
	return &team, nil
}

func (c *Client) GetTeam(ctx context.Context, teamID string) (*Team, error) {
	body, status, err := c.do(ctx, http.MethodGet,
		"/api/v1/teams/"+teamID, nil)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, nil
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var team Team
	if err := json.Unmarshal(body, &team); err != nil {
		return nil, fmt.Errorf("unmarshal team: %w", err)
	}
	return &team, nil
}

func (c *Client) UpdateTeam(ctx context.Context, teamID, name string) (*Team, error) {
	body, status, err := c.do(ctx, http.MethodPut,
		"/api/v1/teams/"+teamID,
		map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var team Team
	if err := json.Unmarshal(body, &team); err != nil {
		return nil, fmt.Errorf("unmarshal team: %w", err)
	}
	return &team, nil
}

func (c *Client) DeleteTeam(ctx context.Context, teamID string) error {
	body, status, err := c.do(ctx, http.MethodDelete,
		"/api/v1/teams/"+teamID, nil)
	if err != nil {
		return err
	}
	if status >= 300 {
		return parseError(status, body)
	}
	return nil
}

func (c *Client) ListTeams(ctx context.Context) ([]Team, error) {
	body, status, err := c.do(ctx, http.MethodGet,
		"/api/v1/accounts/"+c.AccountID+"/teams", nil)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var teams []Team
	if err := json.Unmarshal(body, &teams); err != nil {
		return nil, fmt.Errorf("unmarshal teams: %w", err)
	}
	return teams, nil
}

// ── Monitors ─────────────────────────────────────────────────────────────────

type Monitor struct {
	ID                      string            `json:"id"`
	TeamID                  string            `json:"team_id"`
	Name                    string            `json:"name"`
	URL                     string            `json:"url"`
	Status                  string            `json:"status"`
	TimeoutMs               int64             `json:"timeout_ms"`
	CheckIntervalMs         int64             `json:"check_interval_ms"`
	HTTPMethod              string            `json:"http_method"`
	HTTPVersion             string            `json:"http_version"`
	Headers                 map[string]string `json:"headers"`
	FollowRedirects         bool              `json:"follow_redirects"`
	AllowedStatusCodes      []string          `json:"allowed_status_codes"`
	FailureThreshold        int               `json:"failure_threshold"`
	LatencyThresholdMs      int               `json:"latency_threshold_ms"`
	SSLExpiryEnabled        bool              `json:"ssl_expiry_enabled"`
	SSLExpiryThresholds     []int             `json:"ssl_expiry_thresholds"`
	DomainExpiryEnabled     bool              `json:"domain_expiry_enabled"`
	DomainExpiryThresholds  []int             `json:"domain_expiry_thresholds"`
	UptimeThresholdGood     float64           `json:"uptime_threshold_good"`
	UptimeThresholdDegraded float64           `json:"uptime_threshold_degraded"`
	UptimeThresholdCritical float64           `json:"uptime_threshold_critical"`
	CheckerRegion           string            `json:"checker_region"`
	CheckerCountry          string            `json:"checker_country"`
	ResolveOverrideIP       string            `json:"resolve_override_ip"`
	HealthStatus            string            `json:"health_status"`
	CreatedAt               string            `json:"created_at"`
	UpdatedAt               string            `json:"updated_at"`
}

type CreateMonitorRequest struct {
	Name                    string            `json:"name"`
	URL                     string            `json:"url"`
	CheckIntervalMs         int64             `json:"check_interval_ms,omitempty"`
	TimeoutMs               int64             `json:"timeout_ms,omitempty"`
	HTTPMethod              string            `json:"http_method,omitempty"`
	HTTPVersion             string            `json:"http_version,omitempty"`
	Headers                 map[string]string `json:"headers,omitempty"`
	FollowRedirects         *bool             `json:"follow_redirects,omitempty"`
	AllowedStatusCodes      []string          `json:"allowed_status_codes,omitempty"`
	FailureThreshold        int               `json:"failure_threshold,omitempty"`
	LatencyThresholdMs      int               `json:"latency_threshold_ms,omitempty"`
	SSLExpiryEnabled        *bool             `json:"ssl_expiry_enabled,omitempty"`
	SSLExpiryThresholds     []int             `json:"ssl_expiry_thresholds,omitempty"`
	DomainExpiryEnabled     *bool             `json:"domain_expiry_enabled,omitempty"`
	DomainExpiryThresholds  []int             `json:"domain_expiry_thresholds,omitempty"`
	UptimeThresholdGood     *float64          `json:"uptime_threshold_good,omitempty"`
	UptimeThresholdDegraded *float64          `json:"uptime_threshold_degraded,omitempty"`
	UptimeThresholdCritical *float64          `json:"uptime_threshold_critical,omitempty"`
	CheckerRegion           string            `json:"checker_region,omitempty"`
	CheckerCountry          string            `json:"checker_country,omitempty"`
	ResolveOverrideIP       string            `json:"resolve_override_ip,omitempty"`
	Status                  string            `json:"status,omitempty"`
}

func (c *Client) CreateMonitor(ctx context.Context, teamID string, req *CreateMonitorRequest) (*Monitor, error) {
	body, status, err := c.do(ctx, http.MethodPost,
		"/api/v1/teams/"+teamID+"/monitors", req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var monitor Monitor
	if err := json.Unmarshal(body, &monitor); err != nil {
		return nil, fmt.Errorf("unmarshal monitor: %w", err)
	}
	return &monitor, nil
}

func (c *Client) GetMonitor(ctx context.Context, teamID, monitorID string) (*Monitor, error) {
	body, status, err := c.do(ctx, http.MethodGet,
		"/api/v1/teams/"+teamID+"/monitors/"+monitorID, nil)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, nil
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var monitor Monitor
	if err := json.Unmarshal(body, &monitor); err != nil {
		return nil, fmt.Errorf("unmarshal monitor: %w", err)
	}
	return &monitor, nil
}

func (c *Client) UpdateMonitor(ctx context.Context, teamID, monitorID string, req *CreateMonitorRequest) (*Monitor, error) {
	body, status, err := c.do(ctx, http.MethodPut,
		"/api/v1/teams/"+teamID+"/monitors/"+monitorID, req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var monitor Monitor
	if err := json.Unmarshal(body, &monitor); err != nil {
		return nil, fmt.Errorf("unmarshal monitor: %w", err)
	}
	return &monitor, nil
}

func (c *Client) DeleteMonitor(ctx context.Context, teamID, monitorID string) error {
	body, status, err := c.do(ctx, http.MethodDelete,
		"/api/v1/teams/"+teamID+"/monitors/"+monitorID, nil)
	if err != nil {
		return err
	}
	if status >= 300 {
		return parseError(status, body)
	}
	return nil
}

// ── Checkers ─────────────────────────────────────────────────────────────────

type Checker struct {
	ID      string `json:"id"`
	Region  string `json:"region"`
	Country string `json:"country"`
	IP      string `json:"ip"`
	ASN     any    `json:"asn"`
	Mode    string `json:"mode"`
	Status  string `json:"status"`
}

func (c *Client) ListCheckers(ctx context.Context) ([]Checker, error) {
	body, status, err := c.do(ctx, http.MethodGet,
		"/api/v1/checkers", nil)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var checkers []Checker
	if err := json.Unmarshal(body, &checkers); err != nil {
		return nil, fmt.Errorf("unmarshal checkers: %w", err)
	}
	return checkers, nil
}

// ── Alert Channels ───────────────────────────────────────────────────────────

type AlertChannel struct {
	ID            string            `json:"id"`
	TeamID        string            `json:"team_id"`
	Type          string            `json:"type"`
	Name          string            `json:"name"`
	Config        map[string]string `json:"config"`
	Enabled       bool              `json:"enabled"`
	EmailVerified bool              `json:"email_verified"`
	Scope         string            `json:"scope"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
}

type CreateAlertChannelRequest struct {
	Type    string            `json:"type"`
	Name    string            `json:"name"`
	Config  map[string]string `json:"config"`
	Enabled *bool             `json:"enabled,omitempty"`
}

func (c *Client) CreateAlertChannel(ctx context.Context, teamID string, req *CreateAlertChannelRequest) (*AlertChannel, error) {
	body, status, err := c.do(ctx, http.MethodPost,
		"/api/v1/teams/"+teamID+"/alert-channels", req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var ch AlertChannel
	if err := json.Unmarshal(body, &ch); err != nil {
		return nil, fmt.Errorf("unmarshal alert channel: %w", err)
	}
	return &ch, nil
}

func (c *Client) GetAlertChannel(ctx context.Context, teamID, channelID string) (*AlertChannel, error) {
	// List all and find — the API doesn't have a get-by-id endpoint for team channels.
	channels, err := c.ListAlertChannels(ctx, teamID)
	if err != nil {
		return nil, err
	}
	for _, ch := range channels {
		if ch.ID == channelID {
			return &ch, nil
		}
	}
	return nil, nil
}

func (c *Client) ListAlertChannels(ctx context.Context, teamID string) ([]AlertChannel, error) {
	body, status, err := c.do(ctx, http.MethodGet,
		"/api/v1/teams/"+teamID+"/alert-channels", nil)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var channels []AlertChannel
	if err := json.Unmarshal(body, &channels); err != nil {
		return nil, fmt.Errorf("unmarshal alert channels: %w", err)
	}
	return channels, nil
}

func (c *Client) UpdateAlertChannel(ctx context.Context, teamID, channelID string, req *CreateAlertChannelRequest) (*AlertChannel, error) {
	body, status, err := c.do(ctx, http.MethodPut,
		"/api/v1/teams/"+teamID+"/alert-channels/"+channelID, req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var ch AlertChannel
	if err := json.Unmarshal(body, &ch); err != nil {
		return nil, fmt.Errorf("unmarshal alert channel: %w", err)
	}
	return &ch, nil
}

func (c *Client) DeleteAlertChannel(ctx context.Context, teamID, channelID string) error {
	body, status, err := c.do(ctx, http.MethodDelete,
		"/api/v1/teams/"+teamID+"/alert-channels/"+channelID, nil)
	if err != nil {
		return err
	}
	if status >= 300 {
		return parseError(status, body)
	}
	return nil
}

// ── Monitor-Channel Links ────────────────────────────────────────────────────

func (c *Client) LinkMonitorChannel(ctx context.Context, teamID, monitorID, channelID string) error {
	body, status, err := c.do(ctx, http.MethodPost,
		"/api/v1/teams/"+teamID+"/monitors/"+monitorID+"/alert-channels/"+channelID, nil)
	if err != nil {
		return err
	}
	if status >= 300 {
		return parseError(status, body)
	}
	return nil
}

func (c *Client) UnlinkMonitorChannel(ctx context.Context, teamID, monitorID, channelID string) error {
	body, status, err := c.do(ctx, http.MethodDelete,
		"/api/v1/teams/"+teamID+"/monitors/"+monitorID+"/alert-channels/"+channelID, nil)
	if err != nil {
		return err
	}
	if status >= 300 {
		return parseError(status, body)
	}
	return nil
}

type MonitorChannelsResponse struct {
	ChannelIDs []string `json:"channel_ids"`
}

func (c *Client) ListMonitorChannels(ctx context.Context, teamID, monitorID string) ([]string, error) {
	body, status, err := c.do(ctx, http.MethodGet,
		"/api/v1/teams/"+teamID+"/monitors/"+monitorID+"/alert-channels", nil)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var resp MonitorChannelsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal monitor channels: %w", err)
	}
	return resp.ChannelIDs, nil
}

// ── Status Pages ─────────────────────────────────────────────────────────────

type StatusPage struct {
	ID                   string  `json:"id"`
	AccountID            string  `json:"account_id"`
	Name                 string  `json:"name"`
	Slug                 string  `json:"slug"`
	Description          string  `json:"description"`
	CustomDomain         *string `json:"custom_domain"`
	HasPassword          bool    `json:"has_password"`
	AllowIframe          bool    `json:"allow_iframe"`
	ShowHistoricalUptime bool    `json:"show_historical_uptime"`
	BrandingLogoURL      *string `json:"branding_logo_url"`
	BrandingFaviconURL   *string `json:"branding_favicon_url"`
	BrandingPrimaryColor *string `json:"branding_primary_color"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}

type CreateStatusPageRequest struct {
	Name                 string  `json:"name"`
	Slug                 string  `json:"slug"`
	Description          string  `json:"description,omitempty"`
	CustomDomain         *string `json:"custom_domain,omitempty"`
	Password             *string `json:"password,omitempty"`
	AllowIframe          *bool   `json:"allow_iframe,omitempty"`
	ShowHistoricalUptime *bool   `json:"show_historical_uptime,omitempty"`
	BrandingLogoURL      *string `json:"branding_logo_url,omitempty"`
	BrandingFaviconURL   *string `json:"branding_favicon_url,omitempty"`
	BrandingPrimaryColor *string `json:"branding_primary_color,omitempty"`
}

func (c *Client) CreateStatusPage(ctx context.Context, req *CreateStatusPageRequest) (*StatusPage, error) {
	body, status, err := c.do(ctx, http.MethodPost,
		"/api/v1/accounts/"+c.AccountID+"/status-pages", req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var sp StatusPage
	if err := json.Unmarshal(body, &sp); err != nil {
		return nil, fmt.Errorf("unmarshal status page: %w", err)
	}
	return &sp, nil
}

func (c *Client) GetStatusPage(ctx context.Context, pageID string) (*StatusPage, error) {
	body, status, err := c.do(ctx, http.MethodGet,
		"/api/v1/accounts/"+c.AccountID+"/status-pages/"+pageID, nil)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, nil
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var sp StatusPage
	if err := json.Unmarshal(body, &sp); err != nil {
		return nil, fmt.Errorf("unmarshal status page: %w", err)
	}
	return &sp, nil
}

func (c *Client) UpdateStatusPage(ctx context.Context, pageID string, req *CreateStatusPageRequest) (*StatusPage, error) {
	body, status, err := c.do(ctx, http.MethodPut,
		"/api/v1/accounts/"+c.AccountID+"/status-pages/"+pageID, req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var sp StatusPage
	if err := json.Unmarshal(body, &sp); err != nil {
		return nil, fmt.Errorf("unmarshal status page: %w", err)
	}
	return &sp, nil
}

func (c *Client) DeleteStatusPage(ctx context.Context, pageID string) error {
	body, status, err := c.do(ctx, http.MethodDelete,
		"/api/v1/accounts/"+c.AccountID+"/status-pages/"+pageID, nil)
	if err != nil {
		return err
	}
	if status >= 300 {
		return parseError(status, body)
	}
	return nil
}

// ── Status Page Component Groups ─────────────────────────────────────────────

type ComponentGroup struct {
	ID           string `json:"id"`
	StatusPageID string `json:"status_page_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Position     int    `json:"position"`
	Collapsed    bool   `json:"collapsed"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type CreateComponentGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Position    int    `json:"position"`
	Collapsed   *bool  `json:"collapsed,omitempty"`
}

func (c *Client) CreateComponentGroup(ctx context.Context, pageID string, req *CreateComponentGroupRequest) (*ComponentGroup, error) {
	body, status, err := c.do(ctx, http.MethodPost,
		"/api/v1/accounts/"+c.AccountID+"/status-pages/"+pageID+"/component-groups", req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var g ComponentGroup
	if err := json.Unmarshal(body, &g); err != nil {
		return nil, fmt.Errorf("unmarshal component group: %w", err)
	}
	return &g, nil
}

func (c *Client) GetComponentGroup(ctx context.Context, pageID, groupID string) (*ComponentGroup, error) {
	// List and find by ID.
	body, status, err := c.do(ctx, http.MethodGet,
		"/api/v1/accounts/"+c.AccountID+"/status-pages/"+pageID+"/component-groups", nil)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var groups []ComponentGroup
	if err := json.Unmarshal(body, &groups); err != nil {
		return nil, fmt.Errorf("unmarshal groups: %w", err)
	}
	for _, g := range groups {
		if g.ID == groupID {
			return &g, nil
		}
	}
	return nil, nil
}

func (c *Client) UpdateComponentGroup(ctx context.Context, pageID, groupID string, req *CreateComponentGroupRequest) (*ComponentGroup, error) {
	body, status, err := c.do(ctx, http.MethodPut,
		"/api/v1/accounts/"+c.AccountID+"/status-pages/"+pageID+"/component-groups/"+groupID, req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var g ComponentGroup
	if err := json.Unmarshal(body, &g); err != nil {
		return nil, fmt.Errorf("unmarshal component group: %w", err)
	}
	return &g, nil
}

func (c *Client) DeleteComponentGroup(ctx context.Context, pageID, groupID string) error {
	body, status, err := c.do(ctx, http.MethodDelete,
		"/api/v1/accounts/"+c.AccountID+"/status-pages/"+pageID+"/component-groups/"+groupID, nil)
	if err != nil {
		return err
	}
	if status >= 300 {
		return parseError(status, body)
	}
	return nil
}

// ── Status Page Components ───────────────────────────────────────────────────

type Component struct {
	ID            string `json:"id"`
	StatusPageID  string `json:"status_page_id"`
	GroupID       string `json:"group_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	DisplayUptime bool   `json:"display_uptime"`
	Position      int    `json:"position"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type CreateComponentRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	GroupID       string `json:"group_id,omitempty"`
	DisplayUptime *bool  `json:"display_uptime,omitempty"`
	Position      int    `json:"position"`
}

func (c *Client) CreateComponent(ctx context.Context, pageID string, req *CreateComponentRequest) (*Component, error) {
	body, status, err := c.do(ctx, http.MethodPost,
		"/api/v1/accounts/"+c.AccountID+"/status-pages/"+pageID+"/components", req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var comp Component
	if err := json.Unmarshal(body, &comp); err != nil {
		return nil, fmt.Errorf("unmarshal component: %w", err)
	}
	return &comp, nil
}

func (c *Client) GetComponent(ctx context.Context, pageID, compID string) (*Component, error) {
	body, status, err := c.do(ctx, http.MethodGet,
		"/api/v1/accounts/"+c.AccountID+"/status-pages/"+pageID+"/components", nil)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var comps []Component
	if err := json.Unmarshal(body, &comps); err != nil {
		return nil, fmt.Errorf("unmarshal components: %w", err)
	}
	for _, comp := range comps {
		if comp.ID == compID {
			return &comp, nil
		}
	}
	return nil, nil
}

func (c *Client) UpdateComponent(ctx context.Context, pageID, compID string, req *CreateComponentRequest) (*Component, error) {
	body, status, err := c.do(ctx, http.MethodPut,
		"/api/v1/accounts/"+c.AccountID+"/status-pages/"+pageID+"/components/"+compID, req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var comp Component
	if err := json.Unmarshal(body, &comp); err != nil {
		return nil, fmt.Errorf("unmarshal component: %w", err)
	}
	return &comp, nil
}

func (c *Client) DeleteComponent(ctx context.Context, pageID, compID string) error {
	body, status, err := c.do(ctx, http.MethodDelete,
		"/api/v1/accounts/"+c.AccountID+"/status-pages/"+pageID+"/components/"+compID, nil)
	if err != nil {
		return err
	}
	if status >= 300 {
		return parseError(status, body)
	}
	return nil
}

// ── Watchdogs ────────────────────────────────────────────────────────────────

type Watchdog struct {
	ID                string `json:"id"`
	ComponentID       string `json:"component_id"`
	MonitorID         string `json:"monitor_id"`
	Severity          string `json:"severity"`
	AutoCreate        bool   `json:"auto_create"`
	AutoResolve       bool   `json:"auto_resolve"`
	NotifySubscribers bool   `json:"notify_subscribers"`
	TemplateID        string `json:"template_id"`
	CreatedAt         string `json:"created_at"`
}

type CreateWatchdogRequest struct {
	MonitorID         string `json:"monitor_id"`
	Severity          string `json:"severity"`
	AutoCreate        *bool  `json:"auto_create,omitempty"`
	AutoResolve       *bool  `json:"auto_resolve,omitempty"`
	NotifySubscribers *bool  `json:"notify_subscribers,omitempty"`
	TemplateID        string `json:"template_id,omitempty"`
}

func (c *Client) CreateWatchdog(ctx context.Context, pageID, compID string, req *CreateWatchdogRequest) (*Watchdog, error) {
	body, status, err := c.do(ctx, http.MethodPost,
		"/api/v1/accounts/"+c.AccountID+"/status-pages/"+pageID+"/components/"+compID+"/watchdogs", req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var w Watchdog
	if err := json.Unmarshal(body, &w); err != nil {
		return nil, fmt.Errorf("unmarshal watchdog: %w", err)
	}
	return &w, nil
}

func (c *Client) GetWatchdog(ctx context.Context, pageID, compID, watchdogID string) (*Watchdog, error) {
	body, status, err := c.do(ctx, http.MethodGet,
		"/api/v1/accounts/"+c.AccountID+"/status-pages/"+pageID+"/components/"+compID+"/watchdogs", nil)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var watchdogs []Watchdog
	if err := json.Unmarshal(body, &watchdogs); err != nil {
		return nil, fmt.Errorf("unmarshal watchdogs: %w", err)
	}
	for _, w := range watchdogs {
		if w.ID == watchdogID {
			return &w, nil
		}
	}
	return nil, nil
}

func (c *Client) DeleteWatchdog(ctx context.Context, pageID, compID, watchdogID string) error {
	body, status, err := c.do(ctx, http.MethodDelete,
		"/api/v1/accounts/"+c.AccountID+"/status-pages/"+pageID+"/components/"+compID+"/watchdogs/"+watchdogID, nil)
	if err != nil {
		return err
	}
	if status >= 300 {
		return parseError(status, body)
	}
	return nil
}

// ── External Links ───────────────────────────────────────────────────────────

type ExternalLink struct {
	ID                string `json:"id"`
	TeamID            string `json:"team_id"`
	Name              string `json:"name"`
	URLTemplate       string `json:"url_template"`
	IconURL           string `json:"icon_url"`
	TimeWindowMinutes int    `json:"time_window_minutes"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type CreateExternalLinkRequest struct {
	Name              string `json:"name"`
	URLTemplate       string `json:"url_template"`
	IconURL           string `json:"icon_url,omitempty"`
	TimeWindowMinutes int    `json:"time_window_minutes"`
}

func (c *Client) CreateExternalLink(ctx context.Context, teamID string, req *CreateExternalLinkRequest) (*ExternalLink, error) {
	body, status, err := c.do(ctx, http.MethodPost,
		"/api/v1/teams/"+teamID+"/external-links", req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var link ExternalLink
	if err := json.Unmarshal(body, &link); err != nil {
		return nil, fmt.Errorf("unmarshal external link: %w", err)
	}
	return &link, nil
}

func (c *Client) GetExternalLink(ctx context.Context, teamID, linkID string) (*ExternalLink, error) {
	body, status, err := c.do(ctx, http.MethodGet,
		"/api/v1/teams/"+teamID+"/external-links", nil)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var links []ExternalLink
	if err := json.Unmarshal(body, &links); err != nil {
		return nil, fmt.Errorf("unmarshal external links: %w", err)
	}
	for _, l := range links {
		if l.ID == linkID {
			return &l, nil
		}
	}
	return nil, nil
}

func (c *Client) UpdateExternalLink(ctx context.Context, teamID, linkID string, req *CreateExternalLinkRequest) (*ExternalLink, error) {
	body, status, err := c.do(ctx, http.MethodPut,
		"/api/v1/teams/"+teamID+"/external-links/"+linkID, req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var link ExternalLink
	if err := json.Unmarshal(body, &link); err != nil {
		return nil, fmt.Errorf("unmarshal external link: %w", err)
	}
	return &link, nil
}

func (c *Client) DeleteExternalLink(ctx context.Context, teamID, linkID string) error {
	body, status, err := c.do(ctx, http.MethodDelete,
		"/api/v1/teams/"+teamID+"/external-links/"+linkID, nil)
	if err != nil {
		return err
	}
	if status >= 300 {
		return parseError(status, body)
	}
	return nil
}

// ── PagerDuty Integration ────────────────────────────────────────────────────

type PDIntegration struct {
	ID           string   `json:"id"`
	AccountID    string   `json:"account_id"`
	APIKey       string   `json:"api_key"`
	Region       string   `json:"region"`
	ServiceIDs   []string `json:"service_ids"`
	SyncEnabled  bool     `json:"sync_enabled"`
	SyncError    string   `json:"sync_error"`
	LastSyncedAt *string  `json:"last_synced_at"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

type CreatePDIntegrationRequest struct {
	APIKey     string   `json:"api_key"`
	Region     string   `json:"region"`
	ServiceIDs []string `json:"service_ids,omitempty"`
}

func (c *Client) CreatePDIntegration(ctx context.Context, req *CreatePDIntegrationRequest) (*PDIntegration, error) {
	body, status, err := c.do(ctx, http.MethodPost,
		"/api/v1/accounts/"+c.AccountID+"/integrations/pagerduty", req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var pd PDIntegration
	if err := json.Unmarshal(body, &pd); err != nil {
		return nil, fmt.Errorf("unmarshal pd integration: %w", err)
	}
	return &pd, nil
}

func (c *Client) GetPDIntegration(ctx context.Context) (*PDIntegration, error) {
	body, status, err := c.do(ctx, http.MethodGet,
		"/api/v1/accounts/"+c.AccountID+"/integrations/pagerduty", nil)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, nil
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var pd PDIntegration
	if err := json.Unmarshal(body, &pd); err != nil {
		return nil, fmt.Errorf("unmarshal pd integration: %w", err)
	}
	return &pd, nil
}

func (c *Client) UpdatePDIntegration(ctx context.Context, req map[string]any) (*PDIntegration, error) {
	body, status, err := c.do(ctx, http.MethodPut,
		"/api/v1/accounts/"+c.AccountID+"/integrations/pagerduty", req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var pd PDIntegration
	if err := json.Unmarshal(body, &pd); err != nil {
		return nil, fmt.Errorf("unmarshal pd integration: %w", err)
	}
	return &pd, nil
}

func (c *Client) DeletePDIntegration(ctx context.Context) error {
	body, status, err := c.do(ctx, http.MethodDelete,
		"/api/v1/accounts/"+c.AccountID+"/integrations/pagerduty", nil)
	if err != nil {
		return err
	}
	if status >= 300 {
		return parseError(status, body)
	}
	return nil
}

// ── Cloudflare Integration ───────────────────────────────────────────────────

type CFIntegration struct {
	ID              string  `json:"id"`
	AccountID       string  `json:"account_id"`
	ZoneID          string  `json:"zone_id"`
	ZoneName        string  `json:"zone_name"`
	APIToken        string  `json:"api_token"`
	Enabled         bool    `json:"enabled"`
	SessionError    string  `json:"session_error"`
	LastConnectedAt *string `json:"last_connected_at"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type CreateCFIntegrationRequest struct {
	ZoneID   string `json:"zone_id"`
	ZoneName string `json:"zone_name"`
	APIToken string `json:"api_token"`
}

func (c *Client) CreateCFIntegration(ctx context.Context, req *CreateCFIntegrationRequest) (*CFIntegration, error) {
	body, status, err := c.do(ctx, http.MethodPost,
		"/api/v1/accounts/"+c.AccountID+"/integrations/cloudflare-zone", req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var cf CFIntegration
	if err := json.Unmarshal(body, &cf); err != nil {
		return nil, fmt.Errorf("unmarshal cf integration: %w", err)
	}
	return &cf, nil
}

func (c *Client) ListCFIntegrations(ctx context.Context) ([]CFIntegration, error) {
	body, status, err := c.do(ctx, http.MethodGet,
		"/api/v1/accounts/"+c.AccountID+"/integrations/cloudflare-zone", nil)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var cfs []CFIntegration
	if err := json.Unmarshal(body, &cfs); err != nil {
		return nil, fmt.Errorf("unmarshal cf integrations: %w", err)
	}
	return cfs, nil
}

func (c *Client) GetCFIntegration(ctx context.Context, cfID string) (*CFIntegration, error) {
	cfs, err := c.ListCFIntegrations(ctx)
	if err != nil {
		return nil, err
	}
	for _, cf := range cfs {
		if cf.ID == cfID {
			return &cf, nil
		}
	}
	return nil, nil
}

func (c *Client) DeleteCFIntegration(ctx context.Context, cfID string) error {
	body, status, err := c.do(ctx, http.MethodDelete,
		"/api/v1/accounts/"+c.AccountID+"/integrations/cloudflare-zone/"+cfID, nil)
	if err != nil {
		return err
	}
	if status >= 300 {
		return parseError(status, body)
	}
	return nil
}

// ── Team API Keys ────────────────────────────────────────────────────────────

type TeamAPIKey struct {
	ID        string  `json:"id"`
	TeamID    string  `json:"team_id"`
	Name      string  `json:"name"`
	KeyPrefix string  `json:"key_prefix"`
	CreatedBy string  `json:"created_by"`
	ExpiresAt *string `json:"expires_at"`
	CreatedAt string  `json:"created_at"`
}

type CreateTeamAPIKeyRequest struct {
	Name      string  `json:"name"`
	ExpiresAt *string `json:"expires_at,omitempty"`
}

type CreateTeamAPIKeyResponse struct {
	Key    string     `json:"key"`
	APIKey TeamAPIKey `json:"api_key"`
}

func (c *Client) CreateTeamAPIKey(ctx context.Context, teamID string, req *CreateTeamAPIKeyRequest) (*CreateTeamAPIKeyResponse, error) {
	body, status, err := c.do(ctx, http.MethodPost,
		"/api/v1/teams/"+teamID+"/api-keys", req)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var resp CreateTeamAPIKeyResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal team api key: %w", err)
	}
	return &resp, nil
}

func (c *Client) ListTeamAPIKeys(ctx context.Context, teamID string) ([]TeamAPIKey, error) {
	body, status, err := c.do(ctx, http.MethodGet,
		"/api/v1/teams/"+teamID+"/api-keys", nil)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, parseError(status, body)
	}
	var keys []TeamAPIKey
	if err := json.Unmarshal(body, &keys); err != nil {
		return nil, fmt.Errorf("unmarshal team api keys: %w", err)
	}
	return keys, nil
}

func (c *Client) GetTeamAPIKey(ctx context.Context, teamID, keyID string) (*TeamAPIKey, error) {
	keys, err := c.ListTeamAPIKeys(ctx, teamID)
	if err != nil {
		return nil, err
	}
	for _, k := range keys {
		if k.ID == keyID {
			return &k, nil
		}
	}
	return nil, nil
}

func (c *Client) DeleteTeamAPIKey(ctx context.Context, teamID, keyID string) error {
	body, status, err := c.do(ctx, http.MethodDelete,
		"/api/v1/teams/"+teamID+"/api-keys/"+keyID, nil)
	if err != nil {
		return err
	}
	if status >= 300 {
		return parseError(status, body)
	}
	return nil
}
