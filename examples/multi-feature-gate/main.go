// Multi-Feature Gate
//
// Check multiple features concurrently before starting an operation that
// requires all of them, such as an AI pipeline consuming both API calls
// and image generation credits.
//
//	ROLLOVER_API_KEY=ro_live_... go run ./examples/multi-feature-gate
package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	rollover "github.com/rolloverdotdev/rollover-go"
)

var ro = rollover.New()

// checkAll checks multiple features concurrently and returns the list of
// features that are blocked (empty if all are allowed).
func checkAll(ctx context.Context, wallet string, features []string) []string {
	var mu sync.Mutex
	var blocked []string
	var wg sync.WaitGroup

	for _, f := range features {
		wg.Add(1)
		go func(feature string) {
			defer wg.Done()
			result, err := ro.Check(ctx, wallet, feature)
			if err != nil || !result.Allowed {
				mu.Lock()
				blocked = append(blocked, feature)
				mu.Unlock()
			}
		}(f)
	}

	wg.Wait()
	return blocked
}

// trackAll tracks usage for multiple features concurrently.
func trackAll(ctx context.Context, wallet string, features map[string]int) {
	var wg sync.WaitGroup
	for feature, amount := range features {
		wg.Add(1)
		go func(f string, a int) {
			defer wg.Done()
			if _, err := ro.Track(ctx, wallet, f, a); err != nil {
				log.Printf("rollover: track %s failed: %v", f, err)
			}
		}(feature, amount)
	}
	wg.Wait()
}

func main() {
	ctx := context.Background()
	wallet := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"

	// This operation requires both api-calls and image-gen.
	required := []string{"api-calls", "image-gen"}

	blocked := checkAll(ctx, wallet, required)
	if len(blocked) > 0 {
		fmt.Printf("Blocked on: %s\n", strings.Join(blocked, ", "))
		fmt.Println("Please upgrade your plan to continue.")
		return
	}

	fmt.Println("All features available. Running pipeline...")
	fmt.Println("Pipeline completed.")

	trackAll(ctx, wallet, map[string]int{
		"api-calls": 1,
		"image-gen": 1,
	})
	fmt.Println("Usage tracked for all features.")
}
