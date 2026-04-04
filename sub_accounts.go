package euromail

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// SubAccount represents a sub-account managed by a parent account.
type SubAccount struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Email               string `json:"email"`
	Plan                string `json:"plan"`
	MonthlyQuota        int    `json:"monthly_quota"`
	EmailsSentThisMonth int64  `json:"emails_sent_this_month"`
	ParentAccountID     string `json:"parent_account_id"`
	IsActive            bool   `json:"is_active"`
	CreatedAt           string `json:"created_at"`
}

// CreateSubAccountParams are the parameters for creating a sub-account.
type CreateSubAccountParams struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	MonthlyQuota int    `json:"monthly_quota"`
}

// UpdateSubAccountParams are the parameters for updating a sub-account.
type UpdateSubAccountParams struct {
	Name         *string `json:"name,omitempty"`
	MonthlyQuota *int    `json:"monthly_quota,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

// CreateSubAccount creates a new sub-account under the current parent account.
func (c *Client) CreateSubAccount(ctx context.Context, params CreateSubAccountParams) (*SubAccount, error) {
	wrapper, err := doJSON[dataResponse[SubAccount]](c, ctx, http.MethodPost, "/v1/accounts", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// ListSubAccounts returns all sub-accounts for the current parent account.
func (c *Client) ListSubAccounts(ctx context.Context, params *ListParams) ([]SubAccount, *Pagination, error) {
	q := url.Values{}
	addListParams(q, params)
	path := "/v1/accounts"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	wrapper, err := doJSON[paginatedResponse[SubAccount]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return wrapper.Data, &wrapper.Pagination, nil
}

// GetSubAccount returns a specific sub-account by ID.
func (c *Client) GetSubAccount(ctx context.Context, id string) (*SubAccount, error) {
	wrapper, err := doJSON[dataResponse[SubAccount]](c, ctx, http.MethodGet, fmt.Sprintf("/v1/accounts/%s", id), nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// UpdateSubAccount updates a sub-account's name, quota, or active status.
func (c *Client) UpdateSubAccount(ctx context.Context, id string, params UpdateSubAccountParams) (*SubAccount, error) {
	wrapper, err := doJSON[dataResponse[SubAccount]](c, ctx, http.MethodPatch, fmt.Sprintf("/v1/accounts/%s", id), params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// DeleteSubAccount permanently deletes a sub-account and all its data.
func (c *Client) DeleteSubAccount(ctx context.Context, id string) error {
	resp, err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/v1/accounts/%s", id), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// GetSubAccountAnalytics returns analytics for a specific sub-account.
func (c *Client) GetSubAccountAnalytics(ctx context.Context, id string, aq *AnalyticsQuery) (*AnalyticsSummary, error) {
	q := url.Values{}
	addAnalyticsQuery(q, aq)
	path := fmt.Sprintf("/v1/accounts/%s/analytics", id)
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	wrapper, err := doJSON[AnalyticsSummary](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return &wrapper, nil
}

// GetAggregateAnalytics returns analytics aggregated across the parent account and all sub-accounts.
func (c *Client) GetAggregateAnalytics(ctx context.Context, aq *AnalyticsQuery) (*AnalyticsSummary, error) {
	q := url.Values{}
	addAnalyticsQuery(q, aq)
	path := "/v1/analytics/aggregate"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	wrapper, err := doJSON[AnalyticsSummary](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return &wrapper, nil
}
