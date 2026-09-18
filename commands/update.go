package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/update"
	"github.com/jjuanrivvera/alegra-cli/internal/version"
)

// `update` self-replaces the running binary from the latest GitHub release (checksum-verified).
// It is a LOCAL command — not an Alegra API op — so it is registered WITHOUT openWorld hints
// (the agent guard must not classify it), listed in guardLocalGroups, and kept out of the MCP
// surface by exact path (see mcp.go). This complements the existing `version --check`, which
// only reports availability and points at brew.
func init() {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update alegra to the latest release",
		Long: `Check GitHub for a newer release and, if one exists, download it, verify it
against the release checksums, and replace the running binary in place.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), 60*time.Second)
			defer cancel()
			res := update.NewUpdater(version.Short()).CheckAndUpdate(ctx)
			if res.Error != nil {
				return res.Error
			}
			out := cmd.OutOrStdout()
			if res.Updated {
				fmt.Fprintf(out, "Updated %s → %s. Restart to use the new version.\n", res.FromVersion, res.ToVersion)
			} else {
				fmt.Fprintln(out, "Already on the latest version.")
			}
			return nil
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "check",
		Short: "Check for a newer release without installing it",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
			defer cancel()
			cur := version.Short()
			rel, err := update.NewUpdater(cur).GetLatestRelease(ctx)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Current: %s\nLatest:  %s\n", cur, rel.TagName)
			if update.IsNewer(rel.TagName, cur) {
				fmt.Fprintln(out, "An update is available. Run `alegra update` to install it.")
			} else {
				fmt.Fprintln(out, "You are on the latest version.")
			}
			return nil
		},
	})
	rootCmd.AddCommand(cmd)
}
