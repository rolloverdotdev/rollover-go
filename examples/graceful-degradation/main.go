// Graceful Degradation
//
// Return a helpful 429 response with usage details and an upgrade path when
// a wallet hits its limit, rather than a generic rate limit error.
//
//	ROLLOVER_API_KEY=ro_live_... go run ./examples/graceful-degradation
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	rollover "github.com/rolloverdotdev/rollover-go"
)

var ro = rollover.New()

func main() {
	http.HandleFunc("POST /generate", func(w http.ResponseWriter, r *http.Request) {
		wallet := r.Header.Get("X-Wallet")
		if wallet == "" {
			writeJSON(w, 401, map[string]string{"error": "wallet required"})
			return
		}

		result, err := ro.Check(r.Context(), wallet, "generations")
		if err != nil {
			log.Printf("billing check failed: %v (failing open)", err)
			doGenerate(w)
			return
		}

		if result.Allowed {
			doGenerate(w)
			go ro.Track(context.Background(), wallet, "generations", 1)
			return
		}

		// Limit reached: return a helpful response with upgrade info.
		writeJSON(w, http.StatusTooManyRequests, map[string]any{
			"error":   "generation limit reached",
			"used":    result.Used,
			"limit":   result.Limit,
			"plan":    result.Plan,
			"upgrade": "https://app.example.com/billing/upgrade",
			"message": "You've used all your generations for this period. Upgrade your plan for more.",
		})
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func doGenerate(w http.ResponseWriter) {
	writeJSON(w, 200, map[string]string{"result": "generated content here"})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
