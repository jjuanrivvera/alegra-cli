package commands

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.Bill]{
		Use:         "bills",
		Aliases:     []string{"bill"},
		Short:       "Manage provider bills (facturas de proveedor)",
		Long:        "Manage provider bills (facturas de proveedor) — purchases you owe. Supports attachments, comments, advance application, perceptions/retentions, and importing received Colombian e-invoices by CUFE.",
		New:         func(c *api.Client) *api.Resource[api.Bill] { return c.Bills() },
		Columns:     []string{"id", "date", "dueDate", "status", "total", "balance"},
		OrderFields: []string{"id", "date", "dueDate", "status"},
		ListFilters: append([]listFilter{
			{Flag: "status", Query: "status", Usage: "Filter by status"},
			{Flag: "provider-name", Query: "provider_name", Usage: "Filter by provider name"},
			{Flag: "client-id", Query: "client_id", Usage: "Filter by provider ID"},
		}, dateRangeFilters()...),
		Extra: func(parent *cobra.Command, sp resourceSpec[api.Bill]) {
			parent.AddCommand(NewActionCmd(sp, "close", "close", "Close a bill with pending balance", true))
			parent.AddCommand(NewActionCmd(sp, "comments", "comments", "Add a comment to a bill", true))
			parent.AddCommand(NewActionCmd(sp, "advances", "advances-applied", "Apply provider advances to a bill", true))
			parent.AddCommand(NewActionCmd(sp, "attach", "attachment", "Attach a file to a bill (JSON body with a base64 'file' field)", true))
			parent.AddCommand(NewPutActionCmd(sp, "perceptions", "perceptions", "Replace a bill's perceptions"))
			parent.AddCommand(NewPutActionCmd(sp, "retentions", "retentions", "Replace a bill's retentions"))
			parent.AddCommand(billCommentUpdateCmd(sp))
			parent.AddCommand(billCommentDeleteCmd(sp))
			parent.AddCommand(billAttachmentDeleteCmd(sp))
			parent.AddCommand(NewCollectionActionCmd(sp, "import-by-cufe", "import-by-cufe", "Import a bill by CUFE (Colombia)"))
		},
	})
}

// billCommentUpdateCmd: PUT /bills/{id}/comments/{commentId}.
func billCommentUpdateCmd(sp resourceSpec[api.Bill]) *cobra.Command {
	var bf bodyFlags
	cmd := &cobra.Command{
		Use:   "comment-update <id> <commentId>",
		Short: "Edit a comment on a bill",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := bf.build()
			if err != nil {
				return err
			}
			client, err := getAPIClient(cmd)
			if err != nil {
				return err
			}
			path := fmt.Sprintf("%s/%s/comments/%s", sp.New(client).Path(), url.PathEscape(args[0]), url.PathEscape(args[1]))
			var out map[string]any
			if err := client.PutInto(cmd.Context(), path, body, &out); err != nil {
				return err
			}
			if len(out) == 0 {
				if !flagDryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "OK: updated comment %s on bill %s\n", args[1], args[0])
				}
				return nil
			}
			return render(cmd, out, nil)
		},
	}
	addBodyFlags(cmd, &bf)
	return cmd
}

// billCommentDeleteCmd: DELETE /bills/{id}/comments/{commentId}.
func billCommentDeleteCmd(sp resourceSpec[api.Bill]) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "comment-delete <id> <commentId>",
		Short: "Delete a comment from a bill",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes && !confirm(cmd, fmt.Sprintf("Delete comment %s from bill %s?", args[1], args[0])) {
				return fmt.Errorf("aborted")
			}
			client, err := getAPIClient(cmd)
			if err != nil {
				return err
			}
			path := fmt.Sprintf("%s/%s/comments/%s", sp.New(client).Path(), url.PathEscape(args[0]), url.PathEscape(args[1]))
			if err := client.DeleteInto(cmd.Context(), path, nil); err != nil {
				return err
			}
			if !flagDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted comment %s from bill %s\n", args[1], args[0])
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")
	return cmd
}

// billAttachmentDeleteCmd: DELETE /bills/attachment/{fileId} (collection-level).
func billAttachmentDeleteCmd(sp resourceSpec[api.Bill]) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "attachment-delete <fileId>",
		Short: "Delete a bill attachment by file ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes && !confirm(cmd, fmt.Sprintf("Delete attachment %s?", args[0])) {
				return fmt.Errorf("aborted")
			}
			client, err := getAPIClient(cmd)
			if err != nil {
				return err
			}
			path := fmt.Sprintf("%s/attachment/%s", sp.New(client).Path(), url.PathEscape(args[0]))
			if err := client.DeleteInto(cmd.Context(), path, nil); err != nil {
				return err
			}
			if !flagDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted attachment %s\n", args[0])
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")
	return cmd
}
