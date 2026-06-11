package commands

import (
	"testing"

	"github.com/njayp/ophis"
	"github.com/stretchr/testify/assert"
)

// TestMCPExcludesSetupCommands locks the MCP tool surface to accounting
// operations: local setup/meta commands must not be exposed as tools, while
// resource operations (including destructive ones) must remain reachable.
func TestMCPExcludesSetupCommands(t *testing.T) {
	exclude := ophis.ExcludeCmdsContaining(mcpExcludedCommands...)

	excluded := [][]string{
		{"agent", "guard"}, {"skills", "install"}, {"auth", "login"},
		{"config", "view"}, {"alias", "set"}, {"init"},
	}
	for _, p := range excluded {
		c := find(t, p...)
		assert.False(t, exclude(c), "%v must be excluded from the MCP tool surface", p)
	}

	included := [][]string{
		{"invoices", "void"}, {"invoices", "list"}, {"contacts", "delete"},
		{"catalog"}, {"doctor"}, {"reports", "sales-by-client"},
	}
	for _, p := range included {
		c := find(t, p...)
		assert.True(t, exclude(c), "%v must stay in the MCP tool surface", p)
	}
}
