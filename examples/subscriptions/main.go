// Subscription Lifecycle
//
// Manage the full subscription lifecycle by listing active subscriptions,
// filtering by wallet, and inspecting subscription details.
//
//	ROLLOVER_API_KEY=ro_test_... go run ./examples/subscriptions
package main

import (
	"context"
	"fmt"
	"log"

	rollover "github.com/rolloverdotdev/rollover-go"
)

func main() {
	ro := rollover.New()
	ctx := context.Background()

	// List all active subscriptions.
	subs, err := ro.ListSubscriptions(ctx, &rollover.ListOptions{Status: "active", Limit: 5})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Active subscriptions: %d\n", subs.Total)

	for _, s := range subs.Data {
		fmt.Printf("  %s -> %s (status: %s, ends: %s)\n",
			s.WalletAddress, s.PlanName, s.Status, s.PeriodEnd.Format("2006-01-02"))
	}

	// Filter by wallet.
	if len(subs.Data) > 0 {
		wallet := subs.Data[0].WalletAddress
		filtered, err := ro.ListSubscriptions(ctx, &rollover.ListOptions{Wallet: wallet})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\nSubscriptions for %s: %d\n", wallet[:12]+"...", filtered.Total)
	}
}
