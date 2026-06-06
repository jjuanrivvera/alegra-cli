package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.Currency]{
		Use:     "currencies",
		Aliases: []string{"currency"},
		Short:   "Manage currencies (monedas)",
		Long:    "Manage currencies. Currencies are identified by their ISO code (e.g. USD), not a numeric id; get/update take the code.",
		New:     func(c *api.Client) *api.Resource[api.Currency] { return c.Currencies() },
		Columns: []string{"id", "code", "name", "symbol"},
		ListFilters: []listFilter{
			{Flag: "fields", Query: "fields", Usage: "Extra fields to include: deletable, canBeInactive, autoUpdate"},
		},
	})
}
