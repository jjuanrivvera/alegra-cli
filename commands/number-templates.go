package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.NumberTemplate]{
		Use:         "number-templates",
		Aliases:     []string{"numbering"},
		Short:       "Manage document numberings (numeraciones de facturación)",
		New:         func(c *api.Client) *api.Resource[api.NumberTemplate] { return c.NumberTemplates() },
		Columns:     []string{"id", "name", "prefix", "status", "documentType"},
		OrderFields: []string{"id", "name"},
		ListFilters: []listFilter{
			{Flag: "document-type", Query: "documentType", Usage: "Filter by document type: invoice, estimate, creditNote, debitNote, incomeDebitNote, transactionIn, transactionOut (required by the API)"},
		},
	})
}
