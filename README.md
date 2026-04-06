# Rollover Go SDK

The official Go client for the [Rollover](https://rollover.dev) API, a subscription billing platform built on [x402](https://github.com/coinbase/x402) that settles in USDC on-chain.

## Install

```bash
go get github.com/rolloverdotdev/rollover-go
```

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "log"

    rollover "github.com/rolloverdotdev/rollover-go"
)

func main() {
    ro := rollover.New() // reads ROLLOVER_API_KEY env var

    ctx := context.Background()
    wallet := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"

    // Check if the wallet can use a feature.
    result, err := ro.Check(ctx, wallet, "api-calls")
    if err != nil {
        log.Fatal(err)
    }
    if !result.Allowed {
        fmt.Println("Limit reached")
        return
    }

    // Do your work, then track the usage.
    _, err = ro.Track(ctx, wallet, "api-calls", 1)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Configuration

```go
// Default: reads ROLLOVER_API_KEY from environment
ro := rollover.New()

// Explicit API key
ro := rollover.New(rollover.WithAPIKey("ro_test_..."))

// Custom base URL (for local dev)
ro := rollover.New(rollover.WithBaseURL("http://localhost:9000"))

// Custom HTTP client (for timeouts, transports, or testing)
ro := rollover.New(rollover.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}))
```

The mode (`test` or `live`) is parsed from the API key prefix (`ro_test_` or `ro_live_`). The default HTTP client has a 30-second timeout.

## API

### Core

```go
// Check if a wallet can use a feature.
result, err := ro.Check(ctx, wallet, "api-calls")
// result.Allowed, result.Used, result.Remaining, result.Limit,
// result.Plan, result.CreditBalance, result.CreditCost

// Track usage.
result, err := ro.Track(ctx, wallet, "api-calls", 1)
// result.Allowed, result.Used, result.Remaining, result.CreditBalance

// Track with idempotency key to prevent double-counting.
result, err := ro.Track(ctx, wallet, "api-calls", 1,
    rollover.WithIdempotencyKey("order-12345"),
)
```

### Credits

```go
// Get credit balance.
balance, err := ro.GetCredits(ctx, wallet)
// balance.Wallet, balance.Balance

// Grant credits.
result, err := ro.GrantCredits(ctx, wallet, 500,
    rollover.WithDescription("Welcome bonus"),
)
// result.Balance, result.Granted

// List credit transaction history.
txns, err := ro.ListCreditTransactions(ctx, &rollover.ListOptions{Wallet: wallet})
```

### Plans

```go
// List plans.
plans, err := ro.ListPlans(ctx, &rollover.ListOptions{Limit: 10})

// Get a plan.
plan, err := ro.GetPlan(ctx, "starter")

// Create a plan.
plan, err := ro.CreatePlan(ctx, rollover.CreatePlanParams{
    Slug: "starter", Name: "Starter", PriceUSDC: "9.99", BillingPeriod: "monthly",
})

// Update a plan.
plan, err := ro.UpdatePlan(ctx, "starter", rollover.UpdatePlanParams{
    Name: rollover.Ptr("Starter Plus"),
})

// Archive a plan.
err := ro.ArchivePlan(ctx, "starter")

// Add a feature to a plan.
feature, err := ro.CreateFeature(ctx, "starter", rollover.CreateFeatureParams{
    FeatureSlug: "api-calls", Name: "API Calls", LimitAmount: 10000, ResetPeriod: "monthly",
})

// Update a feature.
feature, err := ro.UpdateFeature(ctx, "starter", "api-calls", rollover.UpdateFeatureParams{
    LimitAmount: rollover.Ptr(20000),
})

// Delete a feature.
err := ro.DeleteFeature(ctx, "starter", "api-calls")

// List public pricing for a pricing page.
plans, err := ro.ListPricing(ctx, "your-org-slug")
```

### Subscriptions

```go
// List subscriptions.
subs, err := ro.ListSubscriptions(ctx, &rollover.ListOptions{
    Wallet: "0xabc...",
    Status: "active",
})

// Get a single subscription.
sub, err := ro.GetSubscription(ctx, subscriptionID)

// Create a subscription (admin).
sub, err := ro.CreateSubscription(ctx, "0xabc...", "starter")

// Cancel a subscription.
sub, err := ro.CancelSubscription(ctx, subscriptionID)
```

### Usage and Analytics

```go
// List usage events.
events, err := ro.ListUsage(ctx, &rollover.ListOptions{
    Wallet:  "0xabc...",
    Feature: "api-calls",
    After:   "2025-01-01T00:00:00Z",
})

// Get analytics stats.
stats, err := ro.GetAnalytics(ctx)
// stats.MRR, stats.ActiveSubs, stats.TotalRevenue, stats.TopFeatures

// List invoices.
invoices, err := ro.ListInvoices(ctx, &rollover.ListOptions{Wallet: "0xabc..."})

// Get organization info.
org, err := ro.GetOrganization(ctx)
```

## Pagination

All list methods accept `*ListOptions` with `Limit` and `Offset` fields. For convenience, the SDK provides two helpers that handle pagination automatically.

```go
// Collect loads all items into a single slice.
all, err := rollover.Collect(ctx, ro.ListUsage, &rollover.ListOptions{
    Feature: "api-calls",
})

// Pages iterates one page at a time without loading everything into memory.
iter := rollover.Pages(ro.ListUsage, &rollover.ListOptions{Feature: "api-calls"})
for iter.Next(ctx) {
    for _, e := range iter.Page().Data {
        fmt.Println(e.FeatureSlug, e.Amount)
    }
}
if err := iter.Err(); err != nil {
    log.Fatal(err)
}
```

## Error handling

Non-2xx responses are returned as `*rollover.Error` with a status code, error code, and message.

```go
result, err := ro.Check(ctx, wallet, "api-calls")
if err != nil {
    var roErr *rollover.Error
    if errors.As(err, &roErr) {
        fmt.Println(roErr.StatusCode, roErr.Code, roErr.Message)

        if roErr.Temporary() {
            fmt.Println("Transient error, safe to retry.")
        }
    }
}

// Or use IsErrorCode for clean checks without type assertions.
if rollover.IsErrorCode(err, rollover.ErrCodeNotFound) {
    fmt.Println("Not found.")
}
```

Error code constants: `ErrCodeInvalidAPIKey`, `ErrCodeUnauthorized`, `ErrCodeRateLimit`, `ErrCodeNotFound`, `ErrCodeInsufficientCredits`, `ErrCodeValidation`.

## Examples

See the [examples](./examples) directory:

- [check-and-track](./examples/check-and-track) - Verify feature access before doing work, then record usage after the operation succeeds
- [middleware](./examples/middleware) - An HTTP middleware that gates endpoints by verifying usage and recording consumption
- [credits](./examples/credits) - Protect expensive operations by requiring an available credit balance
- [metered-api](./examples/metered-api) - Track usage for multiple features across different routes
- [idempotency](./examples/idempotency) - Avoid double-counting in distributed systems by using idempotency keys
- [provisioning](./examples/provisioning) - A complete server-side onboarding flow that creates a plan, subscribes a wallet, and grants credits
- [pricing-page](./examples/pricing-page) - Return plans as JSON for a pricing page, with a single API call fetching each plan and its included features
- [usage-dashboard](./examples/usage-dashboard) - Pull analytics stats and paginated usage events to display in an admin dashboard
- [graceful-degradation](./examples/graceful-degradation) - Return a helpful 429 response with usage details and an upgrade path when a wallet hits its limit
- [multi-feature-gate](./examples/multi-feature-gate) - Check multiple features concurrently before starting an operation that requires all of them
- [credit-topup](./examples/credit-topup) - Monitor a wallet's credit balance and automatically grant more credits when it drops below a threshold
- [subscriptions](./examples/subscriptions) - Manage the full subscription lifecycle with listing, filtering, and inspection
- [batch-usage-report](./examples/batch-usage-report) - Query usage events for a time range with pagination and aggregate totals by feature and wallet
- [error-handling](./examples/error-handling) - Handle API errors by inspecting status codes, error codes, and retryability
- [admin-operations](./examples/admin-operations) - Manage plans, features, subscriptions, invoices, and credit transactions using the admin API
- [webhooks](./examples/webhooks) - Process real-time events from Rollover via webhook

## Docs

Visit [docs.rollover.dev](https://docs.rollover.dev) for guides and API reference.

## License

[MIT](LICENSE)
