package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.Item]{
		Use:         "items",
		Aliases:     []string{"item"},
		Short:       "Manage items (products and services)",
		Long:        "Manage items (products and services) — your catalog. Items carry price, taxes, reference codes, and inventory settings, and become the line items on invoices and bills.",
		New:         func(c *api.Client) *api.Resource[api.Item] { return c.Items() },
		Columns:     []string{"id", "name", "reference", "price", "status"},
		OrderFields: []string{"id", "name", "reference", "price"},
		ListFilters: []listFilter{
			{Flag: "name", Query: "name", Usage: "Filter by name"},
		},
	})
}
