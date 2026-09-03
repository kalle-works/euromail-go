package euromail

import (
	"context"
	"net/http"
	"net/url"
)

// AddSuppression adds an email address to the suppression list.
// If reason is empty, it defaults to "manual".
func (c *Client) AddSuppression(ctx context.Context, email string, reason string) (*Suppression, error) {
	if reason == "" {
		reason = "manual"
	}
	body := map[string]string{
		"email_address": email,
		"reason":        reason,
	}
	wrapper, err := doJSON[dataResponse[Suppression]](c, ctx, http.MethodPost, "/v1/suppressions", body)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// DeleteSuppression removes an email address from the suppression list.
func (c *Client) DeleteSuppression(ctx context.Context, email string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v1/suppressions/"+url.PathEscape(email), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListSuppressions returns a paginated list of suppressions.
func (c *Client) ListSuppressions(ctx context.Context, params *ListParams) ([]Suppression, *Pagination, error) {
	q := url.Values{}
	addListParams(q, params)

	path := "/v1/suppressions"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[Suppression]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}

// ImportSuppressions bulk-adds up to 10,000 email addresses to the
// suppression list in a single request. Invalid addresses are skipped rather
// than failing the whole request; check ImportSuppressionsResult.InvalidAddresses
// to see which ones were rejected.
func (c *Client) ImportSuppressions(ctx context.Context, params ImportSuppressionsParams) (*ImportSuppressionsResult, error) {
	wrapper, err := doJSON[dataResponse[ImportSuppressionsResult]](c, ctx, http.MethodPost, "/v1/suppressions/import", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// ExportSuppressions downloads the account's full suppression list as CSV
// (columns: email_address, reason, created_at). Large lists are streamed by
// the server, so this may take a while for accounts with a long history.
func (c *Client) ExportSuppressions(ctx context.Context) (string, error) {
	return doRawText(c, ctx, http.MethodGet, "/v1/suppressions/export")
}
