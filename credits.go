package rollover

import (
	"context"
	"net/url"
	"time"
)

// GetCredits returns the current credit balance for the given wallet address.
func (c *Client) GetCredits(ctx context.Context, wallet string) (*CreditBalance, error) {
	q := url.Values{}
	q.Set("wallet", wallet)

	var result CreditBalance
	if err := c.get(ctx, "/v1/credits", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GrantCredits adds the specified number of credits to a wallet, with an
// optional description and expiration time.
func (c *Client) GrantCredits(ctx context.Context, wallet string, amount int, opts ...GrantOption) (*GrantResult, error) {
	var cfg grantConfig
	for _, o := range opts {
		o(&cfg)
	}

	body := struct {
		Wallet      string `json:"wallet"`
		Amount      int    `json:"amount"`
		Description string `json:"description,omitempty"`
		ExpiresAt   string `json:"expires_at,omitempty"`
	}{
		Wallet:      wallet,
		Amount:      amount,
		Description: cfg.description,
	}
	if cfg.expiresAt != nil {
		body.ExpiresAt = cfg.expiresAt.Format(time.RFC3339)
	}

	var result GrantResult
	if err := c.post(ctx, "/v1/credits", nil, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListCreditTransactions returns a paginated list of credit ledger entries,
// with optional filtering by wallet address.
func (c *Client) ListCreditTransactions(ctx context.Context, opts *ListOptions) (*Page[CreditTransaction], error) {
	extra := url.Values{}
	if opts != nil {
		setPagination(extra, opts.Limit, opts.Offset)
		setIfNotEmpty(extra, "wallet", opts.Wallet)
	}

	q, err := c.adminQuery(ctx, extra)
	if err != nil {
		return nil, err
	}

	var result Page[CreditTransaction]
	if err := c.get(ctx, "/v1/credits/transactions", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
