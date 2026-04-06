// Batch Usage Report
//
// Query usage events for a time range with pagination and aggregate totals
// by feature and wallet, useful for generating daily or weekly usage digests.
//
//	ROLLOVER_API_KEY=ro_live_... go run ./examples/batch-usage-report
package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	rollover "github.com/rolloverdotdev/rollover-go"
)

func main() {
	ro := rollover.New()
	ctx := context.Background()

	now := time.Now().UTC()
	from := now.Add(-24 * time.Hour).Format(time.RFC3339)
	to := now.Format(time.RFC3339)
	opts := &rollover.ListOptions{After: from, Before: to}

	fmt.Printf("Usage report: %s to %s\n\n", from[:10], to[:10])

	// Collect loads all events into memory at once, handling pagination
	// automatically behind the scenes.
	all, err := rollover.Collect(ctx, ro.ListUsage, opts)
	if err != nil {
		log.Fatal(err)
	}

	byFeature := make(map[string]float64)
	byWallet := make(map[string]float64)
	for _, e := range all {
		amt, _ := strconv.ParseFloat(e.Amount, 64)
		byFeature[e.FeatureSlug] += amt
		byWallet[e.WalletAddress] += amt
	}

	fmt.Printf("Total events: %d\n\n", len(all))

	fmt.Println("By feature:")
	for f, total := range byFeature {
		fmt.Printf("  %-25s %.0f units\n", f, total)
	}

	fmt.Println("\nBy wallet:")
	for w, total := range byWallet {
		fmt.Printf("  %-15s %.0f units\n", shortAddr(w), total)
	}

	// Pages fetches one page at a time, letting you process events as they
	// arrive without holding the full dataset in memory.
	fmt.Println("\nPage-by-page:")
	iter := rollover.Pages(ro.ListUsage, opts)
	pageNum := 0
	for iter.Next(ctx) {
		pageNum++
		fmt.Printf("Page %d: %d events\n", pageNum, len(iter.Page().Data))
	}
	if err := iter.Err(); err != nil {
		log.Fatal(err)
	}
}

func shortAddr(s string) string {
	if len(s) <= 12 {
		return s
	}
	return s[:10] + "..."
}
