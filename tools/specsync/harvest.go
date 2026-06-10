package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// harvest fetches each modern REST reference page, extracts the OpenAPI
// definition each page embeds (a ```json block in the static markdown), and
// reduces it to a slim committed manifest of documented response fields per
// resource. This is what the field-level contract test in commands/ checks the
// CLI's typed structs against.
//
// Pages are cached under .alegra-spec/pages/ (gitignored): re-runs are free
// locally, and CI starts from a clean checkout so the weekly spec-sync job
// always re-fetches.

// pageURLPrefix is a var so tests can point the harvest at a local server.
var pageURLPrefix = "https://developer.alegra.com/reference/"

var jsonBlockRe = regexp.MustCompile("(?s)```json\n(.*?)\n```")

// FieldSchemas maps resource → field → documented JSON type(s). Types seen
// across pages/operations are unioned with "|" (e.g. "integer|string").
type FieldSchemas map[string]map[string]string

func harvest(endpoints []RESTEndpoint) (FieldSchemas, error) {
	if err := os.MkdirAll(filepath.Join(".alegra-spec", "pages"), 0o755); err != nil {
		return nil, err
	}
	// resource → field → set of seen types
	acc := map[string]map[string]map[string]bool{}
	var fetched, failed int
	for _, e := range endpoints {
		body, err := fetchPage(e.Slug)
		if err != nil {
			// A single missing/renamed page must not sink the whole harvest;
			// the committed manifest keeps the last good fields for it.
			fmt.Fprintf(os.Stderr, "specsync: harvest %s: %v (skipped)\n", e.Slug, err)
			failed++
			continue
		}
		fetched++
		spec, err := extractOpenAPI(body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "specsync: harvest %s: %v (skipped)\n", e.Slug, err)
			failed++
			continue
		}
		mergeSpecFields(acc, spec)
	}
	if fetched == 0 {
		return nil, fmt.Errorf("harvested 0 of %d pages — refusing to write an empty schema manifest", len(endpoints))
	}
	if failed > 0 {
		fmt.Fprintf(os.Stderr, "specsync: harvest skipped %d/%d pages\n", failed, len(endpoints))
	}

	out := FieldSchemas{}
	for res, fields := range acc {
		out[res] = map[string]string{}
		for f, types := range fields {
			var ts []string
			for t := range types {
				ts = append(ts, t)
			}
			sort.Strings(ts)
			out[res][f] = strings.Join(ts, "|")
		}
	}
	return out, nil
}

func fetchPage(slug string) ([]byte, error) {
	cache := filepath.Join(".alegra-spec", "pages", slug+".md")
	if b, err := os.ReadFile(cache); err == nil { //nolint:gosec // local cache
		return b, nil
	}
	b, err := fetch(pageURLPrefix + slug + ".md")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(cache, b, 0o644); err != nil { //nolint:gosec // public docs
		return nil, err
	}
	// Be polite to the docs host: ~5 pages/second across the ~230-page sweep.
	time.Sleep(200 * time.Millisecond)
	return b, nil
}

// extractOpenAPI returns the first ```json block that parses as an OpenAPI
// document (the pages embed exactly one, under "# OpenAPI definition").
func extractOpenAPI(page []byte) (map[string]any, error) {
	for _, m := range jsonBlockRe.FindAllSubmatch(page, -1) {
		var doc map[string]any
		if err := json.Unmarshal(m[1], &doc); err != nil {
			continue
		}
		if _, ok := doc["openapi"]; ok {
			return doc, nil
		}
	}
	return nil, fmt.Errorf("no embedded OpenAPI definition found")
}

// mergeSpecFields walks every operation's 200/201 JSON response and
// accumulates the resource object's top-level properties. Pages document
// responses either with a schema or with examples only (e.g. invoices); both
// are reduced. List envelopes are unwrapped: a bare array uses its items; a
// {metadata, data:[...]} wrapper uses the data items.
func mergeSpecFields(acc map[string]map[string]map[string]bool, spec map[string]any) {
	paths, _ := spec["paths"].(map[string]any)
	for path, rawOps := range paths {
		resource := firstSegment(path)
		if resource == "" {
			continue
		}
		ops, _ := rawOps.(map[string]any)
		for _, rawOp := range ops {
			op, isMap := rawOp.(map[string]any)
			if !isMap {
				continue
			}
			responses, _ := op["responses"].(map[string]any)
			for _, code := range []string{"200", "201"} {
				media := responseMedia(responses, code)
				if media == nil {
					continue
				}
				fields := map[string]string{}
				if schema, _ := media["schema"].(map[string]any); schema != nil {
					fields = objectFields(schema, spec)
				}
				if len(fields) == 0 {
					fields = exampleFields(media)
				}
				if isStatusEnvelope(fields) {
					continue
				}
				for name, typ := range fields {
					if typ == "" { // property declared without type information
						continue
					}
					if acc[resource] == nil {
						acc[resource] = map[string]map[string]bool{}
					}
					if acc[resource][name] == nil {
						acc[resource][name] = map[string]bool{}
					}
					acc[resource][name][typ] = true
				}
			}
		}
	}
}

// isStatusEnvelope detects action acknowledgements ({code, message}) that some
// operations (void, attachments) document as their 200 body — they describe the
// action's outcome, not the resource, and would poison the field manifest.
func isStatusEnvelope(fields map[string]string) bool {
	if len(fields) == 0 || len(fields) > 2 {
		return false
	}
	for name := range fields {
		if name != "code" && name != "message" && name != "status" {
			return false
		}
	}
	return true
}

func responseMedia(responses map[string]any, code string) map[string]any {
	resp, _ := responses[code].(map[string]any)
	content, _ := resp["content"].(map[string]any)
	media, _ := content["application/json"].(map[string]any)
	return media
}

// exampleFields reduces a media object's example payload(s) to field names and
// JSON value types, applying the same list/envelope unwrapping as schemas.
func exampleFields(media map[string]any) map[string]string {
	out := map[string]string{}
	if ex, ok := media["example"]; ok {
		maps.Copy(out, exampleObjectFields(ex))
	}
	if examples, ok := media["examples"].(map[string]any); ok {
		for _, rawEx := range examples {
			if ex, isMap := rawEx.(map[string]any); isMap {
				maps.Copy(out, exampleObjectFields(ex["value"]))
			}
		}
	}
	return out
}

func exampleObjectFields(value any) map[string]string {
	switch v := value.(type) {
	case []any: // bare list response: union across the example's items
		out := map[string]string{}
		for _, item := range v {
			maps.Copy(out, exampleObjectFields(item))
		}
		return out
	case map[string]any:
		// Paginated envelope: {metadata?, data: [...]} — an example that is
		// only a data wrapper (with or without metadata) is not the resource.
		if data, ok := v["data"].([]any); ok && len(v) <= 2 {
			_, hasMeta := v["metadata"]
			if len(v) == 1 || hasMeta {
				return exampleObjectFields(data)
			}
		}
		out := map[string]string{}
		for name, val := range v {
			if t := jsonValueType(val); t != "" {
				out[name] = t
			}
		}
		return out
	}
	return nil
}

// jsonValueType names the JSON type of a decoded example value; "" for null
// (no type information).
func jsonValueType(v any) string {
	switch v.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "boolean"
	case map[string]any:
		return "object"
	case []any:
		return "array"
	}
	return ""
}

// objectFields reduces a response schema to the resource object's top-level
// property names and types.
func objectFields(schema map[string]any, spec map[string]any) map[string]string {
	schema = deref(schema, spec)
	if schema == nil {
		return nil
	}
	// Bare list response: [ {resource} ].
	if schema["type"] == "array" {
		items, _ := schema["items"].(map[string]any)
		return objectFields(items, spec)
	}
	props, _ := schema["properties"].(map[string]any)
	// Paginated envelope: {metadata, data: [ {resource} ]}.
	if data, ok := props["data"].(map[string]any); ok {
		if d := deref(data, spec); d != nil && d["type"] == "array" {
			items, _ := d["items"].(map[string]any)
			return objectFields(items, spec)
		}
	}
	out := map[string]string{}
	for name, rawProp := range props {
		prop, isMap := rawProp.(map[string]any)
		if !isMap {
			continue
		}
		out[name] = schemaType(prop, spec)
	}
	return out
}

// schemaType names a property's JSON type, following $ref and oneOf/anyOf.
func schemaType(prop map[string]any, spec map[string]any) string {
	prop = deref(prop, spec)
	if prop == nil {
		return "object"
	}
	if t, ok := prop["type"].(string); ok && t != "" {
		return t
	}
	for _, key := range []string{"oneOf", "anyOf"} {
		if alts, ok := prop[key].([]any); ok {
			types := map[string]bool{}
			for _, alt := range alts {
				if m, isMap := alt.(map[string]any); isMap {
					types[schemaType(m, spec)] = true
				}
			}
			var ts []string
			for t := range types {
				ts = append(ts, t)
			}
			sort.Strings(ts)
			return strings.Join(ts, "|")
		}
	}
	if _, ok := prop["properties"]; ok {
		return "object"
	}
	if _, ok := prop["items"]; ok {
		return "array"
	}
	// A bare {description: ...} property declares no type at all (e.g.
	// nextNumber on numerations); report nothing rather than guess.
	return ""
}

// deref resolves a local "#/components/schemas/X" reference within the page's
// own document (each page embeds a self-contained spec).
func deref(schema map[string]any, spec map[string]any) map[string]any {
	for range 10 { // cycle guard
		ref, ok := schema["$ref"].(string)
		if !ok {
			return schema
		}
		name, ok := strings.CutPrefix(ref, "#/components/schemas/")
		if !ok {
			return nil
		}
		components, _ := spec["components"].(map[string]any)
		schemas, _ := components["schemas"].(map[string]any)
		target, isMap := schemas[name].(map[string]any)
		if !isMap {
			return nil
		}
		schema = target
	}
	return nil
}

func writeFieldSchemas(fs FieldSchemas) error {
	out, err := json.MarshalIndent(fs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("testdata", "spec", "schemas.json"), append(out, '\n'), 0o644) //nolint:gosec // public docs
}
