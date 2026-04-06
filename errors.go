package rollover

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Error code constants returned by the Rollover API.
const (
	ErrCodeInvalidAPIKey       = "invalid_api_key"
	ErrCodeUnauthorized        = "unauthorized"
	ErrCodeRateLimit           = "rate_limit_exceeded"
	ErrCodeNotFound            = "not_found"
	ErrCodeInsufficientCredits = "insufficient_credits"
	ErrCodeValidation          = "validation_error"
)

// Error represents an API error returned by the Rollover server.
type Error struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("rollover: %s (%d): %s", e.Code, e.StatusCode, e.Message)
}

// Temporary returns true if the error is likely transient and the request
// could succeed on retry, such as rate limits (429) or server errors (5xx).
func (e *Error) Temporary() bool {
	return e.StatusCode == http.StatusTooManyRequests || e.StatusCode >= 500
}

// IsErrorCode checks whether an error is a Rollover API error with the given
// error code, unwrapping the error chain as needed.
func IsErrorCode(err error, code string) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

func parseError(statusCode int, body []byte) *Error {
	apiErr := &Error{StatusCode: statusCode}

	if len(body) == 0 {
		apiErr.Code = "http_error"
		apiErr.Message = http.StatusText(statusCode)
		return apiErr
	}

	if err := json.Unmarshal(body, apiErr); err != nil {
		apiErr.Code = "unknown_error"
		apiErr.Message = fmt.Sprintf("unexpected response (HTTP %d): %s", statusCode, string(body))
	}

	return apiErr
}
