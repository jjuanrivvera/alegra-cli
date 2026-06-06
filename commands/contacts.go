package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.Contact]{
		Use:     "contacts",
		Aliases: []string{"contact"},
		Short:   "Manage contacts (clients and providers)",
		Long:    "Create, list, update, and delete Alegra contacts — your clients and providers.",
		New:     func(c *api.Client) *api.Resource[api.Contact] { return c.Contacts() },
		Columns: []string{"id", "name", "identification", "email", "status"},
		OrderFields: []string{
			"id", "name", "email",
		},
		ListFilters: []listFilter{
			{Flag: "identification", Query: "identification", Usage: "Filter by identification number"},
			{Flag: "name", Query: "name", Usage: "Filter by name"},
			{Flag: "email", Query: "email", Usage: "Filter by email"},
			{Flag: "type", Query: "type", Usage: "Filter by type: client or provider"},
		},
	})
}
