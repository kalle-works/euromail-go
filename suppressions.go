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
