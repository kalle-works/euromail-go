package euromail

import (
	"context"
	"net/http"
	"net/url"
)

// ListInboundEmails returns a paginated list of inbound emails.
func (c *Client) ListInboundEmails(ctx context.Context, params *ListParams) ([]InboundEmail, *Pagination, error) {
	q := url.Values{}
	addListParams(q, params)

	path := "/v1/inbound"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[InboundEmail]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}

// GetInboundEmail retrieves a single inbound email by ID.
func (c *Client) GetInboundEmail(ctx context.Context, id string) (*InboundEmail, error) {
	wrapper, err := doJSON[dataResponse[InboundEmail]](c, ctx, http.MethodGet, "/v1/inbound/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// DeleteInboundEmail deletes an inbound email by ID.
func (c *Client) DeleteInboundEmail(ctx context.Context, id string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v1/inbound/"+url.PathEscape(id), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
