package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.CustomField]{
		Use:         "custom-fields",
		Aliases:     []string{"custom-field"},
		Short:       "Manage custom fields (campos adicionales)",
		New:         func(c *api.Client) *api.Resource[api.CustomField] { return c.CustomFields() },
		Columns:     []string{"id", "name", "type", "status"},
		OrderFields: []string{"id", "name"},
		ListFilters: []listFilter{
			{Flag: "name", Query: "name_query", Usage: "Filter by text contained in the field name"},
			{Flag: "description", Query: "description_query", Usage: "Filter by text contained in the field description"},
			{Flag: "resource-type", Query: "resourceType", Usage: "Filter by associated resource type (e.g. item)"},
			{Flag: "type", Query: "type", Usage: "Filter by field type: text, int, date, boolean, optionsList"},
			{Flag: "status", Query: "status", Usage: "Filter by status: active or inactive"},
		},
	})
}
