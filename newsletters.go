package euromail

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// CreateNewsletter creates a new newsletter draft.
func (c *Client) CreateNewsletter(ctx context.Context, params CreateNewsletterParams) (*Newsletter, error) {
	wrapper, err := doJSON[dataResponse[Newsletter]](c, ctx, http.MethodPost, "/v1/newsletters", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// GetNewsletter returns a newsletter by ID.
func (c *Client) GetNewsletter(ctx context.Context, id string) (*Newsletter, error) {
	wrapper, err := doJSON[dataResponse[Newsletter]](c, ctx, http.MethodGet, "/v1/newsletters/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// UpdateNewsletter updates a newsletter draft.
func (c *Client) UpdateNewsletter(ctx context.Context, id string, params UpdateNewsletterParams) (*Newsletter, error) {
	wrapper, err := doJSON[dataResponse[Newsletter]](c, ctx, http.MethodPut, "/v1/newsletters/"+url.PathEscape(id), params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// DeleteNewsletter deletes a newsletter.
func (c *Client) DeleteNewsletter(ctx context.Context, id string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v1/newsletters/"+url.PathEscape(id), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListNewsletters returns newsletters with optional limit/offset.
func (c *Client) ListNewsletters(ctx context.Context, limit, offset *int) ([]Newsletter, error) {
	q := url.Values{}
	if limit != nil {
		q.Set("limit", intToStr(*limit))
	}
	if offset != nil {
		q.Set("offset", intToStr(*offset))
	}
	path := "/v1/newsletters"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	type resp struct {
		Data []Newsletter `json:"data"`
	}
	wrapper, err := doJSON[resp](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return wrapper.Data, nil
}

// SendNewsletter sends a newsletter to its contact list.
func (c *Client) SendNewsletter(ctx context.Context, id string) (*NewsletterSendResponse, error) {
	wrapper, err := doJSON[dataResponse[NewsletterSendResponse]](c, ctx, http.MethodPost, fmt.Sprintf("/v1/newsletters/%s/send", url.PathEscape(id)), struct{}{})
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
