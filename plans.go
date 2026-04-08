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
	if err := c.doRequest(ctx, http.MethodPatch, path.Join("/v1/plans", url.PathEscape(planSlug)), q, params, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ArchivePlan archives a plan by slug.
func (c *Client) ArchivePlan(ctx context.Context, planSlug string) error {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return err
	}
	return c.doRequest(ctx, http.MethodDelete, path.Join("/v1/plans", url.PathEscape(planSlug)), q, nil, nil, nil)
}

// CreateFeature adds a feature to a plan.
func (c *Client) CreateFeature(ctx context.Context, planSlug string, params CreateFeatureParams) (*Feature, error) {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result Feature
	if err := c.post(ctx, path.Join("/v1/plans", url.PathEscape(planSlug), "features"), q, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateFeature updates an existing feature on a plan.
func (c *Client) UpdateFeature(ctx context.Context, planSlug, featureSlug string, params UpdateFeatureParams) (*Feature, error) {
	q, err := c.adminQuery(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result Feature
	if err := c.doRequest(ctx, http.MethodPatch, path.Join("/v1/plans", url.PathEscape(planSlug), "features", url.PathEscape(featureSlug)), q, params, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteFeature removes a feature from a plan.
func (c *Client) DeleteFeature(ctx context.Context, planSlug, featureSlug string) error {
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
