// Usage Middleware
//
// An HTTP middleware that gates endpoints by verifying usage before
// handling the request and recording consumption after a successful response.
//
//	ROLLOVER_API_KEY=ro_test_... go run ./examples/middleware
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	rollover "github.com/rolloverdotdev/rollover-go"
)

var ro = rollover.New()

func usageGate(feature string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wallet := r.Header.Get("X-Wallet")
		if wallet == "" {
			http.Error(w, "missing X-Wallet header", http.StatusBadRequest)
			return
		}

		result, err := ro.Check(r.Context(), wallet, feature)
		if err != nil {
			http.Error(w, "usage check failed", http.StatusInternalServerError)
			return
		}

		if !result.Allowed {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintf(w, `{"error":"limit reached","used":%d,"limit":%d}`, result.Used, result.Limit)
			return
		}

		next(w, r)

		ro.Track(context.Background(), wallet, feature, 1)
	}
}

func handleTranslate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"result":"translated"}`)
}

func main() {
	http.HandleFunc("/v1/translate", usageGate("translations", handleTranslate))

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
