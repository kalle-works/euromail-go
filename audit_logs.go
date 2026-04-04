package euromail

import (
	"context"
	"net/http"
	"net/url"
)

// ListAuditLogs returns a paginated list of audit log entries.
func (c *Client) ListAuditLogs(ctx context.Context, params *ListParams) ([]AuditLog, *Pagination, error) {
	q := url.Values{}
	addListParams(q, params)

	path := "/v1/audit-logs"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[AuditLog]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}
