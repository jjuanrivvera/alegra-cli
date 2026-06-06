package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomFields_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/custom-fields", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","name":"Código de barras","description":null,"defaultValue":null,"resourceType":"item","status":"active","key":"barcode","type":"text","settings":{"isRequired":false,"showInItemVariants":true,"printOnInvoices":null}},
			{"id":2,"name":"Fecha de vencimiento","description":"Indica la fecha en la que se vence el producto","resourceType":"item","status":"active","key":null,"type":"date","settings":{"isRequired":true,"showInItemVariants":null,"printOnInvoices":true}}
		]`))
	})
	items, err := c.CustomFields().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "Código de barras", items[0].Name)
	assert.Equal(t, "text", items[0].Type)
	assert.Equal(t, ID("2"), items[1].ID)
	require.NotNil(t, items[1].Settings)
	assert.True(t, items[1].Settings.IsRequired)
}

func TestCustomFields_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/custom-fields/1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id":"1",
			"name":"Código de barras",
			"description":null,
			"defaultValue":null,
			"resourceType":"item",
			"status":"active",
			"key":"barcode",
			"type":"text",
			"settings":{"isRequired":false,"showInItemVariants":true,"printOnInvoices":null}
		}`))
	})
	cf, err := c.CustomFields().Get(context.Background(), "1")
	require.NoError(t, err)
	require.NotNil(t, cf)
	assert.Equal(t, ID("1"), cf.ID)
	assert.Equal(t, "Código de barras", cf.Name)
	assert.Equal(t, "item", cf.ResourceType)
	assert.Equal(t, "active", cf.Status)
	assert.Equal(t, "barcode", cf.Key)
	require.NotNil(t, cf.Settings)
	assert.True(t, cf.Settings.ShowInItemVariants)
}
