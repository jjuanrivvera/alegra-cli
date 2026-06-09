package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.Retention]{
		Use:         "retentions",
		Aliases:     []string{"retention"},
		Short:       "Manage retentions (withholdings)",
		Long:        "Manage retentions (retenciones / withholdings) — tax withholdings such as retefuente or reteICA applied to documents; definitions are country-specific.",
		New:         func(c *api.Client) *api.Resource[api.Retention] { return c.Retentions() },
		Columns:     []string{"id", "name", "percentage", "status"},
		OrderFields: []string{"id", "name"},
		ListFilters: []listFilter{
			{Flag: "type", Query: "type", Usage: "Filter by retention type (e.g. FUENTE)"},
			{Flag: "status", Query: "status", Usage: "Filter by status: active or inactive"},
		},
	})
}
