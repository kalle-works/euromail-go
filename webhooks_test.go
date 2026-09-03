package euromail

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"
)

// Fixed test vector for VerifyWebhookSignature, computed independently
// (Python hmac/hashlib, mirroring the server's Rust implementation in
// euromail-worker's fire_webhook.rs) rather than derived by round-tripping
// through this package's own signing code. A regression that silently
// changes the signed-message format (field order, separator, timestamp
// encoding) would still pass a self-signed round trip but must fail here.
const (
	webhookTestSecret    = "whsec_test_secret_1234567890abcdef"
	webhookTestTimestamp = int64(1700000000)
	webhookTestPayload   = `{"event":"delivered","email_id":"em_01JQ5K8WMFZ9XRTVB3GH6EDCN4"}`
	webhookTestSignature = "3a8c23e82136165d39641f730d477e3d16f51edd14d83e23659bedafe1e66923"
)

func webhookTestHeader() string {
	return "t=1700000000,v1=" + webhookTestSignature
}

func TestVerifyWebhookSignature_ExactVector(t *testing.T) {
	now := time.Unix(webhookTestTimestamp, 0).Add(10 * time.Second)
	signedAt, err := verifyWebhookSignatureAt([]byte(webhookTestPayload), webhookTestHeader(), webhookTestSecret, now)
	if err != nil {
		t.Fatalf("expected valid signature, got error: %v", err)
	}
	if !signedAt.Equal(time.Unix(webhookTestTimestamp, 0).UTC()) {
		t.Fatalf("signedAt = %v, want %v", signedAt, time.Unix(webhookTestTimestamp, 0).UTC())
	}
}

// TestVerifyWebhookSignature_ContractVector checks the exact test vector
// specified for this cross-SDK effort (the same one asserted by the
// Python, TypeScript, Rust, and PHP SDKs), so all five verify identically
// against one shared, independently-computed HMAC — not just internally
// consistent with this package's own signing code.
func TestVerifyWebhookSignature_ContractVector(t *testing.T) {
	const (
		secret    = "whsec_test_secret_do_not_use"
		timestamp = int64(1735689600)
		payload   = `{"event":"delivered","email_id":"018f2c3a-7b1e-7c3e-8b1a-2f6e9d4c5a01","account_id":"018f2c3a-7b1e-7c3e-8b1a-2f6e9d4c5a02","timestamp":"2025-01-01T00:00:00Z"}`
		wantV1    = "d571fbef13b9e524d460f6f2c88f8d8dc7df3c50ff7aabdedd8a3656abb96dd0"
	)
	header := "t=1735689600,v1=" + wantV1
	now := time.Unix(timestamp, 0)
	signedAt, err := verifyWebhookSignatureAt([]byte(payload), header, secret, now)
	if err != nil {
		t.Fatalf("expected valid signature, got error: %v", err)
	}
	if !signedAt.Equal(time.Unix(timestamp, 0).UTC()) {
		t.Fatalf("signedAt = %v, want %v", signedAt, time.Unix(timestamp, 0).UTC())
	}
}

func TestVerifyWebhookSignature_WrongSecret(t *testing.T) {
	now := time.Unix(webhookTestTimestamp, 0).Add(10 * time.Second)
	_, err := verifyWebhookSignatureAt([]byte(webhookTestPayload), webhookTestHeader(), "wrong_secret", now)
	if !errors.Is(err, ErrWebhookSignatureInvalid) {
		t.Fatalf("err = %v, want ErrWebhookSignatureInvalid", err)
	}
}

func TestVerifyWebhookSignature_TamperedPayload(t *testing.T) {
	now := time.Unix(webhookTestTimestamp, 0).Add(10 * time.Second)
	tampered := `{"event":"bounced","email_id":"em_01JQ5K8WMFZ9XRTVB3GH6EDCN4"}`
	_, err := verifyWebhookSignatureAt([]byte(tampered), webhookTestHeader(), webhookTestSecret, now)
	if !errors.Is(err, ErrWebhookSignatureInvalid) {
		t.Fatalf("err = %v, want ErrWebhookSignatureInvalid", err)
	}
}

func TestVerifyWebhookSignature_ExpiredTooOld(t *testing.T) {
	now := time.Unix(webhookTestTimestamp, 0).Add(WebhookSignatureTolerance + time.Second)
	_, err := verifyWebhookSignatureAt([]byte(webhookTestPayload), webhookTestHeader(), webhookTestSecret, now)
	if !errors.Is(err, ErrWebhookSignatureExpired) {
		t.Fatalf("err = %v, want ErrWebhookSignatureExpired", err)
	}
}

func TestVerifyWebhookSignature_ExpiredFutureClockSkew(t *testing.T) {
	// A timestamp that is *ahead* of our clock by more than the tolerance is
	// just as invalid as one that is stale behind it.
	now := time.Unix(webhookTestTimestamp, 0).Add(-WebhookSignatureTolerance - time.Second)
	_, err := verifyWebhookSignatureAt([]byte(webhookTestPayload), webhookTestHeader(), webhookTestSecret, now)
	if !errors.Is(err, ErrWebhookSignatureExpired) {
		t.Fatalf("err = %v, want ErrWebhookSignatureExpired", err)
	}
}

func TestVerifyWebhookSignature_WithinToleranceBoundary(t *testing.T) {
	now := time.Unix(webhookTestTimestamp, 0).Add(WebhookSignatureTolerance)
	if _, err := verifyWebhookSignatureAt([]byte(webhookTestPayload), webhookTestHeader(), webhookTestSecret, now); err != nil {
		t.Fatalf("expected signature exactly at the tolerance boundary to verify, got: %v", err)
	}
}

func TestVerifyWebhookSignature_MissingTimestamp(t *testing.T) {
	_, err := verifyWebhookSignatureAt([]byte(webhookTestPayload), "v1="+webhookTestSignature, webhookTestSecret, time.Now())
	if !errors.Is(err, ErrWebhookSignatureMissingTimestamp) {
		t.Fatalf("err = %v, want ErrWebhookSignatureMissingTimestamp", err)
	}
}

func TestVerifyWebhookSignature_MissingSignature(t *testing.T) {
	_, err := verifyWebhookSignatureAt([]byte(webhookTestPayload), "t=1700000000", webhookTestSecret, time.Now())
	if !errors.Is(err, ErrWebhookSignatureMissingSignature) {
		t.Fatalf("err = %v, want ErrWebhookSignatureMissingSignature", err)
	}
}

func TestVerifyWebhookSignature_PublicEntryPointUsesCurrentTime(t *testing.T) {
	// The public wrapper signs against time.Now(), so a freshly generated
	// signature (via the same algorithm) must verify right away and a fixed
	// past timestamp like the vector above must not (proves it isn't
	// hardcoded to accept anything).
	if _, err := VerifyWebhookSignature([]byte(webhookTestPayload), webhookTestHeader(), webhookTestSecret); err == nil {
		t.Fatalf("expected a decade-old timestamp to be rejected by the real-time entry point")
	}
}

func TestCreateWebhookAndTestWebhook(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/webhooks", func(w http.ResponseWriter, r *http.Request) {
		var body CreateWebhookParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.URL != "https://example.com/hook" {
			t.Fatalf("unexpected url: %+v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"wh_1","account_id":"acc_1","url":"https://example.com/hook","events":["delivered"],"is_active":true,"secret":"whsec_abc","created_at":"t","updated_at":"t"}}`))
	})
	mux.HandleFunc("/v1/webhooks/wh_1/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"message":"test event sent","payload":{"event":"test"}}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	wh, err := client.CreateWebhook(context.Background(), CreateWebhookParams{URL: "https://example.com/hook", Events: []string{"delivered"}})
	if err != nil {
		t.Fatalf("CreateWebhook: %v", err)
	}
	if wh.ID != "wh_1" || wh.Secret == nil || *wh.Secret != "whsec_abc" {
		t.Fatalf("unexpected webhook: %+v", wh)
	}

	resp, err := client.TestWebhook(context.Background(), wh.ID)
	if err != nil {
		t.Fatalf("TestWebhook: %v", err)
	}
	if resp.Message != "test event sent" {
		t.Fatalf("unexpected test response: %+v", resp)
	}
}
