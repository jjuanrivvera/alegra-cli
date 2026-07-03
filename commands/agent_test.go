package commands

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func bucketTools(cs []guardCmd) map[string]bool {
	m := map[string]bool{}
	for _, c := range cs {
		m[c.tool] = true
	}
	return m
}

func TestClassifyAPICommands(t *testing.T) {
	read, writes, irreversible := classifyAPICommands(RootCmd())

	r, w, irr := bucketTools(read), bucketTools(writes), bucketTools(irreversible)

	// Reads.
	assert.True(t, r["alegra_invoices_list"], "list is read")
	assert.True(t, r["alegra_invoices_get"], "get is read")
	assert.True(t, r["alegra_catalog"], "catalog is read")

	// Irreversible — including compound -delete actions.
	assert.True(t, irr["alegra_invoices_void"], "void is irreversible")
	assert.True(t, irr["alegra_invoices_emit"], "emit is irreversible")
	assert.True(t, irr["alegra_contacts_delete"], "delete is irreversible")
	assert.True(t, irr["alegra_bills_attachment-delete"], "attachment-delete is irreversible")

	// Ordinary writes.
	assert.True(t, w["alegra_contacts_create"], "create is a write")
	assert.True(t, w["alegra_contacts_update"], "update is a write")

	// The hand-built company singleton must be annotated: "company update" is a
	// remote PUT — before the audit fix it carried no openWorldHint and was
	// silently exempted from the guard.
	assert.True(t, r["alegra_company_get"], "company get is read")
	assert.True(t, w["alegra_company_update"], "company update is an approval-gated write")

	// Local utility commands (no API call) must never be classified/gated.
	for _, m := range []map[string]bool{r, w, irr} {
		assert.False(t, m["alegra_agent_guard"], "guard is not an API op")
		assert.False(t, m["alegra_skills_install"], "skills is not an API op")
		assert.False(t, m["alegra_auth_login"], "auth is not an API op")
	}

	// No operation is in two buckets at once.
	for tool := range irr {
		assert.False(t, w[tool], "%s must not be both irreversible and write", tool)
		assert.False(t, r[tool], "%s must not be both irreversible and read", tool)
	}
}

// guardLocalGroups are the top-level command groups that never call the Alegra
// API. Every runnable leaf OUTSIDE these groups must carry openWorldHint so the
// guard classifies it — classification silently skips unannotated commands, so
// a missing hint on an API command exempts it from gating entirely (this is
// exactly how "company update" escaped the guard before the audit).
var guardLocalGroups = map[string]bool{
	"agent":      true,
	"alias":      true,
	"auth":       true,
	"catalog":    true, // embedded reference data; sync-sat downloads SAT (not Alegra) data
	"completion": true,
	"config":     true,
	"init":       true,
	"mcp":        true,
	"skills":     true,
	"help":       true,
}

func TestEveryAPICommandIsAnnotated(t *testing.T) {
	root := RootCmd()
	var walk func(c *cobra.Command, top string)
	walk = func(c *cobra.Command, top string) {
		for _, sub := range c.Commands() {
			name := top
			if name == "" {
				name = sub.Name()
			}
			if sub.Runnable() && !sub.Hidden && sub.Name() != "help" && !guardLocalGroups[name] {
				assert.Equalf(t, "true", sub.Annotations["openWorldHint"],
					"%s is outside the local groups but has no openWorldHint — the agent guard would silently skip it", sub.CommandPath())
			}
			walk(sub, name)
		}
	}
	walk(root, "")
}

func runGuard(t *testing.T, args ...string) string {
	t.Helper()
	var buf bytes.Buffer
	c := newAgentGuardCmd()
	c.SetOut(&buf)
	c.SetErr(&buf)
	c.SetArgs(args)
	require.NoError(t, c.Execute())
	return buf.String()
}

func TestGuard_ClaudeCode(t *testing.T) {
	out := runGuard(t, "--host", "claude-code")
	// The hook is emitted and covers both surfaces.
	assert.Contains(t, out, "alegra-guard.sh")
	assert.Contains(t, out, "permissionDecision")
	assert.Contains(t, out, "void")

	// The settings JSON block (after its file marker) parses and denies the
	// irreversible verbs.
	marker := strings.Index(out, "settings.json")
	require.GreaterOrEqual(t, marker, 0)
	start := strings.Index(out[marker:], "{")
	require.GreaterOrEqual(t, start, 0)
	var settings map[string]any
	require.NoError(t, json.Unmarshal([]byte(out[marker+start:]), &settings), "settings must be valid JSON")
	perms := settings["permissions"].(map[string]any)
	deny := perms["deny"].([]any)
	denyStr := joinAny(deny)
	assert.Contains(t, denyStr, "delete")
	assert.Contains(t, denyStr, "void")
	// create is approval-gated, not hard-denied.
	assert.NotContains(t, denyStr, "create:")
}

func TestGuard_AllWritesFoldsIntoDeny(t *testing.T) {
	out := runGuard(t, "--host", "opencode", "--all-writes")
	assert.Contains(t, out, `"alegra_contacts_create": "deny"`)
	assert.Contains(t, out, `"alegra contacts create": "deny"`)
	assert.NotContains(t, out, `": "ask"`)
}

func TestGuard_Codex(t *testing.T) {
	out := runGuard(t, "--host", "codex")
	assert.Contains(t, out, `sandbox_mode    = "read-only"`)
	assert.Contains(t, out, "approval")
}

func TestGuard_UnknownHost(t *testing.T) {
	c := newAgentGuardCmd()
	c.SetOut(&bytes.Buffer{})
	c.SetArgs([]string{"--host", "nope"})
	assert.Error(t, c.Execute())
}

func joinAny(xs []any) string {
	parts := make([]string, 0, len(xs))
	for _, x := range xs {
		parts = append(parts, x.(string))
	}
	return strings.Join(parts, " ")
}
