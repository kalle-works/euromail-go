package euromail

import (
	"context"
	"net/http"
	"net/url"
)

// CreateWebhook creates a new webhook endpoint.
func (c *Client) CreateWebhook(ctx context.Context, params CreateWebhookParams) (*Webhook, error) {
	wrapper, err := doJSON[dataResponse[Webhook]](c, ctx, http.MethodPost, "/v1/webhooks", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// GetWebhook retrieves a webhook by ID.
func (c *Client) GetWebhook(ctx context.Context, webhookID string) (*Webhook, error) {
	wrapper, err := doJSON[dataResponse[Webhook]](c, ctx, http.MethodGet, "/v1/webhooks/"+url.PathEscape(webhookID), nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// UpdateWebhook updates an existing webhook.
func (c *Client) UpdateWebhook(ctx context.Context, webhookID string, params UpdateWebhookParams) (*Webhook, error) {
	wrapper, err := doJSON[dataResponse[Webhook]](c, ctx, http.MethodPut, "/v1/webhooks/"+url.PathEscape(webhookID), params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// TestWebhook sends a test event to a webhook.
func (c *Client) TestWebhook(ctx context.Context, webhookID string) (*WebhookTestResponse, error) {
	wrapper, err := doJSON[dataResponse[WebhookTestResponse]](c, ctx, http.MethodPost, "/v1/webhooks/"+url.PathEscape(webhookID)+"/test", struct{}{})
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// DeleteWebhook deletes a webhook.
func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v1/webhooks/"+url.PathEscape(webhookID), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListWebhooks returns a paginated list of webhooks.
func (c *Client) ListWebhooks(ctx context.Context, params *ListParams) ([]Webhook, *Pagination, error) {
	q := url.Values{}
	addListParams(q, params)

	path := "/v1/webhooks"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[Webhook]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}
