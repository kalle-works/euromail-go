package euromail

import (
	"context"
	"net/http"
	"net/url"
)

// CreateContactList creates a new contact list.
func (c *Client) CreateContactList(ctx context.Context, params CreateContactListParams) (*ContactList, error) {
	wrapper, err := doJSON[dataResponse[ContactList]](c, ctx, http.MethodPost, "/v1/contact-lists", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// GetContactList retrieves a contact list by ID.
func (c *Client) GetContactList(ctx context.Context, listID string) (*ContactList, error) {
	wrapper, err := doJSON[dataResponse[ContactList]](c, ctx, http.MethodGet, "/v1/contact-lists/"+url.PathEscape(listID), nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// UpdateContactList updates an existing contact list.
func (c *Client) UpdateContactList(ctx context.Context, listID string, params UpdateContactListParams) (*ContactList, error) {
	wrapper, err := doJSON[dataResponse[ContactList]](c, ctx, http.MethodPut, "/v1/contact-lists/"+url.PathEscape(listID), params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// DeleteContactList deletes a contact list.
func (c *Client) DeleteContactList(ctx context.Context, listID string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v1/contact-lists/"+url.PathEscape(listID), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListContactLists returns all contact lists for the account.
func (c *Client) ListContactLists(ctx context.Context) ([]ContactList, error) {
	wrapper, err := doJSON[dataResponse[[]ContactList]](c, ctx, http.MethodGet, "/v1/contact-lists", nil)
	if err != nil {
		return nil, err
	}
	return wrapper.Data, nil
}

// AddContact adds a contact to a contact list.
func (c *Client) AddContact(ctx context.Context, listID string, params AddContactParams) (*Contact, error) {
	wrapper, err := doJSON[dataResponse[Contact]](c, ctx, http.MethodPost, "/v1/contact-lists/"+url.PathEscape(listID)+"/contacts", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// BulkAddContacts adds multiple contacts to a contact list.
func (c *Client) BulkAddContacts(ctx context.Context, listID string, params BulkAddContactsParams) (*BulkAddContactsResponse, error) {
	wrapper, err := doJSON[dataResponse[BulkAddContactsResponse]](c, ctx, http.MethodPost, "/v1/contact-lists/"+url.PathEscape(listID)+"/contacts", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// ListContacts returns a paginated list of contacts in a contact list.
func (c *Client) ListContacts(ctx context.Context, listID string, params *ListContactsParams) ([]Contact, *Pagination, error) {
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

	path := "/v1/contact-lists/" + url.PathEscape(listID) + "/contacts"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	result, err := doJSON[paginatedResponse[Contact]](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, &result.Pagination, nil
}

// RemoveContact removes a contact from a contact list by email address.
func (c *Client) RemoveContact(ctx context.Context, listID string, email string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v1/contact-lists/"+url.PathEscape(listID)+"/contacts/"+url.PathEscape(email), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
