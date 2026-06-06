package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShellSplit(t *testing.T) {
	cases := map[string][]string{
		`invoices list --status open`: {"invoices", "list", "--status", "open"},
		`a "b c" d`:                   {"a", "b c", "d"},
		`x --set 'name=Acme S.A.S'`:   {"x", "--set", "name=Acme S.A.S"},
		`  leading   and   spaces  `:  {"leading", "and", "spaces"},
		`--data '{"a":1,"b":"c d"}'`:  {"--data", `{"a":1,"b":"c d"}`},
	}
	for in, want := range cases {
		got, err := shellSplit(in)
		require.NoError(t, err, in)
		assert.Equal(t, want, got, in)
	}
	_, err := shellSplit(`unterminated "quote`)
	assert.Error(t, err)
}

func TestExpandAliasArgs(t *testing.T) {
	aliases := map[string]string{"unpaid": "invoices list --status open --all"}
	builtin := func(s string) bool { return s == "invoices" || s == "config" }

	// alias expands and preserves trailing args
	got := expandAliasArgs([]string{"alegra", "unpaid", "--client-id", "12"}, aliases, builtin)
	assert.Equal(t, []string{"alegra", "invoices", "list", "--status", "open", "--all", "--client-id", "12"}, got)

	// built-in command is never expanded
	got = expandAliasArgs([]string{"alegra", "invoices", "list"}, aliases, builtin)
	assert.Equal(t, []string{"alegra", "invoices", "list"}, got)

	// unknown token untouched
	got = expandAliasArgs([]string{"alegra", "whatever"}, aliases, builtin)
	assert.Equal(t, []string{"alegra", "whatever"}, got)

	// flags / no args untouched
	assert.Equal(t, []string{"alegra", "--help"}, expandAliasArgs([]string{"alegra", "--help"}, aliases, builtin))
	assert.Equal(t, []string{"alegra"}, expandAliasArgs([]string{"alegra"}, aliases, builtin))
}
