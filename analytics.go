package euromail

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// GetAnalyticsOverview returns summary analytics for the account.
func (c *Client) GetAnalyticsOverview(ctx context.Context, query *AnalyticsQuery) (*AnalyticsSummary, error) {
	q := url.Values{}
	addAnalyticsQuery(q, query)

	path := "/v1/analytics/overview"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[AnalyticsSummary](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAnalyticsTimeseries returns time series analytics data.
func (c *Client) GetAnalyticsTimeseries(ctx context.Context, query *TimeseriesQuery) (*TimeseriesResponse, error) {
	q := url.Values{}
	if query != nil {
		addAnalyticsQuery(q, &query.AnalyticsQuery)
		if query.Metrics != nil {
			q.Set("metrics", *query.Metrics)
		}
	}

	path := "/v1/analytics/timeseries"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[TimeseriesResponse](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAnalyticsDomains returns per-domain analytics data.
func (c *Client) GetAnalyticsDomains(ctx context.Context, query *DomainAnalyticsQuery) (*DomainAnalyticsResponse, error) {
	q := url.Values{}
	if query != nil {
		addAnalyticsQuery(q, &query.AnalyticsQuery)
		if query.Limit != nil {
			q.Set("limit", strconv.Itoa(*query.Limit))
		}
	}

	path := "/v1/analytics/domains"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[DomainAnalyticsResponse](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ExportAnalyticsCSV returns analytics data as a CSV string.
func (c *Client) ExportAnalyticsCSV(ctx context.Context, query *AnalyticsQuery) (string, error) {
	q := url.Values{}
	addAnalyticsQuery(q, query)

	path := "/v1/analytics/export"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	return doRawText(c, ctx, http.MethodGet, path)
}
