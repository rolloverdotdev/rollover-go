// Error Handling
//
// Handle Rollover API errors by inspecting the status code and error code,
// allowing your application to respond differently to authentication failures,
// rate limits, and other error conditions.
//
//	ROLLOVER_API_KEY=ro_test_... go run ./examples/error-handling
package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	rollover "github.com/rolloverdotdev/rollover-go"
)

func main() {
	ro := rollover.New()
	ctx := context.Background()

	result, err := ro.Check(ctx, "0xinvalid", "api-calls")
	if err != nil {
		var roErr *rollover.Error
		if errors.As(err, &roErr) {
			fmt.Printf("API error: %s (status %d)\n", roErr.Message, roErr.StatusCode)

			if roErr.Temporary() {
				fmt.Println("This is a transient error, safe to retry.")
			}
		} else {
			fmt.Printf("Network or other error: %v\n", err)
		}
		return
	}

	fmt.Printf("Allowed: %v\n", result.Allowed)

	// Use IsErrorCode for clean checks without type assertions.
	_, err = ro.GetPlan(ctx, "nonexistent-plan")
	if rollover.IsErrorCode(err, rollover.ErrCodeNotFound) {
		fmt.Println("\nPlan not found (checked via IsErrorCode).")
	}

	// Error code constants work in switch statements too.
	_, err = ro.GrantCredits(ctx, "0xabc", -1)
	if err != nil {
		var roErr *rollover.Error
		if errors.As(err, &roErr) {
			switch roErr.Code {
			case rollover.ErrCodeValidation:
				fmt.Printf("\nValidation error: %s\n", roErr.Message)
			case rollover.ErrCodeUnauthorized:
				fmt.Println("\nCheck your API key.")
			default:
				fmt.Printf("\nUnexpected: %s\n", roErr.Message)
			}
		} else {
			log.Fatal(err)
		}
	}
}
