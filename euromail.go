// Package euromail provides a Go client for the EuroMail email delivery API.
//
// Usage:
//
//	client := euromail.NewClient("em_live_your_api_key")
//	resp, err := client.SendEmail(ctx, euromail.SendEmailParams{
//	    From:    "sender@example.com",
//	    To:      "recipient@example.com",
//	    Subject: euromail.String("Hello"),
//	    HTMLBody: euromail.String("<h1>Hello World</h1>"),
//	})
package euromail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.euromail.dev"
	defaultTimeout = 30 * time.Second
)

// Client is the EuroMail API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Option configures the Client.
type Option func(*Client)

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// NewClient creates a new EuroMail API client.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	if !strings.HasPrefix(c.baseURL, "https://") && !strings.HasPrefix(c.baseURL, "http://localhost") && !strings.HasPrefix(c.baseURL, "http://127.0.0.1") {
		fmt.Fprintln(os.Stderr, "WARNING: EuroMail base URL does not use HTTPS. API keys will be sent in cleartext.")
	}
	return c
}

// String returns a pointer to the given string value. Useful for optional fields.
func String(s string) *string { return &s }

// Int returns a pointer to the given int value. Useful for optional fields.
func Int(i int) *int { return &i }

// Bool returns a pointer to the given bool value. Useful for optional fields.
func Bool(b bool) *bool { return &b }

// paginatedResponse is a generic wrapper for paginated list responses.
type paginatedResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// dataResponse is a generic wrapper for single-item responses.
type dataResponse[T any] struct {
	Data T `json:"data"`
}

func (c *Client) do(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	u := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("euromail: failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reqBody)
	if err != nil {
		return nil, fmt.Errorf("euromail: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("euromail: request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, errorFromResponse(resp)
	}

	return resp, nil
}

func doJSON[T any](c *Client, ctx context.Context, method, path string, body interface{}) (T, error) {
	var zero T
	resp, err := c.do(ctx, method, path, body)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return zero, nil
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return zero, fmt.Errorf("euromail: failed to decode response: %w", err)
	}
	return result, nil
}

func doRawText(c *Client, ctx context.Context, method, path string) (string, error) {
	resp, err := c.do(ctx, method, path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("euromail: failed to read response body: %w", err)
	}
	return string(b), nil
}

func addListParams(q url.Values, params *ListParams) {
	if params == nil {
		return
	}
	if params.Page != nil {
		q.Set("page", strconv.Itoa(*params.Page))
	}
	if params.PerPage != nil {
		q.Set("per_page", strconv.Itoa(*params.PerPage))
	}
}

// FlexFloat64 unmarshals a JSON float that may be encoded as a string (e.g. "0.0").
type FlexFloat64 float64

func (f *FlexFloat64) UnmarshalJSON(data []byte) error {
	var v float64
	if err := json.Unmarshal(data, &v); err == nil {
		*f = FlexFloat64(v)
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("FlexFloat64: cannot unmarshal %s", string(data))
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("FlexFloat64: cannot parse %q: %w", s, err)
	}
	*f = FlexFloat64(v)
	return nil
}

func intToStr(i int) string {
	return strconv.Itoa(i)
}

func addAnalyticsQuery(q url.Values, aq *AnalyticsQuery) {
	if aq == nil {
		return
	}
	if aq.Period != nil {
		q.Set("period", *aq.Period)
	}
	if aq.From != nil {
		q.Set("from", *aq.From)
	}
	if aq.To != nil {
		q.Set("to", *aq.To)
	}
}
