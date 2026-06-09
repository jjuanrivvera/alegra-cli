// Package catalog serves Alegra's per-country reference catalogs (units of
// measure and reference enums such as identification types, tax types, and
// payment methods) from data embedded at build time.
//
// Alegra exposes no public REST endpoint for these — the official MCP serves
// them from its own per-country dataset — so the CLI embeds the same data,
// generated from Alegra's published country parameter pages by tools/catalogen
// (run via `make catalog-sync`).
package catalog

import (
	"embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

//go:embed data/*.json
var dataFS embed.FS

// Entry is one reference value: a code and its human-readable name.
type Entry struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Category is a named group of reference entries (e.g. "Unidades de medida").
type Category struct {
	Key     string  `json:"key"`
	Title   string  `json:"title"`
	Entries []Entry `json:"entries"`
}

// Catalog is one country's full set of reference categories.
type Catalog struct {
	Country    string     `json:"country"`
	Label      string     `json:"label"`
	Source     string     `json:"source"`
	Categories []Category `json:"categories"`
}

// categoryAliases maps friendly English keys to the slugified Spanish category
// keys produced from the source pages.
var categoryAliases = map[string]string{
	"units":                "unidades-de-medida",
	"identification-types": "tipos-de-identificacion",
	"payment-methods":      "formas-de-pago",
	"invoice-types":        "tipos-de-factura",
	"regimes":              "regimenes",
}

// Normalize maps a country string (e.g. the account's lowercased
// applicationVersion, or a user-supplied --country) to a catalog key.
func Normalize(country string) string {
	c := strings.ToLower(strings.TrimSpace(country))
	c = strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ñ", "n", " ", "").Replace(c)
	switch c {
	case "co":
		return "colombia"
	case "mx", "mexico":
		return "mexico"
	case "cr", "costarica":
		return "costarica"
	case "pe", "peru":
		return "peru"
	case "es", "espana", "spain":
		return "spain"
	case "pa", "panama":
		return "panama"
	}
	return c
}

// Available returns the country keys that have an embedded catalog, sorted.
func Available() []string {
	ents, _ := dataFS.ReadDir("data")
	out := make([]string, 0, len(ents))
	for _, e := range ents {
		out = append(out, strings.TrimSuffix(e.Name(), ".json"))
	}
	sort.Strings(out)
	return out
}

// Load returns the catalog for a country, or an error naming the available
// countries when none matches.
func Load(country string) (*Catalog, error) {
	key := Normalize(country)
	if key == "" {
		return nil, fmt.Errorf("no country given: pass --country (%s)", strings.Join(Available(), ", "))
	}
	b, err := dataFS.ReadFile("data/" + key + ".json")
	if err != nil {
		return nil, fmt.Errorf("no reference catalog for country %q (available: %s)", country, strings.Join(Available(), ", "))
	}
	var c Catalog
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("catalog: decoding %s: %w", key, err)
	}
	return &c, nil
}

// Category returns one category by key or friendly alias.
func (c *Catalog) Category(name string) (*Category, bool) {
	want := strings.ToLower(strings.TrimSpace(name))
	if alias, ok := categoryAliases[want]; ok {
		want = alias
	}
	for i := range c.Categories {
		if c.Categories[i].Key == want {
			return &c.Categories[i], true
		}
	}
	return nil, false
}

// CategoryKeys returns the category keys available in this catalog, sorted.
func (c *Catalog) CategoryKeys() []string {
	out := make([]string, 0, len(c.Categories))
	for _, cat := range c.Categories {
		out = append(out, cat.Key)
	}
	sort.Strings(out)
	return out
}
