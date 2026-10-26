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
	// oneOf union is captured as the individual JSON types.
	assert.True(t, acc["taxes"]["id"]["integer"])
	assert.True(t, acc["taxes"]["id"]["string"])

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

func TestMergeSpecFields_ResponseAndExampleRefs(t *testing.T) {
	// Invoices (and credit-notes, estimates, remissions) moved the 200 body
	// behind #/components/responses, whose examples are themselves $refs.
	spec := map[string]any{
		"openapi": "3.0.3",
		"paths": map[string]any{
			"/invoices": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": map[string]any{"$ref": "#/components/responses/getList200"},
					},
				},
			},
		},
		"components": map[string]any{
			"responses": map[string]any{
				"getList200": map[string]any{
					"content": map[string]any{
						"application/json": map[string]any{
							"examples": map[string]any{
								"invoices-obtained": map[string]any{"$ref": "#/components/examples/GetInvoiceList"},
							},
						},
					},
				},
			},
			"examples": map[string]any{
				"GetInvoiceList": map[string]any{
					"summary": "Lista",
					"value": []any{
						map[string]any{"id": "1", "client": map[string]any{"id": "1", "name": "Ana"}, "total": 10.0},
					},
				},
			},
		},
	}
	acc := map[string]map[string]map[string]bool{}
	mergeSpecFields(acc, spec)
	require.Contains(t, acc, "invoices")
	assert.True(t, acc["invoices"]["id"]["string"])
	assert.True(t, acc["invoices"]["client"]["object"])
	assert.True(t, acc["invoices"]["total"]["number"])
	assert.NotContains(t, acc["invoices"], "metadata")
}

func TestMergeSpecFields_SchemaTypeBeatsRequestExample(t *testing.T) {
	// The GET example is the create body (provider id string). The POST
	// response schema types provider as an object and omits dueDate, which
	// the example may still contribute.
	spec := map[string]any{
		"openapi": "3.0.3",
		"paths": map[string]any{
			"/purchase-orders": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": map[string]any{
							"content": map[string]any{
								"application/json": map[string]any{
									"examples": map[string]any{
										"purchaseOrder": map[string]any{
											"value": map[string]any{
												"date":     "2019-02-12",
												"dueDate":  "2019-02-12",
												"provider": "1",
											},
										},
									},
								},
							},
						},
					},
				},
				"post": map[string]any{
					"responses": map[string]any{
						"200": map[string]any{
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{
										"type": "object",
										"properties": map[string]any{
											"id":       map[string]any{"type": "string"},
											"provider": map[string]any{"type": "object"},
											"total":    map[string]any{"type": "integer"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	acc := map[string]map[string]map[string]bool{}
	mergeSpecFields(acc, spec)
	require.Contains(t, acc, "purchase-orders")
	assert.True(t, acc["purchase-orders"]["provider"]["object"])
	assert.False(t, acc["purchase-orders"]["provider"]["string"], "request-shaped example must not widen a schema object")
	assert.True(t, acc["purchase-orders"]["dueDate"]["string"], "example-only fields are kept")
	assert.True(t, acc["purchase-orders"]["id"]["string"])
}

func TestMergeSpecFields_OneOfAllOfAndPathRef(t *testing.T) {
	spec := map[string]any{
		"openapi": "3.0.3",
		"paths": map[string]any{
			"/bills": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": map[string]any{
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{
										"type": "array",
										"items": map[string]any{
											"type": "object",
											"oneOf": []any{
												map[string]any{"$ref": "#/components/schemas/bill"},
												map[string]any{"$ref": "#/components/schemas/colombia-bill"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"/wrapped": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": map[string]any{
							"content": map[string]any{
								"application/json": map[string]any{
									// Example Object placed where a raw example belongs.
									"example": map[string]any{
										"summary": "Nota",
										"value":   map[string]any{"id": "10", "observations": "x"},
									},
								},
							},
						},
					},
				},
			},
			"/debit-notes": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": map[string]any{
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{
										"type":  "array",
										"items": map[string]any{"$ref": "#/paths/~1debit-notes/post/responses/201/content/application~1json/schema"},
									},
								},
							},
						},
					},
				},
				"post": map[string]any{
					"responses": map[string]any{
						"201": map[string]any{
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{
										"type": "object",
										"properties": map[string]any{
											"id":         map[string]any{"type": "string"},
											"costCenter": map[string]any{"type": "object"},
											"bills":      map[string]any{"type": "array"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"components": map[string]any{
			"schemas": map[string]any{
				"bill": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id":       map[string]any{"type": "string"},
						"provider": map[string]any{"type": "object"},
					},
				},
				"colombia-bill": map[string]any{
					"allOf": []any{
						map[string]any{"$ref": "#/components/schemas/bill"},
						map[string]any{
							"type": "object",
							"properties": map[string]any{
								"physicalDocument": map[string]any{"type": "object"},
							},
						},
					},
				},
			},
		},
	}
	acc := map[string]map[string]map[string]bool{}
	mergeSpecFields(acc, spec)

	require.Contains(t, acc, "bills")
	assert.True(t, acc["bills"]["provider"]["object"])
	assert.True(t, acc["bills"]["physicalDocument"]["object"], "allOf country variant fields are unioned")

	require.Contains(t, acc, "debit-notes")
	assert.True(t, acc["debit-notes"]["costCenter"]["object"], "path JSON pointer must resolve")
	assert.True(t, acc["debit-notes"]["bills"]["array"])

	require.Contains(t, acc, "wrapped")
	assert.True(t, acc["wrapped"]["observations"]["string"], "Example Object wrapper must unwrap to its value")
	assert.NotContains(t, acc["wrapped"], "summary")
	assert.NotContains(t, acc["wrapped"], "value")
}

func TestMergeSpecFields_RequestSchemaDoesNotWidenResponse(t *testing.T) {
	// The contact 200 schema $refs the create-body oneOf, where priceList is
	// an id string. Response examples are objects, plus one Chile example that
	// repeats the id. The harvested response type stays object.
	spec := map[string]any{
		"openapi": "3.0.3",
		"paths": map[string]any{
			"/contacts": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": map[string]any{
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{
										"oneOf": []any{
											map[string]any{"$ref": "#/paths/~1contacts/post/requestBody/content/application~1json/schema/oneOf/0"},
										},
									},
									"examples": map[string]any{
										"list": map[string]any{
											"value": []any{map[string]any{
												"id":        "1",
												"priceList": map[string]any{"id": "1", "name": "General"},
											}},
										},
										"chile": map[string]any{
											"value": map[string]any{"id": "1", "priceList": "2"},
										},
									},
								},
							},
						},
					},
				},
				"post": map[string]any{
					"requestBody": map[string]any{
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"oneOf": []any{
										map[string]any{
											"type": "object",
											"properties": map[string]any{
												"priceList": map[string]any{"type": "string"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	acc := map[string]map[string]map[string]bool{}
	mergeSpecFields(acc, spec)
	require.Contains(t, acc, "contacts")
	assert.True(t, acc["contacts"]["priceList"]["object"])
	assert.False(t, acc["contacts"]["priceList"]["string"])
}

func TestDeref_PathPointerAndEncoding(t *testing.T) {
	spec := map[string]any{
		"paths": map[string]any{
			"/purchase-orders/{id}": map[string]any{
				"get": map[string]any{"summary": "one"},
			},
		},
	}
	got := resolvePointer(spec, "#/paths/~1purchase-orders~1%7Bid%7D/get")
	require.NotNil(t, got)
	assert.Equal(t, "one", got["summary"])
	assert.Nil(t, resolvePointer(spec, "#/paths/~1missing/get"))
	assert.Nil(t, resolvePointer(spec, "http://external/ref"))

	indexed := map[string]any{
		"oneOf": []any{
			map[string]any{"type": "string"},
			map[string]any{"type": "object"},
		},
	}
	branch := resolvePointer(indexed, "#/oneOf/1")
	require.NotNil(t, branch)
	assert.Equal(t, "object", branch["type"])
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

// Self-referential schemas used to recurse until the goroutine stack
// overflowed (a fatal error, not a recoverable panic). The walk must cut the
// cycle at the repeated $ref and still report the non-recursive fields.
func TestObjectFields_SelfReferentialSchemas(t *testing.T) {
	ref := func(name string) map[string]any {
		return map[string]any{"$ref": "#/components/schemas/" + name}
	}
	spec := map[string]any{
		"components": map[string]any{
			"schemas": map[string]any{
				// Node is oneOf itself plus a concrete branch.
				"Node": map[string]any{
					"oneOf": []any{
						ref("Node"),
						map[string]any{"type": "object", "properties": map[string]any{
							"id": map[string]any{"type": "integer"},
						}},
					},
				},
				// Ping/Pong cycle through allOf on each side.
				"Ping": map[string]any{
					"allOf": []any{ref("Pong")},
					"properties": map[string]any{
						"ping": map[string]any{"type": "string"},
					},
				},
				"Pong": map[string]any{
					"anyOf": []any{ref("Ping")},
					"properties": map[string]any{
						"pong": map[string]any{"type": "boolean"},
					},
				},
				// A list whose items are the list itself.
				"Loop": map[string]any{"type": "array", "items": ref("Loop")},
				// A tree: children reference the enclosing schema.
				"Tree": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"name":     map[string]any{"type": "string"},
						"parent":   ref("Tree"),
						"children": map[string]any{"type": "array", "items": ref("Tree")},
					},
				},
				// A property that is only a composed cycle.
				"Odd": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"value": ref("OddValue"),
					},
				},
				"OddValue": map[string]any{"oneOf": []any{ref("OddValue"), map[string]any{"type": "string"}}},
			},
		},
	}

	tests := []struct {
		name   string
		schema map[string]any
		want   map[string]string
	}{
		{"oneOf self", ref("Node"), map[string]string{"id": "integer"}},
		{"mutual allOf/anyOf", ref("Ping"), map[string]string{"ping": "string", "pong": "boolean"}},
		{"array of itself", ref("Loop"), nil},
		{"tree", ref("Tree"), map[string]string{"name": "string", "parent": "object", "children": "array"}},
		{"composed property cycle", ref("Odd"), map[string]string{"value": "object|string"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, objectFields(tt.schema, spec))
		})
	}
}

// Deep but acyclic nesting stops at maxSchemaDepth instead of fanning out.
func TestObjectFields_DepthCap(t *testing.T) {
	leaf := map[string]any{"type": "object", "properties": map[string]any{"id": map[string]any{"type": "integer"}}}
	shallow := leaf
	for range 3 {
		shallow = map[string]any{"type": "array", "items": shallow}
	}
	assert.Equal(t, map[string]string{"id": "integer"}, objectFields(shallow, nil))

	deep := leaf
	for range maxSchemaDepth + 1 {
		deep = map[string]any{"type": "array", "items": deep}
	}
	assert.Nil(t, objectFields(deep, nil))
}

func TestResolvePointer_FailureCases(t *testing.T) {
	doc := map[string]any{
		"components": map[string]any{
			"schemas": map[string]any{
				"Str":  "not an object",
				"List": []any{map[string]any{"type": "object"}},
			},
		},
	}
	tests := []struct {
		name string
		ref  string
	}{
		{"empty", ""},
		{"no leading #/", "components/schemas/List"},
		{"external", "other.json#/components/schemas/List"},
		{"missing key", "#/components/schemas/Nope"},
		{"bad percent escape", "#/components/schemas/%zz"},
		{"non-numeric index", "#/components/schemas/List/x"},
		{"negative index", "#/components/schemas/List/-1"},
		{"index out of range", "#/components/schemas/List/1"},
		{"walk through scalar", "#/components/schemas/Str/deeper"},
		{"target not an object", "#/components/schemas/Str"},
		{"target is an array", "#/components/schemas/List"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Nil(t, resolvePointer(doc, tt.ref))
		})
	}
	require.NotNil(t, resolvePointer(doc, "#/components/schemas/List/0"))
}
