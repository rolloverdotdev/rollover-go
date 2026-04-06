package rollover

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestListSubscriptions(t *testing.T) {
	c := testClient(t, orgHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("wallet") != "0xabc" {
			t.Errorf("expected wallet filter")
		}
		if r.URL.Query().Get("status") != "active" {
			t.Errorf("expected status filter")
		}
		w.Write([]byte(`{"data":[{"id":"sub-1","wallet_address":"0xabc","plan_name":"Starter","status":"active"}],"total":1,"limit":20,"offset":0}`))
	}))

	result, err := c.ListSubscriptions(context.Background(), &ListOptions{Wallet: "0xabc", Status: "active"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 || result.Data[0].PlanName != "Starter" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestGetSubscription(t *testing.T) {
	c := testClient(t, orgHandler(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/sub-123") {
			t.Errorf("expected path ending in /sub-123, got %s", r.URL.Path)
		}
		w.Write([]byte(`{"id":"sub-123","wallet_address":"0xabc","plan_name":"Starter","status":"active"}`))
	}))

	result, err := c.GetSubscription(context.Background(), "sub-123")
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "sub-123" {
		t.Errorf("expected sub-123, got %s", result.ID)
	}
}

func TestCreateSubscription(t *testing.T) {
	c := testClient(t, orgHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(201)
		w.Write([]byte(`{"id":"sub-new","wallet_address":"0xabc","plan_name":"Starter","status":"active"}`))
	}))

	result, err := c.CreateSubscription(context.Background(), "0xabc", "starter")
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "sub-new" {
		t.Errorf("expected sub-new, got %s", result.ID)
	}
}

func TestCancelSubscription(t *testing.T) {
	c := testClient(t, orgHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.Write([]byte(`{"id":"sub-123","status":"active","cancel_at_end":true}`))
	}))

	result, err := c.CancelSubscription(context.Background(), "sub-123")
	if err != nil {
		t.Fatal(err)
	}
	if !result.CancelAtEnd {
		t.Error("expected cancel_at_end to be true")
	}
}
