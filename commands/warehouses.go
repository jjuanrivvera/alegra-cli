package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.Warehouse]{
		Use:         "warehouses",
		Aliases:     []string{"warehouse"},
		Short:       "Manage inventory warehouses (bodegas)",
		New:         func(c *api.Client) *api.Resource[api.Warehouse] { return c.Warehouses() },
		Columns:     []string{"id", "name", "address", "isDefault", "status"},
		OrderFields: []string{"id", "name"},
		ListFilters: []listFilter{
			{Flag: "name", Query: "name", Usage: "Filter by warehouse name"},
			{Flag: "status", Query: "status", Usage: "Filter by status: active or inactive"},
		},
	})
}
