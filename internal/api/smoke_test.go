//go:build smoke

// Live smoke tests against a REAL Alegra account (issue #27). They validate that
// the Go structs match what the API returns at runtime — not just what the docs
// say — and surface fields Alegra returns that we don't model.
//
// Gated twice so they never affect `make test` or normal CI:
//   - build tag `smoke` (this file only compiles under `go test -tags smoke`)
//   - ALEGRA_SMOKE_EMAIL / ALEGRA_SMOKE_TOKEN env vars (t.Skip when absent)
//
// Read-only by default (one List/Get per resource). The create→update→delete
// write cycle is additionally gated by ALEGRA_SMOKE_WRITE=1 and limited to safe
// master-data resources (never fiscal documents).
//
//	make smoke          # read-only
//	ALEGRA_SMOKE_WRITE=1 make smoke
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func smokeClient(t *testing.T) *Client {
	t.Helper()
	email := os.Getenv("ALEGRA_SMOKE_EMAIL")
	token := os.Getenv("ALEGRA_SMOKE_TOKEN")
	if email == "" || token == "" {
		t.Skip("live smoke tests: set ALEGRA_SMOKE_EMAIL and ALEGRA_SMOKE_TOKEN to run")
	}
	opts := []Option{WithBasicAuth(email, token)}
	if base := os.Getenv("ALEGRA_SMOKE_BASE_URL"); base != "" {
		opts = append(opts, WithBaseURL(base))
	}
	return New(opts...)
}

// jsonTags returns the top-level json field names a struct type declares.
func jsonTags(rt reflect.Type) map[string]bool {
	for rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	tags := map[string]bool{}
	if rt.Kind() != reflect.Struct {
		return tags
	}
	for i := 0; i < rt.NumField(); i++ {
		name, _, _ := strings.Cut(rt.Field(i).Tag.Get("json"), ",")
		if name != "" && name != "-" {
			tags[name] = true
		}
	}
	return tags
}

// unmodeledKeys reports keys present in the response object but not declared on
// the struct type. Structs are intentionally non-exhaustive, so this is a
// warning signal (fields Alegra returns that we could surface), not a failure.
func unmodeledKeys(obj map[string]json.RawMessage, modeled map[string]bool) []string {
	var extra []string
	for k := range obj {
		if !modeled[k] {
			extra = append(extra, k)
		}
	}
	sort.Strings(extra)
	return extra
}

// firstListObject unwraps a list response (bare array or {data|results|rows|
// subscriptions} wrapper) and returns the first element's keys, or nil if empty.
func firstListObject(raw []byte) map[string]json.RawMessage {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil
	}
	var arr []json.RawMessage
	if trimmed[0] == '[' {
		_ = json.Unmarshal(trimmed, &arr)
	} else {
		var wrapper map[string]json.RawMessage
		_ = json.Unmarshal(trimmed, &wrapper)
		for _, k := range []string{"data", "results", "rows", "subscriptions"} {
			if v, ok := wrapper[k]; ok {
				_ = json.Unmarshal(v, &arr)
				break
			}
		}
	}
	if len(arr) == 0 {
		return nil
	}
	var obj map[string]json.RawMessage
	_ = json.Unmarshal(arr[0], &obj)
	return obj
}

// planGated reports whether err is a 402/403 — the resource isn't in this
// account's plan, which is not a fidelity failure (skip rather than fail).
func planGated(err error) bool {
	if apiErr, ok := AsAPIError(err); ok {
		return apiErr.StatusCode == http.StatusPaymentRequired || apiErr.StatusCode == http.StatusForbidden
	}
	return false
}

// smokeRead does one live List per resource: assert it succeeds, decodes into
// the struct (catches type mismatches like the seller `identification` bug), and
// reports unmodeled fields.
func smokeRead[T any](t *testing.T, c *Client, path string) {
	t.Run(path, func(t *testing.T) {
		raw, err := c.do(context.Background(), http.MethodGet, path, url.Values{"limit": {"1"}}, nil)
		if planGated(err) {
			t.Skipf("/%s not in this account's plan (%v)", path, err)
		}
		require.NoError(t, err, "GET /%s must succeed", path)

		if _, err := decodeList[T](raw); err != nil {
			t.Fatalf("/%s: response does not decode into %T: %v", path, *new(T), err)
		}

		obj := firstListObject(raw)
		if obj == nil {
			t.Logf("/%s: empty list (valid pass)", path)
			return
		}
		if extra := unmodeledKeys(obj, jsonTags(reflect.TypeFor[T]())); len(extra) > 0 {
			t.Logf("⚠ /%s returns %d unmodeled field(s): %v", path, len(extra), extra)
		}
	})
}

// smokeGetObject does one live GET for a singleton resource (e.g. /company): the
// response is a single object rather than a list.
func smokeGetObject[T any](t *testing.T, c *Client, path string) {
	t.Run(path, func(t *testing.T) {
		raw, err := c.do(context.Background(), http.MethodGet, path, nil, nil)
		if planGated(err) {
			t.Skipf("/%s not in this account's plan (%v)", path, err)
		}
		require.NoError(t, err, "GET /%s must succeed", path)

		var typed T
		if err := json.Unmarshal(bytes.TrimSpace(raw), &typed); err != nil {
			t.Fatalf("/%s: response does not decode into %T: %v", path, typed, err)
		}
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(bytes.TrimSpace(raw), &obj); err == nil {
			if extra := unmodeledKeys(obj, jsonTags(reflect.TypeFor[T]())); len(extra) > 0 {
				t.Logf("⚠ /%s returns %d unmodeled field(s): %v", path, len(extra), extra)
			}
		}
	})
}

// TestSmokeReadOnly_AllResources lists every registered resource against the
// live account. Safe to run repeatedly: read-only, and an empty list is a pass.
func TestSmokeReadOnly_AllResources(t *testing.T) {
	c := smokeClient(t)

	smokeRead[AdditionalCharge](t, c, "additional-charges")
	smokeRead[BankAccount](t, c, "bank-accounts")
	smokeRead[Bill](t, c, "bills")
	smokeRead[Category](t, c, "categories")
	smokeRead[Contact](t, c, "contacts")
	smokeRead[CostCenter](t, c, "cost-centers")
	smokeRead[CreditNote](t, c, "credit-notes")
	smokeRead[Currency](t, c, "currencies")
	smokeRead[CustomField](t, c, "custom-fields")
	smokeRead[DebitNote](t, c, "debit-notes")
	smokeRead[Estimate](t, c, "estimates")
	smokeRead[GlobalInvoice](t, c, "global-invoices")
	smokeRead[IncomeDebitNote](t, c, "income-debit-notes")
	smokeRead[InventoryAdjustment](t, c, "inventory-adjustments")
	smokeRead[InventoryAdjustmentNumeration](t, c, "inventory-adjustments/numerations")
	smokeRead[Invoice](t, c, "invoices")
	smokeRead[ItemCategory](t, c, "item-categories")
	smokeRead[Item](t, c, "items")
	smokeRead[Journal](t, c, "journals")
	smokeRead[NumberTemplate](t, c, "number-templates")
	smokeRead[Payment](t, c, "payments")
	smokeRead[PriceList](t, c, "price-lists")
	smokeRead[PurchaseOrder](t, c, "purchase-orders")
	smokeRead[Reconciliation](t, c, "conciliations")
	smokeRead[RecurringInvoice](t, c, "recurring-invoices")
	smokeRead[RecurringPayment](t, c, "recurring-payments")
	smokeRead[Remission](t, c, "remissions")
	smokeRead[Retention](t, c, "retentions")
	smokeRead[Seller](t, c, "sellers")
	smokeRead[Tax](t, c, "taxes")
	smokeRead[Term](t, c, "terms")
	smokeRead[TransportationReceipt](t, c, "transportation-receipts")
	smokeRead[User](t, c, "users")
	smokeRead[VariantAttribute](t, c, "variant-attributes")
	smokeRead[WarehouseTransfer](t, c, "warehouse-transfers")
	smokeRead[Warehouse](t, c, "warehouses")
	smokeRead[WebhookSubscription](t, c, "webhooks/subscriptions")

	// Singleton (single object, not a list).
	smokeGetObject[Company](t, c, "company")
}

// reflectID reads the ID field of a returned record via reflection (the structs
// share an `ID api.ID` field but generics can't reach it directly).
func reflectID(v any) string {
	rv := reflect.Indirect(reflect.ValueOf(v))
	if !rv.IsValid() {
		return ""
	}
	f := rv.FieldByName("ID")
	if !f.IsValid() {
		return ""
	}
	return fmt.Sprint(f.Interface())
}

// smokeWriteCycle exercises create→get→update→delete for one resource, cleaning
// up the throwaway record even if a step fails.
func smokeWriteCycle[T any](t *testing.T, c *Client, path string, create, update map[string]any) {
	t.Run(path, func(t *testing.T) {
		r := NewResource[T](c, path)
		ctx := context.Background()

		created, err := r.Create(ctx, create)
		if planGated(err) {
			t.Skipf("/%s not in this account's plan (%v)", path, err)
		}
		require.NoError(t, err, "create /%s", path)
		require.NotNil(t, created)

		id := reflectID(created)
		require.NotEmpty(t, id, "created /%s record must have an id", path)
		t.Cleanup(func() { _ = r.Delete(ctx, id) }) // best-effort, in case a step below fails

		got, err := r.Get(ctx, id)
		require.NoError(t, err, "get /%s/%s", path, id)
		require.NotNil(t, got)

		_, err = r.Update(ctx, id, update)
		require.NoError(t, err, "update /%s/%s", path, id)

		assert.NoError(t, r.Delete(ctx, id), "delete /%s/%s", path, id)
	})
}

// TestSmokeWriteCycle exercises request bodies on a curated set of SAFE
// master-data resources. Never fiscal documents (invoices/bills/credit-notes),
// which emit to DIAN/SAT. Gated by ALEGRA_SMOKE_WRITE=1.
func TestSmokeWriteCycle(t *testing.T) {
	if os.Getenv("ALEGRA_SMOKE_WRITE") != "1" {
		t.Skip("set ALEGRA_SMOKE_WRITE=1 to run the create→update→delete write cycle")
	}
	c := smokeClient(t)
	now := time.Now()
	tag := fmt.Sprintf("smoke-%d", now.Unix())

	// Contact bodies are country-specific; this is the Colombia shape (NIT
	// identification object + kindOfPerson), matching the smoke account. Alegra's
	// contact update is a full replace (not a patch), so the update body repeats
	// the whole create body. The NIT is clock-derived to avoid collisions.
	contact := func(observations string) map[string]any {
		return map[string]any{
			"name":           tag,
			"identification": map[string]any{"type": "NIT", "number": fmt.Sprintf("9%08d", now.UnixNano()%100000000)},
			"type":           []string{"client"},
			"kindOfPerson":   "LEGAL_ENTITY",
			"observations":   observations,
		}
	}
	smokeWriteCycle[Contact](t, c, "contacts", contact(""), contact("smoke-updated"))

	smokeWriteCycle[Item](t, c, "items",
		map[string]any{"name": tag, "price": 1000},
		map[string]any{"description": "smoke-updated"})

	smokeWriteCycle[Seller](t, c, "sellers",
		map[string]any{"name": tag},
		map[string]any{"observations": "smoke-updated"})
}
