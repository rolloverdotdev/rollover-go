package rollover

import (
	"net/http"
	"time"
)

// Option configures a Client.
type Option func(*config)

type config struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// WithAPIKey sets the API key. If not provided, reads ROLLOVER_API_KEY from env.
func WithAPIKey(key string) Option {
	return func(c *config) { c.apiKey = key }
}

// WithBaseURL overrides the default API base URL (https://api.rollover.dev).
func WithBaseURL(url string) Option {
	return func(c *config) { c.baseURL = url }
}

// WithHTTPClient sets a custom HTTP client for timeouts, transports, or testing.
func WithHTTPClient(client *http.Client) Option {
	return func(c *config) { c.httpClient = client }
}

// GrantOption configures GrantCredits.
type GrantOption func(*grantConfig)

type grantConfig struct {
	description string
	expiresAt   *time.Time
}

// WithDescription sets a description for a credit grant.
func WithDescription(d string) GrantOption {
	return func(c *grantConfig) { c.description = d }
}

// WithExpiresAt sets an expiration time for a credit grant.
func WithExpiresAt(t time.Time) GrantOption {
	return func(c *grantConfig) { c.expiresAt = &t }
}

// TrackOption configures Track.
type TrackOption func(*trackConfig)

type trackConfig struct {
	idempotencyKey string
}

// WithIdempotencyKey sets the idempotency key to prevent double-counting.
func WithIdempotencyKey(key string) TrackOption {
	return func(c *trackConfig) { c.idempotencyKey = key }
}

// ListOptions configures any paginated list method, where nil uses defaults
// and zero-value fields are omitted from the request.
type ListOptions struct {
	Limit   int
	Offset  int
	Wallet  string
	Status  string
	PlanID  string
	Feature string
	After   string // RFC3339 timestamp, e.g. "2025-01-01T00:00:00Z"
	Before  string // RFC3339 timestamp, e.g. "2025-01-02T00:00:00Z"
}
