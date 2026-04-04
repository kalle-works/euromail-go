package euromail

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// EuroMailError represents an API error returned by EuroMail.
type EuroMailError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *EuroMailError) Error() string {
	return fmt.Sprintf("euromail: %d %s: %s", e.Status, e.Code, e.Message)
}

// AuthenticationError indicates an invalid or missing API key (HTTP 401).
type AuthenticationError struct {
	EuroMailError
}

// ValidationError indicates a request validation failure (HTTP 422).
type ValidationError struct {
	EuroMailError
}

// RateLimitError indicates the request was rate limited (HTTP 429).
type RateLimitError struct {
	EuroMailError
	RetryAfter int // seconds until retry is allowed, 0 if unknown
}

// NotFoundError indicates the requested resource was not found (HTTP 404).
type NotFoundError struct {
	EuroMailError
}

// IsAuthenticationError returns true if the error is an authentication error.
func IsAuthenticationError(err error) bool {
	_, ok := err.(*AuthenticationError)
	return ok
}

// IsValidationError returns true if the error is a validation error.
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// IsRateLimitError returns true if the error is a rate limit error.
func IsRateLimitError(err error) bool {
	_, ok := err.(*RateLimitError)
	return ok
}

// IsNotFoundError returns true if the error is a not found error.
func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

func errorFromResponse(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	var apiErr struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &apiErr); err != nil {
		apiErr.Code = "unknown"
		apiErr.Message = http.StatusText(resp.StatusCode)
	}
	if apiErr.Code == "" {
		apiErr.Code = "unknown"
	}
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(resp.StatusCode)
	}

	base := EuroMailError{
		Status:  resp.StatusCode,
		Code:    apiErr.Code,
		Message: apiErr.Message,
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return &AuthenticationError{EuroMailError: base}
	case http.StatusUnprocessableEntity:
		return &ValidationError{EuroMailError: base}
	case http.StatusNotFound:
		return &NotFoundError{EuroMailError: base}
	case http.StatusTooManyRequests:
		retryAfter := 0
		if h := resp.Header.Get("Retry-After"); h != "" {
			if v, err := strconv.Atoi(h); err == nil {
				retryAfter = v
			}
		}
		return &RateLimitError{EuroMailError: base, RetryAfter: retryAfter}
	default:
		return &base
	}
}
