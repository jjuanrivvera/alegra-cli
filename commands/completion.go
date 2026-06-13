package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

// Dynamic shell completion. Cobra ships static command/flag-name completion for
// free; the functions here add live, data-aware suggestions: real record IDs for
// `get`/`update`/`delete`/actions, configured profiles, resource columns, and
// enum filter values. Every completer is best-effort and offline-tolerant — on
// any error (no credentials, network failure, timeout) it yields no suggestions
// rather than blocking the shell.

const (
	// completionLimit caps how many records a live ID completer fetches. One
	// max-size page keeps a <Tab> press snappy while still covering the recent
	// records a user is most likely reaching for.
	completionLimit = 30
	// completionTimeout bounds the single API call behind a <Tab> press so a slow
	// network degrades to "no suggestions" instead of a hung terminal.
	completionTimeout = 4 * time.Second
)

// completeFunc is Cobra's completion callback signature.
type completeFunc = func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective)

// resourceIDCompleter returns a completer that fetches a page of the resource and
// offers each record's id, annotated with a human label (name, number, …) shown
// by description-aware shells (zsh, fish). It completes only the first argument.
func resourceIDCompleter[T any](sp resourceSpec[T]) completeFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Only the <id> positional is completable; later args (if any) are not ids.
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		client, err := getAPIClient(cmd)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		ctx, cancel := context.WithTimeout(cmd.Context(), completionTimeout)
		defer cancel()
		items, err := sp.New(client).List(ctx, api.ListParams{Limit: completionLimit, Extra: url.Values{}})
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return idCompletions(items, toComplete), cobra.ShellCompDirectiveNoFileComp
	}
}

// idCompletions projects records to "id\tlabel" completion entries, filtered to
// those whose id has the typed prefix. Records are normalized through JSON so the
// same logic works for every resource type without per-resource code.
func idCompletions[T any](items []T, toComplete string) []string {
	raw, err := json.Marshal(items)
	if err != nil {
		return nil
	}
	var rows []map[string]any
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil
	}
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		id := scalar(r["id"])
		if id == "" || (toComplete != "" && !strings.HasPrefix(id, toComplete)) {
			continue
		}
		if label := completionLabel(r); label != "" {
			out = append(out, id+"\t"+label)
		} else {
			out = append(out, id)
		}
	}
	return out
}

// labelFields are the record fields, in priority order, used to describe an id in
// completion output (e.g. an invoice's number, a contact's name).
var labelFields = []string{"name", "fullNumber", "number", "idNumber", "email", "description", "reference", "status", "total"}

func completionLabel(r map[string]any) string {
	for _, f := range labelFields {
		if v := scalar(r[f]); v != "" {
			return v
		}
	}
	return ""
}

// scalar renders a JSON-decoded scalar as a string, ignoring nested structures.
func scalar(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case bool:
		if t {
			return "true"
		}
		return "false"
	case float64:
		// Whole numbers render without a trailing ".0" (ids, counts).
		if t == float64(int64(t)) {
			return fmt.Sprintf("%d", int64(t))
		}
		return fmt.Sprintf("%g", t)
	default:
		return ""
	}
}

// fixedCompleter offers a static set of values, filtered by the typed prefix.
func fixedCompleter(values ...string) completeFunc {
	return func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		out := make([]string, 0, len(values))
		for _, v := range values {
			if strings.HasPrefix(v, toComplete) {
				out = append(out, v)
			}
		}
		return out, cobra.ShellCompDirectiveNoFileComp
	}
}

// columnsAnnotation stores a command's known columns (comma-joined) so the single
// global --columns completer can offer resource-specific fields. --columns is a
// persistent root flag shared by every command, so per-command completers can't
// be registered against it (Cobra keys completers by the one shared flag); the
// annotation carries the per-resource context instead.
const columnsAnnotation = "alegra/columns"

// withColumns records a resource's columns on cmd for completeColumns to read.
func withColumns(cmd *cobra.Command, cols []string) {
	if len(cols) == 0 {
		return
	}
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[columnsAnnotation] = strings.Join(cols, ",")
}

// completeColumns completes a comma-separated --columns list one field at a time,
// preserving already-typed segments and suppressing the trailing space so the
// user can keep appending fields. Candidates come from the command's annotation.
func completeColumns(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cols := strings.Split(cmd.Annotations[columnsAnnotation], ",")
	prefix, cur := "", toComplete
	if i := strings.LastIndex(toComplete, ","); i >= 0 {
		prefix, cur = toComplete[:i+1], toComplete[i+1:]
	}
	out := make([]string, 0, len(cols))
	for _, c := range cols {
		if c != "" && strings.HasPrefix(c, cur) {
			out = append(out, prefix+c)
		}
	}
	return out, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

// completeProfiles offers configured profile names, annotated with their email.
func completeProfiles(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cfg, err := config.Load()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)
	out := make([]string, 0, len(names))
	for _, name := range names {
		if !strings.HasPrefix(name, toComplete) {
			continue
		}
		if email := cfg.Profiles[name].Email; email != "" {
			out = append(out, name+"\t"+email)
		} else {
			out = append(out, name)
		}
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

// filterEnum returns the allowed values for a list filter when they can be
// determined: an explicit Values set, or a Usage string that is itself a compact
// enum like "open,closed,draft,void" (comma-separated, no prose/spaces).
func filterEnum(f listFilter) []string {
	if len(f.Values) > 0 {
		return f.Values
	}
	u := strings.TrimSpace(f.Usage)
	if u == "" || !strings.Contains(u, ",") || strings.ContainsAny(u, " \t") {
		return nil
	}
	var vals []string
	for _, p := range strings.Split(u, ",") {
		if p = strings.TrimSpace(p); p != "" {
			vals = append(vals, p)
		}
	}
	return vals
}
