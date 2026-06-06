package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

// listFilter declares a resource-specific list filter mapped to an API query
// parameter (e.g. CLI --status -> ?status=).
type listFilter struct {
	Flag  string
	Query string
	Usage string
}

// resourceSpec declares how to expose one CRUD resource. Concrete resource
// files build one of these and call registerResource in their init().
//
// Only Use, Short, and New are required. Everything else is optional.
type resourceSpec[T any] struct {
	Use     string   // command name, e.g. "contacts"
	Aliases []string // alternate names, e.g. ["contact"]
	Short   string   // one-line description
	Long    string   // optional long description

	// New returns the typed resource handle for a client.
	New func(*api.Client) *api.Resource[T]

	// Columns is the default table/csv column set (JSON keys) for list/get.
	Columns []string

	// OrderFields documents valid --order-field values (shown in help).
	OrderFields []string

	// ListFilters adds resource-specific list query filters.
	ListFilters []listFilter

	// NoCreate/NoUpdate/NoDelete suppress those subcommands for read-only or
	// restricted resources.
	NoCreate bool
	NoUpdate bool
	NoDelete bool

	// Extra adds custom subcommands (void, email, stamp, ...).
	Extra func(parent *cobra.Command, sp resourceSpec[T])
}

// registerResource builds the resource command tree and attaches it to root.
func registerResource[T any](sp resourceSpec[T]) {
	rootCmd.AddCommand(buildResourceCmd(sp))
}

func buildResourceCmd[T any](sp resourceSpec[T]) *cobra.Command {
	cmd := &cobra.Command{
		Use:     sp.Use,
		Aliases: sp.Aliases,
		Short:   sp.Short,
		Long:    sp.Long,
	}
	cmd.AddCommand(newListCmd(sp))
	cmd.AddCommand(newGetCmd(sp))
	if !sp.NoCreate {
		cmd.AddCommand(newCreateCmd(sp))
	}
	if !sp.NoUpdate {
		cmd.AddCommand(newUpdateCmd(sp))
	}
	if !sp.NoDelete {
		cmd.AddCommand(newDeleteCmd(sp))
	}
	if sp.Extra != nil {
		sp.Extra(cmd, sp)
	}
	return cmd
}

// --- list ---

// reservedListFlags are the built-in flag names on every `list` subcommand;
// resource-specific filters that collide with these are skipped.
var reservedListFlags = map[string]bool{
	"start": true, "limit": true, "all": true, "query": true,
	"order-field": true, "order-direction": true,
}

type listFlags struct {
	start    int
	limit    int
	all      bool
	query    string
	orderBy  string
	orderDir string
}

func newListCmd[T any](sp resourceSpec[T]) *cobra.Command {
	var lf listFlags
	filterVals := map[string]*string{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List " + sp.Use,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			params := api.ListParams{
				Start:          lf.start,
				Limit:          lf.limit,
				Query:          lf.query,
				OrderField:     lf.orderBy,
				OrderDirection: lf.orderDir,
				Extra:          url.Values{},
			}
			for _, f := range sp.ListFilters {
				if v := filterVals[f.Query]; v != nil && *v != "" {
					params.Extra.Set(f.Query, *v)
				}
			}
			res := sp.New(client)
			var items []T
			if lf.all {
				items, err = res.ListAll(cmd.Context(), params, 0)
			} else {
				items, err = res.List(cmd.Context(), params)
			}
			if err != nil {
				return err
			}
			return render(cmd, items, sp.Columns)
		},
	}

	fs := cmd.Flags()
	fs.IntVar(&lf.start, "start", 0, "Offset to start from (pagination)")
	fs.IntVar(&lf.limit, "limit", 0, "Max records per page (max 30)")
	fs.BoolVar(&lf.all, "all", false, "Fetch all pages")
	fs.StringVarP(&lf.query, "query", "q", "", "Free-text search")
	orderUsage := "Field to sort by"
	if len(sp.OrderFields) > 0 {
		orderUsage += " (" + strings.Join(sp.OrderFields, ", ") + ")"
	}
	fs.StringVar(&lf.orderBy, "order-field", "", orderUsage)
	fs.StringVar(&lf.orderDir, "order-direction", "", "Sort direction: ASC or DESC")
	// Register resource-specific filters, skipping any that would collide with a
	// built-in list flag or a previously declared filter (defensive: a bad
	// resource definition must never panic the whole CLI at init).
	for _, f := range sp.ListFilters {
		if f.Flag == "" || f.Query == "" || reservedListFlags[f.Flag] || fs.Lookup(f.Flag) != nil {
			continue
		}
		filterVals[f.Query] = fs.String(f.Flag, "", f.Usage)
	}
	return cmd
}

// --- get ---

func newGetCmd[T any](sp resourceSpec[T]) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a single " + singular(sp.Use) + " by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			item, err := sp.New(client).Get(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return render(cmd, item, sp.Columns)
		},
	}
}

// --- create ---

func newCreateCmd[T any](sp resourceSpec[T]) *cobra.Command {
	var bf bodyFlags
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a " + singular(sp.Use),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := bf.build()
			if err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			item, err := sp.New(client).Create(cmd.Context(), body)
			if err != nil {
				return err
			}
			return render(cmd, item, sp.Columns)
		},
	}
	addBodyFlags(cmd, &bf)
	return cmd
}

// --- update ---

func newUpdateCmd[T any](sp resourceSpec[T]) *cobra.Command {
	var bf bodyFlags
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a " + singular(sp.Use) + " by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := bf.build()
			if err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			item, err := sp.New(client).Update(cmd.Context(), args[0], body)
			if err != nil {
				return err
			}
			return render(cmd, item, sp.Columns)
		},
	}
	addBodyFlags(cmd, &bf)
	return cmd
}

// --- delete ---

func newDeleteCmd[T any](sp resourceSpec[T]) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a " + singular(sp.Use) + " by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes && !confirm(cmd, fmt.Sprintf("Delete %s %s?", singular(sp.Use), args[0])) {
				return fmt.Errorf("aborted")
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			if err := sp.New(client).Delete(cmd.Context(), args[0]); err != nil {
				return err
			}
			if !flagDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted %s %s\n", singular(sp.Use), args[0])
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")
	return cmd
}

// --- custom actions ---

// NewActionCmd builds a subcommand that POSTs to /<resource>/<id>/<apiAction>.
func NewActionCmd[T any](sp resourceSpec[T], use, apiAction, short string) *cobra.Command {
	var bf bodyFlags
	cmd := &cobra.Command{
		Use:   use + " <id>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, _ := bf.buildOptional()
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			var out map[string]any
			if err := sp.New(client).Action(cmd.Context(), args[0], apiAction, body, &out); err != nil {
				return err
			}
			if len(out) == 0 {
				if !flagDryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "OK: %s %s\n", apiAction, args[0])
				}
				return nil
			}
			return render(cmd, out, nil)
		},
	}
	addBodyFlags(cmd, &bf)
	return cmd
}

// NewCollectionActionCmd builds a subcommand that POSTs to /<resource>/<apiAction>.
func NewCollectionActionCmd[T any](sp resourceSpec[T], use, apiAction, short string) *cobra.Command {
	var bf bodyFlags
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := bf.build()
			if err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			var out map[string]any
			if err := sp.New(client).CollectionAction(cmd.Context(), apiAction, body, &out); err != nil {
				return err
			}
			return render(cmd, out, nil)
		},
	}
	addBodyFlags(cmd, &bf)
	return cmd
}

// --- request body flags ---

type bodyFlags struct {
	data string
	file string
	sets []string
}

func addBodyFlags(cmd *cobra.Command, bf *bodyFlags) {
	fs := cmd.Flags()
	fs.StringVarP(&bf.data, "data", "d", "", "Request body as a JSON string")
	fs.StringVarP(&bf.file, "file", "f", "", "Read JSON request body from a file (use - for stdin)")
	fs.StringArrayVar(&bf.sets, "set", nil, "Set a top-level field: key=value (value parsed as JSON when valid). Repeatable.")
}

// build returns the JSON body, erroring if nothing was provided.
func (bf bodyFlags) build() (json.RawMessage, error) {
	raw, err := bf.buildOptional()
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, fmt.Errorf("no request body: provide --data, --file, or one or more --set key=value")
	}
	return raw, nil
}

// buildOptional returns the JSON body or nil if none was provided.
func (bf bodyFlags) buildOptional() (json.RawMessage, error) {
	base := map[string]any{}
	provided := false

	switch {
	case bf.file != "":
		provided = true
		var data []byte
		var err error
		if bf.file == "-" {
			data, err = io.ReadAll(os.Stdin)
		} else {
			data, err = os.ReadFile(bf.file) //nolint:gosec // user-specified path
		}
		if err != nil {
			return nil, fmt.Errorf("reading body file: %w", err)
		}
		if err := json.Unmarshal(data, &base); err != nil {
			return nil, fmt.Errorf("parsing JSON body from %s: %w", bf.file, err)
		}
	case bf.data != "":
		provided = true
		if err := json.Unmarshal([]byte(bf.data), &base); err != nil {
			return nil, fmt.Errorf("parsing --data JSON: %w", err)
		}
	}

	for _, kv := range bf.sets {
		provided = true
		k, v, ok := strings.Cut(kv, "=")
		if !ok {
			return nil, fmt.Errorf("invalid --set %q (expected key=value)", kv)
		}
		base[k] = inferValue(v)
	}

	if !provided {
		return nil, nil
	}
	return json.Marshal(base)
}

// inferValue parses v as JSON when possible (number, bool, object, array,
// null), otherwise treats it as a string.
func inferValue(v string) any {
	if v == "" {
		return ""
	}
	switch v {
	case "true":
		return true
	case "false":
		return false
	case "null":
		return nil
	}
	if n, err := strconv.ParseFloat(v, 64); err == nil {
		return n
	}
	if c := v[0]; c == '{' || c == '[' || c == '"' {
		var parsed any
		if err := json.Unmarshal([]byte(v), &parsed); err == nil {
			return parsed
		}
	}
	return v
}

// --- helpers ---

func confirm(cmd *cobra.Command, prompt string) bool {
	fmt.Fprintf(cmd.OutOrStdout(), "%s [y/N]: ", prompt)
	reader := bufio.NewReader(cmd.InOrStdin())
	line, err := reader.ReadString('\n')
	if err != nil && line == "" {
		return false
	}
	line = strings.ToLower(strings.TrimSpace(line))
	return line == "y" || line == "yes"
}

// singular best-effort singularizes a resource name for help text.
func singular(s string) string {
	switch {
	case strings.HasSuffix(s, "ies"):
		return s[:len(s)-3] + "y"
	case strings.HasSuffix(s, "ses"):
		return s[:len(s)-2]
	case strings.HasSuffix(s, "s"):
		return s[:len(s)-1]
	default:
		return s
	}
}
