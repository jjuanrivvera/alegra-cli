package commands

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// guardHookTestCmds is a small fixed blocked set shared by the hook tests that
// don't need the full command tree.
func guardHookTestCmds() []guardCmd {
	return []guardCmd{
		{cli: "invoices void", tool: "alegra_invoices_void", verb: "void"},
		{cli: "invoices emit", tool: "alegra_invoices_emit", verb: "emit"},
		{cli: "contacts delete", tool: "alegra_contacts_delete", verb: "delete"},
	}
}

// TestClaudeHookScript_BashExecution exercises the generated hook script with
// real bash against the full command tree, covering the adversarial cases from
// the agent-guard audit: obfuscation, command-position anchoring, path-invoked
// binaries, and benign commands that merely mention a blocked verb.
func TestClaudeHookScript_BashExecution(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("bash hook tests require a POSIX shell; skipping on windows")
	}
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found in PATH; skipping hook execution tests")
	}
	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq not available")
	}

	// Generate the hook from the real command tree so the blocked_cmds /
	// blocked_tools arrays are fully populated.
	_, _, irreversible := classifyAPICommands(RootCmd())
	hookContent := claudeHookScript(irreversible)
	hookFile := filepath.Join(t.TempDir(), "alegra-guard.sh")
	if err := os.WriteFile(hookFile, []byte(hookContent), 0o755); err != nil {
		t.Fatalf("write hook: %v", err)
	}

	bashPayload := func(command string) string {
		b, _ := json.Marshal(map[string]any{
			"tool_name":  "Bash",
			"tool_input": map[string]any{"command": command},
		})
		return string(b)
	}
	mcpPayload := func(toolName string) string {
		b, _ := json.Marshal(map[string]any{
			"tool_name":  toolName,
			"tool_input": map[string]any{},
		})
		return string(b)
	}

	runHook := func(t *testing.T, payload string) string {
		t.Helper()
		cmd := exec.Command(bash, hookFile)
		cmd.Stdin = strings.NewReader(payload)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("hook script exited non-zero: %v\noutput: %s", err, out)
		}
		return string(out)
	}
	isDenied := func(output string) bool {
		return strings.Contains(output, `"permissionDecision":"deny"`)
	}

	cases := []struct {
		name        string
		payload     string
		wantDenied  bool
		description string
	}{
		{
			name:        "bash_invoices_delete_denied",
			payload:     bashPayload("alegra invoices delete 5"),
			wantDenied:  true,
			description: "direct alegra invoices delete must be denied",
		},
		{
			name:        "bash_invoices_void_denied",
			payload:     bashPayload("alegra invoices void 12"),
			wantDenied:  true,
			description: "alegra invoices void must be denied",
		},
		{
			name:        "bash_invoices_emit_denied",
			payload:     bashPayload("alegra invoices emit 12"),
			wantDenied:  true,
			description: "alegra invoices emit (DIAN fiscal event) must be denied",
		},
		{
			name:        "bash_compound_attachment_delete_denied",
			payload:     bashPayload("alegra bills attachment-delete 3 9"),
			wantDenied:  true,
			description: "compound attachment-delete must be denied",
		},
		{
			name:        "bash_obfuscated_delete_denied",
			payload:     bashPayload(`alegra invoices de""lete 5`),
			wantDenied:  true,
			description: "quote-obfuscated delete must be denied",
		},
		{
			name:        "bash_backslash_obfuscated_denied",
			payload:     bashPayload(`alegra invoices de\lete 5`),
			wantDenied:  true,
			description: "backslash-obfuscated delete must be denied",
		},
		{
			name:        "bash_after_semicolon_denied",
			payload:     bashPayload("echo hi; alegra invoices delete 5"),
			wantDenied:  true,
			description: "blocked command after ; must be denied",
		},
		{
			name:        "bash_after_pipe_denied",
			payload:     bashPayload("echo hi | alegra invoices delete 5"),
			wantDenied:  true,
			description: "blocked command after | must be denied",
		},
		{
			name:        "bash_after_and_denied",
			payload:     bashPayload("true && alegra invoices delete 5"),
			wantDenied:  true,
			description: "blocked command after && must be denied",
		},
		{
			name:        "bash_env_prefix_denied",
			payload:     bashPayload("env X=1 alegra invoices delete 5"),
			wantDenied:  true,
			description: "env-prefixed blocked command must be denied",
		},
		{
			name:        "bash_relative_path_binary_denied",
			payload:     bashPayload("./bin/alegra invoices delete 5"),
			wantDenied:  true,
			description: "path-invoked binary ./bin/alegra must be denied",
		},
		{
			name:        "bash_absolute_path_binary_denied",
			payload:     bashPayload("/usr/local/bin/alegra invoices delete 5"),
			wantDenied:  true,
			description: "path-invoked binary /usr/local/bin/alegra must be denied",
		},
		{
			name:        "bash_read_command_allowed",
			payload:     bashPayload("alegra invoices list --limit 5"),
			wantDenied:  false,
			description: "a read command must NOT be denied",
		},
		{
			name:        "bash_benign_arg_with_blocked_verb_allowed",
			payload:     bashPayload(`alegra invoices create --data "note: delete later"`),
			wantDenied:  false,
			description: "a write whose ARGUMENT contains a blocked verb must NOT be denied",
		},
		{
			name:        "bash_benign_arg_void_word_allowed",
			payload:     bashPayload(`alegra contacts update 5 --set name="Void Setup SAS"`),
			wantDenied:  false,
			description: "an argument containing 'void' must NOT be denied",
		},
		{
			name:        "bash_unrelated_file_allowed",
			payload:     bashPayload("cat something_delete.go"),
			wantDenied:  false,
			description: "non-alegra Bash command mentioning 'delete' must NOT be denied",
		},
		{
			// Accepted safe-direction false positive: quote-stripping (the
			// defense against vo""id-style obfuscation) makes a QUOTED blocked
			// command indistinguishable from a real one, so it is denied.
			name:        "bash_rg_quoting_blocked_cmd_denied_conservative",
			payload:     bashPayload(`rg "alegra invoices delete" src/`),
			wantDenied:  true,
			description: "a quoted blocked command in an argument is denied (conservative; quoting is the obfuscation channel)",
		},
		{
			name:        "bash_other_binary_named_like_alegra_allowed",
			payload:     bashPayload("myalegra invoices delete 5"),
			wantDenied:  false,
			description: "a different binary whose name merely ends in 'alegra' must NOT be denied",
		},
		{
			name:        "bash_glued_separator_denied",
			payload:     bashPayload("alegra invoices void;true"),
			wantDenied:  true,
			description: "shell separator glued directly to the verb must still be denied",
		},
		{
			name:        "bash_glued_pipe_denied",
			payload:     bashPayload("alegra invoices emit|cat"),
			wantDenied:  true,
			description: "pipe glued directly to the verb must still be denied",
		},
		{
			name:        "mcp_invoices_delete_denied",
			payload:     mcpPayload("mcp__alegra__alegra_invoices_delete"),
			wantDenied:  true,
			description: "MCP alegra_invoices_delete tool must be denied",
		},
		{
			name:        "mcp_attachment_delete_denied",
			payload:     mcpPayload("mcp__alegra__alegra_bills_attachment-delete"),
			wantDenied:  true,
			description: "MCP compound attachment-delete tool must be denied",
		},
		{
			name:        "mcp_invoices_list_allowed",
			payload:     mcpPayload("mcp__alegra__alegra_invoices_list"),
			wantDenied:  false,
			description: "MCP alegra_invoices_list tool must NOT be denied",
		},
		{
			name:        "mcp_near_miss_open_allowed",
			payload:     mcpPayload("mcp__alegra__alegra_payments_open"),
			wantDenied:  false,
			description: "MCP payments_open (approval-gated write, not irreversible) must NOT be hook-denied",
		},
		{
			name:        "mcp_other_server_allowed",
			payload:     mcpPayload("mcp__other__somethingelse_delete_local"),
			wantDenied:  false,
			description: "a non-alegra MCP tool must NOT be denied",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			output := runHook(t, tc.payload)
			denied := isDenied(output)
			if denied != tc.wantDenied {
				want := "allowed"
				if tc.wantDenied {
					want = "denied"
				}
				got := "allowed"
				if denied {
					got = "denied"
				}
				t.Errorf("%s: want %s, got %s\noutput: %s", tc.description, want, got, output)
			}
		})
	}
}
