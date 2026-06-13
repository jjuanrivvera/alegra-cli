package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.User]{
		Use:         "users",
		Aliases:     []string{"user"},
		Short:       "Manage account users",
		Long:        "Manage Alegra account users, who are the people with access to the company.",
		New:         func(c *api.Client) *api.Resource[api.User] { return c.Users() },
		Columns:     []string{"id", "name", "email", "role", "status"},
		OrderFields: []string{"id", "name", "email"},
		ListFilters: []listFilter{
			{Flag: "status", Query: "status", Usage: "Filter by status: active or inactive"},
			{Flag: "role", Query: "role", Usage: "Filter by role"},
		},
		Extra: func(parent *cobra.Command, sp resourceSpec[api.User]) {
			// GET /users/self — read-only; opt back in to read-only MCP hints so it
			// isn't defaulted to destructive (M1).
			parent.AddCommand(readOnlyHints(&cobra.Command{
				Use:   "self",
				Short: "Show the currently authenticated user",
				Args:  cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					client, err := getAPIClient(cmd)
					if err != nil {
						return err
					}
					var out api.User
					if err := client.GetInto(cmd.Context(), "users/self", nil, &out); err != nil {
						return err
					}
					return render(cmd, out, nil)
				},
			}))
		},
	})
}
