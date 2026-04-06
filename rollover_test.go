package rollover

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func testClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return New(WithAPIKey("ro_test_key"), WithBaseURL(srv.URL))
}

func TestNewDefaultMode(t *testing.T) {
	c := New(WithAPIKey("ro_test_abc"))
	if c.mode != "test" {
		t.Errorf("expected mode test, got %s", c.mode)
	}

	c = New(WithAPIKey("ro_live_abc"))
	if c.mode != "live" {
		t.Errorf("expected mode live, got %s", c.mode)
	}
}

func TestNewDefaultTimeout(t *testing.T) {
	c := New(WithAPIKey("ro_test_abc"))
	if c.httpClient.Timeout != 30*time.Second {
		t.Errorf("expected 30s timeout, got %s", c.httpClient.Timeout)
	}
}

func TestNewWithHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 5 * time.Second}
	c := New(WithAPIKey("ro_test_abc"), WithHTTPClient(custom))
	if c.httpClient != custom {
		t.Error("expected custom HTTP client")
	}
}

func TestResolveSlugCachesSuccess(t *testing.T) {
	calls := 0
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Write([]byte(`{"slug":"acme"}`))
	})

	ctx := context.Background()
	slug, err := c.resolveSlug(ctx)
	if err != nil || slug != "acme" {
		t.Fatalf("expected acme, got %s (err: %v)", slug, err)
	}

	slug, _ = c.resolveSlug(ctx)
	if slug != "acme" {
		t.Fatalf("expected cached acme, got %s", slug)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestResolveSlugRetriesAfterFailure(t *testing.T) {
	calls := 0
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			w.WriteHeader(500)
			w.Write([]byte(`{"code":"internal","message":"down"}`))
			return
		}
		w.Write([]byte(`{"slug":"acme"}`))
	})

	ctx := context.Background()
	_, err := c.resolveSlug(ctx)
	if err == nil {
		t.Fatal("expected error on first call")
	}

	slug, err := c.resolveSlug(ctx)
	if err != nil || slug != "acme" {
		t.Fatalf("expected acme on retry, got %s (err: %v)", slug, err)
	}
}

func TestAPIKeyHeader(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-API-Key"); got != "ro_test_key" {
			t.Errorf("expected ro_test_key, got %s", got)
		}
		w.Write([]byte(`{"allowed":true}`))
	})

	c.Check(context.Background(), "0xabc", "feature")
}
