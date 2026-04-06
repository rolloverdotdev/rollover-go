// Package rollover provides a Go client for the Rollover API.
//
// Rollover is a subscription billing platform built on x402 that manages
// plans, usage, credits, and recurring billing, settling in USDC on-chain.
//
// Usage:
//
//	ro := rollover.New()
//	result, err := ro.Check(ctx, "0xabc...", "api-calls")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if result.Allowed {
//	    ro.Track(ctx, "0xabc...", "api-calls", 1)
//	}
package rollover

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const defaultBaseURL = "https://api.rollover.dev"

// Client is the Rollover API client.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	mode       string

	slugMu sync.Mutex
	slug   string
}

// New creates a new Rollover client. By default it reads the ROLLOVER_API_KEY
// environment variable. Use WithAPIKey to set it explicitly.
//
//	ro := rollover.New()
//	ro := rollover.New(rollover.WithAPIKey("ro_test_..."))
//	ro := rollover.New(rollover.WithBaseURL("http://localhost:9000"))
func New(opts ...Option) *Client {
	cfg := &config{
		baseURL: defaultBaseURL,
		apiKey:  os.Getenv("ROLLOVER_API_KEY"),
	}
	for _, o := range opts {
		o(cfg)
	}

	mode := "live"
	if strings.HasPrefix(cfg.apiKey, "ro_test_") {
		mode = "test"
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	if cfg.httpClient != nil {
		httpClient = cfg.httpClient
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    strings.TrimRight(cfg.baseURL, "/"),
		apiKey:     cfg.apiKey,
		mode:       mode,
	}
}

func (c *Client) resolveSlug(ctx context.Context) (string, error) {
	c.slugMu.Lock()
	defer c.slugMu.Unlock()

	if c.slug != "" {
		return c.slug, nil
	}

	var org Organization
	if err := c.get(ctx, "/v1/organization", nil, &org); err != nil {
		return "", err
	}
	c.slug = org.Slug
	return c.slug, nil
}

func (c *Client) adminQuery(ctx context.Context, extra url.Values) (url.Values, error) {
	slug, err := c.resolveSlug(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving org slug: %w", err)
	}
	q := url.Values{}
	q.Set("slug", slug)
	q.Set("mode", c.mode)
	for k, vs := range extra {
		for _, v := range vs {
			q.Set(k, v)
		}
	}
	return q, nil
}

func (c *Client) get(ctx context.Context, path string, query url.Values, dest any) error {
	return c.doRequest(ctx, http.MethodGet, path, query, nil, nil, dest)
}

func (c *Client) post(ctx context.Context, path string, query url.Values, body any, dest any) error {
	return c.doRequest(ctx, http.MethodPost, path, query, body, nil, dest)
}

func (c *Client) doRequest(ctx context.Context, method, path string, query url.Values, body any, headers http.Header, dest any) error {
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, vs := range headers {
		for _, v := range vs {
			req.Header.Set(k, v)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseError(resp.StatusCode, respBody)
	}

	if dest != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, dest); err != nil {
			return fmt.Errorf("parsing response: %w", err)
		}
	}

	return nil
}

func setPagination(q url.Values, limit, offset int) {
	if limit > 0 {
		q.Set("limit", fmt.Sprint(limit))
	}
	if offset > 0 {
		q.Set("offset", fmt.Sprint(offset))
	}
}

func setIfNotEmpty(q url.Values, key, value string) {
	if value != "" {
		q.Set(key, value)
	}
}
