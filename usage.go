package rollover

import (
	"context"
	"net/http"
	"net/url"
)

// Check returns whether a wallet is allowed to use a feature.
func (c *Client) Check(ctx context.Context, wallet, feature string) (*CheckResult, error) {
	q := url.Values{}
	q.Set("wallet", wallet)
	q.Set("feature", feature)

	var result CheckResult
	if err := c.get(ctx, "/v1/check", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Track records a usage event for the given wallet and feature, deducting from
// the wallet's quota or credits as configured on the plan.
func (c *Client) Track(ctx context.Context, wallet, feature string, amount int, opts ...TrackOption) (*TrackResult, error) {
	var cfg trackConfig
	for _, o := range opts {
		o(&cfg)
	}

	body := struct {
		Wallet  string `json:"wallet"`
		Feature string `json:"feature"`
		Amount  int    `json:"amount"`
	}{wallet, feature, amount}

	var headers http.Header
	if cfg.idempotencyKey != "" {
		headers = http.Header{}
		headers.Set("Idempotency-Key", cfg.idempotencyKey)
	}

	var result TrackResult
	if err := c.doRequest(ctx, http.MethodPost, "/v1/track", nil, body, headers, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListUsage returns a paginated list of usage events, with optional filters
// for wallet, feature, and time range.
func (c *Client) ListUsage(ctx context.Context, opts *ListOptions) (*Page[UsageEvent], error) {
	extra := url.Values{}
	if opts != nil {
		setPagination(extra, opts.Limit, opts.Offset)
		setIfNotEmpty(extra, "wallet", opts.Wallet)
		setIfNotEmpty(extra, "feature", opts.Feature)
		setIfNotEmpty(extra, "after", opts.After)
		setIfNotEmpty(extra, "before", opts.Before)
	}

	q, err := c.adminQuery(ctx, extra)
	if err != nil {
		return nil, err
	}

	var result Page[UsageEvent]
	if err := c.get(ctx, "/v1/usage", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
