// Provision a Customer
//
// A complete server-side onboarding flow that creates a plan with features,
// subscribes a wallet, and grants welcome credits for a new customer.
//
//	ROLLOVER_API_KEY=ro_test_... go run ./examples/provisioning
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	rollover "github.com/rolloverdotdev/rollover-go"
)

func main() {
	ro := rollover.New()
	ctx := context.Background()

	slug := fmt.Sprintf("starter-%d", time.Now().UnixNano()%100000)

	// 1. Create a plan.
	plan, err := ro.CreatePlan(ctx, rollover.CreatePlanParams{
		Slug:          slug,
		Name:          "Starter",
		PriceUSDC:     "9.99",
		BillingPeriod: "monthly",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created plan: %s (%s)\n", plan.Name, plan.Slug)

	// 2. Link features to the plan. New feature slugs are auto-created in the catalog
	// as metered features.
	link, err := ro.LinkFeature(ctx, slug, rollover.LinkFeatureParams{
		FeatureSlug: "api-calls",
		LimitAmount: 10000,
		ResetPeriod: "monthly",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Linked feature: %s (limit: %d)\n", link.Feature.Slug, link.LimitAmount)

	// 3. Subscribe a wallet.
	wallet := fmt.Sprintf("0x%040x", time.Now().UnixNano())
	sub, err := ro.CreateSubscription(ctx, wallet, slug)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Subscribed %s to %s (status: %s)\n", wallet[:12]+"...", sub.PlanName, sub.Status)

	// 4. Grant welcome credits.
	grant, err := ro.GrantCredits(ctx, wallet, 500, rollover.WithDescription("Welcome bonus"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Granted 500 credits (balance: %d)\n", grant.Balance)

	// Cleanup.
	ro.ArchivePlan(ctx, slug)
}
