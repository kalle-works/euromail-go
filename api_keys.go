package euromail

import (
	"context"
	"fmt"
	"net/http"
)

// CreateApiKey creates a new API key with optional scopes.
func (c *Client) CreateApiKey(ctx context.Context, params CreateApiKeyParams) (*ApiKeyCreated, error) {
	wrapper, err := doJSON[dataResponse[ApiKeyCreated]](c, ctx, http.MethodPost, "/v1/api-keys", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// ListApiKeys returns all API keys for the account.
func (c *Client) ListApiKeys(ctx context.Context) ([]ApiKey, error) {
	wrapper, err := doJSON[dataResponse[[]ApiKey]](c, ctx, http.MethodGet, "/v1/api-keys", nil)
	if err != nil {
		return nil, err
	}
	return wrapper.Data, nil
}

// DeleteApiKey revokes an API key.
func (c *Client) DeleteApiKey(ctx context.Context, id string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v1/api-keys/"+id, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// CreateSubAccountApiKey creates an API key for a sub-account.
func (c *Client) CreateSubAccountApiKey(ctx context.Context, subAccountID string, params CreateApiKeyParams) (*ApiKeyCreated, error) {
	wrapper, err := doJSON[dataResponse[ApiKeyCreated]](c, ctx, http.MethodPost, fmt.Sprintf("/v1/accounts/%s/api-keys", subAccountID), params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
