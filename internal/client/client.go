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
	defer resp.Body.Close()

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
	ID                    string            `json:"id"`
	TeamID                string            `json:"team_id"`
	Name                  string            `json:"name"`
	URL                   string            `json:"url"`
	Status                string            `json:"status"`
	TimeoutMs             int64             `json:"timeout_ms"`
	CheckIntervalMs       int64             `json:"check_interval_ms"`
	HTTPMethod            string            `json:"http_method"`
	HTTPVersion           string            `json:"http_version"`
	Headers               map[string]string `json:"headers"`
	FollowRedirects       bool              `json:"follow_redirects"`
	AllowedStatusCodes    []string          `json:"allowed_status_codes"`
	FailureThreshold      int               `json:"failure_threshold"`
	LatencyThresholdMs    int               `json:"latency_threshold_ms"`
	SSLExpiryEnabled      bool              `json:"ssl_expiry_enabled"`
	SSLExpiryThresholds   []int             `json:"ssl_expiry_thresholds"`
	DomainExpiryEnabled   bool              `json:"domain_expiry_enabled"`
	DomainExpiryThresholds []int            `json:"domain_expiry_thresholds"`
	UptimeThresholdGood   float64           `json:"uptime_threshold_good"`
	UptimeThresholdDegraded float64         `json:"uptime_threshold_degraded"`
	UptimeThresholdCritical float64         `json:"uptime_threshold_critical"`
	CheckerRegion         string            `json:"checker_region"`
	CheckerCountry        string            `json:"checker_country"`
	ResolveOverrideIP     string            `json:"resolve_override_ip"`
	HealthStatus          string            `json:"health_status"`
	CreatedAt             string            `json:"created_at"`
	UpdatedAt             string            `json:"updated_at"`
}

type CreateMonitorRequest struct {
	Name                   string            `json:"name"`
	URL                    string            `json:"url"`
	CheckIntervalMs        int64             `json:"check_interval_ms,omitempty"`
	TimeoutMs              int64             `json:"timeout_ms,omitempty"`
	HTTPMethod             string            `json:"http_method,omitempty"`
	HTTPVersion            string            `json:"http_version,omitempty"`
	Headers                map[string]string `json:"headers,omitempty"`
	FollowRedirects        *bool             `json:"follow_redirects,omitempty"`
	AllowedStatusCodes     []string          `json:"allowed_status_codes,omitempty"`
	FailureThreshold       int               `json:"failure_threshold,omitempty"`
	LatencyThresholdMs     int               `json:"latency_threshold_ms,omitempty"`
	SSLExpiryEnabled       *bool             `json:"ssl_expiry_enabled,omitempty"`
	SSLExpiryThresholds    []int             `json:"ssl_expiry_thresholds,omitempty"`
	DomainExpiryEnabled    *bool             `json:"domain_expiry_enabled,omitempty"`
	DomainExpiryThresholds []int             `json:"domain_expiry_thresholds,omitempty"`
	UptimeThresholdGood    *float64          `json:"uptime_threshold_good,omitempty"`
	UptimeThresholdDegraded *float64         `json:"uptime_threshold_degraded,omitempty"`
	UptimeThresholdCritical *float64         `json:"uptime_threshold_critical,omitempty"`
	CheckerRegion          string            `json:"checker_region,omitempty"`
	CheckerCountry         string            `json:"checker_country,omitempty"`
	ResolveOverrideIP      string            `json:"resolve_override_ip,omitempty"`
	Status                 string            `json:"status,omitempty"`
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
	ASN     string `json:"asn"`
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
		"/api/v1/teams/"+teamID+"/monitors/"+monitorID+"/channels/"+channelID, nil)
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
		"/api/v1/teams/"+teamID+"/monitors/"+monitorID+"/channels/"+channelID, nil)
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
		"/api/v1/teams/"+teamID+"/monitors/"+monitorID+"/channels", nil)
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
