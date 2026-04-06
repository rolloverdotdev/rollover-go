package rollover

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestParseErrorJSON(t *testing.T) {
	err := parseError(400, []byte(`{"code":"validation_error","message":"invalid wallet"}`))
	if err.StatusCode != 400 {
		t.Errorf("expected 400, got %d", err.StatusCode)
	}
	if err.Code != "validation_error" {
		t.Errorf("expected validation_error, got %s", err.Code)
	}
	if err.Message != "invalid wallet" {
		t.Errorf("expected invalid wallet, got %s", err.Message)
	}
}

func TestParseErrorHTML(t *testing.T) {
	err := parseError(502, []byte("<html>Bad Gateway</html>"))
	if err.Code != "unknown_error" {
		t.Errorf("expected unknown_error, got %s", err.Code)
	}
}

func TestParseErrorEmpty(t *testing.T) {
	err := parseError(503, nil)
	if err.Code != "http_error" {
		t.Errorf("expected http_error, got %s", err.Code)
	}
	if err.Message != "Service Unavailable" {
		t.Errorf("expected Service Unavailable, got %s", err.Message)
	}
}

func TestTemporary(t *testing.T) {
	tests := []struct {
		status   int
		expected bool
	}{
		{429, true},
		{500, true},
		{502, true},
		{400, false},
		{401, false},
		{404, false},
	}
	for _, tt := range tests {
		err := &Error{StatusCode: tt.status}
		if err.Temporary() != tt.expected {
			t.Errorf("Temporary() for %d: expected %v", tt.status, tt.expected)
		}
	}
}

func TestIsErrorCode(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"code":"not_found","message":"plan not found"}`))
	})

	_, err := c.Check(context.Background(), "0xabc", "feature")
	if !IsErrorCode(err, ErrCodeNotFound) {
		t.Error("expected IsErrorCode to match not_found")
	}
	if IsErrorCode(err, ErrCodeRateLimit) {
		t.Error("expected IsErrorCode to not match rate_limit_exceeded")
	}
}

func TestErrorAsUnwrap(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"code":"unauthorized","message":"bad key"}`))
	})

	_, err := c.Check(context.Background(), "0xabc", "feature")
	var roErr *Error
	if !errors.As(err, &roErr) {
		t.Fatal("expected errors.As to work")
	}
	if roErr.StatusCode != 401 {
		t.Errorf("expected 401, got %d", roErr.StatusCode)
	}
}
