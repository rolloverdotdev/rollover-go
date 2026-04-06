# Examples

Practical integration patterns for the Rollover Go SDK, with each example provided as a standalone program you can run directly.

```bash
export ROLLOVER_API_KEY="ro_test_..."
go run ./examples/check-and-track/
```

---

## Check and Track

The core Rollover pattern is to verify a wallet has feature access before doing any work, then record usage after the operation succeeds.

```go
ro := rollover.New()

result, err := ro.Check(ctx, wallet, "api-calls")
if err != nil {
    log.Fatal(err)
}
if !result.Allowed {
    fmt.Printf("Limit reached. %d/%d used.\n", result.Used, result.Limit)
    return
}

// Do your work...

ro.Track(ctx, wallet, "api-calls", 1)
```

**Full example:** [check-and-track](./check-and-track/main.go)

---

## HTTP Middleware

An HTTP middleware that gates endpoints by verifying usage before handling the request and recording consumption after a successful response.

```go
func usageGate(feature string, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        wallet := r.Header.Get("X-Wallet")
        result, _ := ro.Check(r.Context(), wallet, feature)
        if !result.Allowed {
            http.Error(w, "rate limited", http.StatusTooManyRequests)
            return
        }
        next(w, r)
        ro.Track(context.Background(), wallet, feature, 1)
    }
}
```

**Full example:** [middleware](./middleware/main.go)

---

## Credit-Gated Access

Protect expensive operations by requiring an available credit balance, with credits automatically deducted according to the feature's `credit_cost`.

```go
balance, _ := ro.GetCredits(ctx, wallet)
fmt.Printf("Credits: %d\n", balance.Balance)

grant, _ := ro.GrantCredits(ctx, wallet, 500, rollover.WithDescription("Welcome bonus"))
fmt.Printf("Granted %d, balance: %d\n", grant.Granted, grant.Balance)

result, _ := ro.Check(ctx, wallet, "image-gen")
if result.Allowed {
    ro.Track(ctx, wallet, "image-gen", 1)
}
```

**Full example:** [credits](./credits/main.go)

---

## Metered API Server

Track usage for multiple features across different routes, with each route mapped to a Rollover feature.

```go
http.HandleFunc("POST /v1/translate", metered("translations", translateHandler))
http.HandleFunc("POST /v1/summarize", metered("summaries", summarizeHandler))
http.HandleFunc("POST /v1/embeddings", metered("embeddings", embeddingsHandler))
```

**Full example:** [metered-api](./metered-api/main.go)

---

## Idempotent Tracking

Avoid double-counting in distributed systems by using idempotency keys, where the same `Idempotency-Key` always produces the same result.

```go
ro.Track(ctx, wallet, "api-calls", 1, rollover.WithIdempotencyKey("order-12345"))
```

**Full example:** [idempotency](./idempotency/main.go)

---

## Provision a Customer

A complete server-side onboarding flow that creates a plan with features, subscribes a wallet, and grants welcome credits for a new customer.

```go
plan, _ := ro.CreatePlan(ctx, rollover.CreatePlanParams{
    Slug: "starter", Name: "Starter", PriceUSDC: "9.99", BillingPeriod: "monthly",
})

ro.CreateFeature(ctx, "starter", rollover.CreateFeatureParams{
    FeatureSlug: "api-calls", Name: "API Calls", LimitAmount: 10000, ResetPeriod: "monthly",
})

ro.CreateSubscription(ctx, wallet, "starter")
ro.GrantCredits(ctx, wallet, 500, rollover.WithDescription("Welcome bonus"))
```

**Full example:** [provisioning](./provisioning/main.go)

---

## Pricing Page

Return your plans as JSON for a pricing page, with a single API call fetching each plan and its included features.

```go
plans, _ := ro.ListPricing(ctx, "your-org-slug")
json.NewEncoder(w).Encode(plans)
```

**Full example:** [pricing-page](./pricing-page/main.go)

---

## Usage Dashboard

Pull analytics stats and paginated usage events to display in an admin dashboard, combining MRR, active subscriptions, and event history.

```go
stats, _ := ro.GetAnalytics(ctx)
fmt.Printf("MRR: $%s, Active subs: %d\n", stats.MRR, stats.ActiveSubs)

events, _ := ro.ListUsage(ctx, &rollover.ListOptions{Limit: 10})
for _, e := range events.Data {
    fmt.Printf("%s  %s  %s units\n", e.WalletAddress, e.FeatureSlug, e.Amount)
}
```

**Full example:** [usage-dashboard](./usage-dashboard/main.go)

---

## Graceful Degradation

Return a helpful 429 response with usage details and an upgrade path when a wallet hits its limit, rather than a generic rate limit error.

```go
result, _ := ro.Check(ctx, wallet, "generations")
if !result.Allowed {
    writeJSON(w, 429, map[string]any{
        "error": "limit reached",
        "used":  result.Used,
        "limit": result.Limit,
        "plan":  result.Plan,
    })
}
```

**Full example:** [graceful-degradation](./graceful-degradation/main.go)

---

## Multi-Feature Gate

Check multiple features concurrently before starting an operation that requires all of them, such as an AI pipeline consuming both API calls and image generation credits.

```go
required := []string{"api-calls", "image-gen"}
blocked := checkAll(ctx, wallet, required)
if len(blocked) > 0 {
    fmt.Printf("Blocked on: %s\n", strings.Join(blocked, ", "))
    return
}
trackAll(ctx, wallet, map[string]int{"api-calls": 1, "image-gen": 1})
```

**Full example:** [multi-feature-gate](./multi-feature-gate/main.go)

---

## Credit Top-Up

Monitor a wallet's credit balance on an interval and automatically grant more credits when the balance drops below a configured threshold.

```go
balance, _ := ro.GetCredits(ctx, wallet)
if balance.Balance < threshold {
    ro.GrantCredits(ctx, wallet, topUpAmount,
        rollover.WithDescription("Auto top-up: balance below threshold"),
    )
}
```

**Full example:** [credit-topup](./credit-topup/main.go)

---

## Subscription Lifecycle

Manage the full subscription lifecycle by listing active subscriptions, filtering by wallet, and inspecting subscription details.

```go
subs, _ := ro.ListSubscriptions(ctx, &rollover.ListOptions{Status: "active"})
for _, s := range subs.Data {
    fmt.Printf("%s -> %s (status: %s)\n", s.WalletAddress, s.PlanName, s.Status)
}
```

**Full example:** [subscriptions](./subscriptions/main.go)

---

## Batch Usage Report

Query usage events for a time range with pagination and aggregate totals by feature and wallet, useful for generating daily or weekly usage digests.

```go
page, _ := ro.ListUsage(ctx, &rollover.ListOptions{
    After:  "2025-01-01T00:00:00Z",
    Before: "2025-01-02T00:00:00Z",
    Limit:  100,
})
```

**Full example:** [batch-usage-report](./batch-usage-report/main.go)

---

## Error Handling

Handle Rollover API errors by inspecting the status code and error code, allowing your application to respond differently to authentication failures, rate limits, and other error conditions.

```go
result, err := ro.Check(ctx, wallet, "api-calls")
if err != nil {
    var roErr *rollover.Error
    if errors.As(err, &roErr) {
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

**Full example:** [error-handling](./error-handling/main.go)

---

## Admin Operations

Manage plans, features, subscriptions, invoices, and credit transactions using the admin API, covering the full set of operations available to API key holders.

```go
ro.UpdatePlan(ctx, "starter", rollover.UpdatePlanParams{Name: rollover.Ptr("Starter (Updated)")})
ro.UpdateFeature(ctx, "starter", "api-calls", rollover.UpdateFeatureParams{LimitAmount: rollover.Ptr(20000)})
ro.DeleteFeature(ctx, "starter", "old-feature")

sub, _ := ro.GetSubscription(ctx, subscriptionID)
txns, _ := ro.ListCreditTransactions(ctx, &rollover.ListOptions{Wallet: wallet})
invoices, _ := ro.ListInvoices(ctx, &rollover.ListOptions{Wallet: wallet})
```

**Full example:** [admin-operations](./admin-operations/main.go)

---

## Webhook Receiver

Process real-time events from Rollover by registering a webhook URL in the dashboard and handling events as they arrive.

```go
http.HandleFunc("POST /webhooks/rollover", func(w http.ResponseWriter, r *http.Request) {
    var event WebhookEvent
    json.NewDecoder(r.Body).Decode(&event)
    switch event.Type {
    case "subscription.created":
        // handle new subscription
    case "subscription.canceled":
        // handle cancellation
    }
    w.WriteHeader(http.StatusOK)
})
```

**Full example:** [webhooks](./webhooks/main.go)
