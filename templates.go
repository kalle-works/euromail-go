package euromail

import (
	"context"
	"net/http"
	"net/url"
)

// CreateTemplate creates a new email template.
func (c *Client) CreateTemplate(ctx context.Context, params CreateTemplateParams) (*Template, error) {
	wrapper, err := doJSON[dataResponse[Template]](c, ctx, http.MethodPost, "/v1/templates", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// GetTemplate retrieves a template by ID.
func (c *Client) GetTemplate(ctx context.Context, templateID string) (*Template, error) {
	wrapper, err := doJSON[dataResponse[Template]](c, ctx, http.MethodGet, "/v1/templates/"+url.PathEscape(templateID), nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// UpdateTemplate updates an existing template.
func (c *Client) UpdateTemplate(ctx context.Context, templateID string, params UpdateTemplateParams) (*Template, error) {
	wrapper, err := doJSON[dataResponse[Template]](c, ctx, http.MethodPut, "/v1/templates/"+url.PathEscape(templateID), params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// DeleteTemplate deletes a template.
func (c *Client) DeleteTemplate(ctx context.Context, templateID string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v1/templates/"+url.PathEscape(templateID), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListTemplates returns a paginated list of templates.
func (c *Client) ListTemplates(ctx context.Context, params *ListParams) ([]Template, *Pagination, error) {
	q := url.Values{}
	addListParams(q, params)

	path := "/v1/templates"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[Template]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}
