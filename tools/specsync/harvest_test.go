package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// harvestPageFixture mirrors a /reference/<slug>.md page: prose, then the
// embedded OpenAPI inside a ```json block. It exercises $ref resolution, the
// bare-array list shape, the {metadata,data} envelope, and oneOf types.
const harvestPageFixture = "# Listado de impuestos\n\nProse.\n\n# OpenAPI definition\n\n```json\n" + `{
  "openapi": "3.0.3",
  "paths": {
    "/taxes": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {"type": "array", "items": {"$ref": "#/components/schemas/Tax"}}
              }
            }
          }
        }
      },
      "post": {
        "responses": {
          "201": {
            "content": {
              "application/json": {"schema": {"$ref": "#/components/schemas/Tax"}}
            }
          }
        }
      }
    },
    "/invoices": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "metadata": {"type": "object"},
                    "data": {"type": "array", "items": {"$ref": "#/components/schemas/Invoice"}}
                  }
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Tax": {
        "type": "object",
        "properties": {
          "id": {"oneOf": [{"type": "integer"}, {"type": "string"}]},
          "name": {"type": "string"},
          "percentage": {"type": "number"}
        }
      },
      "Invoice": {
        "type": "object",
        "properties": {
          "id": {"type": "integer"},
          "total": {"type": "number"},
          "client": {"$ref": "#/components/schemas/Client"}
        }
      },
      "Client": {"type": "object", "properties": {"id": {"type": "integer"}}}
    }
  }
}` + "\n```\n"

func TestExtractOpenAPI(t *testing.T) {
	spec, err := extractOpenAPI([]byte(harvestPageFixture))
	require.NoError(t, err)
	assert.Contains(t, spec, "paths")

	_, err = extractOpenAPI([]byte("# A page\n\n```json\n{\"not\": \"openapi\"}\n```\n"))
	assert.ErrorContains(t, err, "no embedded OpenAPI")
}

func TestMergeSpecFields(t *testing.T) {
	spec, err := extractOpenAPI([]byte(harvestPageFixture))
	require.NoError(t, err)

	acc := map[string]map[string]map[string]bool{}
	mergeSpecFields(acc, spec)

	// Bare-array list and $ref'd 201 detail both reduce to Tax's fields.
	require.Contains(t, acc, "taxes")
	assert.Contains(t, acc["taxes"], "percentage")
	assert.True(t, acc["taxes"]["name"]["string"])
	// oneOf union is captured as a single joined type token.
	assert.True(t, acc["taxes"]["id"]["integer|string"])

	// The {metadata,data} envelope is unwrapped to the resource object.
	require.Contains(t, acc, "invoices")
	assert.Contains(t, acc["invoices"], "total")
	assert.NotContains(t, acc["invoices"], "metadata", "envelope fields must not leak into the resource")
	assert.NotContains(t, acc["invoices"], "data")
	// A $ref'd nested object reads as type object.
	assert.True(t, acc["invoices"]["client"]["object"])
}

func TestHarvest_EndToEnd(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/get_taxes.md":
			_, _ = w.Write([]byte(harvestPageFixture))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	orig := pageURLPrefix
	pageURLPrefix = srv.URL + "/"
	t.Cleanup(func() { pageURLPrefix = orig })
	t.Chdir(t.TempDir()) // .alegra-spec/pages and testdata/spec land here

	endpoints := []RESTEndpoint{
		{Method: "GET", Path: "/taxes", Slug: "get_taxes"},
		{Method: "GET", Path: "/gone", Slug: "get_gone"}, // 404s must not sink the harvest
	}
	fs, err := harvest(endpoints)
	require.NoError(t, err)
	require.Contains(t, fs, "taxes")
	assert.Equal(t, "integer|string", fs["taxes"]["id"])
	assert.Equal(t, "number", fs["taxes"]["percentage"])

	// Second run is served from the page cache (server can be gone).
	srv.Close()
	again, err := harvest(endpoints[:1])
	require.NoError(t, err)
	assert.Equal(t, fs["taxes"], again["taxes"])

	require.NoError(t, os.MkdirAll("testdata/spec", 0o755))
	require.NoError(t, writeFieldSchemas(fs))
	written, err := os.ReadFile("testdata/spec/schemas.json")
	require.NoError(t, err)
	assert.Contains(t, string(written), `"percentage": "number"`)
}

func TestHarvest_AllPagesFailingIsFatal(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	orig := pageURLPrefix
	pageURLPrefix = srv.URL + "/"
	t.Cleanup(func() { pageURLPrefix = orig })
	t.Chdir(t.TempDir())

	_, err := harvest([]RESTEndpoint{{Slug: "get_nope"}})
	assert.ErrorContains(t, err, "harvested 0")
}

func TestIsStatusEnvelope(t *testing.T) {
	assert.True(t, isStatusEnvelope(map[string]string{"code": "number", "message": "string"}))
	assert.True(t, isStatusEnvelope(map[string]string{"message": "string"}))
	assert.False(t, isStatusEnvelope(map[string]string{"id": "number", "message": "string"}))
	assert.False(t, isStatusEnvelope(map[string]string{}))
}

func TestDeref_CycleAndMissing(t *testing.T) {
	spec := map[string]any{
		"components": map[string]any{
			"schemas": map[string]any{
				"A": map[string]any{"$ref": "#/components/schemas/B"},
				"B": map[string]any{"$ref": "#/components/schemas/A"},
			},
		},
	}
	assert.Nil(t, deref(map[string]any{"$ref": "#/components/schemas/A"}, spec), "ref cycle must terminate")
	assert.Nil(t, deref(map[string]any{"$ref": "#/components/schemas/Nope"}, spec))
	assert.Nil(t, deref(map[string]any{"$ref": "http://external/ref"}, spec))
}
