package euromail

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// CreateSignupForm creates a new signup form.
func (c *Client) CreateSignupForm(ctx context.Context, params CreateSignupFormParams) (*SignupForm, error) {
	wrapper, err := doJSON[dataResponse[SignupForm]](c, ctx, http.MethodPost, "/v1/signup-forms", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// ListSignupForms returns all signup forms for the account.
func (c *Client) ListSignupForms(ctx context.Context) ([]SignupForm, error) {
	wrapper, err := doJSON[dataResponse[[]SignupForm]](c, ctx, http.MethodGet, "/v1/signup-forms", nil)
	if err != nil {
		return nil, err
	}
	return wrapper.Data, nil
}

// GetSignupForm retrieves a signup form by ID.
func (c *Client) GetSignupForm(ctx context.Context, id string) (*SignupForm, error) {
	wrapper, err := doJSON[dataResponse[SignupForm]](c, ctx, http.MethodGet, "/v1/signup-forms/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// UpdateSignupForm updates an existing signup form.
func (c *Client) UpdateSignupForm(ctx context.Context, id string, params UpdateSignupFormParams) (*SignupForm, error) {
	wrapper, err := doJSON[dataResponse[SignupForm]](c, ctx, http.MethodPut, "/v1/signup-forms/"+url.PathEscape(id), params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// DeleteSignupForm deletes a signup form.
func (c *Client) DeleteSignupForm(ctx context.Context, id string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v1/signup-forms/"+url.PathEscape(id), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ToggleSignupForm toggles a signup form's active status.
func (c *Client) ToggleSignupForm(ctx context.Context, id string) (*SignupForm, error) {
	wrapper, err := doJSON[dataResponse[SignupForm]](c, ctx, http.MethodPost, fmt.Sprintf("/v1/signup-forms/%s/toggle", url.PathEscape(id)), struct{}{})
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
