package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.InventoryAdjustmentNumeration]{
		Use:     "inventory-adjustment-numerations",
		Aliases: []string{"inventory-adjustment-numeration"},
		Short:   "Manage inventory adjustment numerations",
		New: func(c *api.Client) *api.Resource[api.InventoryAdjustmentNumeration] {
			return c.InventoryAdjustmentNumerations()
		},
		Columns:     []string{"id", "name", "prefix", "isDefault", "status"},
		OrderFields: []string{"id", "name"},
		ListFilters: []listFilter{
			{Flag: "name", Query: "name", Usage: "Filter by numerations whose name contains the text"},
			{Flag: "status", Query: "status", Usage: "Filter by status: active or inactive"},
			{Flag: "auto-increment", Query: "autoIncrement", Usage: "Filter by autoincremental (true) or non-autoincremental (false) numerations"},
		},
	})
}
