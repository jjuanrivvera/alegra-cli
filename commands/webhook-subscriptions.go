package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.WebhookSubscription]{
		Use:     "webhook-subscriptions",
		Aliases: []string{"webhooks"},
		Short:   "Manage webhook subscriptions (event notifications)",
		Long: "Manage webhook subscriptions: Alegra POSTs to your URL whenever an event fires " +
			"(e.g. new-invoice, edit-client, delete-item). Create one per event+URL with " +
			`--set event=new-invoice --set url=https://your.app/hook.`,
		New:     func(c *api.Client) *api.Resource[api.WebhookSubscription] { return c.WebhookSubscriptions() },
		Columns: []string{"id", "event", "url", "status"},
	})
}
