package rollover

import "time"

// CheckResult is returned by Check.
type CheckResult struct {
	Allowed       bool   `json:"allowed"`
	Used          int    `json:"used"`
	Remaining     int    `json:"remaining"`
	Limit         int    `json:"limit"`
	Plan          string `json:"plan"`
	CreditBalance int    `json:"credit_balance"`
	CreditCost    int    `json:"credit_cost"`
}

// TrackResult is returned by Track.
type TrackResult struct {
	Allowed       bool `json:"allowed"`
	Used          int  `json:"used"`
	Remaining     int  `json:"remaining"`
	CreditBalance int  `json:"credit_balance"`
}

// CreditBalance is returned by GetCredits.
type CreditBalance struct {
	Wallet  string `json:"wallet"`
	Balance int    `json:"balance"`
}

// GrantResult is returned by GrantCredits.
type GrantResult struct {
	Balance int `json:"balance"`
	Granted int `json:"granted"`
}

// Plan represents a billing plan.
type Plan struct {
	ID               string    `json:"id"`
	Slug             string    `json:"slug"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	PriceUSDC        string    `json:"price_usdc"`
	SetupFeeUSDC     string    `json:"setup_fee_usdc"`
	BillingPeriod    string    `json:"billing_period"`
	TrialDays        int       `json:"trial_days"`
	AutoAssign       bool      `json:"auto_assign"`
	IsActive         bool      `json:"is_active"`
	SortOrder        int       `json:"sort_order"`
	Subscribers      int       `json:"subscribers"`
	Features         []Feature `json:"features"`
	Metadata         any       `json:"metadata"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	LastSubscribedAt time.Time `json:"last_subscribed_at"`
}

// Feature represents a metered feature on a plan.
type Feature struct {
	ID           string `json:"id"`
	FeatureSlug  string `json:"feature_slug"`
	Name         string `json:"name"`
	LimitAmount  int    `json:"limit_amount"`
	ResetPeriod  string `json:"reset_period"`
	CreditCost   int    `json:"credit_cost"`
	OveragePrice string `json:"overage_price"`
	Weight       string `json:"weight"`
}

// Subscription represents a wallet's subscription to a plan.
type Subscription struct {
	ID            string    `json:"id"`
	WalletAddress string    `json:"wallet_address"`
	PlanID        string    `json:"plan_id"`
	PlanName      string    `json:"plan_name"`
	Status        string    `json:"status"`
	Mode          string    `json:"mode"`
	PeriodStart   time.Time `json:"period_start"`
	PeriodEnd     time.Time `json:"period_end"`
	TrialEnd      time.Time `json:"trial_end"`
	CancelAtEnd   bool      `json:"cancel_at_end"`
	Metadata      any       `json:"metadata"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// UsageEvent represents a single usage tracking event.
type UsageEvent struct {
	ID             string    `json:"id"`
	WalletAddress  string    `json:"wallet_address"`
	FeatureSlug    string    `json:"feature_slug"`
	Amount         string    `json:"amount"`
	SubscriptionID string    `json:"subscription_id"`
	RecordedAt     time.Time `json:"recorded_at"`
}

// Organization represents the org associated with the API key.
type Organization struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	Logo       string    `json:"logo"`
	WebhookURL string    `json:"webhook_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Page is a paginated list response.
type Page[T any] struct {
	Data   []T `json:"data"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// CreatePlanParams are the parameters for CreatePlan.
type CreatePlanParams struct {
	Slug          string `json:"slug"`
	Name          string `json:"name"`
	PriceUSDC     string `json:"price_usdc"`
	Description   string `json:"description,omitempty"`
	BillingPeriod string `json:"billing_period,omitempty"`
	SetupFeeUSDC  string `json:"setup_fee_usdc,omitempty"`
	TrialDays     int    `json:"trial_days,omitempty"`
	AutoAssign    bool   `json:"auto_assign,omitempty"`
	SortOrder     int    `json:"sort_order,omitempty"`
}

// CreateFeatureParams are the parameters for CreateFeature.
type CreateFeatureParams struct {
	FeatureSlug  string `json:"feature_slug"`
	Name         string `json:"name"`
	LimitAmount  int    `json:"limit_amount,omitempty"`
	ResetPeriod  string `json:"reset_period,omitempty"`
	CreditCost   int    `json:"credit_cost,omitempty"`
	OveragePrice string `json:"overage_price,omitempty"`
	Weight       string `json:"weight,omitempty"`
}

// UpdatePlanParams are the parameters for UpdatePlan. Only non-nil fields are
// sent to the server, allowing explicit zero values like false or 0.
type UpdatePlanParams struct {
	Name          *string `json:"name,omitempty"`
	Description   *string `json:"description,omitempty"`
	PriceUSDC     *string `json:"price_usdc,omitempty"`
	SetupFeeUSDC  *string `json:"setup_fee_usdc,omitempty"`
	BillingPeriod *string `json:"billing_period,omitempty"`
	TrialDays     *int    `json:"trial_days,omitempty"`
	AutoAssign    *bool   `json:"auto_assign,omitempty"`
	IsActive      *bool   `json:"is_active,omitempty"`
	SortOrder     *int    `json:"sort_order,omitempty"`
}

// UpdateFeatureParams are the parameters for UpdateFeature. Only non-nil fields
// are sent to the server, allowing explicit zero values.
type UpdateFeatureParams struct {
	Name         *string `json:"name,omitempty"`
	LimitAmount  *int    `json:"limit_amount,omitempty"`
	ResetPeriod  *string `json:"reset_period,omitempty"`
	CreditCost   *int    `json:"credit_cost,omitempty"`
	OveragePrice *string `json:"overage_price,omitempty"`
	Weight       *string `json:"weight,omitempty"`
}

// Ptr returns a pointer to the given value, useful for setting fields on
// update param structs.
func Ptr[T any](v T) *T { return &v }

// AnalyticsStats contains high-level analytics for the organization.
type AnalyticsStats struct {
	MRR            string        `json:"mrr"`
	ActiveSubs     int           `json:"active_subs"`
	TotalRevenue   string        `json:"total_revenue"`
	TopFeatures    []TopFeature  `json:"top_features"`
	RecentActivity []RecentEvent `json:"recent_activity"`
}

// TopFeature represents a feature ranked by total usage.
type TopFeature struct {
	FeatureSlug string `json:"feature_slug"`
	TotalUsed   int    `json:"total_used"`
}

// RecentEvent represents a recent usage event in the activity feed.
type RecentEvent struct {
	WalletAddress string    `json:"wallet_address"`
	FeatureSlug   string    `json:"feature_slug"`
	Amount        string    `json:"amount"`
	RecordedAt    time.Time `json:"recorded_at"`
}

// CreditTransaction represents a single credit ledger entry.
type CreditTransaction struct {
	ID             string    `json:"id"`
	WalletAddress  string    `json:"wallet_address"`
	Amount         int       `json:"amount"`
	Type           string    `json:"type"`
	Description    string    `json:"description"`
	Mode           string    `json:"mode"`
	SubscriptionID string    `json:"subscription_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// Invoice represents a billing invoice.
type Invoice struct {
	ID             string    `json:"id"`
	WalletAddress  string    `json:"wallet_address"`
	SubscriptionID string    `json:"subscription_id"`
	Status         string    `json:"status"`
	BaseAmount     string    `json:"base_amount"`
	OverageAmount  string    `json:"overage_amount"`
	TotalAmount    string    `json:"total_amount"`
	TxHash         string    `json:"tx_hash"`
	PeriodStart    time.Time `json:"period_start"`
	PeriodEnd      time.Time `json:"period_end"`
	SettledAt      time.Time `json:"settled_at"`
	CreatedAt      time.Time `json:"created_at"`
}
