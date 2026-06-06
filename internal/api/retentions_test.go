package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetentions_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/retentions", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","idRetentionReference":"1","name":"Arrendamiento de bienes muebles","calculatedBy":"percentage","percentage":4,"type":"FUENTE","status":"active"},
			{"id":2,"name":"Arrendamiento de bienes raíces","percentage":3.5,"status":"active"}
		]`))
	})

	items, err := c.Retentions().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "Arrendamiento de bienes muebles", items[0].Name)
	assert.Equal(t, "percentage", items[0].CalculatedBy)
	assert.Equal(t, Money(4), items[0].Percentage)
	// Numeric id is normalized to a string.
	assert.Equal(t, ID("2"), items[1].ID)
	assert.Equal(t, Money(3.5), items[1].Percentage)
}

func TestRetentions_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/retentions/1", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":1,"idRetentionReference":"2","name":"Arrendamiento de bienes muebles","calculatedBy":"percentage","percentage":4,"type":"FUENTE","status":"active"}`))
	})
	ret, err := c.Retentions().Get(context.Background(), "1")
	require.NoError(t, err)
	assert.Equal(t, ID("1"), ret.ID)
	assert.Equal(t, ID("2"), ret.IDRetentionReference)
	assert.Equal(t, "FUENTE", ret.Type)
	assert.Equal(t, "active", ret.Status)
}
