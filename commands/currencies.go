package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.Currency]{
		Use:     "currencies",
		Aliases: []string{"currency"},
		Short:   "Manage currencies (monedas)",
		Long:    "Manage currencies. Currencies are identified by their ISO code (e.g. USD), not a numeric id; get/update take the code.",
		New:     func(c *api.Client) *api.Resource[api.Currency] { return c.Currencies() },
		// Currencies are keyed by ISO code, not a numeric id; the struct has no
		// id field, so listing "id" rendered an empty leading column (L6).
		Columns: []string{"code", "name", "symbol", "status"},
		ListFilters: []listFilter{
			{Flag: "fields", Query: "fields", Usage: "Extra fields to include: deletable, canBeInactive, autoUpdate"},
		},
	})
}
