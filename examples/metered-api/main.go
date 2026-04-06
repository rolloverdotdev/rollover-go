// Metered API Server
//
// Track usage for multiple features across different routes, with each
// route mapped to a Rollover feature.
//
//	ROLLOVER_API_KEY=ro_test_... go run ./examples/metered-api
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	rollover "github.com/rolloverdotdev/rollover-go"
)

var ro = rollover.New()

func metered(feature string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wallet := r.Header.Get("X-Wallet")

		result, err := ro.Check(r.Context(), wallet, feature)
		if err != nil || !result.Allowed {
			http.Error(w, "rate limited", http.StatusTooManyRequests)
			return
		}

		handler(w, r)

		ro.Track(context.Background(), wallet, feature, 1)
	}
}

func main() {
	http.HandleFunc("POST /v1/translate", metered("translations", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"text":"translated"}`)
	}))
	http.HandleFunc("POST /v1/summarize", metered("summaries", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"text":"summarized"}`)
	}))
	http.HandleFunc("POST /v1/embeddings", metered("embeddings", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"embeddings":[0.1,0.2]}`)
	}))

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
