// Multi-Feature Gate
//
// Check multiple features in one call before starting an operation that
// requires all of them, such as an AI pipeline consuming both API calls
// and image generation credits.
//
//	ROLLOVER_API_KEY=ro_live_... go run ./examples/multi-feature-gate
package main

import (
	"context"
	"fmt"
	"strings"

	rollover "github.com/rolloverdotdev/rollover-go"
)

func main() {
	ctx := context.Background()
	ro := rollover.New()
	wallet := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"

	// CheckBatch resolves the subscription once and answers for every feature
	// in a single request. Supplying Amount makes Allowed reflect whether N
	// units would succeed, not just whether any quota remains.
	gate, err := ro.CheckBatch(ctx, wallet, []rollover.BatchCheckItem{
		{Feature: "api-calls", Amount: rollover.Ptr(int64(1))},
		{Feature: "image-gen", Amount: rollover.Ptr(int64(1))},
	})
	if err != nil {
		fmt.Printf("checkBatch failed: %v\n", err)
		return
	}

	var blocked []string
	for _, r := range gate.Results {
		if !r.Allowed {
			blocked = append(blocked, r.Feature)
		}
	}
	if len(blocked) > 0 {
		fmt.Printf("Blocked on: %s\n", strings.Join(blocked, ", "))
		fmt.Println("Please upgrade your plan to continue.")
		return
	}

	fmt.Println("All features available. Running pipeline...")
	fmt.Println("Pipeline completed.")

	// TrackBatch records every event in one call and groups the resulting
	// usage_events rows under a shared BatchID. AtomicityAllOrNothing rolls
	// the whole batch back if any event would block.
	result, err := ro.TrackBatch(ctx, wallet, []rollover.BatchTrackEvent{
		{Feature: "api-calls", Amount: 1},
		{Feature: "image-gen", Amount: 1},
	}, rollover.AtomicityAllOrNothing)
	if err != nil {
		fmt.Printf("trackBatch failed: %v\n", err)
		return
	}
	fmt.Printf("Usage tracked (batch %s).\n", result.BatchID)
}
