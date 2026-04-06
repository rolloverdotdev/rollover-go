package rollover

import (
	"context"
	"net/url"
)

// ListInvoices returns a paginated list of invoices, with optional filtering
// by wallet address and status.
func (c *Client) ListInvoices(ctx context.Context, opts *ListOptions) (*Page[Invoice], error) {
	extra := url.Values{}
	if opts != nil {
		setPagination(extra, opts.Limit, opts.Offset)
		setIfNotEmpty(extra, "wallet", opts.Wallet)
		setIfNotEmpty(extra, "status", opts.Status)
	}

	q, err := c.adminQuery(ctx, extra)
	if err != nil {
		return nil, err
	}

	var result Page[Invoice]
	if err := c.get(ctx, "/v1/invoices", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
