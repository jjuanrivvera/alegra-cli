package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.Tax]{
		Use:         "taxes",
		Aliases:     []string{"tax"},
		Short:       "Manage taxes (e.g. IVA)",
		New:         func(c *api.Client) *api.Resource[api.Tax] { return c.Taxes() },
		Columns:     []string{"id", "name", "percentage", "status", "type"},
		OrderFields: []string{"id", "name", "percentage"},
		ListFilters: []listFilter{
			{Flag: "type", Query: "type", Usage: "Filter by tax type (e.g. IVA)"},
			{Flag: "status", Query: "status", Usage: "Filter by status: active or inactive"},
		},
	})
}
