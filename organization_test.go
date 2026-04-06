package rollover

import (
	"context"
	"net/http"
	"testing"
)

func TestGetOrganization(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Write([]byte(`{"id":"org-1","name":"Acme","slug":"acme","webhook_url":"https://example.com/hook"}`))
	})

	result, err := c.GetOrganization(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result.Slug != "acme" || result.Name != "Acme" {
		t.Errorf("unexpected result: %+v", result)
	}
}
