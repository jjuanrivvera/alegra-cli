package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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
	// Survives across pages: a later example must not forget that an earlier
	// response schema already typed the field.
	fromSchema := map[string]map[string]bool{}
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
		mergeSpecFieldsInto(acc, fromSchema, spec)
	}
	preferStructuredExampleTypes(acc, fromSchema)
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
//
// A field typed by a response schema keeps that type. Example payloads fill
// fields the schema doesn't mention, but they must not widen a schema type:
// several purchase-order pages publish the create/update body (provider as an
// id string) as the GET example next to a response schema that types provider
// as an object.
func mergeSpecFields(acc map[string]map[string]map[string]bool, spec map[string]any) {
	fromSchema := map[string]map[string]bool{}
	mergeSpecFieldsInto(acc, fromSchema, spec)
	preferStructuredExampleTypes(acc, fromSchema)
}

func mergeSpecFieldsInto(acc map[string]map[string]map[string]bool, fromSchema map[string]map[string]bool, spec map[string]any) {
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
				media := responseMedia(responses, code, spec)
				if media == nil {
					continue
				}
				fields := map[string]string{}
				schemaTyped := false
				if schema, _ := media["schema"].(map[string]any); schema != nil {
					fields = objectFields(schema, spec)
					schemaTyped = len(fields) > 0
				}
				if len(fields) == 0 {
					fields = exampleFields(media, spec)
				}
				if isStatusEnvelope(fields) {
					continue
				}
				for name, typ := range fields {
					recordField(acc, fromSchema, resource, name, typ, schemaTyped)
				}
			}
		}
	}
}

// preferStructuredExampleTypes drops primitive example types for a field that
// any example typed as object or array. Schema-sourced types are left alone,
// so a real oneOf (string in one country, object in another) still records
// both.
func preferStructuredExampleTypes(acc map[string]map[string]map[string]bool, fromSchema map[string]map[string]bool) {
	for resource, fields := range acc {
		for name, types := range fields {
			if fromSchema[resource][name] {
				continue
			}
			if !types["object"] && !types["array"] {
				continue
			}
			for _, primitive := range []string{"string", "number", "integer", "boolean"} {
				delete(types, primitive)
			}
		}
	}
}

// recordField adds a documented type for a resource field. Once any response
// schema has typed the field, later example-only types are ignored, and a
// schema type replaces example types recorded earlier in the walk.
func recordField(acc map[string]map[string]map[string]bool, fromSchema map[string]map[string]bool, resource, name, typ string, schema bool) {
	if typ == "" || name == "" { // property declared without type information
		return
	}
	if acc[resource] == nil {
		acc[resource] = map[string]map[string]bool{}
	}
	if fromSchema[resource] == nil {
		fromSchema[resource] = map[string]bool{}
	}
	if schema {
		if !fromSchema[resource][name] {
			acc[resource][name] = map[string]bool{}
			fromSchema[resource][name] = true
		}
	} else if fromSchema[resource][name] {
		return
	}
	if acc[resource][name] == nil {
		acc[resource][name] = map[string]bool{}
	}
	for _, part := range strings.Split(typ, "|") {
		if part != "" {
			acc[resource][name][part] = true
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

func responseMedia(responses map[string]any, code string, spec map[string]any) map[string]any {
	resp, _ := responses[code].(map[string]any)
	if resp == nil {
		return nil
	}
	// List pages (invoices, credit-notes, estimates, remissions) now point the
	// 200 response at #/components/responses/... instead of inlining it.
	if _, ok := resp["$ref"]; ok {
		resp = deref(resp, spec)
		if resp == nil {
			return nil
		}
	}
	content, _ := resp["content"].(map[string]any)
	media, _ := content["application/json"].(map[string]any)
	return media
}

// exampleFields reduces a media object's example payload(s) to field names and
// JSON value types, applying the same list/envelope unwrapping as schemas.
func exampleFields(media map[string]any, spec map[string]any) map[string]string {
	out := map[string]string{}
	if ex, ok := media["example"]; ok {
		mergeFieldTypes(out, exampleObjectFields(unwrapExample(ex, spec)))
	}
	if examples, ok := media["examples"].(map[string]any); ok {
		for _, rawEx := range examples {
			// Union across examples. Copying would let one request-shaped
			// example (priceList: "2") replace the object form from another,
			// and map iteration made that choice unstable.
			mergeFieldTypes(out, exampleObjectFields(unwrapExample(rawEx, spec)))
		}
	}
	return out
}

func mergeFieldTypes(dst, src map[string]string) {
	for name, typ := range src {
		dst[name] = unionType(dst[name], typ)
	}
}

// unwrapExample follows local $refs and OpenAPI Example Objects ({summary,
// value} or a $ref to one) down to the payload the docs are illustrating.
func unwrapExample(v any, spec map[string]any) any {
	for range 10 {
		m, ok := v.(map[string]any)
		if !ok {
			return v
		}
		if _, has := m["$ref"]; has {
			target := deref(m, spec)
			if target == nil {
				return nil
			}
			v = target
			continue
		}
		if isExampleObject(m) {
			v = m["value"]
			continue
		}
		return m
	}
	return nil
}

// isExampleObject reports whether m is an OpenAPI Example Object rather than
// a resource payload that happens to contain a "value" field.
func isExampleObject(m map[string]any) bool {
	if _, ok := m["value"]; !ok {
		return false
	}
	for k := range m {
		switch k {
		case "summary", "description", "value", "externalValue":
		default:
			return false
		}
	}
	return true
}

func exampleObjectFields(value any) map[string]string {
	switch v := value.(type) {
	case []any: // bare list response: union across the example's items
		out := map[string]string{}
		for _, item := range v {
			mergeFieldTypes(out, exampleObjectFields(item))
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
	// allOf/oneOf/anyOf before the array arm: a node can be `type: object`
	// with oneOf of country variants (bills) or a oneOf of a bare list and a
	// {metadata, data} envelope (stamp documents) and have no properties of
	// its own.
	if fields, ok := composedFields(schema, spec); ok {
		return fields
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

// composedFields unions the object fields of allOf/oneOf/anyOf branches.
// ok is false when the schema isn't composed, so the caller can keep walking.
func composedFields(schema map[string]any, spec map[string]any) (map[string]string, bool) {
	var alts []any
	for _, key := range []string{"allOf", "oneOf", "anyOf"} {
		if a, ok := schema[key].([]any); ok {
			alts = append(alts, a...)
		}
	}
	if len(alts) == 0 {
		return nil, false
	}
	out := map[string]string{}
	for _, alt := range alts {
		m, ok := alt.(map[string]any)
		if !ok {
			continue
		}
		for name, typ := range objectFields(m, spec) {
			out[name] = unionType(out[name], typ)
		}
	}
	if props, ok := schema["properties"].(map[string]any); ok {
		for name, raw := range props {
			prop, isMap := raw.(map[string]any)
			if !isMap {
				continue
			}
			out[name] = unionType(out[name], schemaType(prop, spec))
		}
	}
	return out, true
}

func unionType(a, b string) string {
	if a == "" {
		return b
	}
	if b == "" || a == b {
		return a
	}
	set := map[string]bool{}
	for _, part := range strings.Split(a+"|"+b, "|") {
		if part != "" {
			set[part] = true
		}
	}
	ts := make([]string, 0, len(set))
	for t := range set {
		ts = append(ts, t)
	}
	sort.Strings(ts)
	return strings.Join(ts, "|")
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
	for _, key := range []string{"oneOf", "anyOf", "allOf"} {
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

// deref resolves a local JSON pointer ("#/components/schemas/X",
// "#/components/responses/...", "#/paths/~1debit-notes/...") within the page's
// own document. Each reference page embeds a self-contained spec, but the
// pointer is no longer always a component schema: list responses and examples
// are shared via components, and some operations point at a sibling path.
func deref(schema map[string]any, spec map[string]any) map[string]any {
	for range 10 { // cycle guard
		ref, ok := schema["$ref"].(string)
		if !ok {
			return schema
		}
		// Contact (and similar) response schemas $ref the POST requestBody,
		// where relations are id strings. That is the write shape, not the
		// response; refuse it so the response examples supply the fields.
		if strings.Contains(ref, "/requestBody/") {
			return nil
		}
		target := resolvePointer(spec, ref)
		if target == nil {
			return nil
		}
		schema = target
	}
	return nil
}

// resolvePointer walks a "#/a/b" JSON pointer. Tokens may be OpenAPI-escaped
// (~1 for /, ~0 for ~) and percent-encoded ({id} is often %7Bid%7D).
func resolvePointer(doc map[string]any, ref string) map[string]any {
	if !strings.HasPrefix(ref, "#/") {
		return nil
	}
	var cur any = doc
	for _, raw := range strings.Split(ref[2:], "/") {
		token, err := url.PathUnescape(raw)
		if err != nil {
			return nil
		}
		// RFC 6901: ~1 is "/" and ~0 is "~", applied in that order so ~01 stays "~1".
		token = strings.ReplaceAll(token, "~1", "/")
		token = strings.ReplaceAll(token, "~0", "~")
		var next any
		switch node := cur.(type) {
		case map[string]any:
			var ok bool
			next, ok = node[token]
			if !ok {
				return nil
			}
		case []any:
			// Pointers such as .../oneOf/0 index the branch.
			i, err := strconv.Atoi(token)
			if err != nil || i < 0 || i >= len(node) {
				return nil
			}
			next = node[i]
		default:
			return nil
		}
		cur = next
	}
	target, ok := cur.(map[string]any)
	if !ok {
		return nil
	}
	return target
}

func writeFieldSchemas(fs FieldSchemas) error {
	out, err := json.MarshalIndent(fs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("testdata", "spec", "schemas.json"), append(out, '\n'), 0o644) //nolint:gosec // public docs
}
