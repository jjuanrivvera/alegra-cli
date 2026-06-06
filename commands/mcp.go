package commands

import "github.com/njayp/ophis"

// init registers the `alegra mcp` command, which exposes the entire CLI command
// tree as a Model Context Protocol server so AI agents can drive Alegra.
func init() {
	rootCmd.AddCommand(ophis.Command(&ophis.Config{
		ToolNamePrefix: "alegra",
		Selectors: []ophis.Selector{
			{
				// Never expose credential-related inherited flags as MCP inputs.
				InheritedFlagSelector: ophis.ExcludeFlags("show-token", "profile"),
			},
		},
	}))
}
