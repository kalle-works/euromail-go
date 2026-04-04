package euromail

import (
	"context"
	"net/http"
)

// ListPlans returns available billing plans.
func (c *Client) ListPlans(ctx context.Context) ([]BillingPlan, error) {
	type resp struct {
		Data []BillingPlan `json:"data"`
	}
	wrapper, err := doJSON[resp](c, ctx, http.MethodGet, "/v1/billing/plans", nil)
	if err != nil {
		return nil, err
	}
	return wrapper.Data, nil
}

// GetSubscription returns the current billing subscription.
func (c *Client) GetSubscription(ctx context.Context) (*Subscription, error) {
	wrapper, err := doJSON[dataResponse[Subscription]](c, ctx, http.MethodGet, "/v1/billing/subscription", nil)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// CreateCheckout creates a Stripe checkout session for upgrading.
func (c *Client) CreateCheckout(ctx context.Context, params CheckoutParams) (*CheckoutResponse, error) {
	wrapper, err := doJSON[dataResponse[CheckoutResponse]](c, ctx, http.MethodPost, "/v1/billing/checkout", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// CreateBillingPortal creates a Stripe billing portal session.
func (c *Client) CreateBillingPortal(ctx context.Context, params PortalParams) (*PortalResponse, error) {
	wrapper, err := doJSON[dataResponse[PortalResponse]](c, ctx, http.MethodPost, "/v1/billing/portal", params)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
