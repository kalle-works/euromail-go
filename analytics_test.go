package euromail

import (
	"context"
	"net/http"
	"testing"
)

// A domain whose mail had no tracking has no rate: nil, not 0.
func TestGetAnalyticsDomainsReadsAMissingRate(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/analytics/domains", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"domain":"example.com","sent":500,"delivered":490,"bounced":10,"opens_tracked":0,"clicks_tracked":100,"open_rate":null,"click_rate":"12.5"}],"period":{"from":"2026-02-07","to":"2026-03-07"}}`))
	})

	client, srv := newTestClient(t, mux)
	defer srv.Close()

	result, err := client.GetAnalyticsDomains(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetAnalyticsDomains: %v", err)
	}
	d := result.Data[0]
	if d.OpenRate != nil {
		t.Fatalf("open_rate = %v, want nil", *d.OpenRate)
	}
	if d.ClickRate == nil || *d.ClickRate != 12.5 {
		t.Fatalf("click_rate = %v, want 12.5", d.ClickRate)
	}
	if d.ClicksTracked != 100 {
		t.Fatalf("clicks_tracked = %d, want 100", d.ClicksTracked)
	}
}
