// Check and Track
//
// The core Rollover pattern is to verify a wallet has feature access before
// doing any work, then record usage after the operation succeeds.
//
//	ROLLOVER_API_KEY=ro_test_... go run ./examples/check-and-track
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
	wallet := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"

	result, err := ro.Check(ctx, wallet, "api-calls")
	if err != nil {
		log.Fatal(err)
	}

	if !result.Allowed {
		fmt.Printf("Limit reached. %d/%d used.\n", result.Used, result.Limit)
		return
	}

	fmt.Printf("Access granted. %d/%d remaining.\n", result.Remaining, result.Limit)

	// Do your work here...

	track, err := ro.Track(ctx, wallet, "api-calls", 1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Tracked. %d used, %d remaining.\n", track.Used, track.Remaining)
}
