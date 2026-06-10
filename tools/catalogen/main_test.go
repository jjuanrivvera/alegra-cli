package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sample = "# Colombia\n\n" +
	"## Unidades de medida\n\nA continuación las unidades.\n\n" +
	"| Key     | Nombre   |\n| ------- | -------- |\n| unit    | Unidad   |\n| service | Servicio |\n\n" +
	"## Tipos de identificación\n\n" +
	"| Id      | Descripción |\n| ------- | ----------- |\n| NIT\\_X  | NIT escaped |\n| CC      | Cédula      |\n\n" +
	"## Ubicaciones\n\n" +
	"| Prov | Ciudades                          |\n| ---- | --------------------------------- |\n| A    | uno <br /> dos <br /> tres        |\n"

func TestParse(t *testing.T) {
	cats := parse(sample)
	// The geo "Ubicaciones" section's only row is a <br />-joined list, which is
	// filtered out, leaving the category empty and therefore dropped.
	if len(cats) != 2 {
		t.Fatalf("want 2 categories, got %d: %+v", len(cats), cats)
	}

	if cats[0].Key != "unidades-de-medida" || len(cats[0].Entries) != 2 {
		t.Fatalf("units category wrong: %+v", cats[0])
	}
	if cats[0].Entries[0] != (entry{Code: "unit", Name: "Unidad"}) {
		t.Fatalf("first unit wrong: %+v", cats[0].Entries[0])
	}

	if cats[1].Key != "tipos-de-identificacion" {
		t.Fatalf("second category key wrong: %q", cats[1].Key)
	}
	// Markdown escaping is stripped.
	if cats[1].Entries[0].Code != "NIT_X" {
		t.Fatalf("escape not stripped: %q", cats[1].Entries[0].Code)
	}
}

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Unidades de medida":                  "unidades-de-medida",
		"Tipos de identificación":             "tipos-de-identificacion",
		"Estado de emisión (emission_status)": "estado-de-emision",
		"Regímenes":                           "regimenes",
	}
	for in, want := range cases {
		if got := slugify(in); got != want {
			t.Errorf("slugify(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRun_EndToEnd(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, ".md") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(sample))
	}))
	defer srv.Close()
	orig := base
	base = srv.URL + "/"
	t.Cleanup(func() { base = orig })

	outDir := t.TempDir()
	if err := run(outDir); err != nil {
		t.Fatalf("run: %v", err)
	}

	// One committed-shape JSON per source country, parseable and populated.
	for _, s := range sources {
		data, err := os.ReadFile(filepath.Join(outDir, s.key+".json"))
		if err != nil {
			t.Fatalf("%s: %v", s.key, err)
		}
		var cat catalog
		if err := json.Unmarshal(data, &cat); err != nil {
			t.Fatalf("%s: invalid JSON: %v", s.key, err)
		}
		if cat.Country != s.key || cat.Label != s.label || len(cat.Categories) != 2 {
			t.Fatalf("%s: wrong catalog: %+v", s.key, cat)
		}
	}
}

func TestRun_FetchErrorAborts(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	orig := base
	base = srv.URL + "/"
	t.Cleanup(func() { base = orig })

	err := run(t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "HTTP 404") {
		t.Fatalf("want HTTP 404 error, got %v", err)
	}
}

func TestRun_EmptyParseRefusesToWrite(t *testing.T) {
	// A page that parses to zero entries must abort: the JSON is embedded in
	// the binary, so writing it would silently hollow out a country's catalog.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("# Colombia\n\nprose only, no tables\n"))
	}))
	defer srv.Close()
	orig := base
	base = srv.URL + "/"
	t.Cleanup(func() { base = orig })

	outDir := t.TempDir()
	err := run(outDir)
	if err == nil || !strings.Contains(err.Error(), "0 entries") {
		t.Fatalf("want 0-entries error, got %v", err)
	}
	files, _ := os.ReadDir(outDir)
	if len(files) != 0 {
		t.Fatalf("no catalog files should be written, found %d", len(files))
	}
}
