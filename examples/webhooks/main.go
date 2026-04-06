// Webhook Receiver
//
// Process real-time events from Rollover by registering a webhook URL
// in the dashboard and handling events as they arrive.
//
//	go run ./examples/webhooks
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type WebhookEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type SubscriptionData struct {
	WalletAddress string `json:"wallet_address"`
	PlanName      string `json:"plan_name"`
	Status        string `json:"status"`
}

func main() {
	http.HandleFunc("POST /webhooks/rollover", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		var event WebhookEvent
		if err := json.Unmarshal(body, &event); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		switch event.Type {
		case "subscription.created":
			var data SubscriptionData
			json.Unmarshal(event.Data, &data)
			fmt.Printf("New subscription: %s -> %s\n", data.WalletAddress, data.PlanName)

		case "subscription.canceled":
			var data SubscriptionData
			json.Unmarshal(event.Data, &data)
			fmt.Printf("Canceled: %s from %s\n", data.WalletAddress, data.PlanName)

		default:
			fmt.Printf("Received event: %s\n", event.Type)
		}

		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Webhook receiver on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
