package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNumberTemplates_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/number-templates", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","name":"Principal","prefix":null,"status":"active","documentType":"invoice","isDefault":true,"nextInvoiceNumber":125,"maxInvoiceNumber":1000,"isElectronic":true},
			{"id":2,"name":"Cotización","status":"active","documentType":"estimate","nextInvoiceNumber":2}
		]`))
	})

	items, err := c.NumberTemplates().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "Principal", items[0].Name)
	assert.Equal(t, "invoice", items[0].DocumentType)
	assert.True(t, items[0].IsElectronic)
	assert.Equal(t, Int(125), items[0].NextInvoiceNumber)
	// Numeric id is normalized to a string.
	assert.Equal(t, ID("2"), items[1].ID)
	assert.Equal(t, "estimate", items[1].DocumentType)
}

func TestNumberTemplates_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/number-templates/1", r.URL.Path)
		_, _ = w.Write([]byte(`{
			"id":1,
			"name":"Principal",
			"status":"active",
			"documentType":"invoice",
			"resolutionNumber":"1234567890",
			"startDate":"2025-01-01",
			"endDate":"2025-12-31",
			"minInvoiceNumber":1,
			"maxInvoiceNumber":1000
		}`))
	})
	tmpl, err := c.NumberTemplates().Get(context.Background(), "1")
	require.NoError(t, err)
	assert.Equal(t, ID("1"), tmpl.ID)
	assert.Equal(t, "Principal", tmpl.Name)
	assert.Equal(t, "1234567890", tmpl.ResolutionNumber)
	assert.Equal(t, "2025-01-01", tmpl.StartDate)
	assert.Equal(t, Int(1000), tmpl.MaxInvoiceNumber)
}
