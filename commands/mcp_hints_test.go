package commands

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// find walks the registered command tree by path (e.g. "invoices", "void").
func find(t *testing.T, path ...string) *cobra.Command {
	t.Helper()
	cur := RootCmd()
	for _, name := range path {
		var next *cobra.Command
		for _, c := range cur.Commands() {
			if c.Name() == name {
				next = c
				break
			}
		}
		require.NotNilf(t, next, "command %v not found (missing %q)", path, name)
		cur = next
	}
	return cur
}

// TestMCPHints asserts the MCP tool annotations `alegra mcp` exposes (via ophis)
// so hosts that honor them gate writes. Read commands must be read-only; writes
// and custom actions must not be read-only, and the destructive ones must say so.
func TestMCPHints(t *testing.T) {
	readOnly := []([]string){
		{"contacts", "list"}, {"contacts", "get"}, {"contacts", "export"},
		{"invoices", "list"}, {"items", "stock"},
		{"catalog"}, {"catalog", "product-keys"},
		{"reports", "sales-by-client"}, {"reports", "account-statement"},
		{"doctor"}, {"version"},
		// Read-only custom GET subcommands that must opt back in to read-only
		// hints rather than defaulting to destructive (M1).
		{"journals", "balance"}, {"users", "self"}, {"categories", "settings"},
	}
	for _, p := range readOnly {
		c := find(t, p...)
		assert.Equal(t, "true", c.Annotations["readOnlyHint"], "%v should be read-only", p)
		assert.NotEqual(t, "true", c.Annotations["destructiveHint"], "%v must not be destructive", p)
	}

	destructive := []([]string){
		{"contacts", "delete"}, {"contacts", "update"},
		{"invoices", "void"}, {"invoices", "emit"},
		// The PUT counterpart of categories settings must stay destructive (M1).
		{"categories", "set-settings"},
	}
	for _, p := range destructive {
		c := find(t, p...)
		assert.Equal(t, "true", c.Annotations["destructiveHint"], "%v should be destructive", p)
		assert.NotEqual(t, "true", c.Annotations["readOnlyHint"], "%v must not be read-only", p)
	}

	// A write that is not data-destructive (create) is not read-only, so an MCP
	// host's default (destructive) still gates it.
	create := find(t, "contacts", "create")
	assert.NotEqual(t, "true", create.Annotations["readOnlyHint"], "create must not be read-only")

	// Annotating a command must never clobber the completion columns annotation.
	list := find(t, "contacts", "list")
	assert.NotEmpty(t, list.Annotations["alegra/columns"], "columns annotation must survive MCP hints")
}
