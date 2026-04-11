package euromail

import (
	"context"
	"net/http"
)

// GenerateInsights triggers an AI-powered operational insights report for this account.
//
// Analyses the last 7 days of operational data (bounce rates, complaint rates,
// throttled domains, FBL reports) and returns a structured report with up to three
// findings, each with severity, area, observation, and recommendation.
//
// Requires the account:admin scope. At most one report can be generated per account
// per 24-hour period.
func (c *Client) GenerateInsights(ctx context.Context) (*InsightReport, error) {
	report, err := doJSON[InsightReport](c, ctx, http.MethodPost, "/v1/insights/generate", struct{}{})
	if err != nil {
		return nil, err
	}
	return &report, nil
}
