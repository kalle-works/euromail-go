package euromail

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestSendEmailIncludesNewFields(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/emails", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body["transactional"] != false {
			t.Fatalf("transactional not forwarded: %+v", body)
		}
		if body["stream"] != "marketing" {
			t.Fatalf("stream not forwarded: %+v", body)
		}
		if body["send_at"] != "2026-08-01T00:00:00Z" {
			t.Fatalf("send_at not forwarded: %+v", body)
		}
		if body["tracking"] != true {
			t.Fatalf("tracking not forwarded: %+v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"email_1","message_id":"<m>","status":"queued","to":"c@d.com","sandbox":false,"scheduled_at":null,"created_at":"t"}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	_, err := client.SendEmail(context.Background(), SendEmailParams{
		From:          "a@b.com",
		To:            ToRecipient("c@d.com"),
		Subject:       String("Marketing update"),
		TextBody:      String("News"),
		Transactional: Bool(false),
		Stream:        String("marketing"),
		SendAt:        String("2026-08-01T00:00:00Z"),
		Tracking:      Bool(true),
	})
	if err != nil {
		t.Fatalf("SendEmail: %v", err)
	}
}

func TestSendBroadcastIncludesTrackingAndTransactional(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/emails/broadcast", func(w http.ResponseWriter, r *http.Request) {
		var body BroadcastParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Transactional == nil || *body.Transactional != true {
			t.Fatalf("transactional not forwarded: %+v", body)
		}
		if body.Tracking == nil || *body.Tracking != false {
			t.Fatalf("tracking not forwarded: %+v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"operation_id":"op_1","total_recipients":3,"message":"queued"}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	resp, err := client.SendBroadcast(context.Background(), BroadcastParams{
		ContactListID: "cl_001",
		FromAddress:   "sender@example.com",
		Subject:       String("Migration notice"),
		TextBody:      String("We moved!"),
		Transactional: Bool(true),
		Tracking:      Bool(false),
	})
	if err != nil {
		t.Fatalf("SendBroadcast: %v", err)
	}
	if resp.TotalRecipients != 3 {
		t.Fatalf("unexpected total_recipients: %d", resp.TotalRecipients)
	}
}
