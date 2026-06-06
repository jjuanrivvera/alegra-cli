package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebitNotes_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/debit-notes", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"26","date":"2020-05-13","total":12000,"balance":7000,"client":{"id":"92","name":"11111"}},
			{"id":2,"date":"2020-04-29","total":"6000","client":{"id":13,"name":"123123213"}}
		]`))
	})

	notes, err := c.DebitNotes().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, notes, 2)
	assert.Equal(t, ID("26"), notes[0].ID)
	assert.Equal(t, Money(12000), notes[0].Total)
	require.NotNil(t, notes[0].Client)
	assert.Equal(t, "11111", notes[0].Client.Name)
	// Numeric id and numeric-string total are normalized.
	assert.Equal(t, ID("2"), notes[1].ID)
	assert.Equal(t, Money(6000), notes[1].Total)
}

func TestDebitNotes_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/debit-notes/26", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id":26,
			"date":"2020-05-13",
			"total":12000,
			"totalApplied":5000,
			"numberTemplate":{"number":"123213","documentType":"debitNote"},
			"items":[{"id":3,"name":"123213","price":6000,"quantity":"1.00","total":6000}],
			"refunds":[{"id":"147","number":"139","amount":5000,"type":"in","status":"open"}]
		}`))
	})

	note, err := c.DebitNotes().Get(context.Background(), "26")
	require.NoError(t, err)
	assert.Equal(t, ID("26"), note.ID)
	assert.Equal(t, Money(5000), note.TotalApplied)
	require.NotNil(t, note.NumberTemplate)
	assert.Equal(t, "debitNote", note.NumberTemplate.DocumentType)
	require.Len(t, note.Items, 1)
	assert.Equal(t, "123213", note.Items[0].Name)
	require.Len(t, note.Refunds, 1)
	assert.Equal(t, "open", note.Refunds[0].Status)
}
