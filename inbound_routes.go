package euromail

import (
	"context"
	"net/http"
	"net/url"
)

// CreateInboundRoute creates a new inbound routing rule.
func (c *Client) CreateInboundRoute(ctx context.Context, params CreateInboundRouteParams) (*InboundRoute, error) {
	wrapper, err := doJSON[dataResponse[InboundRoute]](c, ctx, http.MethodPost, "/v1/inbound-routes", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// GetInboundRoute retrieves an inbound route by ID.
func (c *Client) GetInboundRoute(ctx context.Context, id string) (*InboundRoute, error) {
	wrapper, err := doJSON[dataResponse[InboundRoute]](c, ctx, http.MethodGet, "/v1/inbound-routes/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// UpdateInboundRoute updates an existing inbound route.
func (c *Client) UpdateInboundRoute(ctx context.Context, id string, params UpdateInboundRouteParams) (*InboundRoute, error) {
	wrapper, err := doJSON[dataResponse[InboundRoute]](c, ctx, http.MethodPut, "/v1/inbound-routes/"+url.PathEscape(id), params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// DeleteInboundRoute deletes an inbound route by ID.
func (c *Client) DeleteInboundRoute(ctx context.Context, id string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v1/inbound-routes/"+url.PathEscape(id), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListInboundRoutes returns a paginated list of inbound routes.
func (c *Client) ListInboundRoutes(ctx context.Context, params *ListParams) ([]InboundRoute, *Pagination, error) {
	q := url.Values{}
	addListParams(q, params)

	path := "/v1/inbound-routes"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[InboundRoute]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}
