package euromail

import (
	"context"
	"net/http"
)

// ValidateEmail validates an email address.
func (c *Client) ValidateEmail(ctx context.Context, email string) (*EmailValidation, error) {
	result, err := doJSON[EmailValidation](c, ctx, http.MethodPost, "/v1/validate", map[string]string{"email": email})
	if err != nil {
		return nil, err
	}
	return &result, nil
}
