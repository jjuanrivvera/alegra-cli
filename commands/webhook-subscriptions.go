package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.WebhookSubscription]{
		Use:     "webhook-subscriptions",
		Aliases: []string{"webhooks"},
		Short:   "Manage webhook subscriptions (event notifications)",
		New:     func(c *api.Client) *api.Resource[api.WebhookSubscription] { return c.WebhookSubscriptions() },
		Columns: []string{"id", "event", "url", "status"},
	})
}
