// Package output renders API results in json, yaml, csv, or table form. It
// works uniformly across every resource type by normalizing values through JSON
// into maps, so no per-resource formatting code is required.
package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
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
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	defer enc.Close()
	return enc.Encode(normalized)
}

func renderCSV(w io.Writer, data any, columns []string) error {
	rows := toRows(data)
	if len(rows) == 0 {
		return nil
	}
	cols := resolveColumns(rows, columns)
	cw := csv.NewWriter(w)
	defer cw.Flush()
	if err := cw.Write(cols); err != nil {
		return err
	}
	for _, row := range rows {
		rec := make([]string, len(cols))
		for i, c := range cols {
			rec[i] = scalarString(row[c])
		}
		if err := cw.Write(rec); err != nil {
			return err
		}
	}
	return cw.Error()
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
	cols := resolveColumns(rows, columns)
	tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, strings.Join(upper(cols), "\t"))
	for _, row := range rows {
		cells := make([]string, len(cols))
		for i, c := range cols {
			cells[i] = truncate(scalarString(row[c]), 48)
		}
		fmt.Fprintln(tw, strings.Join(cells, "\t"))
	}
	return tw.Flush()
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

// resolveColumns returns the explicit columns (filtered to ones present) or
// derives a scalar-only column set ordered by preference.
func resolveColumns(rows []map[string]any, columns []string) []string {
	if len(columns) > 0 {
		return columns
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
	if len(cols) > 10 {
		cols = cols[:10]
	}
	return cols
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

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 1 {
		return s[:n]
	}
	return s[:n-1] + "…"
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
