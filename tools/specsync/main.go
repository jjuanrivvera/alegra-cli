// Command specsync reconstructs a manifest of Alegra's documented API surface
// from the official LLM index (https://developer.alegra.com/llms.txt) and writes
// it to testdata/spec/endpoints.json (committed) plus the raw index to
// .alegra-spec/ (gitignored).
//
// Why the index and not the pages: the /reference/*.md pages are JS-rendered
// (ReadMe.io), so their embedded OpenAPI is not present in the static HTML and
// can't be scraped without a browser. llms.txt is plain text and lists every
// reference page; its slugs encode method+path for the modern REST pages and
// tool names for the MCP pages. That makes a reproducible, network-light
// inventory used by `spec-check` to catch resource drift.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

const llmsURL = "https://developer.alegra.com/llms.txt"

// Manifest is the committed inventory of the documented API surface.
type Manifest struct {
	Source        string         `json:"source"`
	RESTEndpoints []RESTEndpoint `json:"restEndpoints"`
	MCPTools      []string       `json:"mcpTools"`
	// LegacySlugs are reference pages whose slug isn't a modern `method_path`
	// REST slug or an MCP `group__tool` slug (e.g. createAdditionalCharge,
	// listReconciliations). They still document real endpoints.
	LegacySlugs []string `json:"legacySlugs"`
	// Resources is the set of unique top-level resource tokens seen across the
	// REST slugs and MCP tools — the granularity spec-check matches CLI resources
	// against (paths can't always be derived authoritatively from a slug).
	Resources []string `json:"resources"`
}

// RESTEndpoint is a documented REST page: its HTTP method and a best-effort path
// derived from the page slug (modern `method_path` slugs only).
type RESTEndpoint struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Slug   string `json:"slug"`
}

var (
	linkRe   = regexp.MustCompile(`\(https://developer\.alegra\.com/reference/([^)]+)\.md\)`)
	methodRe = regexp.MustCompile(`^(get|post|put|patch|delete)_(.+)$`)
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "specsync:", err)
		os.Exit(1)
	}
}

func run() error {
	body, err := fetch(llmsURL)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(".alegra-spec", 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(".alegra-spec/llms.txt", body, 0o644); err != nil { //nolint:gosec // public docs
		return err
	}

	seenSlug := map[string]bool{}
	var rest []RESTEndpoint
	var mcp []string
	var legacy []string
	resources := map[string]bool{}

	for _, m := range linkRe.FindAllStringSubmatch(string(body), -1) {
		slug := m[1]
		if seenSlug[slug] {
			continue
		}
		seenSlug[slug] = true

		switch {
		case strings.Contains(slug, "__"): // MCP tool, e.g. sellers__getSellers
			mcp = append(mcp, slug)
			resources[strings.SplitN(slug, "__", 2)[0]] = true
		case methodRe.MatchString(slug): // modern REST, e.g. post_webhooks-subscriptions
			sm := methodRe.FindStringSubmatch(slug)
			path := "/" + strings.ReplaceAll(sm[2], "-id", "/{id}")
			rest = append(rest, RESTEndpoint{Method: strings.ToUpper(sm[1]), Path: path, Slug: slug})
			resources[firstSegment(path)] = true
		default: // legacy/catalog slug (e.g. listreconciliations, createAdditionalCharge)
			legacy = append(legacy, slug)
		}
	}

	sort.Slice(rest, func(i, j int) bool {
		if rest[i].Path != rest[j].Path {
			return rest[i].Path < rest[j].Path
		}
		return rest[i].Method < rest[j].Method
	})
	sort.Strings(mcp)
	sort.Strings(legacy)
	man := Manifest{Source: llmsURL, RESTEndpoints: rest, MCPTools: mcp, LegacySlugs: legacy, Resources: sortedKeys(resources)}

	if err := os.MkdirAll("testdata/spec", 0o755); err != nil {
		return err
	}
	out, err := json.MarshalIndent(man, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile("testdata/spec/endpoints.json", append(out, '\n'), 0o644); err != nil { //nolint:gosec // public docs
		return err
	}
	fmt.Printf("specsync: %d REST endpoints, %d MCP tools, %d resources → testdata/spec/endpoints.json\n",
		len(rest), len(mcp), len(man.Resources))
	return nil
}

func fetch(url string) ([]byte, error) {
	c := &http.Client{Timeout: 30 * time.Second}
	resp, err := c.Get(url) //nolint:gosec,noctx // fixed public URL
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %d", url, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func firstSegment(path string) string {
	first, _, _ := strings.Cut(strings.Trim(path, "/"), "/")
	return first
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		if k != "" {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}
