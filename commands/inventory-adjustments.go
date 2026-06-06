package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.InventoryAdjustment]{
		Use:         "inventory-adjustments",
		Aliases:     []string{"inventory-adjustment"},
		Short:       "Manage inventory adjustments (manual stock corrections)",
		New:         func(c *api.Client) *api.Resource[api.InventoryAdjustment] { return c.InventoryAdjustments() },
		Columns:     []string{"id", "date", "status", "observations"},
		OrderFields: []string{"id", "number", "date", "observations", "warehouse_name"},
		ListFilters: []listFilter{
			{Flag: "number", Query: "number", Usage: "Filter by adjustment number/id"},
			{Flag: "date", Query: "date", Usage: "Filter by date (YYYY-MM-DD)"},
			{Flag: "warehouse-id", Query: "warehouse_id", Usage: "Filter by warehouse ID"},
		},
	})
}
