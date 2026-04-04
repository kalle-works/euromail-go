package euromail

import (
	"context"
	"net/http"
	"net/url"
)

// ListOperations returns a paginated list of async operations.
func (c *Client) ListOperations(ctx context.Context, params *ListParams) ([]Operation, *Pagination, error) {
	q := url.Values{}
	addListParams(q, params)
	path := "/v1/operations"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	result, err := doJSON[paginatedResponse[Operation]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}

// GetOperation returns a specific operation by ID.
func (c *Client) GetOperation(ctx context.Context, id string) (*Operation, error) {
	wrapper, err := doJSON[dataResponse[Operation]](c, ctx, http.MethodGet, "/v1/operations/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
