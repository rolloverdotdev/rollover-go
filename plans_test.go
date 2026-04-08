package rollover

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func orgHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/v1/organization") {
			w.Write([]byte(`{"slug":"acme"}`))
			return
		}
		next(w, r)
	}
}

func TestListPlans(t *testing.T) {
	c := testClient(t, orgHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("slug") != "acme" {
			t.Errorf("expected slug acme")
		}
		if r.URL.Query().Get("limit") != "5" {
			t.Errorf("expected limit 5")
		}
		w.Write([]byte(`{"data":[{"id":"1","slug":"starter","name":"Starter","price_usdc":"9.99","billing_period":"monthly"}],"total":1,"limit":5,"offset":0}`))
	}))

	result, err := c.ListPlans(context.Background(), &ListOptions{Limit: 5})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 || result.Data[0].Slug != "starter" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestGetPlan(t *testing.T) {
	c := testClient(t, orgHandler(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/starter") {
			t.Errorf("expected path ending in /starter, got %s", r.URL.Path)
		}
		w.Write([]byte(`{"id":"1","slug":"starter","name":"Starter","price_usdc":"9.99","billing_period":"monthly"}`))
	}))

	result, err := c.GetPlan(context.Background(), "starter")
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Starter" {
		t.Errorf("expected Starter, got %s", result.Name)
	}
}

func TestUpdatePlanPointerFields(t *testing.T) {
	c := testClient(t, orgHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		json.Unmarshal(body, &req)
		if req["name"] != "Updated" {
			t.Errorf("expected name Updated, got %v", req["name"])
		}
		if _, exists := req["description"]; exists {
			t.Error("expected description to be omitted")
		}
		w.Write([]byte(`{"id":"1","slug":"starter","name":"Updated","price_usdc":"9.99","billing_period":"monthly"}`))
	}))

	result, err := c.UpdatePlan(context.Background(), "starter", UpdatePlanParams{
		Name: Ptr("Updated"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Updated" {
		t.Errorf("expected Updated, got %s", result.Name)
	}
}

func TestUpdateFeatureUsesPatch(t *testing.T) {
	c := testClient(t, orgHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/starter/features/api-calls") {
			t.Errorf("expected feature path, got %s", r.URL.Path)
		}
		w.Write([]byte(`{"id":"f-1","feature_slug":"api-calls","name":"API Calls","limit_amount":20000,"reset_period":"monthly"}`))
	}))

	result, err := c.UpdateFeature(context.Background(), "starter", "api-calls", UpdateFeatureParams{
		LimitAmount: Ptr(20000),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.LimitAmount != 20000 {
		t.Errorf("expected limit_amount 20000, got %d", result.LimitAmount)
	}
}

func TestListPricingNoAuth(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/v1/pricing/") {
			t.Errorf("expected pricing path, got %s", r.URL.Path)
		}
		w.Write([]byte(`[{"id":"1","slug":"starter","name":"Starter","price_usdc":"9.99","billing_period":"monthly"}]`))
	})

	result, err := c.ListPricing(context.Background(), "acme")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 plan, got %d", len(result))
	}
}

func TestPathEscaping(t *testing.T) {
	c := testClient(t, orgHandler(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.RawPath, "plan%2Fslash") {
			t.Errorf("expected escaped slash in path, got %s", r.URL.RawPath)
		}
		w.Write([]byte(`{"id":"1","slug":"plan","name":"Plan","price_usdc":"0","billing_period":"monthly"}`))
	}))

	c.GetPlan(context.Background(), "plan/slash")
}
