package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.AdditionalCharge]{
		Use:         "additional-charges",
		Aliases:     []string{"additional-charge"},
		Short:       "Manage additional charges (tips and parafiscal contributions)",
		New:         func(c *api.Client) *api.Resource[api.AdditionalCharge] { return c.AdditionalCharges() },
		Columns:     []string{"id", "name", "percentage", "status"},
		OrderFields: []string{"id", "name", "status"},
		ListFilters: []listFilter{
			{Flag: "status", Query: "status", Usage: "Filter by status: active or inactive"},
			{Flag: "name", Query: "name", Usage: "Filter by name (partial match)"},
		},
	})
}
