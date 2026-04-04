package euromail

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// ListDeadLetters returns dead letter messages.
// The count parameter controls how many items to return (default: 50, max: 200).
func (c *Client) ListDeadLetters(ctx context.Context, count *int) (*DeadLetterListResponse, error) {
	q := url.Values{}
	if count != nil {
		q.Set("count", strconv.Itoa(*count))
	}

	path := "/v1/dead-letters"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[DeadLetterListResponse](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RetryDeadLetter re-enqueues a dead letter message for retry.
func (c *Client) RetryDeadLetter(ctx context.Context, streamID string) (*MessageResponse, error) {
	result, err := doJSON[MessageResponse](c, ctx, http.MethodPost, "/v1/dead-letters/"+url.PathEscape(streamID)+"/retry", struct{}{})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteDeadLetter permanently removes a dead letter message.
func (c *Client) DeleteDeadLetter(ctx context.Context, streamID string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v1/dead-letters/"+url.PathEscape(streamID), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
