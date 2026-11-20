package euromail

import (
	"context"
	"io"
	"net/http"
	"testing"
)

func TestGetDomainReadsDetectedProviderBesideRecords(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/domains/dom_1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"dom_1","domain":"example.com","dns_records":{"spf":{"type":"TXT","host":"example.com","value":"v=spf1 include:euromail.dev ~all"}},"detected_provider":"cloudflare","created_at":"2026-10-28T00:00:00Z","updated_at":"2026-10-28T00:00:00Z"}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	d, err := client.GetDomain(context.Background(), "dom_1")
	if err != nil {
		t.Fatalf("GetDomain: %v", err)
	}
	if d.DetectedProvider == nil || *d.DetectedProvider != "cloudflare" {
		t.Fatalf("detected_provider = %v, want cloudflare", d.DetectedProvider)
	}
	if got := d.DnsRecords["spf"].Type; got != "TXT" {
		t.Fatalf("dns_records[spf].type = %q, want TXT", got)
	}
}

func TestUpdateDomainPatchesOnlyRequireTLSAndReturnsIt(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/domains/dom_1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %s, want PATCH", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != `{"require_tls":true}` {
			t.Errorf("body = %s, want {\"require_tls\":true}", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"dom_1","domain":"example.com","require_tls":true}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	on := true
	d, err := client.UpdateDomain(context.Background(), "dom_1", UpdateDomainParams{RequireTLS: &on})
	if err != nil {
		t.Fatalf("UpdateDomain: %v", err)
	}
	if !d.RequireTLS {
		t.Fatal("RequireTLS = false, want true")
	}
}

func TestDomainWithoutRequireTLSParsesAsFalse(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/domains/dom_1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"dom_1","domain":"example.com"}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	d, err := client.GetDomain(context.Background(), "dom_1")
	if err != nil {
		t.Fatalf("GetDomain: %v", err)
	}
	if d.RequireTLS {
		t.Fatal("RequireTLS = true, want false when the field is absent")
	}
}
