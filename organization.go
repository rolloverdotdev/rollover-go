package rollover

import "context"

// GetOrganization returns the organization associated with the API key.
func (c *Client) GetOrganization(ctx context.Context) (*Organization, error) {
	var result Organization
	if err := c.get(ctx, "/v1/organization", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
