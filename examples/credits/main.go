// Credit-Gated Access
//
// Protect expensive operations by requiring an available credit balance,
// with credits automatically deducted according to the feature's credit_cost.
//
//	ROLLOVER_API_KEY=ro_test_... go run ./examples/credits
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

	// Check credit balance.
	balance, err := ro.GetCredits(ctx, wallet)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Credit balance: %d\n", balance.Balance)

	// Grant credits.
	grant, err := ro.GrantCredits(ctx, wallet, 500,
		rollover.WithDescription("Welcome bonus"),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Granted %d credits. New balance: %d\n", grant.Granted, grant.Balance)

	// Check if the wallet can use a credit-gated feature.
	result, err := ro.Check(ctx, wallet, "image-gen")
	if err != nil {
		log.Fatal(err)
	}

	if !result.Allowed {
		fmt.Printf("Not enough credits. Balance: %d, cost: %d\n",
			result.CreditBalance, result.CreditCost)
		return
	}

	// Do the expensive work, then track usage.
	track, err := ro.Track(ctx, wallet, "image-gen", 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Tracked. Credits remaining: %d\n", track.CreditBalance)
}
