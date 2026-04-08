package rollover

import "context"

// GetAnalytics returns high-level analytics stats for the organization,
// including MRR, active subscribers, top features, and recent activity.
func (c *Client) GetAnalytics(ctx context.Context) (*AnalyticsStats, error) {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result AnalyticsStats
	if err := c.get(ctx, "/v1/analytics", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
