package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

// companyColumns are the default columns rendered for the company singleton.
var companyColumns = []string{"name", "identification", "email", "regime", "applicationVersion"}

func init() {
	companyCmd := &cobra.Command{
		Use:   "company",
		Short: "View and update the account's company (empresa)",
		Long: `Manage the singleton company registered on your Alegra account.

The Alegra API exposes the company only at GET /company and PUT /company; there
is no list or per-id access. Use "company get" to view it and "company update"
to edit it with a JSON body.`,
	}

	companyCmd.AddCommand(newCompanyGetCmd())
	companyCmd.AddCommand(newCompanyUpdateCmd())

	rootCmd.AddCommand(companyCmd)
}

func newCompanyGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Show the company information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := getAPIClient(cmd)
			if err != nil {
				return err
			}
			var company api.Company
			if err := client.GetInto(cmd.Context(), "company", nil, &company); err != nil {
				return err
			}
			return render(cmd, company, companyColumns)
		},
	}
}

func newCompanyUpdateCmd() *cobra.Command {
	var bf bodyFlags
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the company information",
		Long: `Update the company with a JSON body.

Provide the body with --data '<json>', --file <path> (use - for stdin), or one
or more --set key=value pairs.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := bf.build()
			if err != nil {
				return err
			}
			client, err := getAPIClient(cmd)
			if err != nil {
				return err
			}
			var company api.Company
			if err := client.PutInto(cmd.Context(), "company", body, &company); err != nil {
				return err
			}
			if flagDryRun {
				return nil
			}
			return render(cmd, company, companyColumns)
		},
	}
	addBodyFlags(cmd, &bf)
	return cmd
}
