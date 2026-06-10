package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

// contractAllowlist lists struct fields that are real in live responses but
// absent from the published reference docs (observed during development with
// real accounts, or returned only for certain countries/plans). Keep entries
// scoped resource.field so a genuinely removed field still fails the test.
var contractAllowlist = map[string]bool{
	// Alegra's docs omit `status` from several documents' response examples
	// even though every live list/get returns it.
	"estimates.status":          true,
	"income-debit-notes.status": true,
	"recurring-invoices.status": true,
	"recurring-payments.status": true,
	"terms.status":              true,
	"webhooks.status":           true,
	"debit-notes.status":        true,
	// Documented request-side or live-observed fields the response docs skip.
	"estimates.priceList":              true,
	"item-categories.parent":           true, // tree shape returned by ?format=tree
	"item-categories.children":         true,
	"number-templates.subDocumentType": true,
	"number-templates.branchOffice":    true,
	"terms.deletable":                  true,
	"transportation-receipts.seller":   true,
	"users.permissions":                true, // returned with ?fields=permissions
	"income-debit-notes.costCenter":    true, // documented array|object; struct uses flexible Refs
	"credit-notes.termsConditions":     true,
	"currencies.deletable":             true,
	"currencies.canBeInactive":         true,
	"currencies.autoUpdate":            true,
	// One docs example renders client as a display string; live responses are
	// {id,name} objects.
	"credit-notes.client": true,
	// Declared in the docs without any type (description only), so the harvest
	// cannot record it.
	"inventory-adjustments.nextNumber": true,
}

// jsonTypeFor maps a struct field's Go type to the documented JSON types it
// can faithfully decode. The CLI's flexible types intentionally absorb the
// API's representation drift, so they accept several.
func jsonTypeFor(t reflect.Type) map[string]bool {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t {
	case reflect.TypeFor[api.ID](), reflect.TypeFor[api.Int](), reflect.TypeFor[api.Money]():
		return map[string]bool{"string": true, "number": true, "integer": true, "boolean": false}
	case reflect.TypeFor[api.StringOrSlice]():
		return map[string]bool{"string": true, "array": true}
	case reflect.TypeFor[api.Refs]():
		return map[string]bool{"object": true, "array": true}
	case reflect.TypeFor[json.RawMessage]():
		return nil // deliberately shapeless — accepts anything
	}
	switch t.Kind() {
	case reflect.String:
		return map[string]bool{"string": true}
	case reflect.Bool:
		return map[string]bool{"boolean": true}
	case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return map[string]bool{"number": true, "integer": true}
	case reflect.Slice, reflect.Array:
		return map[string]bool{"array": true}
	case reflect.Struct, reflect.Map, reflect.Interface:
		return map[string]bool{"object": true}
	default:
		return nil
	}
}

// TestContract_StructFieldsAreDocumented is the field-level companion to the
// resource-level spec test: every json tag on a registered resource's typed
// struct must exist in the documented response fields harvested into
// testdata/spec/schemas.json, with a compatible JSON type. It catches phantom
// fields (typos, removed API fields) that resource-level matching cannot.
//
// One-directional by design: the API documents more fields than the structs
// declare (unknown JSON is intentionally ignored), so missing struct fields
// are fine — undocumented struct fields are not. Refresh with `make spec-sync`.
func TestContract_StructFieldsAreDocumented(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "testdata", "spec", "schemas.json"))
	require.NoError(t, err, "schemas manifest missing — run `make spec-sync`")

	var schemas map[string]map[string]string
	require.NoError(t, json.Unmarshal(data, &schemas))
	require.NotEmpty(t, schemas, "schemas manifest empty — run `make spec-sync`")

	// Index documented resources by the same normalized token spec_test uses,
	// absorbing slug/path formatting differences.
	docFields := map[string]map[string]string{}
	for res, fields := range schemas {
		docFields[specNorm(res)] = fields
	}

	checked := 0
	for _, contractFn := range resourceContractFns {
		path, typ := contractFn()
		fields, documented := docFields[specNorm(firstPathSegment(path))]
		if !documented {
			continue // resource has no harvested schema (docs lack schema+examples)
		}
		// Sanity guard: when almost nothing overlaps, the harvested fields
		// describe something other than this resource object (e.g. only an
		// action's response was documented) — comparing would only produce
		// noise, not drift signal.
		if overlap(typ, fields) < 3 {
			t.Logf("%s: skipped — only %d struct fields appear in the harvested schema", path, overlap(typ, fields))
			continue
		}
		checked++
		for i := range typ.NumField() {
			f := typ.Field(i)
			tag, _, _ := strings.Cut(f.Tag.Get("json"), ",")
			if tag == "" || tag == "-" {
				continue
			}
			key := firstPathSegment(path) + "." + tag
			if contractAllowlist[key] {
				continue
			}
			docType, ok := fields[tag]
			if !ok {
				t.Errorf("%s: struct %s field %q (json %q) is not in the documented response fields — phantom field, or docs drifted (refresh with `make spec-sync`, or allowlist %q with a justification)",
					path, typ.Name(), f.Name, tag, key)
				continue
			}
			accepts := jsonTypeFor(f.Type)
			if accepts == nil {
				continue
			}
			for _, dt := range strings.Split(docType, "|") {
				if !accepts[dt] {
					t.Errorf("%s: field %q is documented as %s but struct %s declares %s (cannot decode %s)",
						path, tag, docType, typ.Name(), f.Type, dt)
					break
				}
			}
		}
	}
	require.Greater(t, checked, 5, "almost no resources were contract-checked — schemas.json likely degraded; run `make spec-sync`")
	t.Logf("contract-checked %d resources against documented schemas", checked)
}

func firstPathSegment(path string) string {
	first, _, _ := strings.Cut(strings.Trim(path, "/"), "/")
	return first
}

// overlap counts how many of the struct's json tags appear in the documented
// fields.
func overlap(typ reflect.Type, fields map[string]string) int {
	n := 0
	for i := range typ.NumField() {
		tag, _, _ := strings.Cut(typ.Field(i).Tag.Get("json"), ",")
		if _, ok := fields[tag]; ok && tag != "" {
			n++
		}
	}
	return n
}
