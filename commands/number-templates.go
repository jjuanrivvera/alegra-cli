package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.NumberTemplate]{
		Use:         "number-templates",
		Aliases:     []string{"numbering"},
		Short:       "Manage document numberings (numeraciones de facturación)",
		Long:        "Manage document numberings (numeraciones de facturación) — the prefixes, ranges, and DIAN/SAT resolutions that number invoices and other documents.",
		New:         func(c *api.Client) *api.Resource[api.NumberTemplate] { return c.NumberTemplates() },
		Columns:     []string{"id", "name", "prefix", "status", "documentType"},
		OrderFields: []string{"id", "name"},
		ListFilters: []listFilter{
			{Flag: "document-type", Query: "documentType", Usage: "Filter by document type: invoice, estimate, creditNote, debitNote, incomeDebitNote, transactionIn, transactionOut (required by the API)"},
		},
	})
}
