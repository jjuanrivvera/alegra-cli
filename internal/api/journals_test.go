package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJournals_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/journals", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"35","date":"2022-10-25","status":"open","total":150000.55}
		]`))
	})
	items, err := c.Journals().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ID("35"), items[0].ID)
	assert.Equal(t, "open", items[0].Status)
	assert.Equal(t, Money(150000.55), items[0].Total)
}

func TestJournals_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/journals/35", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id":"35",
			"date":"2022-10-25",
			"observations":"Ajuste por anticipo",
			"status":"open",
			"total":150000.55,
			"client":{"id":"3232","name":"Nombre de prueba"},
			"entries":[
				{"id":"4","description":"Cuentas por cobrar","debit":150000.55,"credit":0}
			]
		}`))
	})
	item, err := c.Journals().Get(context.Background(), "35")
	require.NoError(t, err)
	assert.Equal(t, ID("35"), item.ID)
	assert.Equal(t, "2022-10-25", item.Date)
	require.NotNil(t, item.Client)
	assert.Equal(t, "Nombre de prueba", item.Client.Name)
	require.Len(t, item.Entries, 1)
	assert.Equal(t, Money(150000.55), item.Entries[0].Debit)
}
