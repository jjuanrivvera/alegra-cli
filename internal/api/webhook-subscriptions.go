package api

// WebhookSubscription is a company subscription to an Alegra webhook event:
// Alegra POSTs to URL whenever Event fires (e.g. new-invoice, edit-client).
type WebhookSubscription struct {
	ID     ID     `json:"id,omitempty"`
	Event  string `json:"event,omitempty"`
	URL    string `json:"url,omitempty"`
	Status string `json:"status,omitempty"`
}

// WebhookSubscriptions returns a typed handle to the /webhooks/subscriptions
// resource (list/get/create/update/delete).
func (c *Client) WebhookSubscriptions() *Resource[WebhookSubscription] {
	return NewResource[WebhookSubscription](c, "webhooks/subscriptions")
}
