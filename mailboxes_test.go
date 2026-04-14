package euromail

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, handler http.Handler) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	c := NewClient("em_test_key", WithBaseURL(srv.URL))
	return c, srv
}

func TestCreateMailbox(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer em_test_key" {
			t.Fatalf("unexpected auth header: %q", got)
		}
		var body CreateMailboxParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.DisplayName == nil || *body.DisplayName != "Support" {
			t.Fatalf("display_name not forwarded: %+v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"mbx_1","account_id":"acc_1","local_part":"agent","domain":"example.com","address":"agent@example.com","display_name":"Support","created_at":"2026-04-13T00:00:00Z"}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	mb, err := client.CreateMailbox(context.Background(), CreateMailboxParams{DisplayName: String("Support")})
	if err != nil {
		t.Fatalf("CreateMailbox: %v", err)
	}
	if mb.ID != "mbx_1" || mb.Address != "agent@example.com" {
		t.Fatalf("unexpected mailbox: %+v", mb)
	}
}

func TestListMailboxes(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("limit"); got != "10" {
			t.Fatalf("limit=%q", got)
		}
		if got := r.URL.Query().Get("offset"); got != "10" {
			t.Fatalf("offset=%q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"mbx_1","account_id":"a","local_part":"x","domain":"d","address":"x@d","display_name":null,"created_at":"t"}],"pagination":{"page":2,"per_page":10,"total":11,"total_pages":2}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	mbs, pag, err := client.ListMailboxes(context.Background(), &ListParams{Page: Int(2), PerPage: Int(10)})
	if err != nil {
		t.Fatalf("ListMailboxes: %v", err)
	}
	if len(mbs) != 1 || pag.Total != 11 {
		t.Fatalf("unexpected result: %+v %+v", mbs, pag)
	}
}

func TestWaitForNextMessage_200(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes/mbx_1/messages/next", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("timeout"); got != "5" {
			t.Fatalf("timeout=%q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"msg_1","mailbox_id":"mbx_1","account_id":"a","message_id":null,"mail_from":"from@x","from_header":null,"reply_to":null,"subject":"hi","text_body":null,"html_body":null,"size_bytes":0,"thread_id":null,"labels":[],"read_at":null,"created_at":"t"},"lease_token":"lt_1","lease_expires_at":"t2"}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	leased, err := client.WaitForNextMessage(context.Background(), "mbx_1", Int(5))
	if err != nil {
		t.Fatalf("WaitForNextMessage: %v", err)
	}
	if leased == nil || leased.LeaseToken != "lt_1" || leased.Data.ID != "msg_1" {
		t.Fatalf("unexpected: %+v", leased)
	}
}

func TestWaitForNextMessage_408(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes/mbx_1/messages/next", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusRequestTimeout)
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	leased, err := client.WaitForNextMessage(context.Background(), "mbx_1", nil)
	if err != nil {
		t.Fatalf("expected nil error on 408, got %v", err)
	}
	if leased != nil {
		t.Fatalf("expected nil LeasedMessage on 408, got %+v", leased)
	}
}

func TestAckMessage(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes/mbx_1/messages/msg_1/ack", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method=%s", r.Method)
		}
		raw, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(raw), `"lease_token":"lt_1"`) {
			t.Fatalf("body missing lease_token: %s", raw)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	if err := client.AckMessage(context.Background(), "mbx_1", "msg_1", "lt_1"); err != nil {
		t.Fatalf("AckMessage: %v", err)
	}
}
