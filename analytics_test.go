package rollover

import (
	"context"
	"net/http"
	"testing"
)

func TestGetAnalytics(t *testing.T) {
	c := testClient(t, orgHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("slug") != "acme" {
			t.Errorf("expected slug acme")
		}
		if r.URL.Query().Get("mode") != "test" {
			t.Errorf("expected mode test")
		}
		w.Write([]byte(`{"mrr":"99.99","active_subs":10,"total_revenue":"500.00","top_features":[{"feature_slug":"api-calls","total_used":1000}]}`))
	}))

	result, err := c.GetAnalytics(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result.MRR != "99.99" {
		t.Errorf("expected MRR 99.99, got %s", result.MRR)
	}
	if result.ActiveSubs != 10 {
		t.Errorf("expected 10 active subs, got %d", result.ActiveSubs)
	}
	if len(result.TopFeatures) != 1 || result.TopFeatures[0].TotalUsed != 1000 {
		t.Errorf("unexpected top features: %+v", result.TopFeatures)
	}
}
