// Package output renders API results in json, yaml, csv, or table form. It
// works uniformly across every resource type by normalizing values through JSON
// into maps, so no per-resource formatting code is required.
package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

// Format is an output format identifier.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
	FormatCSV   Format = "csv"
)

// Valid reports whether f is a supported format.
func (f Format) Valid() bool {
	switch f {
	case FormatTable, FormatJSON, FormatYAML, FormatCSV:
		return true
	}
	return false
}

// preferredOrder lists columns shown first (when present) for readable tables.
var preferredOrder = []string{
	"id", "idNumber", "number", "fullNumber", "numberTemplate",
	"name", "identification", "email", "date", "dueDate",
	"status", "type", "total", "balance", "price", "quantity",
}

// Render writes data to w in the requested format. columns optionally fixes the
// table/csv column set and order (by JSON key); when empty, columns are derived.
func Render(w io.Writer, data any, format Format, columns []string) error {
	switch format {
	case FormatJSON:
		return renderJSON(w, data)
	case FormatYAML:
		return renderYAML(w, data)
	case FormatCSV:
		return renderCSV(w, data, columns)
	case FormatTable, "":
		return renderTable(w, data, columns)
	default:
		return fmt.Errorf("unsupported output format %q (use table, json, yaml, or csv)", format)
	}
}

func renderJSON(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(data)
}

func renderYAML(w io.Writer, data any) error {
	// Round-trip through JSON so json tags drive key names, then YAML-encode.
	normalized, err := toGeneric(data)
	if err != nil {
		return err
	}
	// toGeneric decodes numbers as json.Number (a string type); yaml.v3 would
	// emit those quoted ("total: \"7\""), unlike the JSON renderer. Convert them
	// to native int64/float64 so YAML numbers stay numbers (M3).
	normalized = nativeNumbers(normalized)
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	defer enc.Close()
	return enc.Encode(normalized)
}

// nativeNumbers walks a generic value and replaces every json.Number with a
// native int64 (when integral) or float64, so YAML/serializers that special-case
// json.Number as a string don't quote numeric output.
func nativeNumbers(v any) any {
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			t[k] = nativeNumbers(val)
		}
		return t
	case []any:
		for i, val := range t {
			t[i] = nativeNumbers(val)
		}
		return t
	case json.Number:
		s := t.String()
		if i, err := t.Int64(); err == nil {
			return i
		}
		// A big integer that overflows int64 must NOT become a lossy float64 —
		// that would silently corrupt a large Money/ledger total, the precision
		// hazard the Money type exists to avoid. Keep its exact text; only
		// genuinely fractional/scientific values fall through to float64.
		if !strings.ContainsAny(s, ".eE") {
			return s
		}
		if f, err := t.Float64(); err == nil {
			return f
		}
		return s
	default:
		return v
	}
}

func renderCSV(w io.Writer, data any, columns []string) error {
	rows := toRows(data)
	if len(rows) == 0 {
		// A single record (e.g. `get`) is an object, not an array; emit it as one
		// CSV row rather than nothing, matching renderTable's single-object path.
		if obj, ok := toObject(data); ok {
			rows = []map[string]any{obj}
		}
	}
	if len(rows) == 0 {
		return nil
	}
	cols, dropped := resolveColumns(rows, columns)
	warnDroppedColumns(dropped)
	cw := csv.NewWriter(w)
	defer cw.Flush()
	if err := cw.Write(cols); err != nil {
		return err
	}
	for _, row := range rows {
		rec := make([]string, len(cols))
		for i, c := range cols {
			rec[i] = csvField(scalarString(row[c]))
		}
		if err := cw.Write(rec); err != nil {
			return err
		}
	}
	return cw.Error()
}

// csvField neutralizes spreadsheet formula injection (CWE-1236). A cell whose
// first significant character is =, +, @, or - is interpreted as a formula by
// Excel/LibreOffice/Sheets, so we prefix it with a single quote to force text.
// A leading '-' is only dangerous when the value is not a real, finite negative
// number, so numeric cells (e.g. "-42.5") pass through untouched while
// "-Infinity"/"-NaN" — which ParseFloat accepts but a sheet shouldn't evaluate —
// are escaped. The trigger can hide behind leading whitespace/newlines, so we
// look past those, and a cell that itself begins with a control character is
// escaped too.
func csvField(s string) string {
	if s == "" {
		return s
	}
	trimmed := strings.TrimLeft(s, " \t\r\n")
	if trimmed == "" {
		return s
	}
	guard := false
	switch trimmed[0] {
	case '=', '+', '@':
		guard = true
	case '-':
		if f, err := strconv.ParseFloat(trimmed, 64); err != nil || math.IsInf(f, 0) || math.IsNaN(f) {
			guard = true
		}
	}
	if c := s[0]; c == '\t' || c == '\r' || c == '\n' {
		guard = true
	}
	if guard {
		return "'" + s
	}
	return s
}

func renderTable(w io.Writer, data any, columns []string) error {
	rows := toRows(data)
	if len(rows) == 0 {
		// Single object? Render as key/value.
		if obj, ok := toObject(data); ok {
			return renderKeyValue(w, obj, columns)
		}
		_, err := fmt.Fprintln(w, "(no results)")
		return err
	}
	if len(rows) == 1 {
		return renderKeyValue(w, rows[0], columns)
	}
	cols, dropped := resolveColumns(rows, columns)
	warnDroppedColumns(dropped)
	tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, strings.Join(upper(cols), "\t"))
	clipped := false
	for _, row := range rows {
		cells := make([]string, len(cols))
		for i, c := range cols {
			full := scalarString(row[c])
			if len([]rune(full)) > 48 {
				clipped = true
			}
			cells[i] = truncate(full, 48)
		}
		fmt.Fprintln(tw, strings.Join(cells, "\t"))
	}
	if err := tw.Flush(); err != nil {
		return err
	}
	// The multi-row table clips wide cells for readability; a single detail view
	// shows them in full. Flag the difference on stderr so a user who sees "…"
	// knows where to get the untruncated value (stdout stays clean for pipes).
	if clipped {
		fmt.Fprintln(os.Stderr, "note: some cells were truncated for display; use --output json (or csv) for full values")
	}
	return nil
}

func renderKeyValue(w io.Writer, obj map[string]any, columns []string) error {
	tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
	keys := columns
	if len(keys) == 0 {
		keys = orderedKeys(obj, true)
	}
	for _, k := range keys {
		v, ok := obj[k]
		if !ok {
			continue
		}
		fmt.Fprintf(tw, "%s\t%s\n", upperOne(k), valueString(v))
	}
	return tw.Flush()
}

// --- normalization helpers ---

func toGeneric(data any) (any, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var out any
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.UseNumber()
	if err := dec.Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func toRows(data any) []map[string]any {
	g, err := toGeneric(data)
	if err != nil {
		return nil
	}
	arr, ok := g.([]any)
	if !ok {
		return nil
	}
	rows := make([]map[string]any, 0, len(arr))
	for _, item := range arr {
		if m, ok := item.(map[string]any); ok {
			rows = append(rows, m)
		}
	}
	return rows
}

func toObject(data any) (map[string]any, bool) {
	g, err := toGeneric(data)
	if err != nil {
		return nil, false
	}
	m, ok := g.(map[string]any)
	return m, ok
}

// maxAutoColumns caps the auto-detected column set so wide resources stay
// readable; explicit --columns selections are never capped.
const maxAutoColumns = 10

// resolveColumns returns the explicit columns verbatim — honoring the caller's
// exact selection and order, including keys that may be absent from some rows —
// or, when none are given, derives a scalar-only column set ordered by
// preference, plus how many detected columns were dropped by the cap.
func resolveColumns(rows []map[string]any, columns []string) ([]string, int) {
	if len(columns) > 0 {
		return columns, 0
	}
	seen := map[string]bool{}
	for _, row := range rows {
		for k, v := range row {
			if isScalar(v) {
				seen[k] = true
			}
		}
	}
	cols := orderKeys(seen)
	dropped := 0
	if len(cols) > maxAutoColumns {
		dropped = len(cols) - maxAutoColumns
		cols = cols[:maxAutoColumns]
	}
	return cols, dropped
}

// warnDroppedColumns surfaces the auto-detect cap on stderr; stdout must stay
// clean for piping (a note inside CSV/table data would corrupt it).
func warnDroppedColumns(dropped int) {
	if dropped > 0 {
		fmt.Fprintf(os.Stderr, "note: %d more column(s) detected but not shown; use --columns to choose, or --output json for everything\n", dropped)
	}
}

func orderedKeys(obj map[string]any, includeAll bool) []string {
	seen := map[string]bool{}
	for k, v := range obj {
		if includeAll || isScalar(v) {
			seen[k] = true
		}
	}
	return orderKeys(seen)
}

func orderKeys(seen map[string]bool) []string {
	var out []string
	for _, k := range preferredOrder {
		if seen[k] {
			out = append(out, k)
			delete(seen, k)
		}
	}
	var rest []string
	for k := range seen {
		rest = append(rest, k)
	}
	sort.Strings(rest)
	return append(out, rest...)
}

func isScalar(v any) bool {
	switch v.(type) {
	case nil, bool, string, float64, json.Number:
		return true
	default:
		return false
	}
}

func scalarString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case bool:
		return strconv.FormatBool(t)
	case json.Number:
		return t.String()
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		return ""
	}
}

// valueString renders any value, compacting nested structures to JSON.
func valueString(v any) string {
	if isScalar(v) {
		return scalarString(v)
	}
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return truncate(string(b), 80)
}

// truncate shortens s to at most n display characters, appending an ellipsis.
// It counts and slices by rune, not byte, so it never splits a multi-byte UTF-8
// character (common in accented Spanish data) into invalid output.
func truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n == 1 {
		return string(r[:1])
	}
	return string(r[:n-1]) + "…"
}

func upper(cols []string) []string {
	out := make([]string, len(cols))
	for i, c := range cols {
		out[i] = upperOne(c)
	}
	return out
}

func upperOne(c string) string {
	return strings.ToUpper(c)
}
