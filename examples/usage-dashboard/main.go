// Usage Dashboard
//
// Pull analytics stats and paginated usage events to display in an admin
// dashboard, combining MRR, active subscriptions, and event history.
//
//	ROLLOVER_API_KEY=ro_live_... go run ./examples/usage-dashboard
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

	// 1. Fetch high-level analytics.
	stats, err := ro.GetAnalytics(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Dashboard")
	fmt.Printf("MRR:           $%s\n", stats.MRR)
	fmt.Printf("Active subs:   %d\n", stats.ActiveSubs)
	fmt.Printf("Total revenue: $%s\n", stats.TotalRevenue)

	if len(stats.TopFeatures) > 0 {
		fmt.Println("\nTop features:")
		for _, f := range stats.TopFeatures {
			fmt.Printf("  %-20s %d events\n", f.FeatureSlug, f.TotalUsed)
		}
	}

	// 2. Fetch recent usage events.
	events, err := ro.ListUsage(ctx, &rollover.ListOptions{Limit: 10})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nRecent events (showing %d of %d):\n", len(events.Data), events.Total)
	for _, e := range events.Data {
		fmt.Printf("  %s  %-15s  %s units  %s\n",
			shortAddr(e.WalletAddress),
			e.FeatureSlug,
			e.Amount,
			e.RecordedAt.Format("15:04:05"),
		)
	}
}

func shortAddr(s string) string {
	if len(s) <= 12 {
		return s
	}
	return s[:10] + "..."
}
