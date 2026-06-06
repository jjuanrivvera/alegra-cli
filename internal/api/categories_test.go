package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategories_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/categories", r.URL.Path)
		_, _ = w.Write([]byte(`[{"id":3,"idParent":1,"name":"Ventas","code":"4135","type":"income","nature":"credit","status":"active"}]`))
	})
	items, err := c.Categories().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ID("3"), items[0].ID)
	assert.Equal(t, ID("1"), items[0].IDParent)
	assert.Equal(t, "Ventas", items[0].Name)
	assert.Equal(t, "income", items[0].Type)
	assert.Equal(t, "credit", items[0].Nature)
}

func TestCategories_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/categories/3", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":3,"idParent":1,"name":"Ventas","description":"Ingresos principales","code":"","type":"income","nature":"credit","blocked":"no","readOnly":false,"categoryRule":{"id":17,"name":"Patrimonio","key":"equity"}}`))
	})
	cat, err := c.Categories().Get(context.Background(), "3")
	require.NoError(t, err)
	require.NotNil(t, cat)
	assert.Equal(t, ID("3"), cat.ID)
	assert.Equal(t, "Ventas", cat.Name)
	assert.Equal(t, "income", cat.Type)
	assert.False(t, cat.ReadOnly)
	require.NotNil(t, cat.Rule)
	assert.Equal(t, "equity", cat.Rule.Key)
}
