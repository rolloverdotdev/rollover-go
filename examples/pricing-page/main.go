// Pricing Page
//
// Return your plans as JSON for a pricing page, with a single API call
// fetching each plan and its included features.
//
//	ROLLOVER_API_KEY=ro_live_... go run ./examples/pricing-page
//	curl http://localhost:8080/pricing
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	rollover "github.com/rolloverdotdev/rollover-go"
)

var ro = rollover.New()

func main() {
	org, err := ro.GetOrganization(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("GET /pricing", func(w http.ResponseWriter, r *http.Request) {
		plans, err := ro.ListPricing(r.Context(), org.Slug)
		if err != nil {
			http.Error(w, `{"error":"failed to load pricing"}`, http.StatusServiceUnavailable)
			return
		}

		type feature struct {
			Name  string `json:"name"`
			Limit int    `json:"limit,omitempty"`
		}
		type plan struct {
			Name      string    `json:"name"`
			Slug      string    `json:"slug"`
			PriceUSDC string    `json:"price_usdc"`
			Period    string    `json:"billing_period"`
			TrialDays int       `json:"trial_days,omitempty"`
			Features  []feature `json:"features"`
		}

		out := make([]plan, 0, len(plans))
		for _, p := range plans {
			pp := plan{
				Name:      p.Name,
				Slug:      p.Slug,
				PriceUSDC: p.PriceUSDC,
				Period:    p.BillingPeriod,
				TrialDays: p.TrialDays,
			}
			for _, f := range p.Features {
				name := ""
				if f.Feature != nil {
					name = f.Feature.Name
				}
				pp.Features = append(pp.Features, feature{Name: name, Limit: f.LimitAmount})
			}
			out = append(out, pp)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	})

	fmt.Println("Pricing server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
