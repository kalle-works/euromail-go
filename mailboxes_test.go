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

func TestReplyToMessage(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes/mbx_1/messages/msg_1/reply", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method=%s", r.Method)
		}
		raw, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(raw), `"text_body":"thanks"`) {
			t.Fatalf("body missing text_body: %s", raw)
		}
		if strings.Contains(string(raw), "html_body") {
			t.Fatalf("absent html_body should be omitted: %s", raw)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"data":{"id":"em_1","status":"queued","message_id":"<a@b>","to":"user@x","subject":"Re: hi"}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	res, err := client.ReplyToMessage(context.Background(), "mbx_1", "msg_1", ReplyToMessageParams{TextBody: String("thanks")})
	if err != nil {
		t.Fatalf("ReplyToMessage: %v", err)
	}
	if res.ID != "em_1" || res.Status != "queued" || res.Subject != "Re: hi" {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestListMailboxThreads(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes/mbx_1/threads", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("limit"); got != "10" {
			t.Fatalf("limit=%q", got)
		}
		if got := r.URL.Query().Get("offset"); got != "10" {
			t.Fatalf("offset=%q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"msg_1","mailbox_id":"mbx_1","account_id":"a","message_id":null,"mail_from":"f@x","from_header":null,"reply_to":null,"subject":"hi","text_body":null,"html_body":null,"size_bytes":0,"thread_id":"th_1","labels":[],"read_at":null,"created_at":"t"}],"pagination":{"page":2,"per_page":10,"total":11,"total_pages":2}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	msgs, pag, err := client.ListMailboxThreads(context.Background(), "mbx_1", &ListParams{Page: Int(2), PerPage: Int(10)})
	if err != nil {
		t.Fatalf("ListMailboxThreads: %v", err)
	}
	if len(msgs) != 1 || msgs[0].ThreadID == nil || *msgs[0].ThreadID != "th_1" || pag.Total != 11 {
		t.Fatalf("unexpected result: %+v %+v", msgs, pag)
	}
}

func TestGetMailboxThread(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes/mbx_1/threads/th_1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"msg_1","mailbox_id":"mbx_1","account_id":"a","message_id":null,"mail_from":"f@x","from_header":null,"reply_to":null,"subject":"hi","text_body":null,"html_body":null,"size_bytes":0,"thread_id":"th_1","labels":[],"read_at":null,"created_at":"t","in_reply_to":null},{"id":"msg_2","mailbox_id":"mbx_1","account_id":"a","message_id":null,"mail_from":"f@x","from_header":null,"reply_to":null,"subject":"re","text_body":null,"html_body":null,"size_bytes":0,"thread_id":"th_1","labels":[],"read_at":null,"created_at":"t2","in_reply_to":"<a@b>"}],"pagination":{"page":1,"per_page":100,"total":2,"total_pages":1}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	msgs, pag, err := client.GetMailboxThread(context.Background(), "mbx_1", "th_1", nil)
	if err != nil {
		t.Fatalf("GetMailboxThread: %v", err)
	}
	if len(msgs) != 2 || pag.Total != 2 {
		t.Fatalf("unexpected result: %+v %+v", msgs, pag)
	}
	if msgs[1].InReplyTo == nil || *msgs[1].InReplyTo != "<a@b>" {
		t.Fatalf("in_reply_to not decoded: %+v", msgs[1])
	}
}

func TestSearchMailboxMessages(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes/mbx_1/messages/search", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q"); got != "hello world" {
			t.Fatalf("q=%q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"msg_1","mailbox_id":"mbx_1","account_id":"a","message_id":null,"mail_from":"f@x","from_header":null,"reply_to":null,"subject":"hello","text_body":null,"html_body":null,"size_bytes":0,"thread_id":null,"labels":[],"read_at":null,"created_at":"t"}],"pagination":{"page":1,"per_page":20,"total":1,"total_pages":1}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	msgs, _, err := client.SearchMailboxMessages(context.Background(), "mbx_1", "hello world", nil)
	if err != nil {
		t.Fatalf("SearchMailboxMessages: %v", err)
	}
	if len(msgs) != 1 || msgs[0].ID != "msg_1" {
		t.Fatalf("unexpected result: %+v", msgs)
	}
}

func TestUpdateMessageLabels(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes/mbx_1/messages/msg_1/labels", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method=%s", r.Method)
		}
		raw, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(raw), `"labels":["urgent","vip"]`) {
			t.Fatalf("body missing labels: %s", raw)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"labels":["urgent","vip"]}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	labels, err := client.UpdateMessageLabels(context.Background(), "mbx_1", "msg_1", []string{"urgent", "vip"})
	if err != nil {
		t.Fatalf("UpdateMessageLabels: %v", err)
	}
	if len(labels) != 2 || labels[0] != "urgent" {
		t.Fatalf("unexpected labels: %+v", labels)
	}
}

func TestGetMessageAttachmentURLs(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes/mbx_1/messages/msg_1/attachments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"filename":"a.pdf","content_type":"application/pdf","size":1024,"url":"https://s3/a.pdf","expires_in_seconds":3600}]}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	atts, err := client.GetMessageAttachmentURLs(context.Background(), "mbx_1", "msg_1")
	if err != nil {
		t.Fatalf("GetMessageAttachmentURLs: %v", err)
	}
	if len(atts) != 1 || atts[0].Filename != "a.pdf" || atts[0].Size != 1024 || atts[0].ExpiresInSeconds != 3600 {
		t.Fatalf("unexpected attachments: %+v", atts)
	}
}

func TestGetMessageAttachmentURLs_FallbackMetadata(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes/mbx_1/messages/msg_1/attachments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Attachment bytes were never persisted: raw stored metadata, no url/expiry.
		_, _ = w.Write([]byte(`{"data":[{"filename":"a.pdf","content_type":"application/pdf"}]}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	atts, err := client.GetMessageAttachmentURLs(context.Background(), "mbx_1", "msg_1")
	if err != nil {
		t.Fatalf("GetMessageAttachmentURLs: %v", err)
	}
	if len(atts) != 1 || atts[0].URL != "" || atts[0].ExpiresInSeconds != 0 {
		t.Fatalf("expected empty url/expiry in fallback: %+v", atts)
	}
}

func TestListMailboxContacts(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes/mbx_1/contacts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"email":"a@x","display_name":"Alice","message_count":3,"last_seen":"2026-04-14T00:00:00Z"}],"pagination":{"page":1,"per_page":20,"total":1,"total_pages":1}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	contacts, _, err := client.ListMailboxContacts(context.Background(), "mbx_1", nil)
	if err != nil {
		t.Fatalf("ListMailboxContacts: %v", err)
	}
	if len(contacts) != 1 || contacts[0].Email != "a@x" || contacts[0].MessageCount != 3 {
		t.Fatalf("unexpected contacts: %+v", contacts)
	}
	if contacts[0].DisplayName == nil || *contacts[0].DisplayName != "Alice" {
		t.Fatalf("display_name not decoded: %+v", contacts[0])
	}
}

func TestGetMailboxAnalytics(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes/mbx_1/analytics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"total_messages":42,"unread_messages":5,"total_threads":10,"messages_today":2,"messages_this_week":9}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	a, err := client.GetMailboxAnalytics(context.Background(), "mbx_1")
	if err != nil {
		t.Fatalf("GetMailboxAnalytics: %v", err)
	}
	if a.TotalMessages != 42 || a.UnreadMessages != 5 || a.MessagesThisWeek != 9 {
		t.Fatalf("unexpected analytics: %+v", a)
	}
}

func TestUpdateAutoResponder(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/agent-mailboxes/mbx_1/auto-responder", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("method=%s", r.Method)
		}
		raw, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(raw), `"enabled":true`) {
			t.Fatalf("body missing enabled: %s", raw)
		}
		if !strings.Contains(string(raw), `"rules":[{"match":"all"`) {
			t.Fatalf("body missing rules: %s", raw)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"auto_responder_enabled":true,"auto_responder_rules":[{"match":"all","action":{"reply_text":"hi"}}]}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	cfg, err := client.UpdateAutoResponder(context.Background(), "mbx_1", UpdateAutoResponderParams{
		Enabled: Bool(true),
		Rules:   json.RawMessage(`[{"match":"all","action":{"reply_text":"hi"}}]`),
	})
	if err != nil {
		t.Fatalf("UpdateAutoResponder: %v", err)
	}
	if !cfg.AutoResponderEnabled {
		t.Fatalf("expected enabled: %+v", cfg)
	}
}
