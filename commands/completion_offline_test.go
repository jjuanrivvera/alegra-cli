package commands

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

// The ID completer is best-effort: any failure to reach the API must yield no
// suggestions (never an error or a hung shell). These drive the real __complete
// path through the failure branches.
func TestResourceIDCompleter_OfflineTolerant(t *testing.T) {
	t.Run("no credentials yields no suggestions", func(t *testing.T) {
		keyring.MockInit() // empty keyring: no token for any profile
		t.Setenv(config.EnvEmail, "")
		t.Setenv(config.EnvToken, "")
		t.Setenv(config.EnvProfile, "")
		t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))

		out, err := runRoot(t, "__complete", "contacts", "get", "")
		require.NoError(t, err)
		assert.NotRegexp(t, `(?m)^\d`, out, "no id completions when unauthenticated")
	})

	t.Run("API error yields no suggestions", func(t *testing.T) {
		keyring.MockInit()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		t.Cleanup(srv.Close)
		t.Setenv(config.EnvBaseURL, srv.URL)
		t.Setenv(config.EnvEmail, "e@x.com")
		t.Setenv(config.EnvToken, "tok")
		t.Setenv(config.EnvProfile, "")
		t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))

		out, err := runRoot(t, "__complete", "contacts", "get", "")
		require.NoError(t, err)
		assert.NotRegexp(t, `(?m)^\d`, out, "no id completions when the API errors")
	})
}

// idCompletions normalizes records via JSON; a non-object element type can't be
// projected to id/label and must degrade to no suggestions rather than panic.
func TestIDCompletions_NonObject(t *testing.T) {
	assert.Nil(t, idCompletions([]string{"a", "b"}, ""))
}

// withColumns is a no-op for an empty column set (no annotation is attached).
func TestWithColumns_Empty(t *testing.T) {
	cmd := &cobra.Command{}
	withColumns(cmd, nil)
	_, ok := cmd.Annotations[columnsAnnotation]
	assert.False(t, ok)
}

func TestCompleteProfiles_OfflineBranches(t *testing.T) {
	t.Run("unreadable config yields no suggestions", func(t *testing.T) {
		bad := filepath.Join(t.TempDir(), "config.yaml")
		require.NoError(t, os.WriteFile(bad, []byte("foo: bar: baz\n"), 0o600)) // invalid YAML
		t.Setenv(config.EnvConfig, bad)

		got, dir := completeProfiles(nil, nil, "")
		assert.Nil(t, got)
		assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, dir)
	})

	t.Run("profile without email completes the bare name", func(t *testing.T) {
		t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))
		cfg := config.New()
		cfg.SetProfile(&config.Profile{Name: "noemail"})
		require.NoError(t, cfg.Save())

		got, _ := completeProfiles(nil, nil, "")
		assert.Equal(t, []string{"noemail"}, got)
	})
}
