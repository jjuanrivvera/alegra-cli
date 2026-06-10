package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreditNotes_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/credit-notes", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","date":"2020-02-17","status":"open","total":1000},
			{"id":2,"date":"2020-03-01","status":"closed","total":"500.50"}
		]`))
	})

	notes, err := c.CreditNotes().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, notes, 2)
	assert.Equal(t, ID("1"), notes[0].ID)
	assert.Equal(t, "open", notes[0].Status)
	assert.Equal(t, Money("1000"), notes[0].Total)
	// Numeric id is normalized to a string and numeric-string money parses.
	assert.Equal(t, ID("2"), notes[1].ID)
	assert.Equal(t, Money("500.50"), notes[1].Total)
}

func TestCreditNotes_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/credit-notes/1", r.URL.Path)
		_, _ = w.Write([]byte(`{
			"id":1,
			"date":"2020-02-17",
			"status":"open",
			"total":1000,
			"client":{"id":"21","name":"Usuario Alegra"},
			"items":[{"id":"2","name":"item 2","price":1000,"quantity":"1.00","total":1000}],
			"invoices":[{"id":"69","number":"69","amount":100,"total":12122,"balance":11922}]
		}`))
	})
	note, err := c.CreditNotes().Get(context.Background(), "1")
	require.NoError(t, err)
	assert.Equal(t, ID("1"), note.ID)
	require.NotNil(t, note.Client)
	assert.Equal(t, "Usuario Alegra", note.Client.Name)
	require.Len(t, note.Items, 1)
	assert.Equal(t, "item 2", note.Items[0].Name)
	require.Len(t, note.Invoices, 1)
	assert.Equal(t, ID("69"), note.Invoices[0].ID)
}
