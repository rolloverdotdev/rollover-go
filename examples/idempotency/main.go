// Idempotent Tracking
//
// Avoid double-counting in distributed systems by using idempotency keys,
// where the same Idempotency-Key always produces the same result.
//
//	ROLLOVER_API_KEY=ro_test_... go run ./examples/idempotency
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

	// Use a deterministic key tied to the operation being tracked.
	key := "order-12345-image-gen"

	// First call records the usage.
	r1, err := ro.Track(ctx, wallet, "api-calls", 1, rollover.WithIdempotencyKey(key))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("First:  used=%d remaining=%d\n", r1.Used, r1.Remaining)

	// Second call with same key returns the cached result.
	r2, err := ro.Track(ctx, wallet, "api-calls", 1, rollover.WithIdempotencyKey(key))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Second: used=%d remaining=%d (same as first, not double-counted)\n", r2.Used, r2.Remaining)
}
