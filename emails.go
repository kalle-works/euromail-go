package euromail

import (
	"context"
	"net/http"
	"net/url"
)

// SendEmail sends a single email.
func (c *Client) SendEmail(ctx context.Context, params SendEmailParams) (*SendEmailResponse, error) {
	wrapper, err := doJSON[dataResponse[SendEmailResponse]](c, ctx, http.MethodPost, "/v1/emails", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// SendBatch sends a batch of emails.
func (c *Client) SendBatch(ctx context.Context, params SendBatchParams) (*SendBatchResponse, error) {
	result, err := doJSON[SendBatchResponse](c, ctx, http.MethodPost, "/v1/emails/batch", params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetEmail retrieves a single email by ID, including its delivery events.
func (c *Client) GetEmail(ctx context.Context, emailID string) (*EmailDetail, error) {
	wrapper, err := doJSON[dataResponse[EmailDetail]](c, ctx, http.MethodGet, "/v1/emails/"+url.PathEscape(emailID), nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// CancelScheduledEmail cancels a scheduled or queued email.
func (c *Client) CancelScheduledEmail(ctx context.Context, emailID string) (*SendEmailResponse, error) {
	wrapper, err := doJSON[dataResponse[SendEmailResponse]](c, ctx, http.MethodPost, "/v1/emails/"+url.PathEscape(emailID)+"/cancel", struct{}{})
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// SendBroadcast sends an email to all active subscribers in a contact list.
func (c *Client) SendBroadcast(ctx context.Context, params BroadcastParams) (*BroadcastResponse, error) {
	wrapper, err := doJSON[dataResponse[BroadcastResponse]](c, ctx, http.MethodPost, "/v1/emails/broadcast", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// GetEmailLinks returns per-link click statistics for a sent email.
func (c *Client) GetEmailLinks(ctx context.Context, emailID string) ([]LinkClickStat, error) {
	wrapper, err := doJSON[dataResponse[[]LinkClickStat]](c, ctx, http.MethodGet, "/v1/emails/"+url.PathEscape(emailID)+"/links", nil)
	if err != nil {
		return nil, err
	}
	return wrapper.Data, nil
}

// ListEmails returns a paginated list of emails.
func (c *Client) ListEmails(ctx context.Context, params *ListEmailsParams) ([]Email, *Pagination, error) {
	q := url.Values{}
	if params != nil {
		if params.Page != nil {
			q.Set("page", intToStr(*params.Page))
		}
		if params.PerPage != nil {
			q.Set("per_page", intToStr(*params.PerPage))
		}
		if params.Status != nil {
			q.Set("status", *params.Status)
		}
	}

	path := "/v1/emails"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[Email]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}
