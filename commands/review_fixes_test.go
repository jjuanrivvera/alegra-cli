package commands

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

// TestDryRun_RedactsCredentialsInCurl pins T3: the dry-run curl is wired through
// cmd.OutOrStdout() (so it's capturable) and never prints the live API token.
func TestDryRun_RedactsCredentialsInCurl(t *testing.T) {
	srv := actionServer(t)
	setActionEnv(t, srv)
	t.Setenv(config.EnvToken, "s3cr3t-zzz") // distinctive: won't collide with "token"/"redacted"

	out, err := runRoot(t, "contacts", "list", "--dry-run")
	require.NoError(t, err)
	assert.Contains(t, out, "curl -X GET")
	assert.Contains(t, out, "Authorization: <redacted>")
	assert.NotContains(t, out, "s3cr3t-zzz", "the API token must never appear in dry-run output")

	// Even with --show-token, Basic auth shows only the scheme, not the secret.
	out, err = runRoot(t, "contacts", "list", "--dry-run", "--show-token")
	require.NoError(t, err)
	assert.Contains(t, out, "Basic <base64(email:token)>")
	assert.NotContains(t, out, "s3cr3t-zzz")
}

// TestExportDryRun_DoesNotTruncateFile pins H1: `export --out <file> --dry-run`
// must not create or truncate the target.
func TestExportDryRun_DoesNotTruncateFile(t *testing.T) {
	srv := actionServer(t)
	setActionEnv(t, srv)

	outFile := filepath.Join(t.TempDir(), "important.csv")
	const original = "DO NOT DELETE\n"
	require.NoError(t, os.WriteFile(outFile, []byte(original), 0o600))

	_, err := runRoot(t, "contacts", "export", "--out", outFile, "--dry-run")
	require.NoError(t, err)

	data, rerr := os.ReadFile(outFile)
	require.NoError(t, rerr)
	assert.Equal(t, original, string(data), "dry-run export must leave the existing file untouched")
}

// TestImport_NegativePath pins T5: when a row fails, the import reports the
// failure and returns a non-nil error instead of a success-shaped tally.
func TestImport_NegativePath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`{"message":"invalid"}`))
			return
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	t.Cleanup(srv.Close)
	setActionEnv(t, srv)

	csvFile := filepath.Join(t.TempDir(), "rows.csv")
	require.NoError(t, os.WriteFile(csvFile, []byte("name\nBeta\n"), 0o600))

	out, err := runRoot(t, "contacts", "import", "-f", csvFile)
	require.Error(t, err, "an import where every row fails must return an error")
	assert.Contains(t, out, "Imported 0, failed 1")
	assert.Contains(t, out, "FAILED")
}

// --- T4: the generated Claude Code guard hook is executed, not just string-matched ---

func runGuardHook(t *testing.T, scriptPath, payload string, env []string) string {
	t.Helper()
	cmd := exec.Command("bash", scriptPath)
	cmd.Stdin = strings.NewReader(payload)
	if env != nil {
		cmd.Env = env
	}
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "guard hook should exit 0; output: %s", out)
	return string(out)
}

func TestClaudeGuardHook_DenyAndAllow(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("bash not available")
	}
	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq not available")
	}

	script := claudeHookScript([]string{"void", "delete", "emit"})
	scriptPath := filepath.Join(t.TempDir(), "alegra-guard.sh")
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0o755))

	cases := []struct {
		name    string
		payload string
		deny    bool
	}{
		{"bash void", `{"tool_name":"Bash","tool_input":{"command":"alegra invoices void 5"}}`, true},
		{"bash quote-obfuscated void", `{"tool_name":"Bash","tool_input":{"command":"alegra invoices vo\"\"id 5"}}`, true},
		{"bash delete", `{"tool_name":"Bash","tool_input":{"command":"alegra contacts delete 3"}}`, true},
		{"bash list", `{"tool_name":"Bash","tool_input":{"command":"alegra contacts list"}}`, false},
		{"mcp void", `{"tool_name":"mcp__server_alegra_invoices_void"}`, true},
		{"mcp list", `{"tool_name":"mcp__server_alegra_contacts_list"}`, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := runGuardHook(t, scriptPath, c.payload, nil)
			if c.deny {
				assert.Contains(t, out, `"permissionDecision":"deny"`, "expected deny for %s", c.payload)
			} else {
				assert.NotContains(t, out, "deny", "expected allow for %s", c.payload)
			}
		})
	}
}

// TestClaudeGuardHook_FailsSafeWithoutJq pins M4: without jq the hook must still
// deny an irreversible verb (fail safe), not fall through to allow.
func TestClaudeGuardHook_FailsSafeWithoutJq(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("bash not available")
	}
	// Build a PATH that has the coreutils the no-jq branch needs but NOT jq.
	binDir := t.TempDir()
	for _, tool := range []string{"cat", "tr", "grep"} {
		p, lerr := exec.LookPath(tool)
		if lerr != nil {
			t.Skipf("%s not available", tool)
		}
		require.NoError(t, os.Symlink(p, filepath.Join(binDir, tool)))
	}

	script := claudeHookScript([]string{"void", "delete", "emit"})
	scriptPath := filepath.Join(t.TempDir(), "alegra-guard.sh")
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0o755))

	// PATH=binDir hides jq (it isn't symlinked there) while keeping cat/tr/grep.
	env := []string{"PATH=" + binDir}
	deny := runGuardHook(t, scriptPath, `{"tool_name":"Bash","tool_input":{"command":"alegra invoices void 5"}}`, env)
	assert.Contains(t, deny, `"permissionDecision":"deny"`,
		"without jq, an irreversible verb must still be denied (fail safe, not open)")

	allow := runGuardHook(t, scriptPath, `{"tool_name":"Bash","tool_input":{"command":"alegra contacts list"}}`, env)
	assert.NotContains(t, allow, "deny", "a read command must still be allowed in the no-jq path")
}
