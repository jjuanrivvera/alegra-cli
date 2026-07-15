package commands

import (
	"slices"

	"github.com/njayp/ophis"
	"github.com/spf13/cobra"
)

// mcpExcludedCommands are local setup and meta commands that are not accounting
// operations, so they stay out of the agent-facing MCP tool surface (an agent
// has no business calling `alegra agent guard` or `alegra auth login`).
// mcp/completion/help are already excluded by ophis itself.
//
// These are matched as path substrings (ExcludeCmdsContaining) so a whole
// subtree — e.g. every `auth …` leaf — is dropped by its parent name. Exact-path
// matching would leak the subcommands back into the surface. The trade-off is
// that a future resource whose name contains one of these words would be
// over-excluded; that only ever removes a tool (never leaks one), and
// TestMCPExcludesSetupCommands guards the surface either way.
var mcpExcludedCommands = []string{"agent", "skills", "auth", "config", "alias", "init"}

// selfUpdatePaths are the EXACT command paths of the self-update command. `update` replaces
// the running binary, so it must never be an MCP tool — but it can't go in mcpExcludedCommands,
// whose substring match would also drop every `<resource> update` action. Excluded by exact
// path instead.
var selfUpdatePaths = []string{"alegra update", "alegra update check"}

// init registers the `alegra mcp` command, which exposes the CLI's accounting
// operations as a Model Context Protocol server so AI agents can drive Alegra.
func init() {
	excludeMeta := ophis.ExcludeCmdsContaining(mcpExcludedCommands...)
	rootCmd.AddCommand(ophis.Command(&ophis.Config{
		ToolNamePrefix: "alegra",
		Selectors: []ophis.Selector{
			{
				CmdSelector: func(cmd *cobra.Command) bool {
					return excludeMeta(cmd) && !slices.Contains(selfUpdatePaths, cmd.CommandPath())
				},
				InheritedFlagSelector: ophis.ExcludeFlags("show-token", "profile"),
			},
		},
	}))
}
