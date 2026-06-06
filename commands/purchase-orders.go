package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.PurchaseOrder]{
		Use:         "purchase-orders",
		Aliases:     []string{"purchase-order"},
		Short:       "Manage purchase orders (órdenes de compra)",
		New:         func(c *api.Client) *api.Resource[api.PurchaseOrder] { return c.PurchaseOrders() },
		Columns:     []string{"id", "date", "status", "total"},
		OrderFields: []string{"id", "name", "date", "deliveryDate", "status"},
		ListFilters: []listFilter{
			{Flag: "status", Query: "status", Usage: "Filter by status"},
			{Flag: "provider-name", Query: "provider_name", Usage: "Filter by provider name"},
			{Flag: "client-id", Query: "client_id", Usage: "Filter by client ID"},
			{Flag: "warehouse-id", Query: "warehouse_id", Usage: "Filter by warehouse ID"},
			{Flag: "cost-center-id", Query: "costCenter_id", Usage: "Filter by cost center ID (null/none/0 for none)"},
			{Flag: "number", Query: "number", Usage: "Filter by number"},
			{Flag: "date", Query: "date", Usage: "Filter by date (YYYY-MM-DD)"},
			{Flag: "delivery-date", Query: "deliveryDate", Usage: "Filter by delivery date (YYYY-MM-DD)"},
			{Flag: "currency", Query: "currency", Usage: "Filter by currency code"},
		},
		Extra: func(parent *cobra.Command, sp resourceSpec[api.PurchaseOrder]) {
			parent.AddCommand(NewActionCmd(sp, "void", "void", "Void a purchase order"))
			parent.AddCommand(NewActionCmd(sp, "email", "email", "Email a purchase order", true))
			parent.AddCommand(NewActionCmd(sp, "comments", "comments", "Add a comment", true))
		},
	})
}
