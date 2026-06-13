package commands

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/njayp/ophis"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
	"github.com/jjuanrivvera/alegra-cli/internal/config"
	"github.com/jjuanrivvera/alegra-cli/internal/output"
)

// listFilter declares a resource-specific list filter mapped to an API query
// parameter (e.g. CLI --status -> ?status=).
type listFilter struct {
	Flag  string
	Query string
	Usage string
	// Values optionally enumerates the allowed values for shell completion. When
	// empty, a compact comma-separated Usage (e.g. "open,closed,void") is treated
	// as the value set automatically.
	Values []string
}

// dateRangeFilters returns the standard Alegra date/dueDate range filters shared
// by transactional documents (invoices, bills, payments, ...). Append with
// `append(dateRangeFilters(), ...)` in a resource's ListFilters.
func dateRangeFilters() []listFilter {
	return []listFilter{
		{Flag: "date-after", Query: "date_after", Usage: "On/after this date (YYYY-MM-DD)"},
		{Flag: "date-before", Query: "date_before", Usage: "On/before this date (YYYY-MM-DD)"},
		{Flag: "due-after", Query: "dueDate_after", Usage: "Due on/after this date (YYYY-MM-DD)"},
		{Flag: "due-before", Query: "dueDate_before", Usage: "Due on/before this date (YYYY-MM-DD)"},
	}
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
	if sp.New != nil {
		sp := sp
		resourcePathFns = append(resourcePathFns, func() string { return sp.New(pathProbe).Path() })
		// Record the typed struct for the field-level contract test against
		// the documented schemas (testdata/spec/schemas.json).
		resourceContractFns = append(resourceContractFns, func() (string, reflect.Type) {
			return sp.New(pathProbe).Path(), reflect.TypeFor[T]()
		})
	}
}

// resourceContractFns lazily yields each registered resource's API collection
// path and Go struct type.
var resourceContractFns []func() (string, reflect.Type)

// pathProbe is a credential-free client used only to read a resource's
// collection path (NewResource stores the path; no request is made).
var pathProbe = api.New()

// resourcePathFns lazily yields each registered resource's API collection path.
var resourcePathFns []func() string

// RegisteredResourcePaths returns the API collection path of every registered
// resource (e.g. "contacts", "conciliations", "webhooks/subscriptions"). It is
// used by the spec-coverage test to compare the CLI against the documented API
// manifest (testdata/spec/endpoints.json).
func RegisteredResourcePaths() []string {
	out := make([]string, 0, len(resourcePathFns))
	for _, f := range resourcePathFns {
		out = append(out, f())
	}
	return out
}

func buildResourceCmd[T any](sp resourceSpec[T]) *cobra.Command {
	cmd := &cobra.Command{
		Use:     sp.Use,
		Aliases: sp.Aliases,
		Short:   sp.Short,
		Long:    sp.Long,
	}
	// Tag each subcommand with MCP tool annotations (read-only vs destructive).
	// `alegra mcp` (via ophis) exposes these on every tool, so MCP hosts that
	// honor them gate writes automatically — see docs/user-guide/agent-safety.md.
	cmd.AddCommand(readOnlyHints(newListCmd(sp)))
	cmd.AddCommand(readOnlyHints(newGetCmd(sp)))
	cmd.AddCommand(readOnlyHints(newExportCmd(sp)))
	if !sp.NoCreate {
		cmd.AddCommand(writeHints(newCreateCmd(sp)))
		cmd.AddCommand(writeHints(newImportCmd(sp)))
	}
	if !sp.NoUpdate {
		cmd.AddCommand(destructiveHints(newUpdateCmd(sp)))
	}
	if !sp.NoDelete {
		cmd.AddCommand(destructiveHints(newDeleteCmd(sp)))
	}
	if sp.Extra != nil {
		sp.Extra(cmd, sp)
	}
	// Custom actions (void, emit, stamp, …) registered by Extra mutate state.
	// Any that didn't annotate themselves are treated as destructive so a host
	// gates them by default; the rare read-only custom command can opt back in
	// by calling readOnlyHints in its resource file.
	for _, c := range cmd.Commands() {
		if !hasMCPHints(c) {
			destructiveHints(c)
		}
	}
	return cmd
}

// readOnlyHints / writeHints / destructiveHints merge MCP tool annotations onto
// a command (consumed by `alegra mcp` through ophis). They merge rather than
// overwrite so they never clobber other annotations (e.g. completion columns).
func readOnlyHints(c *cobra.Command) *cobra.Command {
	return mcpHints(c, ophis.AnnotationReadOnly, "true", ophis.AnnotationIdempotent, "true", ophis.AnnotationOpenWorld, "true")
}

func writeHints(c *cobra.Command) *cobra.Command {
	return mcpHints(c, ophis.AnnotationOpenWorld, "true")
}

func destructiveHints(c *cobra.Command) *cobra.Command {
	return mcpHints(c, ophis.AnnotationDestructive, "true", ophis.AnnotationOpenWorld, "true")
}

func mcpHints(c *cobra.Command, kv ...string) *cobra.Command {
	if c.Annotations == nil {
		c.Annotations = map[string]string{}
	}
	for i := 0; i+1 < len(kv); i += 2 {
		c.Annotations[kv[i]] = kv[i+1]
	}
	return c
}

// hasMCPHints reports whether a command already carries any MCP tool annotation.
func hasMCPHints(c *cobra.Command) bool {
	for _, k := range []string{ophis.AnnotationReadOnly, ophis.AnnotationDestructive, ophis.AnnotationOpenWorld} {
		if _, ok := c.Annotations[k]; ok {
			return true
		}
	}
	return false
}

// --- list ---

// droppedListFilters records resource filters that were skipped at registration
// because of an empty definition or a flag collision. The CLI must never panic
// at init over a bad resource definition, so the mistake is surfaced by a
// registry test (TestNoListFiltersAreSilentlyDropped) failing CI instead.
var droppedListFilters []string

// reservedListFlags are the built-in flag names on every `list` subcommand;
// resource-specific filters that collide with these are skipped.
var reservedListFlags = map[string]bool{
	"start": true, "limit": true, "all": true, "query": true,
	"order-field": true, "order-direction": true, "count": true, "param": true,
	"since": true, "until": true,
}

type listFlags struct {
	start    int
	limit    int
	all      bool
	count    bool
	query    string
	orderBy  string
	orderDir string
	params   []string
	since    string
	until    string
}

func newListCmd[T any](sp resourceSpec[T]) *cobra.Command {
	var lf listFlags
	filterVals := map[string]*string{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List " + sp.Use,
		Args:    cobra.NoArgs,
		Example: listExample(sp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := getAPIClient(cmd)
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
			// Arbitrary raw query params (escape hatch for any Alegra filter
			// not exposed as a typed flag), applied first so curated filters win.
			for _, p := range lf.params {
				k, v, ok := strings.Cut(p, "=")
				if !ok {
					return fmt.Errorf("invalid --param %q (expected key=value)", p)
				}
				params.Extra.Set(k, v)
			}
			// Natural date range (resolves to Alegra's date_after/date_before).
			now := time.Now().UTC()
			if lf.since != "" {
				d, derr := parseDateExpr(lf.since, now)
				if derr != nil {
					return fmt.Errorf("--since: %w", derr)
				}
				params.Extra.Set("date_after", d)
			}
			if lf.until != "" {
				d, derr := parseDateExpr(lf.until, now)
				if derr != nil {
					return fmt.Errorf("--until: %w", derr)
				}
				params.Extra.Set("date_before", d)
			}
			for _, f := range sp.ListFilters {
				if v := filterVals[f.Flag]; v != nil && *v != "" {
					params.Extra.Set(f.Query, *v)
				}
			}

			res := sp.New(client)
			if lf.count {
				total, err := res.Count(cmd.Context(), params)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), total)
				return nil
			}

			var items []T
			if lf.all {
				var truncated bool
				items, truncated, err = res.ListAllChecked(cmd.Context(), params, 0)
				if err == nil && truncated {
					fmt.Fprintln(cmd.ErrOrStderr(), "warning: result stopped at the page cap and may be truncated; narrow with filters or a date range and re-run.")
				}
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
	fs.BoolVar(&lf.count, "count", false, "Print only the total number of matching records")
	fs.StringVarP(&lf.query, "query", "q", "", "Free-text search")
	fs.StringVar(&lf.since, "since", "", "Start of date range (YYYY-MM-DD, today, this-month, last-month, 7d, 3m, ...)")
	fs.StringVar(&lf.until, "until", "", "End of date range (same formats as --since)")
	fs.StringArrayVar(&lf.params, "param", nil, "Arbitrary API query parameter: key=value (repeatable; e.g. --param date_after=2026-01-01)")
	orderUsage := "Field to sort by"
	if len(sp.OrderFields) > 0 {
		orderUsage += " (" + strings.Join(sp.OrderFields, ", ") + ")"
	}
	fs.StringVar(&lf.orderBy, "order-field", "", orderUsage)
	fs.StringVar(&lf.orderDir, "order-direction", "", "Sort direction: ASC or DESC")
	// Register resource-specific filters, skipping any that would collide with a
	// built-in list flag or a previously declared filter (defensive: a bad
	// resource definition must never panic the whole CLI at init). Skips are
	// recorded so the registry test fails CI instead of losing filters silently.
	for _, f := range sp.ListFilters {
		if f.Flag == "" || f.Query == "" {
			droppedListFilters = append(droppedListFilters, fmt.Sprintf("%s: filter %+v has an empty flag or query", sp.Use, f))
			continue
		}
		if reservedListFlags[f.Flag] || fs.Lookup(f.Flag) != nil {
			droppedListFilters = append(droppedListFilters, fmt.Sprintf("%s: --%s collides with a built-in or duplicate flag", sp.Use, f.Flag))
			continue
		}
		// Key by the (unique) flag, not the query: two filters sharing a query
		// would otherwise overwrite each other here and silently drop one (L9).
		filterVals[f.Flag] = fs.String(f.Flag, "", f.Usage)
		// Enum-valued filters (e.g. --status open|closed|void) complete their values.
		if vals := filterEnum(f); len(vals) > 0 {
			_ = cmd.RegisterFlagCompletionFunc(f.Flag, fixedCompleter(vals...))
		}
	}
	withColumns(cmd, sp.Columns)
	_ = cmd.RegisterFlagCompletionFunc("order-direction", fixedCompleter("ASC", "DESC"))
	if len(sp.OrderFields) > 0 {
		_ = cmd.RegisterFlagCompletionFunc("order-field", fixedCompleter(sp.OrderFields...))
	}
	return cmd
}

// listExample builds a baseline help example for a resource's list command.
func listExample[T any](sp resourceSpec[T]) string {
	ex := "  alegra " + sp.Use + " list\n" +
		"  alegra " + sp.Use + " list --limit 30 --all -o json\n" +
		"  alegra " + sp.Use + " list --count"
	if len(sp.ListFilters) > 0 {
		ex += "\n  alegra " + sp.Use + " list --" + sp.ListFilters[0].Flag + " <value>"
	}
	ex += "\n  alegra " + sp.Use + " list --param <api_param>=<value>"
	return ex
}

// --- export (CSV/JSON of all pages) ---

func newExportCmd[T any](sp resourceSpec[T]) *cobra.Command {
	var (
		outFile string
		format  string
		params  []string
		query   string
	)
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export all " + sp.Use + " to CSV or JSON",
		Long:  "Fetch every page of " + sp.Use + " and write them to a file (or stdout).",
		Example: "  alegra " + sp.Use + " export > " + sp.Use + ".csv\n" +
			"  alegra " + sp.Use + " export --format json --out " + sp.Use + ".json\n" +
			"  alegra " + sp.Use + " export --param status=open",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := getAPIClient(cmd)
			if err != nil {
				return err
			}
			p := api.ListParams{Query: query, Extra: url.Values{}}
			for _, kv := range params {
				k, v, ok := strings.Cut(kv, "=")
				if !ok {
					return fmt.Errorf("invalid --param %q (expected key=value)", kv)
				}
				p.Extra.Set(k, v)
			}
			items, truncated, err := sp.New(client).ListAllChecked(cmd.Context(), p, 0)
			if err != nil {
				return err
			}
			// In dry-run the client already printed the equivalent curl; never
			// create/truncate the output file or write an (empty) export. os.Create
			// would otherwise clobber an existing file even though no request ran.
			if flagDryRun {
				return nil
			}
			if truncated {
				fmt.Fprintln(cmd.ErrOrStderr(), "warning: export stopped at the page cap and may be incomplete; narrow with --param filters and re-run.")
			}
			w := cmd.OutOrStdout()
			if outFile != "" {
				f, ferr := os.Create(outFile) //nolint:gosec // user-specified path
				if ferr != nil {
					return ferr
				}
				defer f.Close()
				w = f
			}
			cols := flagColumns
			if len(cols) == 0 {
				cols = sp.Columns
			}
			if err := output.Render(w, items, output.Format(format), cols); err != nil {
				return err
			}
			if outFile != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Exported %d %s to %s\n", len(items), sp.Use, outFile)
			}
			return nil
		},
	}
	fs := cmd.Flags()
	fs.StringVar(&outFile, "out", "", "Write to this file (default: stdout)")
	fs.StringVar(&format, "format", "csv", "Export format: csv or json")
	fs.StringArrayVar(&params, "param", nil, "Filter by an API query parameter: key=value (repeatable)")
	fs.StringVarP(&query, "query", "q", "", "Free-text search")
	_ = cmd.RegisterFlagCompletionFunc("format", fixedCompleter("csv", "json"))
	withColumns(cmd, sp.Columns)
	return cmd
}

// --- import (create from CSV) ---

func newImportCmd[T any](sp resourceSpec[T]) *cobra.Command {
	var (
		file    string
		mapping []string
		sets    []string
	)
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Bulk-create " + sp.Use + " from a CSV file",
		Long: `Create one ` + singular(sp.Use) + ` per CSV row.

The header row names the fields; use --map to rename columns to API fields and
dotted paths for nested objects (e.g. --map 'NIT=identification.number'). Apply
constant fields to every row with --set. Rows are processed independently;
failures are reported and do not stop the run.`,
		Example: "  alegra " + sp.Use + " import --file rows.csv\n" +
			"  alegra contacts import -f clients.csv \\\n" +
			"    --map 'Name=name,NIT=identification.number' \\\n" +
			"    --set 'identification.type=NIT' --set 'type=[\"client\"]'",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if file == "" {
				return fmt.Errorf("--file is required")
			}
			colMap, err := parseMapping(mapping)
			if err != nil {
				return fmt.Errorf("--map: %w", err)
			}
			defaults, err := parseKeyVals(sets)
			if err != nil {
				return fmt.Errorf("--set: %w", err)
			}
			f, err := os.Open(file) //nolint:gosec // user-specified path
			if err != nil {
				return err
			}
			defer f.Close()

			records, err := csv.NewReader(f).ReadAll()
			if err != nil {
				return fmt.Errorf("reading CSV: %w", err)
			}
			if len(records) < 2 {
				return fmt.Errorf("CSV has no data rows")
			}
			header := records[0]

			client, err := getAPIClient(cmd)
			if err != nil {
				return err
			}
			res := sp.New(client)

			var created, failed int
			out := cmd.OutOrStdout()
			for i, row := range records[1:] {
				body := map[string]any{}
				for k, v := range defaults {
					setDotPath(body, k, inferValue(v))
				}
				for j, cell := range row {
					if j >= len(header) || cell == "" {
						continue
					}
					field := header[j]
					if mapped, ok := colMap[field]; ok {
						field = mapped
					}
					setDotPath(body, field, inferValue(cell))
				}
				if flagDryRun {
					raw, merr := json.Marshal(body)
					if merr != nil {
						failed++
						fmt.Fprintf(cmd.ErrOrStderr(), "[row %d] FAILED: cannot encode body: %v\n", i+1, merr)
						continue
					}
					fmt.Fprintf(out, "[row %d] would create: %s\n", i+1, raw)
					continue
				}
				if _, cerr := res.Create(cmd.Context(), body); cerr != nil {
					failed++
					fmt.Fprintf(cmd.ErrOrStderr(), "[row %d] FAILED: %v\n", i+1, cerr)
					continue
				}
				created++
				fmt.Fprintf(out, "[row %d] created\n", i+1)
			}
			if !flagDryRun {
				fmt.Fprintf(out, "Imported %d, failed %d\n", created, failed)
			}
			if failed > 0 {
				return fmt.Errorf("%d row(s) failed", failed)
			}
			return nil
		},
	}
	fs := cmd.Flags()
	fs.StringVarP(&file, "file", "f", "", "CSV file to import (required)")
	fs.StringArrayVar(&mapping, "map", nil, "Map a CSV column to a field path: column=field.path (repeatable)")
	fs.StringArrayVar(&sets, "set", nil, "Constant field applied to every row: key=value (repeatable)")
	return cmd
}

// setDotPath assigns val into a (possibly nested) map following a dotted path.
func setDotPath(root map[string]any, path string, val any) {
	parts := strings.Split(path, ".")
	m := root
	for i, p := range parts {
		if i == len(parts)-1 {
			m[p] = val
			return
		}
		next, ok := m[p].(map[string]any)
		if !ok {
			next = map[string]any{}
			m[p] = next
		}
		m = next
	}
}

// parseMapping parses --map entries, where each entry may hold several
// comma-separated `column=field` pairs (field paths never contain commas).
func parseMapping(entries []string) (map[string]string, error) {
	var pairs []string
	for _, e := range entries {
		pairs = append(pairs, strings.Split(e, ",")...)
	}
	return parseKeyVals(pairs)
}

// parseKeyVals turns ["a=b","c=d"] into a map, erroring on malformed entries.
func parseKeyVals(pairs []string) (map[string]string, error) {
	out := map[string]string{}
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			return nil, fmt.Errorf("invalid %q (expected key=value)", p)
		}
		out[k] = v
	}
	return out, nil
}

// --- get ---

func newGetCmd[T any](sp resourceSpec[T]) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "get <id>",
		Short:             "Get a single " + singular(sp.Use) + " by ID",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: resourceIDCompleter(sp),
		Example:           "  alegra " + sp.Use + " get <id>\n  alegra " + sp.Use + " get <id> -o json",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getAPIClient(cmd)
			if err != nil {
				return err
			}
			item, err := sp.New(client).Get(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			// Pass nil columns so a single record shows all of its fields, not
			// just the terse list columns.
			return render(cmd, item, nil)
		},
	}
	// --columns completes against this resource's known fields.
	withColumns(cmd, sp.Columns)
	return cmd
}

// --- create ---

func newCreateCmd[T any](sp resourceSpec[T]) *cobra.Command {
	var bf bodyFlags
	var country string
	var noValidate bool
	var draft bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a " + singular(sp.Use),
		Long: "Create a " + singular(sp.Use) + ".\n\n" +
			"Provide the body with --file <path> (recommended for nested documents),\n" +
			"--data '<json>', or one or more --set key=value pairs for flat fields.\n" +
			"The body is pre-flight validated for your country; use --no-validate to skip.",
		Example: "  alegra " + sp.Use + " create -f " + singular(sp.Use) + ".json\n" +
			"  alegra " + sp.Use + " create --set name=\"Example\"\n" +
			"  echo '{...}' | alegra " + sp.Use + " create -f -",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := bf.build()
			if err != nil {
				return err
			}
			// --draft: never emit electronically — drop any stamp instruction.
			if draft {
				if m, ok := bodyToMap(body); ok {
					delete(m, "stamp")
					body, err = json.Marshal(m)
					if err != nil {
						return fmt.Errorf("re-encoding body after stripping stamp: %w", err)
					}
				}
			}
			if !noValidate {
				if m, ok := bodyToMap(body); ok {
					if problems := validateForCreate(sp.Use, resolveCountry(country), m); len(problems) > 0 {
						return formatValidationError(sp.Use, resolveCountry(country), problems)
					}
				}
			}
			client, err := getAPIClient(cmd)
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
	cmd.Flags().StringVar(&country, "country", "", "Country for pre-flight validation (default: auto-detected from the account)")
	cmd.Flags().BoolVar(&noValidate, "no-validate", false, "Skip client-side pre-flight validation")
	cmd.Flags().BoolVar(&draft, "draft", false, "Create as a draft (strip any electronic stamp from the body)")
	return cmd
}

// resolveCountry picks the validation country. The platform is the source of
// truth, so precedence is: explicit --country flag > the account's auto-detected
// country (cached per profile by `auth login`/`doctor`) > the global offline
// hint (`config set-country`).
func resolveCountry(override string) string {
	if override != "" {
		return strings.ToLower(override)
	}
	cfg, err := config.Load()
	if err != nil {
		return ""
	}
	if p := cfg.Profile(cfg.ActiveProfileName(flagProfile)); p.Country != "" {
		return strings.ToLower(p.Country)
	}
	if cfg.Settings != nil {
		return strings.ToLower(cfg.Settings.Country)
	}
	return ""
}

// --- update ---

func newUpdateCmd[T any](sp resourceSpec[T]) *cobra.Command {
	var bf bodyFlags
	cmd := &cobra.Command{
		Use:               "update <id>",
		Short:             "Update a " + singular(sp.Use) + " by ID",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: resourceIDCompleter(sp),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := bf.build()
			if err != nil {
				return err
			}
			client, err := getAPIClient(cmd)
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
		Use:               "delete <id>",
		Short:             "Delete a " + singular(sp.Use) + " by ID",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: resourceIDCompleter(sp),
		RunE: func(cmd *cobra.Command, args []string) error {
			// In dry-run, skip the interactive confirmation: there is no real
			// deletion to confirm, and prompting would block (or abort on EOF) when
			// the user just wants to preview the request.
			if !yes && !flagDryRun && !confirm(cmd, fmt.Sprintf("Delete %s %s?", singular(sp.Use), args[0])) {
				return fmt.Errorf("aborted")
			}
			client, err := getAPIClient(cmd)
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
// Pass bodyRequired=true for actions whose API call needs a JSON body (e.g.
// email, comments, transfer) so the CLI fails fast with a clear message instead
// of letting the server return an opaque error.
func NewActionCmd[T any](sp resourceSpec[T], use, apiAction, short string, bodyRequired ...bool) *cobra.Command {
	var bf bodyFlags
	required := len(bodyRequired) > 0 && bodyRequired[0]
	cmd := &cobra.Command{
		Use:               use + " <id>",
		Short:             short,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: resourceIDCompleter(sp),
		RunE: func(cmd *cobra.Command, args []string) error {
			var body json.RawMessage
			var err error
			if required {
				body, err = bf.build()
			} else {
				body, err = bf.buildOptional()
			}
			if err != nil {
				return err
			}
			client, err := getAPIClient(cmd)
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

// NewPutActionCmd builds a subcommand that PUTs a body to
// /<resource>/<id>/<apiAction> (e.g. replacing a bill's perceptions or
// retentions). A request body is required.
func NewPutActionCmd[T any](sp resourceSpec[T], use, apiAction, short string) *cobra.Command {
	var bf bodyFlags
	cmd := &cobra.Command{
		Use:               use + " <id>",
		Short:             short,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: resourceIDCompleter(sp),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := bf.build()
			if err != nil {
				return err
			}
			client, err := getAPIClient(cmd)
			if err != nil {
				return err
			}
			path := sp.New(client).Path() + "/" + url.PathEscape(args[0]) + "/" + apiAction
			var out map[string]any
			if err := client.PutInto(cmd.Context(), path, body, &out); err != nil {
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
			client, err := getAPIClient(cmd)
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
	fs.StringArrayVar(&bf.sets, "set", nil, "Set a top-level field: key=value (value parsed as JSON when valid). Repeatable. For nested documents (e.g. invoice items[]) use --file.")
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
	if n, ok := numberIfLossless(v); ok {
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

// numberIfLossless returns v as a numeric value (int64 or float64) only when the
// number formats back to the original text exactly. This keeps identifiers that
// merely look numeric as strings — leading zeros ("007123456"), values past
// 2^53, a leading "+", underscores, "NaN"/"Inf", and scientific notation that
// would rewrite the digits — which is critical for an accounting tool: those are
// the shapes of NITs/RFCs, document numbers, and bank-account numbers that must
// reach the API byte-for-byte. (Finding H2.)
func numberIfLossless(v string) (any, bool) {
	if i, err := strconv.ParseInt(v, 10, 64); err == nil {
		if strconv.FormatInt(i, 10) == v {
			return i, true
		}
		return nil, false // e.g. "007" parses to 7 but isn't the same text
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return nil, false
		}
		if strconv.FormatFloat(f, 'g', -1, 64) == v {
			return f, true
		}
	}
	return nil, false
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
