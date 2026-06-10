package catalog

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SAT product keys (México's c_ClaveProdServ, CFDI 4.0) are the one catalog the
// CLI does not embed: ~52k entries (~7MB source) would bloat the binary for a
// dataset only Mexican accounts need. They are government data published by the
// SAT, not by Alegra — Alegra exposes no REST endpoint for them and the official
// MCP's product-keys tool is unreachable — so the CLI syncs them on demand from
// phpcfdi/resources-sat-catalogs (Unlicense, actively maintained mirror of the
// SAT's catCFDI publication) into a local cache shared across profiles.
var (
	SATDataURL    = "https://raw.githubusercontent.com/phpcfdi/resources-sat-catalogs/master/database/data/cfdi_40_productos_servicios.sql"
	SATVersionURL = "https://raw.githubusercontent.com/phpcfdi/resources-sat-catalogs/master/database/version.txt"
)

// SATEntry is one product/service key from the SAT c_ClaveProdServ catalog.
type SATEntry struct {
	Code string `json:"code"`
	Name string `json:"name"`
	// Similar lists alternative names the SAT publishes for the key (search
	// matches against it; e.g. "Caballos" → "Equinos, Potras, Yeguas...").
	Similar string `json:"similar,omitempty"`
}

// SATCatalog is the cached SAT product-keys catalog plus its provenance.
type SATCatalog struct {
	Version   string     `json:"version"`
	Source    string     `json:"source"`
	FetchedAt time.Time  `json:"fetchedAt"`
	Entries   []SATEntry `json:"entries"`
}

// SATPath returns the cache file inside dir. The cache is shared across
// profiles: the catalog is global SAT data, not account data.
func SATPath(dir string) string {
	return filepath.Join(dir, "sat-product-keys.json.gz")
}

// SATCached reports whether a synced catalog exists in dir.
func SATCached(dir string) bool {
	_, err := os.Stat(SATPath(dir))
	return err == nil
}

// SyncSAT downloads the SAT product-keys catalog and writes the local cache,
// returning the cached catalog metadata (entries are populated). Re-running is
// idempotent: it re-fetches and overwrites atomically.
func SyncSAT(ctx context.Context, dir string) (*SATCatalog, error) {
	version, err := fetchURL(ctx, SATVersionURL)
	if err != nil {
		return nil, fmt.Errorf("fetch catalog version: %w", err)
	}
	data, err := fetchURL(ctx, SATDataURL)
	if err != nil {
		return nil, fmt.Errorf("fetch catalog data: %w", err)
	}
	entries, err := parseSATSQL(data)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("parsed 0 entries — source format likely changed; refusing to write an empty catalog")
	}

	cat := &SATCatalog{
		Version:   strings.TrimSpace(string(version)),
		Source:    SATDataURL,
		FetchedAt: time.Now().UTC(),
		Entries:   entries,
	}
	if err := writeSATCache(dir, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

// LoadSAT reads the cached catalog from dir.
func LoadSAT(dir string) (*SATCatalog, error) {
	f, err := os.Open(SATPath(dir))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("SAT product-keys catalog not synced — run `alegra catalog sync-sat`")
		}
		return nil, err
	}
	defer func() { _ = f.Close() }()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("corrupt catalog cache (re-run `alegra catalog sync-sat`): %w", err)
	}
	defer func() { _ = gz.Close() }()
	var cat SATCatalog
	if err := json.NewDecoder(gz).Decode(&cat); err != nil {
		return nil, fmt.Errorf("corrupt catalog cache (re-run `alegra catalog sync-sat`): %w", err)
	}
	return &cat, nil
}

// SearchSAT returns entries whose code, name, or similar-names list contains
// the query, case- and diacritic-insensitively. An empty query matches all;
// limit <= 0 means no limit.
func SearchSAT(cat *SATCatalog, query string, limit int) []SATEntry {
	q := foldSpanish(query)
	var out []SATEntry
	for _, e := range cat.Entries {
		if q == "" || strings.Contains(e.Code, q) ||
			strings.Contains(foldSpanish(e.Name), q) ||
			strings.Contains(foldSpanish(e.Similar), q) {
			out = append(out, e)
			if limit > 0 && len(out) >= limit {
				break
			}
		}
	}
	return out
}

// foldSpanish lowercases and strips the Spanish diacritics that appear in the
// catalog, so "Visón" matches "vison".
func foldSpanish(s string) string {
	return strings.NewReplacer(
		"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u", "ñ", "n",
	).Replace(strings.ToLower(s))
}

func writeSATCache(dir string, cat *SATCatalog) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	// Write-to-temp + rename so a failed sync can't corrupt an existing cache.
	tmp, err := os.CreateTemp(dir, "sat-product-keys-*.tmp")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmp.Name()) }()
	gz := gzip.NewWriter(tmp)
	if err := json.NewEncoder(gz).Encode(cat); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := gz.Close(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), SATPath(dir))
}

// parseSATSQL extracts (id, texto, similares) from the phpcfdi SQLite dump:
// one `INSERT INTO cfdi_40_productos_servicios VALUES('id','texto',...);` per
// line, nine positional values, strings quoted with doubled single quotes
// as the escape.
func parseSATSQL(data []byte) ([]SATEntry, error) {
	var entries []SATEntry
	sc := bufio.NewScanner(strings.NewReader(string(data)))
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		rest, found := strings.CutPrefix(line, "INSERT INTO cfdi_40_productos_servicios VALUES(")
		if !found {
			continue
		}
		vals, err := parseSQLValues(strings.TrimSuffix(rest, ");"))
		if err != nil {
			return nil, fmt.Errorf("malformed insert line: %w", err)
		}
		if len(vals) != 9 {
			return nil, fmt.Errorf("expected 9 values per row, got %d: %q", len(vals), line)
		}
		entries = append(entries, SATEntry{Code: vals[0], Name: vals[1], Similar: vals[8]})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// parseSQLValues splits a SQLite VALUES tuple body into its fields, honoring
// single-quoted strings with doubled-quote escapes and bare tokens (numbers).
func parseSQLValues(s string) ([]string, error) {
	var vals []string
	for i := 0; i < len(s); {
		if s[i] == '\'' {
			var b strings.Builder
			i++
			for {
				j := strings.IndexByte(s[i:], '\'')
				if j < 0 {
					return nil, fmt.Errorf("unterminated string at %q", s)
				}
				b.WriteString(s[i : i+j])
				i += j + 1
				if i < len(s) && s[i] == '\'' { // doubled quote = literal quote
					b.WriteByte('\'')
					i++
					continue
				}
				break
			}
			vals = append(vals, b.String())
		} else {
			j := strings.IndexByte(s[i:], ',')
			if j < 0 {
				vals = append(vals, strings.TrimSpace(s[i:]))
				i = len(s)
			} else {
				vals = append(vals, strings.TrimSpace(s[i:i+j]))
				i += j
			}
		}
		if i < len(s) {
			if s[i] != ',' {
				return nil, fmt.Errorf("expected ',' at byte %d of %q", i, s)
			}
			i++
		}
	}
	return vals, nil
}

func fetchURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	c := &http.Client{Timeout: 2 * time.Minute}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %d", url, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
