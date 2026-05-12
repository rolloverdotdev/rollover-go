package rollover

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
)

// newIdempotencyKey returns a 128-bit random hex string for use as an Idempotency-Key
// header value. The server only requires that the key be opaque and stable across retries,
// so we avoid pulling in a UUID dependency.
func newIdempotencyKey() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

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

	key := cfg.idempotencyKey
	if key == "" {
		key = newIdempotencyKey()
	}
	headers := http.Header{}
	headers.Set("Idempotency-Key", key)

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
