package rollover

import (
	"context"
	"net/http"
	"net/url"
	"path"
)

// ListPlans returns a paginated list of plans for the organization associated
// with the API key.
func (c *Client) ListPlans(ctx context.Context, opts *ListOptions) (*Page[Plan], error) {
	extra := url.Values{}
	if opts != nil {
		setPagination(extra, opts.Limit, opts.Offset)
	}

	q, err := c.adminQuery(ctx, extra)
	if err != nil {
		return nil, err
	}

	var result Page[Plan]
	if err := c.get(ctx, "/v1/plans", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPlan returns a single plan identified by its slug, including any
// features attached to it.
func (c *Client) GetPlan(ctx context.Context, planSlug string) (*Plan, error) {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result Plan
	if err := c.get(ctx, path.Join("/v1/plans", url.PathEscape(planSlug)), q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreatePlan creates a new plan.
func (c *Client) CreatePlan(ctx context.Context, params CreatePlanParams) (*Plan, error) {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result Plan
	if err := c.post(ctx, "/v1/plans", q, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdatePlan updates an existing plan's metadata.
func (c *Client) UpdatePlan(ctx context.Context, planSlug string, params UpdatePlanParams) (*Plan, error) {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result Plan
	if err := c.doRequest(ctx, http.MethodPut, path.Join("/v1/plans", url.PathEscape(planSlug)), q, params, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ArchivePlan archives a plan by slug, hiding it from new subscribers while existing
// subscribers keep their current subscription on the revision they signed up on.
func (c *Client) ArchivePlan(ctx context.Context, planSlug string) error {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return err
	}
	return c.doRequest(ctx, http.MethodDelete, path.Join("/v1/plans", url.PathEscape(planSlug)), q, nil, nil, nil)
}

// DeletePlan hard removes a plan and all of its revisions from the org. The server returns
// 409 plan_has_subscriptions if any subscription past or present references the plan, so
// reach for ArchivePlan when the plan has ever had a subscriber.
func (c *Client) DeletePlan(ctx context.Context, planSlug string) error {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return err
	}
	q.Set("hard", "true")
	return c.doRequest(ctx, http.MethodDelete, path.Join("/v1/plans", url.PathEscape(planSlug)), q, nil, nil, nil)
}

// LinkFeature attaches a catalog feature to a plan. If params.FeatureSlug names a feature
// that does not yet exist in the org catalog, the server creates one as a metered feature.
func (c *Client) LinkFeature(ctx context.Context, planSlug string, params LinkFeatureParams) (*PlanFeature, error) {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result PlanFeature
	if err := c.post(ctx, path.Join("/v1/plans", url.PathEscape(planSlug), "features"), q, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdatePlanFeature edits the limits or policy on an existing plan-feature link.
func (c *Client) UpdatePlanFeature(ctx context.Context, planSlug, featureSlug string, params UpdatePlanFeatureParams) (*PlanFeature, error) {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result PlanFeature
	if err := c.doRequest(ctx, http.MethodPut, path.Join("/v1/plans", url.PathEscape(planSlug), "features", url.PathEscape(featureSlug)), q, params, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UnlinkFeature detaches a feature from a plan. The catalog feature itself is unaffected.
func (c *Client) UnlinkFeature(ctx context.Context, planSlug, featureSlug string) error {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return err
	}
	return c.doRequest(ctx, http.MethodDelete, path.Join("/v1/plans", url.PathEscape(planSlug), "features", url.PathEscape(featureSlug)), q, nil, nil, nil)
}

// ListPricing returns the active plans for a given org slug, intended for
// rendering a public pricing page. This endpoint does not require authentication.
func (c *Client) ListPricing(ctx context.Context, orgSlug string) ([]Plan, error) {
	var result []Plan
	if err := c.get(ctx, path.Join("/v1/pricing", url.PathEscape(orgSlug)), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}
