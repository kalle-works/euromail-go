package euromail

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// CreateWebhook creates a new webhook endpoint.
func (c *Client) CreateWebhook(ctx context.Context, params CreateWebhookParams) (*Webhook, error) {
	wrapper, err := doJSON[dataResponse[Webhook]](c, ctx, http.MethodPost, "/v1/webhooks", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// GetWebhook retrieves a webhook by ID.
func (c *Client) GetWebhook(ctx context.Context, webhookID string) (*Webhook, error) {
	wrapper, err := doJSON[dataResponse[Webhook]](c, ctx, http.MethodGet, "/v1/webhooks/"+url.PathEscape(webhookID), nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// UpdateWebhook updates an existing webhook.
func (c *Client) UpdateWebhook(ctx context.Context, webhookID string, params UpdateWebhookParams) (*Webhook, error) {
	wrapper, err := doJSON[dataResponse[Webhook]](c, ctx, http.MethodPut, "/v1/webhooks/"+url.PathEscape(webhookID), params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// TestWebhook sends a test event to a webhook.
func (c *Client) TestWebhook(ctx context.Context, webhookID string) (*WebhookTestResponse, error) {
	wrapper, err := doJSON[dataResponse[WebhookTestResponse]](c, ctx, http.MethodPost, "/v1/webhooks/"+url.PathEscape(webhookID)+"/test", struct{}{})
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// DeleteWebhook deletes a webhook.
func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v1/webhooks/"+url.PathEscape(webhookID), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListWebhooks returns a paginated list of webhooks.
func (c *Client) ListWebhooks(ctx context.Context, params *ListParams) ([]Webhook, *Pagination, error) {
	q := url.Values{}
	addListParams(q, params)

	path := "/v1/webhooks"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[Webhook]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}

// WebhookSignatureTolerance is the maximum age of a webhook timestamp that
// VerifyWebhookSignature accepts, matching the server's own replay window.
const WebhookSignatureTolerance = 5 * time.Minute

// Sentinel errors returned by VerifyWebhookSignature. Use errors.Is to check
// for a specific failure reason.
var (
	// ErrWebhookSignatureMissingTimestamp is returned when the signature
	// header has no "t=" component.
	ErrWebhookSignatureMissingTimestamp = errors.New("euromail: webhook signature header is missing a timestamp")
	// ErrWebhookSignatureMissingSignature is returned when the signature
	// header has no "v1=" component.
	ErrWebhookSignatureMissingSignature = errors.New("euromail: webhook signature header is missing a v1 signature")
	// ErrWebhookSignatureExpired is returned when the timestamp in the
	// header is outside WebhookSignatureTolerance of the current time,
	// which also rejects replayed requests.
	ErrWebhookSignatureExpired = errors.New("euromail: webhook signature timestamp is outside the tolerance window")
	// ErrWebhookSignatureInvalid is returned when none of the v1 signatures
	// in the header match the payload signed with the given secret.
	ErrWebhookSignatureInvalid = errors.New("euromail: webhook signature is invalid")
)

// VerifyWebhookSignature verifies the X-Euromail-Signature header EuroMail
// attaches to every webhook delivery and returns the timestamp it was signed
// at.
//
// The header has the Stripe-style format "t=<unix_ts>,v1=<hex_hmac>". The
// signed message is "<unix_ts>.<raw_request_body>", HMAC-SHA256'd with the
// webhook's signing secret (the value returned as Webhook.Secret when the
// webhook was created). Comparison is constant-time (hmac.Equal) to avoid
// leaking the valid signature through response-timing side channels, and the
// timestamp must be within WebhookSignatureTolerance of now to reject replayed
// deliveries.
//
// payload must be the exact, unmodified request body bytes — re-marshaling
// JSON before verifying will usually change byte-for-byte output and make a
// genuine signature fail to verify.
func VerifyWebhookSignature(payload []byte, sigHeader string, secret string) (time.Time, error) {
	return verifyWebhookSignatureAt(payload, sigHeader, secret, time.Now())
}

func verifyWebhookSignatureAt(payload []byte, sigHeader string, secret string, now time.Time) (time.Time, error) {
	var timestamp int64
	haveTimestamp := false
	var signatures []string

	for _, part := range strings.Split(sigHeader, ",") {
		part = strings.TrimSpace(part)
		switch {
		case strings.HasPrefix(part, "t="):
			ts, err := strconv.ParseInt(strings.TrimPrefix(part, "t="), 10, 64)
			if err != nil {
				continue
			}
			timestamp = ts
			haveTimestamp = true
		case strings.HasPrefix(part, "v1="):
			signatures = append(signatures, strings.TrimPrefix(part, "v1="))
		}
	}

	if !haveTimestamp {
		return time.Time{}, ErrWebhookSignatureMissingTimestamp
	}
	if len(signatures) == 0 {
		return time.Time{}, ErrWebhookSignatureMissingSignature
	}

	signedAt := time.Unix(timestamp, 0).UTC()
	if age := now.Sub(signedAt); age > WebhookSignatureTolerance || age < -WebhookSignatureTolerance {
		return time.Time{}, ErrWebhookSignatureExpired
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strconv.FormatInt(timestamp, 10)))
	mac.Write([]byte("."))
	mac.Write(payload)
	expected := mac.Sum(nil)

	for _, sigHex := range signatures {
		got, err := hex.DecodeString(sigHex)
		if err != nil {
			continue
		}
		if hmac.Equal(got, expected) {
			return signedAt, nil
		}
	}

	return time.Time{}, ErrWebhookSignatureInvalid
}
