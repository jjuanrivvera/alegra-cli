package catalog

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// satSQLFixture mirrors the phpcfdi dump format byte-for-byte, including the
// doubled-single-quote escape (Levi's) and unquoted integer fields.
const satSQLFixture = `PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
INSERT INTO cfdi_40_productos_servicios VALUES('01010101','No existe en el catálogo','','','','2022-01-01','','','Público en general');
INSERT INTO cfdi_40_productos_servicios VALUES('10101506','Caballos','','','','2022-01-01','',1,'Equinos, Potrancas, Potros, Yeguas');
INSERT INTO cfdi_40_productos_servicios VALUES('10101504','Visón','','','','2022-01-01','',1,'');
INSERT INTO cfdi_40_productos_servicios VALUES('53101601','Pantalones Levi''s','','','','2022-01-01','',1,'Mezclilla');
COMMIT;
`

func TestParseSATSQL(t *testing.T) {
	entries, err := parseSATSQL([]byte(satSQLFixture))
	require.NoError(t, err)
	require.Len(t, entries, 4)
	assert.Equal(t, SATEntry{Code: "01010101", Name: "No existe en el catálogo", Similar: "Público en general"}, entries[0])
	assert.Equal(t, "Pantalones Levi's", entries[3].Name, "doubled-quote escape must unescape")
}

func TestParseSATSQL_MalformedRow(t *testing.T) {
	_, err := parseSATSQL([]byte(`INSERT INTO cfdi_40_productos_servicios VALUES('only','two');`))
	assert.ErrorContains(t, err, "expected 9 values")

	_, err = parseSATSQL([]byte(`INSERT INTO cfdi_40_productos_servicios VALUES('unterminated);`))
	assert.Error(t, err)
}

func TestSyncSAT_AndLoad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/version.txt":
			_, _ = w.Write([]byte("10.7.20260603\n"))
		case "/data.sql":
			_, _ = w.Write([]byte(satSQLFixture))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	origData, origVersion := SATDataURL, SATVersionURL
	SATDataURL, SATVersionURL = srv.URL+"/data.sql", srv.URL+"/version.txt"
	t.Cleanup(func() { SATDataURL, SATVersionURL = origData, origVersion })

	dir := t.TempDir()
	assert.False(t, SATCached(dir))

	cat, err := SyncSAT(t.Context(), dir)
	require.NoError(t, err)
	assert.Equal(t, "10.7.20260603", cat.Version)
	assert.Len(t, cat.Entries, 4)
	assert.True(t, SATCached(dir))

	loaded, err := LoadSAT(dir)
	require.NoError(t, err)
	assert.Equal(t, cat.Version, loaded.Version)
	assert.Equal(t, cat.Entries, loaded.Entries)

	// Idempotent: re-sync overwrites cleanly.
	_, err = SyncSAT(t.Context(), dir)
	require.NoError(t, err)
}

func TestSyncSAT_EmptySourceFailsLoudly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/version.txt" {
			_, _ = w.Write([]byte("v"))
			return
		}
		_, _ = w.Write([]byte("-- a dump with no recognizable inserts\n"))
	}))
	defer srv.Close()
	origData, origVersion := SATDataURL, SATVersionURL
	SATDataURL, SATVersionURL = srv.URL+"/data.sql", srv.URL+"/version.txt"
	t.Cleanup(func() { SATDataURL, SATVersionURL = origData, origVersion })

	_, err := SyncSAT(t.Context(), t.TempDir())
	assert.ErrorContains(t, err, "0 entries")
}

func TestLoadSAT_NotSynced(t *testing.T) {
	_, err := LoadSAT(t.TempDir())
	assert.ErrorContains(t, err, "sync-sat")
}

func TestSearchSAT(t *testing.T) {
	entries, err := parseSATSQL([]byte(satSQLFixture))
	require.NoError(t, err)
	cat := &SATCatalog{Entries: entries}

	// Diacritic-insensitive on name.
	got := SearchSAT(cat, "vison", 0)
	require.Len(t, got, 1)
	assert.Equal(t, "10101504", got[0].Code)

	// Matches the similar-names list.
	got = SearchSAT(cat, "yeguas", 0)
	require.Len(t, got, 1)
	assert.Equal(t, "Caballos", got[0].Name)

	// Code prefix.
	assert.Len(t, SearchSAT(cat, "101015", 0), 2)

	// Limit honored; empty query matches all.
	assert.Len(t, SearchSAT(cat, "", 2), 2)
	assert.Len(t, SearchSAT(cat, "", 0), 4)
}
