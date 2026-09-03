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
	Type    string `json:"type,omitempty"`
	Message string `json:"message"`
	// DocsURL points to the relevant documentation page, when the API supplies
	// one. Empty for infrastructure errors (database/redis/internal) that have
	// no actionable docs page.
	DocsURL string `json:"docs_url,omitempty"`
	// RequestID is the value of the response's X-Request-Id header, when
	// present. Useful to hand to support when reporting an issue.
	RequestID string `json:"-"`
}

func (e *EuroMailError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("euromail: %d %s: %s (request_id=%s)", e.Status, e.Code, e.Message, e.RequestID)
	}
	return fmt.Sprintf("euromail: %d %s: %s", e.Status, e.Code, e.Message)
}

// AuthenticationError indicates an invalid or missing API key (HTTP 401).
type AuthenticationError struct {
	EuroMailError
}

// ValidationError indicates a request validation failure. The API returns
// this both as HTTP 422 (e.g. AttachmentTooLarge-style semantic failures)
// and as HTTP 400 (e.g. malformed send_at), always with code
// "VALIDATION_ERROR" / type "validation_error" — classification below keys
// off that, not the status code alone.
type ValidationError struct {
	EuroMailError
}

// RateLimitError indicates the request was rate limited (HTTP 429), including
// both plain rate limiting and quota-exceeded responses — both carry a
// Retry-After header the caller should honor.
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

// apiErrorBody mirrors the shape of an EuroMail API error response:
//
//	{"error": {"type": "...", "code": "...", "message": "...", "docs_url": "..."}}
//
// A flat (non-nested) body is also accepted as a fallback for forward/backward
// compatibility with older error shapes.
type apiErrorBody struct {
	Error *struct {
		Type    string `json:"type"`
		Code    string `json:"code"`
		Message string `json:"message"`
		DocsURL string `json:"docs_url"`
	} `json:"error"`
	Code    string `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
	DocsURL string `json:"docs_url"`
}

func errorFromResponse(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	var parsed apiErrorBody
	_ = json.Unmarshal(body, &parsed)

	code, errType, message, docsURL := parsed.Code, parsed.Type, parsed.Message, parsed.DocsURL
	if parsed.Error != nil {
		code, errType, message, docsURL = parsed.Error.Code, parsed.Error.Type, parsed.Error.Message, parsed.Error.DocsURL
	}
	if code == "" {
		code = "unknown"
	}
	if message == "" {
		message = http.StatusText(resp.StatusCode)
	}

	base := EuroMailError{
		Status:    resp.StatusCode,
		Code:      code,
		Type:      errType,
		Message:   message,
		DocsURL:   docsURL,
		RequestID: resp.Header.Get("X-Request-Id"),
	}

	// Validation failures are identified by code/type, not status: the API
	// returns them as both 400 (AppError::Validation, e.g. a bad send_at) and
	// 422 (ApiError::UnprocessableEntity). Checking status alone silently
	// downgrades 400 validation failures to a bare EuroMailError.
	if code == "VALIDATION_ERROR" || errType == "validation_error" {
		return &ValidationError{EuroMailError: base}
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
