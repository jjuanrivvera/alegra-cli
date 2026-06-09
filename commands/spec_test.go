package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// specManifest mirrors testdata/spec/endpoints.json (see tools/specsync).
type specManifest struct {
	RESTEndpoints []struct {
		Slug string `json:"slug"`
	} `json:"restEndpoints"`
	MCPTools    []string `json:"mcpTools"`
	LegacySlugs []string `json:"legacySlugs"`
}

// specNorm folds a resource/slug to a comparable token: lowercase, separators
// removed, trailing plural dropped. This absorbs the doc surfaces' formatting
// differences (REST `bank-accounts` vs MCP `bankAccount`, `/`-vs-`-` in slugs).
func specNorm(s string) string {
	s = strings.ToLower(strings.NewReplacer("-", "", "_", "", "/", "").Replace(s))
	return strings.TrimSuffix(s, "s")
}

// TestSpecManifest_CLIResourcesAreDocumented is a network-free guardrail: every
// registered CLI resource must appear in the official API surface captured in
// testdata/spec/endpoints.json (REST slugs + MCP tools). It fails CI if a
// resource is undocumented (e.g. a typo'd path or a removed endpoint).
//
// Limitation: the docs are JS-rendered, so authoritative paths can't be scraped;
// matching is by normalized resource name, which confirms a resource is
// documented but won't catch a subtle path variant (that is covered by live
// verification in code and the weekly spec-sync drift job). Refresh with
// `make spec-sync`.
func TestSpecManifest_CLIResourcesAreDocumented(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "testdata", "spec", "endpoints.json"))
	require.NoError(t, err, "manifest missing — run `make spec-sync`")

	var man specManifest
	require.NoError(t, json.Unmarshal(data, &man))
	require.NotEmpty(t, man.RESTEndpoints, "manifest has no REST endpoints — run `make spec-sync`")

	var hay strings.Builder
	for _, e := range man.RESTEndpoints {
		hay.WriteString(specNorm(e.Slug))
		hay.WriteByte(' ')
	}
	for _, tool := range man.MCPTools {
		hay.WriteString(specNorm(tool))
		hay.WriteByte(' ')
	}
	for _, slug := range man.LegacySlugs {
		hay.WriteString(specNorm(slug))
		hay.WriteByte(' ')
	}
	documented := hay.String()

	for _, p := range RegisteredResourcePaths() {
		if !strings.Contains(documented, specNorm(p)) {
			t.Errorf("CLI resource %q is not documented in testdata/spec/endpoints.json; "+
				"verify it against the official API docs (or refresh with `make spec-sync`)", p)
		}
	}
}
