package euromail

import (
	"context"
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
