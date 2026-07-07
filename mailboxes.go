package euromail

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// CreateMailbox creates a new agent mailbox. When params.LocalPart is nil the
// server generates a random local-part; when params.DomainID is nil the
// account's default inbound domain is used.
func (c *Client) CreateMailbox(ctx context.Context, params CreateMailboxParams) (*AgentMailbox, error) {
	wrapper, err := doJSON[dataResponse[AgentMailbox]](c, ctx, http.MethodPost, "/v1/agent-mailboxes", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// ListMailboxes returns a paginated list of agent mailboxes for the account.
// ListParams.Page and ListParams.PerPage are mapped to the API's limit/offset
// query parameters (limit = per_page, offset = (page-1)*per_page).
func (c *Client) ListMailboxes(ctx context.Context, params *ListParams) ([]AgentMailbox, *Pagination, error) {
	q := url.Values{}
	if params != nil {
		perPage := 0
		if params.PerPage != nil {
			perPage = *params.PerPage
			q.Set("limit", intToStr(perPage))
		}
		if params.Page != nil && *params.Page > 0 {
			offset := 0
			if perPage > 0 {
				offset = (*params.Page - 1) * perPage
			}
			q.Set("offset", intToStr(offset))
		}
	}

	path := "/v1/agent-mailboxes"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[AgentMailbox]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}

// mailboxListQuery maps ListParams.Page/PerPage to the limit/offset query
// parameters used by the mailbox list endpoints, mirroring ListMailboxes
// (limit = per_page, offset = (page-1)*per_page).
func mailboxListQuery(params *ListParams) url.Values {
	q := url.Values{}
	if params == nil {
		return q
	}
	perPage := 0
	if params.PerPage != nil {
		perPage = *params.PerPage
		q.Set("limit", intToStr(perPage))
	}
	if params.Page != nil && *params.Page > 0 {
		offset := 0
		if perPage > 0 {
			offset = (*params.Page - 1) * perPage
		}
		q.Set("offset", intToStr(offset))
	}
	return q
}

// GetMailbox retrieves a single agent mailbox by ID.
func (c *Client) GetMailbox(ctx context.Context, id string) (*AgentMailbox, error) {
	wrapper, err := doJSON[dataResponse[AgentMailbox]](c, ctx, http.MethodGet, "/v1/agent-mailboxes/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// DeleteMailbox permanently deletes an agent mailbox and all its messages.
func (c *Client) DeleteMailbox(ctx context.Context, id string) error {
	_, err := doJSON[struct{}](c, ctx, http.MethodDelete, "/v1/agent-mailboxes/"+url.PathEscape(id), nil)
	return err
}

// ListMailboxMessages returns messages stored in a mailbox. Use params.Status
// to filter (e.g. "unread", "read"), and Limit/Offset for pagination.
func (c *Client) ListMailboxMessages(ctx context.Context, mailboxID string, params *ListMailboxMessagesParams) ([]MailboxMessage, *Pagination, error) {
	q := url.Values{}
	if params != nil {
		if params.Status != nil {
			q.Set("status", *params.Status)
		}
		if params.Limit != nil {
			q.Set("limit", intToStr(*params.Limit))
		}
		if params.Offset != nil {
			q.Set("offset", intToStr(*params.Offset))
		}
	}

	path := "/v1/agent-mailboxes/" + url.PathEscape(mailboxID) + "/messages"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[MailboxMessage]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}

// WaitForNextMessage long-polls for the next unprocessed message in a mailbox,
// waiting up to timeoutSecs seconds (server default if nil, typically 30).
// On success it returns a LeasedMessage with a lease token that must be passed
// to AckMessage or NackMessage. When the timeout elapses with no message
// available the server responds with 408 and this method returns (nil, nil).
//
// Any other non-2xx status is converted to an error using the shared
// errorFromResponse helper.
func (c *Client) WaitForNextMessage(ctx context.Context, mailboxID string, timeoutSecs *int) (*LeasedMessage, error) {
	q := url.Values{}
	if timeoutSecs != nil {
		q.Set("timeout", strconv.Itoa(*timeoutSecs))
	}
	path := "/v1/agent-mailboxes/" + url.PathEscape(mailboxID) + "/messages/next"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	// We cannot use c.do here because it converts 408 into an error. Build the
	// request inline so we can inspect the status code directly; reuse the
	// same auth + header conventions.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("euromail: failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("euromail: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusRequestTimeout {
		// Drain body so the connection can be reused.
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, nil
	}

	if resp.StatusCode >= 400 {
		return nil, errorFromResponse(resp)
	}

	var leased LeasedMessage
	if err := json.NewDecoder(resp.Body).Decode(&leased); err != nil {
		return nil, fmt.Errorf("euromail: failed to decode response: %w", err)
	}
	return &leased, nil
}

// DeleteMailboxMessage permanently deletes a message from a mailbox.
func (c *Client) DeleteMailboxMessage(ctx context.Context, mailboxID, messageID string) error {
	_, err := doJSON[struct{}](c, ctx,
		http.MethodDelete,
		"/v1/agent-mailboxes/"+url.PathEscape(mailboxID)+"/messages/"+url.PathEscape(messageID),
		nil,
	)
	return err
}

// AckMessage acknowledges successful processing of a leased message. After
// ack the message will not be redelivered.
func (c *Client) AckMessage(ctx context.Context, mailboxID, messageID, leaseToken string) error {
	body := map[string]string{"lease_token": leaseToken}
	_, err := doJSON[struct{}](c, ctx,
		http.MethodPost,
		"/v1/agent-mailboxes/"+url.PathEscape(mailboxID)+"/messages/"+url.PathEscape(messageID)+"/ack",
		body,
	)
	return err
}

// NackMessage releases a leased message back to the queue for redelivery.
// Use this when processing failed and you want another worker (or a retry) to
// pick it up.
func (c *Client) NackMessage(ctx context.Context, mailboxID, messageID, leaseToken string) error {
	body := map[string]string{"lease_token": leaseToken}
	_, err := doJSON[struct{}](c, ctx,
		http.MethodPost,
		"/v1/agent-mailboxes/"+url.PathEscape(mailboxID)+"/messages/"+url.PathEscape(messageID)+"/nack",
		body,
	)
	return err
}

// ReplyToMessage sends a reply to a received message, threaded onto the original
// conversation. At least one of params.TextBody or params.HTMLBody must be set;
// the server returns 400 if both are absent.
func (c *Client) ReplyToMessage(ctx context.Context, mailboxID, messageID string, params ReplyToMessageParams) (*MailboxReplyResult, error) {
	wrapper, err := doJSON[dataResponse[MailboxReplyResult]](c, ctx,
		http.MethodPost,
		"/v1/agent-mailboxes/"+url.PathEscape(mailboxID)+"/messages/"+url.PathEscape(messageID)+"/reply",
		params,
	)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// ListMailboxThreads returns the latest message of each conversation thread in a
// mailbox. ListParams.Page and ListParams.PerPage are mapped to the API's
// limit/offset query parameters.
func (c *Client) ListMailboxThreads(ctx context.Context, mailboxID string, params *ListParams) ([]MailboxMessage, *Pagination, error) {
	q := mailboxListQuery(params)
	path := "/v1/agent-mailboxes/" + url.PathEscape(mailboxID) + "/threads"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[MailboxMessage]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}

// GetMailboxThread returns every message in a single conversation thread, ordered
// chronologically ascending. ListParams.Page and ListParams.PerPage are mapped to
// the API's limit/offset query parameters.
func (c *Client) GetMailboxThread(ctx context.Context, mailboxID, threadID string, params *ListParams) ([]MailboxMessage, *Pagination, error) {
	q := mailboxListQuery(params)
	path := "/v1/agent-mailboxes/" + url.PathEscape(mailboxID) + "/threads/" + url.PathEscape(threadID)
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[MailboxMessage]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}

// SearchMailboxMessages performs a full-text search over the messages in a
// mailbox. query is required (1-500 characters); the server returns 400 for an
// empty or over-long query. ListParams.Page and ListParams.PerPage are mapped to
// the API's limit/offset query parameters.
func (c *Client) SearchMailboxMessages(ctx context.Context, mailboxID, query string, params *ListParams) ([]MailboxMessage, *Pagination, error) {
	q := mailboxListQuery(params)
	q.Set("q", query)
	path := "/v1/agent-mailboxes/" + url.PathEscape(mailboxID) + "/messages/search?" + q.Encode()

	result, err := doJSON[paginatedResponse[MailboxMessage]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}

// UpdateMessageLabels replaces the labels on a message with the provided set
// (this is a full replace, not a merge). Each label must be 1-64 characters of
// alphanumerics, dashes, or underscores, and at most 50 labels are allowed. It
// returns the stored labels.
func (c *Client) UpdateMessageLabels(ctx context.Context, mailboxID, messageID string, labels []string) ([]string, error) {
	body := map[string][]string{"labels": labels}
	wrapper, err := doJSON[dataResponse[struct {
		Labels []string `json:"labels"`
	}]](c, ctx,
		http.MethodPut,
		"/v1/agent-mailboxes/"+url.PathEscape(mailboxID)+"/messages/"+url.PathEscape(messageID)+"/labels",
		body,
	)
	if err != nil {
		return nil, err
	}
	return wrapper.Data.Labels, nil
}

// GetMessageAttachmentURLs returns the attachments of a message. When the
// attachment bytes were persisted, each entry carries a pre-signed download URL
// (see MailboxAttachmentURL for the fallback behavior when they were not).
func (c *Client) GetMessageAttachmentURLs(ctx context.Context, mailboxID, messageID string) ([]MailboxAttachmentURL, error) {
	wrapper, err := doJSON[dataResponse[[]MailboxAttachmentURL]](c, ctx,
		http.MethodGet,
		"/v1/agent-mailboxes/"+url.PathEscape(mailboxID)+"/messages/"+url.PathEscape(messageID)+"/attachments",
		nil,
	)
	if err != nil {
		return nil, err
	}
	return wrapper.Data, nil
}

// ListMailboxContacts returns the distinct correspondents that have sent mail to
// a mailbox. ListParams.Page and ListParams.PerPage are mapped to the API's
// limit/offset query parameters.
func (c *Client) ListMailboxContacts(ctx context.Context, mailboxID string, params *ListParams) ([]MailboxContact, *Pagination, error) {
	q := mailboxListQuery(params)
	path := "/v1/agent-mailboxes/" + url.PathEscape(mailboxID) + "/contacts"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[MailboxContact]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}

// GetMailboxAnalytics returns aggregate message counts for a mailbox.
func (c *Client) GetMailboxAnalytics(ctx context.Context, mailboxID string) (*MailboxAnalytics, error) {
	wrapper, err := doJSON[dataResponse[MailboxAnalytics]](c, ctx,
		http.MethodGet,
		"/v1/agent-mailboxes/"+url.PathEscape(mailboxID)+"/analytics",
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// UpdateAutoResponder configures the mailbox auto-responder. Fields left nil in
// params are unchanged. It returns the resulting auto-responder configuration.
func (c *Client) UpdateAutoResponder(ctx context.Context, mailboxID string, params UpdateAutoResponderParams) (*AutoResponderConfig, error) {
	wrapper, err := doJSON[dataResponse[AutoResponderConfig]](c, ctx,
		http.MethodPatch,
		"/v1/agent-mailboxes/"+url.PathEscape(mailboxID)+"/auto-responder",
		params,
	)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
