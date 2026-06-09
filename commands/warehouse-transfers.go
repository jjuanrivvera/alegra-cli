package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.WarehouseTransfer]{
		Use:         "warehouse-transfers",
		Aliases:     []string{"warehouse-transfer"},
		Short:       "Manage inventory transfers between warehouses",
		Long:        "Manage warehouse transfers — movements of item stock from one warehouse (bodega) to another.",
		New:         func(c *api.Client) *api.Resource[api.WarehouseTransfer] { return c.WarehouseTransfers() },
		Columns:     []string{"id", "date", "observations"},
		OrderFields: []string{"id", "date"},
		ListFilters: []listFilter{
			{Flag: "date", Query: "date", Usage: "Filter by date (YYYY-MM-DD)"},
			{Flag: "item-id", Query: "item_id", Usage: "Filter transfers involving a specific item ID"},
		},
	})
}
