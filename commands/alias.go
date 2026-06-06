package commands

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

func init() {
	aliasCmd := &cobra.Command{
		Use:   "alias",
		Short: "Save and manage command aliases",
		Long: `Aliases are shortcuts for longer commands, stored in your config.

  alegra alias set unpaid "invoices list --status open --all"
  alegra unpaid --client-id 12     # expands, then appends your extra args

Aliases never shadow built-in commands.`,
	}
	aliasCmd.AddCommand(newAliasSetCmd(), newAliasListCmd(), newAliasRemoveCmd())
	rootCmd.AddCommand(aliasCmd)
}

func newAliasSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <expansion>",
		Short: "Create or update an alias",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, expansion := args[0], args[1]
			if strings.HasPrefix(name, "-") {
				return fmt.Errorf("alias name must not start with '-'")
			}
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if isBuiltinCommand(name) {
				return fmt.Errorf("%q is a built-in command and cannot be aliased", name)
			}
			if cfg.Aliases == nil {
				cfg.Aliases = map[string]string{}
			}
			cfg.Aliases[name] = expansion
			if err := cfg.Save(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Saved alias %q → %s\n", name, expansion)
			return nil
		},
	}
}

func newAliasListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List aliases",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if len(cfg.Aliases) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no aliases)")
				return nil
			}
			names := make([]string, 0, len(cfg.Aliases))
			for n := range cfg.Aliases {
				names = append(names, n)
			}
			sort.Strings(names)
			for _, n := range names {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", n, cfg.Aliases[n])
			}
			return nil
		},
	}
}

func newAliasRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove <name>",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove an alias",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if _, ok := cfg.Aliases[args[0]]; !ok {
				return fmt.Errorf("no such alias %q", args[0])
			}
			delete(cfg.Aliases, args[0])
			if err := cfg.Save(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed alias %q\n", args[0])
			return nil
		},
	}
}

// isBuiltinCommand reports whether name is a registered top-level command or one
// of its aliases.
func isBuiltinCommand(name string) bool {
	for _, c := range rootCmd.Commands() {
		if c.Name() == name || slices.Contains(c.Aliases, name) {
			return true
		}
	}
	return false
}

// expandAliasArgs rewrites os.Args so a leading alias token is replaced by its
// expansion, with any trailing user args preserved. It is a no-op when the first
// token is a flag, a built-in command, or not a known alias.
func expandAliasArgs(argv []string, aliases map[string]string, isBuiltin func(string) bool) []string {
	if len(argv) < 2 {
		return argv
	}
	first := argv[1]
	if first == "" || strings.HasPrefix(first, "-") || isBuiltin(first) {
		return argv
	}
	expansion, ok := aliases[first]
	if !ok {
		return argv
	}
	toks, err := shellSplit(expansion)
	if err != nil || len(toks) == 0 {
		return argv
	}
	out := make([]string, 0, len(argv)+len(toks))
	out = append(out, argv[0])
	out = append(out, toks...)
	out = append(out, argv[2:]...)
	return out
}

// shellSplit splits a string into tokens, honoring single and double quotes.
func shellSplit(s string) ([]string, error) {
	var tokens []string
	var cur strings.Builder
	inSingle, inDouble, hasToken := false, false, false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case inSingle:
			if c == '\'' {
				inSingle = false
			} else {
				cur.WriteByte(c)
			}
		case inDouble:
			if c == '"' {
				inDouble = false
			} else {
				cur.WriteByte(c)
			}
		case c == '\'':
			inSingle, hasToken = true, true
		case c == '"':
			inDouble, hasToken = true, true
		case c == ' ' || c == '\t' || c == '\n':
			if hasToken {
				tokens = append(tokens, cur.String())
				cur.Reset()
				hasToken = false
			}
		default:
			cur.WriteByte(c)
			hasToken = true
		}
	}
	if inSingle || inDouble {
		return nil, fmt.Errorf("unterminated quote in alias expansion")
	}
	if hasToken {
		tokens = append(tokens, cur.String())
	}
	return tokens, nil
}
