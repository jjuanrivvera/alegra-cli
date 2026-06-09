package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.Seller]{
		Use:     "sellers",
		Aliases: []string{"seller"},
		Short:   "Manage sellers (vendedores)",
		Long: "Manage sellers (vendedores) — salespeople you can assign to sales documents and " +
			"roll up in `alegra reports sales-by-seller`.",
		New:     func(c *api.Client) *api.Resource[api.Seller] { return c.Sellers() },
		Columns: []string{"id", "name", "identification", "status"},
		ListFilters: []listFilter{
			{Flag: "status", Query: "status", Usage: "Filter by status: active or inactive"},
		},
	})
}
