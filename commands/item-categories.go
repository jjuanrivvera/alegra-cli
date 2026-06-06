package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.ItemCategory]{
		Use:     "item-categories",
		Aliases: []string{"item-category"},
		Short:   "Manage item categories",
		Long:    "Create, list, update, and delete Alegra item categories used to group products and services.",
		New:     func(c *api.Client) *api.Resource[api.ItemCategory] { return c.ItemCategories() },
		Columns: []string{"id", "name", "description", "status"},
		OrderFields: []string{
			"name",
		},
		ListFilters: []listFilter{
			{Flag: "name", Query: "name", Usage: "Filter by name"},
			{Flag: "description", Query: "description", Usage: "Filter by description"},
			{Flag: "status", Query: "status", Usage: "Filter by status: active or inactive"},
		},
	})
}
