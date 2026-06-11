package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// irreversibleVerbs are the resource actions that cannot be undone (fiscal
// emission, voiding, deletion). `alegra agent guard` hard-blocks these by
// default; ordinary writes (create/update/import) only require approval.
var irreversibleVerbs = map[string]bool{
	"delete": true, "void": true, "emit": true,
	"stamp": true, "close": true, "cancel": true, "annul": true,
}

// isIrreversibleVerb reports whether a subcommand name is an irreversible
// action, handling compound names like "attachment-delete" / "comment-delete".
func isIrreversibleVerb(verb string) bool {
	for _, tok := range strings.Split(verb, "-") {
		if irreversibleVerbs[tok] {
			return true
		}
	}
	return false
}

// guardCmd is one API operation the guard config targets.
type guardCmd struct {
	cli  string // CLI path without the root, e.g. "invoices void"
	tool string // MCP tool name, e.g. "alegra_invoices_void"
	verb string // last path segment, e.g. "void"
}

func init() {
	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Helpers for running alegra under an AI agent",
		Long:  "Helpers for running alegra under an AI agent (Claude Code, Codex, OpenCode, …).",
	}
	agentCmd.AddCommand(newAgentGuardCmd())
	rootCmd.AddCommand(readOnlyHints(agentCmd))
}

func newAgentGuardCmd() *cobra.Command {
	var host string
	var allWrites bool
	var write bool

	cmd := &cobra.Command{
		Use:   "guard",
		Short: "Generate agent-safety config that blocks destructive alegra operations",
		Long: `guard generates the permission rules and hooks that stop an AI agent from
running destructive alegra operations, derived from the live command tree so the
list is always complete.

By default it hard-blocks the irreversible actions (delete, void, emit, stamp,
close, cancel) and makes ordinary writes (create, update, import) require
approval; read operations stay allowed. Pass --all-writes to block writes too.

Output is printed for review by default; pass --write to install it. See the
Agent Safety guide: https://jjuanrivvera.github.io/alegra-cli/user-guide/agent-safety/`,
		Args: cobra.NoArgs,
		Example: "  alegra agent guard --host claude-code\n" +
			"  alegra agent guard --host codex\n" +
			"  alegra agent guard --host opencode --all-writes\n" +
			"  alegra agent guard --host claude-code --write",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, writes, irreversible := classifyAPICommands(rootCmd)
			g := guardPlan{
				irreversible: irreversible,
				writes:       writes,
				allWrites:    allWrites,
			}
			switch host {
			case "claude-code", "claude":
				return emitClaudeCode(cmd, g, write)
			case "codex":
				return emitCodex(cmd, g, write)
			case "opencode":
				return emitOpenCode(cmd, g, write)
			case "":
				return fmt.Errorf("--host is required (claude-code, codex, or opencode)")
			default:
				return fmt.Errorf("unknown host %q (use claude-code, codex, or opencode)", host)
			}
		},
	}
	cmd.Flags().StringVar(&host, "host", "", "Target agent host: claude-code, codex, opencode")
	cmd.Flags().BoolVar(&allWrites, "all-writes", false, "Also block create/update/import (default: those require approval)")
	cmd.Flags().BoolVar(&write, "write", false, "Write the config/hook files instead of printing them")
	_ = cmd.RegisterFlagCompletionFunc("host", fixedCompleter("claude-code", "codex", "opencode"))
	return cmd
}

// guardPlan is the classified set of operations the generators turn into config.
type guardPlan struct {
	irreversible []guardCmd
	writes       []guardCmd
	allWrites    bool // fold writes into the hard-blocked set
}

// blocked returns the operations to hard-block (irreversible, plus writes when
// --all-writes); asked returns the ones that only need approval.
func (g guardPlan) blocked() []guardCmd {
	if g.allWrites {
		return append(append([]guardCmd{}, g.irreversible...), g.writes...)
	}
	return g.irreversible
}

func (g guardPlan) asked() []guardCmd {
	if g.allWrites {
		return nil
	}
	return g.writes
}

// blockedVerbs returns the distinct verbs in the hard-block set (for regexes).
func (g guardPlan) blockedVerbs() []string {
	return distinctVerbs(g.blocked())
}

// classifyAPICommands walks the command tree and buckets the operations that hit
// the Alegra API (those carry the openWorldHint annotation, which excludes local
// utility commands like auth/config/skills/agent). Read operations carry
// readOnlyHint; the irreversible verbs are hard-blocked; the rest are writes.
func classifyAPICommands(root *cobra.Command) (read, writes, irreversible []guardCmd) {
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		for _, sub := range c.Commands() {
			if sub.Runnable() && !sub.Hidden && sub.Name() != "help" {
				gc := guardCmd{
					cli:  strings.TrimPrefix(sub.CommandPath(), root.Name()+" "),
					tool: strings.ReplaceAll(sub.CommandPath(), " ", "_"),
					verb: sub.Name(),
				}
				switch {
				case sub.Annotations["openWorldHint"] != "true":
					// Local/utility command (not an API operation) — never gated.
				case sub.Annotations["readOnlyHint"] == "true":
					read = append(read, gc)
				case isIrreversibleVerb(sub.Name()):
					irreversible = append(irreversible, gc)
				default:
					writes = append(writes, gc)
				}
			}
			walk(sub)
		}
	}
	walk(root)
	sortGuard(read)
	sortGuard(writes)
	sortGuard(irreversible)
	return read, writes, irreversible
}

func sortGuard(cs []guardCmd) {
	sort.Slice(cs, func(i, j int) bool { return cs[i].tool < cs[j].tool })
}

func distinctVerbs(cs []guardCmd) []string {
	seen := map[string]bool{}
	var out []string
	for _, c := range cs {
		if !seen[c.verb] {
			seen[c.verb] = true
			out = append(out, c.verb)
		}
	}
	sort.Strings(out)
	return out
}

// writeOrPrint either writes content to path (creating parent dirs) when write
// is set and the file does not already exist, or prints it to the command's
// output with a header. It never overwrites an existing file.
func writeOrPrint(cmd *cobra.Command, write bool, path, content string, perm os.FileMode) error {
	out := cmd.OutOrStdout()
	if !write {
		fmt.Fprintf(out, "# ----- %s -----\n%s\n", path, content)
		return nil
	}
	if _, err := os.Stat(path); err == nil {
		fmt.Fprintf(out, "# %s already exists — review and merge manually:\n%s\n", path, content)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(content), perm); err != nil {
		return err
	}
	fmt.Fprintf(out, "wrote %s\n", path)
	return nil
}
