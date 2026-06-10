package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncomeDebitNotes_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/income-debit-notes", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","date":"2015-11-15","status":"open","total":119000},
			{"id":2,"date":"2015-12-01","status":"closed","total":"500.50"}
		]`))
	})

	notes, err := c.IncomeDebitNotes().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, notes, 2)
	assert.Equal(t, ID("1"), notes[0].ID)
	assert.Equal(t, "open", notes[0].Status)
	assert.Equal(t, Money("119000"), notes[0].Total)
	// Numeric id is normalized to a string and numeric-string money parses.
	assert.Equal(t, ID("2"), notes[1].ID)
	assert.Equal(t, Money("500.50"), notes[1].Total)
}

func TestIncomeDebitNotes_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/income-debit-notes/1", r.URL.Path)
		_, _ = w.Write([]byte(`{
			"id":1,
			"date":"2015-11-15",
			"subtotal":10000,
			"tax":19000,
			"total":119000,
			"balance":119000,
			"type":"VALUE_CHANGE",
			"client":{"id":"2","name":"Coorporación Alegrate"},
			"numberTemplate":{"id":"1","prefix":"NDI","number":1,"isElectronicResolution":true},
			"items":[{"id":"1","name":"Producto 1","price":10000,"quantity":"1","total":10000,"tax":[{"id":"3","name":"IVA","percentage":19}]}],
			"costCenter":[{"id":"1","name":"Centro 1"}]
		}`))
	})
	note, err := c.IncomeDebitNotes().Get(context.Background(), "1")
	require.NoError(t, err)
	assert.Equal(t, ID("1"), note.ID)
	assert.Equal(t, "VALUE_CHANGE", note.Type)
	assert.Equal(t, Money("119000"), note.Total)
	require.NotNil(t, note.Client)
	assert.Equal(t, "Coorporación Alegrate", note.Client.Name)
	require.NotNil(t, note.NumberTemplate)
	assert.Equal(t, "NDI", note.NumberTemplate.Prefix)
	require.Len(t, note.Items, 1)
	assert.Equal(t, "Producto 1", note.Items[0].Name)
	require.Len(t, note.Items[0].Tax, 1)
	assert.Equal(t, "IVA", note.Items[0].Tax[0].Name)
	require.Len(t, note.CostCenter, 1)
	assert.Equal(t, "Centro 1", note.CostCenter[0].Name)
}
