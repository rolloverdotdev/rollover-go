// Credit Top-Up
//
// Monitor a wallet's credit balance on an interval and automatically grant
// more credits when the balance drops below a configured threshold.
//
//	ROLLOVER_API_KEY=ro_live_... go run ./examples/credit-topup
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	rollover "github.com/rolloverdotdev/rollover-go"
)

const (
	lowBalanceThreshold = 100
	topUpAmount         = 500
	checkInterval       = 30 * time.Second
)

var ro = rollover.New()

func main() {
	ctx := context.Background()
	wallet := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"

	fmt.Printf("Monitoring %s (threshold: %d, top-up: %d)\n",
		wallet[:12]+"...", lowBalanceThreshold, topUpAmount)

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	checkAndTopUp(ctx, wallet)
	for range ticker.C {
		checkAndTopUp(ctx, wallet)
	}
}

func checkAndTopUp(ctx context.Context, wallet string) {
	balance, err := ro.GetCredits(ctx, wallet)
	if err != nil {
		log.Printf("failed to check balance: %v", err)
		return
	}

	fmt.Printf("[%s] balance: %d", time.Now().Format("15:04:05"), balance.Balance)

	if balance.Balance < lowBalanceThreshold {
		fmt.Printf(" (low! granting %d credits)\n", topUpAmount)

		grant, err := ro.GrantCredits(ctx, wallet, topUpAmount,
			rollover.WithDescription("Auto top-up: balance below threshold"),
		)
		if err != nil {
			log.Printf("  top-up failed: %v", err)
			return
		}
		fmt.Printf("  new balance: %d\n", grant.Balance)
	} else {
		fmt.Println(" (ok)")
	}
}
