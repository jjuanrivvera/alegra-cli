package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.PriceList]{
		Use:         "price-lists",
		Aliases:     []string{"price-list"},
		Short:       "Manage price lists",
		New:         func(c *api.Client) *api.Resource[api.PriceList] { return c.PriceLists() },
		Columns:     []string{"id", "name", "type", "status"},
		OrderFields: []string{"id", "name"},
		// Free-text search is already provided by the built-in --query flag.
	})
}
