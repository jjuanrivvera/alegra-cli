package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.VariantAttribute]{
		Use:         "variant-attributes",
		Aliases:     []string{"variant-attribute"},
		Short:       "Manage item variant attributes (e.g. Color, Talla)",
		New:         func(c *api.Client) *api.Resource[api.VariantAttribute] { return c.VariantAttributes() },
		Columns:     []string{"id", "name", "status"},
		OrderFields: []string{"name"},
		ListFilters: []listFilter{
			{Flag: "name", Query: "name", Usage: "Filter variant attributes by name"},
			{Flag: "options", Query: "options", Usage: "Filter variant attributes by options"},
		},
	})
}
