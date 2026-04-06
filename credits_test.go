package rollover

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestGetCredits(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("wallet") != "0xabc" {
			t.Errorf("expected wallet 0xabc")
		}
		w.Write([]byte(`{"wallet":"0xabc","balance":250}`))
	})

	result, err := c.GetCredits(context.Background(), "0xabc")
	if err != nil {
		t.Fatal(err)
	}
	if result.Balance != 250 {
		t.Errorf("expected balance 250, got %d", result.Balance)
	}
}

func TestGrantCredits(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Wallet      string `json:"wallet"`
			Amount      int    `json:"amount"`
			Description string `json:"description"`
		}
		json.Unmarshal(body, &req)
		if req.Wallet != "0xabc" || req.Amount != 500 || req.Description != "bonus" {
			t.Errorf("unexpected body: %+v", req)
		}
		w.Write([]byte(`{"balance":750,"granted":500}`))
	})

	result, err := c.GrantCredits(context.Background(), "0xabc", 500, WithDescription("bonus"))
	if err != nil {
		t.Fatal(err)
	}
	if result.Granted != 500 || result.Balance != 750 {
		t.Errorf("unexpected result: %+v", result)
	}
}
