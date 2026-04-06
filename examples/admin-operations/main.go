// Admin Operations
//
// Manage plans, features, subscriptions, invoices, and credit transactions
// using the admin API, covering the full set of operations available to
// API key holders beyond the core check and track workflow.
//
//	ROLLOVER_API_KEY=ro_test_... go run ./examples/admin-operations
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

	slug := fmt.Sprintf("admin-demo-%d", time.Now().UnixNano()%100000)

	// Create a plan.
	plan, err := ro.CreatePlan(ctx, rollover.CreatePlanParams{
		Slug:          slug,
		Name:          "Admin Demo",
		PriceUSDC:     "19.99",
		BillingPeriod: "monthly",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created plan: %s\n", plan.Name)

	// Update the plan.
	updated, err := ro.UpdatePlan(ctx, slug, rollover.UpdatePlanParams{
		Name:        rollover.Ptr("Admin Demo (Updated)"),
		Description: rollover.Ptr("Updated via SDK"),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated plan: %s\n", updated.Name)

	// Add a feature.
	feature, err := ro.CreateFeature(ctx, slug, rollover.CreateFeatureParams{
		FeatureSlug: "requests",
		Name:        "API Requests",
		LimitAmount: 5000,
		ResetPeriod: "monthly",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Added feature: %s (limit: %d)\n", feature.FeatureSlug, feature.LimitAmount)

	// Update the feature.
	updatedFeature, err := ro.UpdateFeature(ctx, slug, "requests", rollover.UpdateFeatureParams{
		LimitAmount: rollover.Ptr(10000),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated feature limit: %d\n", updatedFeature.LimitAmount)


	// Subscribe a wallet and inspect the subscription.
	wallet := fmt.Sprintf("0x%040x", time.Now().UnixNano())
	sub, err := ro.CreateSubscription(ctx, wallet, slug)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Subscribed: %s (status: %s)\n", wallet[:12]+"...", sub.Status)

	fetched, err := ro.GetSubscription(ctx, sub.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Fetched subscription: plan=%s, period ends %s\n",
		fetched.PlanName, fetched.PeriodEnd.Format("2006-01-02"))

	// Grant credits and list transactions.
	ro.GrantCredits(ctx, wallet, 100, rollover.WithDescription("Demo grant"))
	txns, err := ro.ListCreditTransactions(ctx, &rollover.ListOptions{Wallet: wallet})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Credit transactions: %d\n", txns.Total)
	for _, tx := range txns.Data {
		fmt.Printf("  %s: %d credits (%s)\n", tx.Type, tx.Amount, tx.Description)
	}

	// List invoices.
	invoices, err := ro.ListInvoices(ctx, &rollover.ListOptions{Wallet: wallet})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Invoices: %d\n", invoices.Total)

	// Cleanup.
	ro.DeleteFeature(ctx, slug, "requests")
	ro.ArchivePlan(ctx, slug)
	fmt.Println("Cleaned up.")
}
