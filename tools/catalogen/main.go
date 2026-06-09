// Command catalogen regenerates the embedded per-country reference catalogs
// (internal/catalog/data/*.json) from Alegra's published country parameter pages
// (https://developer.alegra.com/reference/<country>.md). These pages are the
// source for the official MCP's units / reference-enum tools; Alegra exposes no
// public REST endpoint for them, so the CLI embeds the documented catalog data.
//
// Run via `make catalog-sync`. The generated JSON is committed; the build embeds
// it with go:embed.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type entry struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type category struct {
	Key     string  `json:"key"`
	Title   string  `json:"title"`
	Entries []entry `json:"entries"`
}

type catalog struct {
	Country    string     `json:"country"`
	Label      string     `json:"label"`
	Source     string     `json:"source"`
	Categories []category `json:"categories"`
}

type source struct {
	key   string // normalized id used by the CLI (lowercased applicationVersion)
	label string
	slug  string // developer.alegra.com reference slug
}

var sources = []source{
	{"colombia", "Colombia", "colombia"},
	{"mexico", "México", "méxico"},
	{"costarica", "Costa Rica", "costa-rica"},
	{"peru", "Perú", "perú"},
	{"spain", "España", "españa"},
	{"panama", "Panamá", "panamá"},
}

const base = "https://developer.alegra.com/reference/"

var (
	headingRe = regexp.MustCompile(`^##\s+(.+?)\s*$`)
	parenRe   = regexp.MustCompile(`\s*\(.*?\)\s*`)
	sepRe     = regexp.MustCompile(`^[\s:|-]+$`)
)

func main() {
	outDir := "internal/catalog/data"
	if len(os.Args) > 1 {
		outDir = os.Args[1]
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fatal(err)
	}
	client := &http.Client{Timeout: 60 * time.Second}
	for _, s := range sources {
		body, err := fetch(client, base+url.PathEscape(s.slug)+".md")
		if err != nil {
			fatal(fmt.Errorf("%s: %w", s.slug, err))
		}
		cats := parse(body)
		cat := catalog{Country: s.key, Label: s.label, Source: base + s.slug + ".md", Categories: cats}
		data, err := json.MarshalIndent(cat, "", "  ")
		if err != nil {
			fatal(err)
		}
		path := filepath.Join(outDir, s.key+".json")
		if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
			fatal(err)
		}
		total := 0
		for _, c := range cats {
			total += len(c.Entries)
		}
		fmt.Printf("%-10s %2d categories, %5d entries -> %s\n", s.key, len(cats), total, path)
	}
}

// parse extracts each "## Section" heading and the first markdown table that
// follows it into a category of {code, name} entries.
func parse(md string) []category {
	lines := strings.Split(md, "\n")
	var cats []category
	seen := map[string]bool{}
	for i := 0; i < len(lines); i++ {
		m := headingRe.FindStringSubmatch(lines[i])
		if m == nil {
			continue
		}
		title := strings.TrimSpace(strings.ReplaceAll(m[1], `\`, ""))
		// Advance to the first table row under this heading (skip prose).
		j := i + 1
		for j < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[j]), "|") {
			if headingRe.MatchString(lines[j]) {
				break // next section with no table
			}
			j++
		}
		if j >= len(lines) || !strings.HasPrefix(strings.TrimSpace(lines[j]), "|") {
			continue
		}
		entries := readTable(lines, j)
		if len(entries) == 0 {
			continue
		}
		key := slugify(title)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		cats = append(cats, category{Key: key, Title: title, Entries: entries})
		i = j
	}
	return cats
}

// readTable reads a contiguous markdown table starting at lines[start], skipping
// the header and separator rows, returning {code=col0, name=col1} per data row.
func readTable(lines []string, start int) []entry {
	var entries []entry
	row := 0
	for k := start; k < len(lines); k++ {
		t := strings.TrimSpace(lines[k])
		if !strings.HasPrefix(t, "|") {
			break
		}
		row++
		if row == 1 || sepRe.MatchString(t) { // header or |---|--- separator
			continue
		}
		cells := splitRow(t)
		if len(cells) < 2 || cells[0] == "" {
			continue
		}
		// Skip multi-value / geo cells (e.g. a province row whose "name" is an
		// HTML <br />-joined list of every city): not a flat code→name enum, and
		// they bloat the embed. Keeps the catalog to clean reference enums.
		if strings.Contains(cells[1], "<br") || len(cells[1]) > 200 || len(cells[0]) > 120 {
			continue
		}
		entries = append(entries, entry{Code: cells[0], Name: cells[1]})
	}
	return entries
}

func splitRow(t string) []string {
	t = strings.Trim(t, "|")
	parts := strings.Split(t, "|")
	out := make([]string, len(parts))
	for i, p := range parts {
		// Drop markdown escaping (e.g. NOT\_DOMICILIED -> NOT_DOMICILIED).
		out[i] = strings.TrimSpace(strings.ReplaceAll(p, `\`, ""))
	}
	return out
}

func slugify(s string) string {
	s = parenRe.ReplaceAllString(s, " ")
	s = strings.ToLower(strings.TrimSpace(s))
	repl := strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ñ", "n", "ü", "u")
	s = repl.Replace(s)
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		default:
			if !prevDash {
				b.WriteRune('-')
				prevDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

func fetch(c *http.Client, u string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "alegra-cli-catalogen")
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	return string(b), err
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "catalogen:", err)
	os.Exit(1)
}
