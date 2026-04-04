package euromail

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetAccount returns the authenticated account details.
func (c *Client) GetAccount(ctx context.Context) (*Account, error) {
	wrapper, err := doJSON[dataResponse[Account]](c, ctx, http.MethodGet, "/v1/account", nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// ExportAccount returns a GDPR data export of the account as a JSON string.
// This endpoint is rate limited to 1 request per hour.
func (c *Client) ExportAccount(ctx context.Context) (string, error) {
	return doRawText(c, ctx, http.MethodGet, "/v1/account/export")
}

// DeleteAccount requests permanent deletion of the account (GDPR).
// This is irreversible and requires the X-Confirm-Delete header.
func (c *Client) DeleteAccount(ctx context.Context) (*AccountDeleteResponse, error) {
	u := c.baseURL + "/v1/account"

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, fmt.Errorf("euromail: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Confirm-Delete", "DELETE")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("euromail: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errorFromResponse(resp)
	}

	var wrapper dataResponse[AccountDeleteResponse]
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("euromail: failed to decode response: %w", err)
	}
	return &wrapper.Data, nil
}
