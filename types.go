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

// Plan represents a billing plan. Pricing fields are hydrated from the latest revision on the
// server side, so PriceUSDC and friends still read directly off the struct even though the
// canonical pricing now lives in plan_revisions. LatestRevisionID identifies which revision
// new subscribers will be pinned to.
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
	IsArchived       bool      `json:"is_archived"`
	LatestRevisionID string    `json:"latest_revision_id"`
	SortOrder        int       `json:"sort_order"`
	Subscribers      int           `json:"subscribers"`
	Features         []PlanFeature `json:"features"`
	Metadata         any       `json:"metadata"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	LastSubscribedAt time.Time `json:"last_subscribed_at"`
}

// FeatureType is the canonical kind of feature in the org catalog. Boolean is an access
// flag, metered is a usage counter, credit is a pooled balance fed by metered features,
// and static is a non-consumptive numeric cap the consumer app enforces itself.
type FeatureType string

const (
	FeatureTypeBoolean FeatureType = "boolean"
	FeatureTypeMetered FeatureType = "metered"
	FeatureTypeCredit  FeatureType = "credit"
	FeatureTypeStatic  FeatureType = "static"
)

// Policy controls what happens when a subscriber hits the plan-feature limit. HardBlock
// rejects the request, SoftWarn lets it through (metered/credit only) for cycle-end
// reconciliation, and Hide treats the feature as not present at all.
type Policy string

const (
	PolicyHardBlock Policy = "hard_block"
	PolicySoftWarn  Policy = "soft_warn"
	PolicyHide      Policy = "hide"
)

// Feature is an org-scoped catalog feature. Plans reference features through PlanFeature
// link rows; the catalog row owns the canonical slug, display name, and type.
type Feature struct {
	ID   string      `json:"id"`
	Slug string      `json:"slug"`
	Name string      `json:"name"`
	Type FeatureType `json:"type"`
}

// PlanFeature is one feature linked to one plan, carrying the plan-specific limits and
// the policy that controls what happens when a subscriber hits them. The nested Feature
// pointer is populated on responses with the catalog row this link points to.
type PlanFeature struct {
	ID           string   `json:"id"`
	LimitAmount  int      `json:"limit_amount"`
	ResetPeriod  string   `json:"reset_period"`
	OveragePrice string   `json:"overage_price"`
	Weight       string   `json:"weight"`
	CreditCost   int      `json:"credit_cost"`
	Policy       Policy   `json:"policy"`
	Feature      *Feature `json:"feature,omitempty"`
}

// Subscription represents a wallet's subscription to a plan. PlanRevisionID pins the
// subscription to the pricing revision it signed up on, so renewals charge the same price
// even after the plan's price is edited.
type Subscription struct {
	ID             string    `json:"id"`
	WalletAddress  string    `json:"wallet_address"`
	PlanID         string    `json:"plan_id"`
	PlanRevisionID string    `json:"plan_revision_id,omitempty"`
	PlanName       string    `json:"plan_name"`
	Status         string    `json:"status"`
	BillingPeriod  string    `json:"billing_period"`
	Mode           string    `json:"mode"`
	PeriodStart    time.Time `json:"period_start"`
	PeriodEnd      time.Time `json:"period_end"`
	TrialEnd       time.Time `json:"trial_end"`
	CancelAtEnd    bool      `json:"cancel_at_end"`
	Metadata       any       `json:"metadata"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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

// LinkFeatureParams are the parameters for LinkFeature, which attaches a catalog feature
// to a plan. Supply either FeatureID or FeatureSlug; if FeatureSlug names a feature that
// does not yet exist in the org catalog, the server creates one as a metered feature.
// Policy defaults to hard_block on the server when omitted; soft_warn requires a metered
// or credit feature.
type LinkFeatureParams struct {
	FeatureID    string `json:"feature_id,omitempty"`
	FeatureSlug  string `json:"feature_slug,omitempty"`
	LimitAmount  int    `json:"limit_amount,omitempty"`
	ResetPeriod  string `json:"reset_period,omitempty"`
	CreditCost   int    `json:"credit_cost,omitempty"`
	OveragePrice string `json:"overage_price,omitempty"`
	Weight       string `json:"weight,omitempty"`
	Policy       Policy `json:"policy,omitempty"`
}

// UpdatePlanParams are the parameters for UpdatePlan. Only non-nil fields are
// sent to the server, allowing explicit zero values like false or 0.
// Setting any pricing field (PriceUSDC, BillingPeriod, TrialDays, SetupFeeUSDC) on the
// server creates a new plan revision instead of mutating the existing one, so existing
// subscribers stay pinned to the price they signed up on.
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

// UpdatePlanFeatureParams are the parameters for UpdatePlanFeature, which edits one
// plan-feature link. Only non-nil fields are sent to the server.
type UpdatePlanFeatureParams struct {
	LimitAmount  *int    `json:"limit_amount,omitempty"`
	ResetPeriod  *string `json:"reset_period,omitempty"`
	CreditCost   *int    `json:"credit_cost,omitempty"`
	OveragePrice *string `json:"overage_price,omitempty"`
	Weight       *string `json:"weight,omitempty"`
	Policy       *Policy `json:"policy,omitempty"`
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

// Invoice represents a billing invoice. ChainID and Mode identify which chain the invoice
// settled on and which environment it belongs to.
type Invoice struct {
	ID             string    `json:"id"`
	WalletAddress  string    `json:"wallet_address"`
	SubscriptionID string    `json:"subscription_id"`
	Mode           string    `json:"mode"`
	ChainID        string    `json:"chain_id"`
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

// Chain represents a payment destination chain configured on an organization for a given mode.
type Chain struct {
	ID               string    `json:"id"`
	OrgID            string    `json:"org_id"`
	Mode             string    `json:"mode"`
	ChainID          string    `json:"chain_id"`
	PayToAddress     string    `json:"pay_to_address"`
	StablecoinSymbol string    `json:"stablecoin_symbol"`
	Enabled          bool      `json:"enabled"`
	Priority         int       `json:"priority"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CreateChainParams are the parameters for CreateChain.
type CreateChainParams struct {
	ChainID          string `json:"chain_id"`
	PayToAddress     string `json:"pay_to_address"`
	StablecoinSymbol string `json:"stablecoin_symbol,omitempty"`
	Priority         int    `json:"priority,omitempty"`
}

// UpdateChainParams are the parameters for UpdateChain; only non-nil fields are sent so
// you can leave any field out to keep its current value.
type UpdateChainParams struct {
	PayToAddress     *string `json:"pay_to_address,omitempty"`
	StablecoinSymbol *string `json:"stablecoin_symbol,omitempty"`
	Enabled          *bool   `json:"enabled,omitempty"`
	Priority         *int    `json:"priority,omitempty"`
}
