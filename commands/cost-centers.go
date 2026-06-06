package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.CostCenter]{
		Use:         "cost-centers",
		Aliases:     []string{"cost-center"},
		Short:       "Manage cost centers",
		New:         func(c *api.Client) *api.Resource[api.CostCenter] { return c.CostCenters() },
		Columns:     []string{"id", "name", "code", "status"},
		OrderFields: []string{"id", "name", "code", "status"},
		ListFilters: []listFilter{
			{Flag: "status", Query: "status", Usage: "Filter by status: active or inactive"},
			{Flag: "name", Query: "name", Usage: "Filter by name (partial match)"},
		},
	})
}
