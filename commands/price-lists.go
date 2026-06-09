package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.PriceList]{
		Use:         "price-lists",
		Aliases:     []string{"price-list"},
		Short:       "Manage price lists",
		Long:        "Manage price lists (listas de precios) — named pricing tiers (e.g. wholesale, retail) that items carry and documents apply.",
		New:         func(c *api.Client) *api.Resource[api.PriceList] { return c.PriceLists() },
		Columns:     []string{"id", "name", "type", "status"},
		OrderFields: []string{"id", "name"},
		// Free-text search is already provided by the built-in --query flag.
	})
}
