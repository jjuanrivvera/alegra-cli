package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaxes_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/taxes", r.URL.Path)
		_, _ = w.Write([]byte(`[{"id":"1","name":"IVA","percentage":19,"status":"active","type":"IVA"}]`))
	})
	items, err := c.Taxes().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "IVA", items[0].Name)
	assert.Equal(t, Money("19"), items[0].Percentage)
	assert.Equal(t, "active", items[0].Status)
}

func TestTaxes_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/taxes/1", r.URL.Path)
		_, _ = w.Write([]byte(`{
			"id":"1",
			"name":"IVA",
			"percentage":19,
			"status":"active",
			"type":"IVA",
			"categoryFavorable":{"id":"5011","name":"Impuesto a las ventas a favor","type":"asset","categoryRule":{"id":"44","name":"IVA a favor","key":"IVA_IN_FAVOR_COL"}}
		}`))
	})
	tax, err := c.Taxes().Get(context.Background(), "1")
	require.NoError(t, err)
	assert.Equal(t, ID("1"), tax.ID)
	assert.Equal(t, "IVA", tax.Name)
	require.NotNil(t, tax.CategoryFavorable)
	assert.Equal(t, ID("5011"), tax.CategoryFavorable.ID)
	require.NotNil(t, tax.CategoryFavorable.CategoryRule)
	assert.Equal(t, "IVA_IN_FAVOR_COL", tax.CategoryFavorable.CategoryRule.Key)
}
