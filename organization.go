package rollover

import (
	"context"
	"net/http"
	"net/url"
	"path"
)

// GetOrganization returns the organization associated with the API key.
func (c *Client) GetOrganization(ctx context.Context) (*Organization, error) {
	var result Organization
	if err := c.get(ctx, "/v1/organization", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListChains returns every payment chain configured for the API key's org and mode,
// including disabled ones, ordered by priority so the first enabled chain is the one
// subscribers settle to.
func (c *Client) ListChains(ctx context.Context) ([]Chain, error) {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return nil, err
	}
	var result []Chain
	if err := c.get(ctx, "/v1/organization/chains", q, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateChain adds a new payment destination chain. Use this to accept payments on additional
// networks, or to set up your live mode payout address before issuing live API keys.
// Do not call this with a chain id outside the supported networks catalog; the server returns
// 400 unsupported_chain. Test mode keys can only add testnets and live mode keys only mainnets.
func (c *Client) CreateChain(ctx context.Context, params CreateChainParams) (*Chain, error) {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return nil, err
	}
	var result Chain
	if err := c.post(ctx, "/v1/organization/chains", q, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateChain edits an existing chain's address, stablecoin, enabled flag, or priority,
// sending only the fields set in params so the rest stay at their current values.
func (c *Client) UpdateChain(ctx context.Context, id string, params UpdateChainParams) (*Chain, error) {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return nil, err
	}
	var result Chain
	if err := c.doRequest(ctx, http.MethodPut, path.Join("/v1/organization/chains", url.PathEscape(id)), q, params, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteChain removes a chain so subscribers can no longer pay on it; when this was the only
// enabled chain, paid flows fail with no_chain_configured until another is added.
func (c *Client) DeleteChain(ctx context.Context, id string) error {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return err
	}
	return c.doRequest(ctx, http.MethodDelete, path.Join("/v1/organization/chains", url.PathEscape(id)), q, nil, nil, nil)
}
