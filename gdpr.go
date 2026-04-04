package euromail

import (
	"context"
	"net/http"
	"net/url"
)

// GdprExport exports all data for an email address (GDPR).
func (c *Client) GdprExport(ctx context.Context, email string) (*GdprExportResponse, error) {
	path := "/v1/gdpr/export?" + url.Values{"email": {email}}.Encode()
	result, err := doJSON[GdprExportResponse](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GdprErase erases all data for an email address (GDPR).
func (c *Client) GdprErase(ctx context.Context, email string) (*GdprEraseResponse, error) {
	path := "/v1/gdpr/erase?" + url.Values{"email": {email}}.Encode()
	result, err := doJSON[GdprEraseResponse](c, ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
