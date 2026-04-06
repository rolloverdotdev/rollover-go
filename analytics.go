package rollover

import "context"

// GetAnalytics returns high-level analytics stats for the organization,
// including MRR, active subscribers, top features, and recent activity.
func (c *Client) GetAnalytics(ctx context.Context) (*AnalyticsStats, error) {
	var result AnalyticsStats
	if err := c.get(ctx, "/v1/analytics", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
