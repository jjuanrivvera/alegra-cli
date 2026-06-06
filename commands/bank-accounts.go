package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.BankAccount]{
		Use:         "bank-accounts",
		Aliases:     []string{"bank-account"},
		Short:       "Manage bank accounts (bank, credit card, and cash accounts)",
		New:         func(c *api.Client) *api.Resource[api.BankAccount] { return c.BankAccounts() },
		Columns:     []string{"id", "name", "type", "status"},
		OrderFields: []string{"id", "date"},
		ListFilters: []listFilter{
			{Flag: "include-inactive", Query: "includeInactive", Usage: "Include inactive bank accounts"},
			{Flag: "include-balance", Query: "includeBalance", Usage: "Include account balances in the response"},
		},
		Extra: func(parent *cobra.Command, sp resourceSpec[api.BankAccount]) {
			parent.AddCommand(NewActionCmd(sp, "transfer", "transfer", "Transfer funds to another bank account"))
		},
	})
}
