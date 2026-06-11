package commands

import "github.com/njayp/ophis"

// mcpExcludedCommands are local setup and meta commands that are not accounting
// operations, so they stay out of the agent-facing MCP tool surface (an agent
// has no business calling `alegra agent guard` or `alegra auth login`).
// mcp/completion/help are already excluded by ophis itself.
var mcpExcludedCommands = []string{"agent", "skills", "auth", "config", "alias", "init"}

// init registers the `alegra mcp` command, which exposes the CLI's accounting
// operations as a Model Context Protocol server so AI agents can drive Alegra.
func init() {
	rootCmd.AddCommand(ophis.Command(&ophis.Config{
		ToolNamePrefix: "alegra",
		Selectors: []ophis.Selector{
			{
				CmdSelector:           ophis.ExcludeCmdsContaining(mcpExcludedCommands...),
				InheritedFlagSelector: ophis.ExcludeFlags("show-token", "profile"),
			},
		},
	}))
}
