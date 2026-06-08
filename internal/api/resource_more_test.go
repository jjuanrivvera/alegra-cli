package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResource_ListAll_Paginates(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("start") {
		case "0":
			_, _ = w.Write([]byte(`[{"id":"1"},{"id":"2"}]`)) // full page (limit=2) → keep going
		default:
			_, _ = w.Write([]byte(`[{"id":"3"}]`)) // short page → stop
		}
	})
	res := NewResource[map[string]any](c, "contacts")
	all, err := res.ListAll(context.Background(), ListParams{Limit: 2}, 0)
	require.NoError(t, err)
	assert.Len(t, all, 3)
}

func TestResource_Update(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/contacts/5", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":"5","name":"Updated"}`))
	})
	res := NewResource[map[string]any](c, "contacts")
	out, err := res.Update(context.Background(), "5", map[string]any{"name": "Updated"})
	require.NoError(t, err)
	assert.Equal(t, "Updated", (*out)["name"])
}

func TestResource_Action(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/invoices/9/void", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":"9","status":"void"}`))
	})
	res := NewResource[map[string]any](c, "invoices")
	var out map[string]any
	require.NoError(t, res.Action(context.Background(), "9", "void", nil, &out))
	assert.Equal(t, "void", out["status"])
}

func TestResource_CollectionAction(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/bills/import-by-cufe", r.URL.Path)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	res := NewResource[map[string]any](c, "bills")
	var out map[string]any
	require.NoError(t, res.CollectionAction(context.Background(), "import-by-cufe", map[string]any{"cufe": "x"}, &out))
	assert.Equal(t, true, out["ok"])
}

func TestResource_Raw(t *testing.T) {
	c := newTestClient(t, func(http.ResponseWriter, *http.Request) {})
	res := NewResource[map[string]any](c, "company")
	assert.Same(t, c, res.Raw())
	assert.Equal(t, "company", res.Path())
}

func TestClient_PostIntoAndPutInto(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, []string{http.MethodPost, http.MethodPut}, r.Method)
		_, _ = w.Write([]byte(`{"echo":"` + r.Method + `"}`))
	})
	var out map[string]any
	require.NoError(t, c.PostInto(context.Background(), "global-invoices", map[string]any{"a": 1}, &out))
	assert.Equal(t, http.MethodPost, out["echo"])

	require.NoError(t, c.PutInto(context.Background(), "company", map[string]any{"name": "X"}, &out))
	assert.Equal(t, http.MethodPut, out["echo"])
}
