package euromail

import (
	"context"
	"net/http"
	"net/url"
)

// AddDomain registers a new sending domain.
func (c *Client) AddDomain(ctx context.Context, domain string) (*Domain, error) {
	wrapper, err := doJSON[dataResponse[Domain]](c, ctx, http.MethodPost, "/v1/domains", map[string]string{"domain": domain})
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// GetDomain retrieves a domain by ID.
func (c *Client) GetDomain(ctx context.Context, domainID string) (*Domain, error) {
	wrapper, err := doJSON[dataResponse[Domain]](c, ctx, http.MethodGet, "/v1/domains/"+url.PathEscape(domainID), nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// VerifyDomain triggers DNS verification for a domain.
func (c *Client) VerifyDomain(ctx context.Context, domainID string) (*DomainVerificationResult, error) {
	wrapper, err := doJSON[dataResponse[DomainVerificationResult]](c, ctx, http.MethodPost, "/v1/domains/"+url.PathEscape(domainID)+"/verify", struct{}{})
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// DeleteDomain removes a domain.
func (c *Client) DeleteDomain(ctx context.Context, domainID string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v1/domains/"+url.PathEscape(domainID), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// SetTrackingDomain sets a vanity tracking domain.
func (c *Client) SetTrackingDomain(ctx context.Context, domainID, trackingDomain string) (*TrackingDomainResponse, error) {
	result, err := doJSON[TrackingDomainResponse](c, ctx, http.MethodPut, "/v1/domains/"+url.PathEscape(domainID)+"/tracking-domain", map[string]string{"tracking_domain": trackingDomain})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// VerifyTrackingDomain verifies the CNAME record for a tracking domain.
func (c *Client) VerifyTrackingDomain(ctx context.Context, domainID string) (*TrackingDomainVerification, error) {
	result, err := doJSON[TrackingDomainVerification](c, ctx, http.MethodPost, "/v1/domains/"+url.PathEscape(domainID)+"/verify-tracking", struct{}{})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveTrackingDomain removes the vanity tracking domain.
func (c *Client) RemoveTrackingDomain(ctx context.Context, domainID string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v1/domains/"+url.PathEscape(domainID)+"/tracking-domain", nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListDomains returns a paginated list of domains.
func (c *Client) ListDomains(ctx context.Context, params *ListParams) ([]Domain, *Pagination, error) {
	q := url.Values{}
	addListParams(q, params)

	path := "/v1/domains"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[Domain]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}
