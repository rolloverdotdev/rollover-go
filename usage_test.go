package rollover

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestCheck(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("wallet") != "0xabc" {
			t.Errorf("expected wallet 0xabc, got %s", r.URL.Query().Get("wallet"))
		}
		if r.URL.Query().Get("feature") != "api-calls" {
			t.Errorf("expected feature api-calls, got %s", r.URL.Query().Get("feature"))
		}
		w.Write([]byte(`{"allowed":true,"used":5,"remaining":95,"limit":100,"plan":"starter","credit_balance":50,"credit_cost":1}`))
	})

	result, err := c.Check(context.Background(), "0xabc", "api-calls")
	if err != nil {
		t.Fatal(err)
	}
	if !result.Allowed {
		t.Error("expected allowed")
	}
	if result.Used != 5 {
		t.Errorf("expected used 5, got %d", result.Used)
	}
	if result.Remaining != 95 {
		t.Errorf("expected remaining 95, got %d", result.Remaining)
	}
	if result.Limit != 100 {
		t.Errorf("expected limit 100, got %d", result.Limit)
	}
	if result.Plan != "starter" {
		t.Errorf("expected plan starter, got %s", result.Plan)
	}
}

func TestCheckMissingOptionalFields(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"allowed":false}`))
	})

	result, err := c.Check(context.Background(), "0xabc", "api-calls")
	if err != nil {
		t.Fatal(err)
	}
	if result.Allowed {
		t.Error("expected not allowed")
	}
	if result.Used != 0 || result.Remaining != 0 || result.Limit != 0 {
		t.Error("expected zero defaults for missing fields")
	}
}

func TestTrack(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content type")
		}

		body, _ := io.ReadAll(r.Body)
		var req struct {
			Wallet  string `json:"wallet"`
			Feature string `json:"feature"`
			Amount  int    `json:"amount"`
		}
		json.Unmarshal(body, &req)
		if req.Wallet != "0xabc" || req.Feature != "api-calls" || req.Amount != 3 {
			t.Errorf("unexpected body: %+v", req)
		}

		w.Write([]byte(`{"allowed":true,"used":8,"remaining":92}`))
	})

	result, err := c.Track(context.Background(), "0xabc", "api-calls", 3)
	if err != nil {
		t.Fatal(err)
	}
	if result.Used != 8 || result.Remaining != 92 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestTrackWithIdempotencyKey(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Idempotency-Key"); got != "order-123" {
			t.Errorf("expected idempotency key order-123, got %s", got)
		}
		w.Write([]byte(`{"allowed":true,"used":1,"remaining":99}`))
	})

	_, err := c.Track(context.Background(), "0xabc", "api-calls", 1, WithIdempotencyKey("order-123"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestTrackAutoIdempotencyKey(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("Idempotency-Key")
		if got == "" {
			t.Errorf("expected auto-generated idempotency key, got empty header")
		}
		if len(got) != 32 {
			t.Errorf("expected 32-char hex key, got %q (len %d)", got, len(got))
		}
		w.Write([]byte(`{"allowed":true,"used":1,"remaining":99}`))
	})

	_, err := c.Track(context.Background(), "0xabc", "api-calls", 1)
	if err != nil {
		t.Fatal(err)
	}
}
