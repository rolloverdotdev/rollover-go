package rollover

import (
	"context"
	"net/http"
	"net/url"
	"path"
)

// ListSubscriptions returns a paginated list of subscriptions, with optional
// filters for wallet address, status, and plan ID.
func (c *Client) ListSubscriptions(ctx context.Context, opts *ListOptions) (*Page[Subscription], error) {
	extra := url.Values{}
	if opts != nil {
		setPagination(extra, opts.Limit, opts.Offset)
		setIfNotEmpty(extra, "wallet", opts.Wallet)
		setIfNotEmpty(extra, "status", opts.Status)
		setIfNotEmpty(extra, "plan_id", opts.PlanID)
	}

	q, err := c.adminQuery(ctx, extra)
	if err != nil {
		return nil, err
	}

	var result Page[Subscription]
	if err := c.get(ctx, "/v1/subscriptions", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSubscription returns a single subscription by ID.
func (c *Client) GetSubscription(ctx context.Context, subscriptionID string) (*Subscription, error) {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result Subscription
	if err := c.get(ctx, path.Join("/v1/subscriptions", url.PathEscape(subscriptionID)), q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateSubscription creates an admin-initiated subscription that assigns a
// wallet to the specified plan without requiring on-chain payment.
func (c *Client) CreateSubscription(ctx context.Context, wallet, planSlug string) (*Subscription, error) {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return nil, err
	}

	body := struct {
		WalletAddress string `json:"wallet_address"`
		PlanSlug      string `json:"plan_slug"`
	}{wallet, planSlug}

	var result Subscription
	if err := c.post(ctx, "/v1/subscriptions", q, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelSubscription cancels the subscription with the given ID, marking it
// to expire at the end of the current billing period.
func (c *Client) CancelSubscription(ctx context.Context, subscriptionID string) (*Subscription, error) {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result Subscription
	if err := c.doRequest(ctx, http.MethodDelete, path.Join("/v1/subscriptions", url.PathEscape(subscriptionID)), q, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
