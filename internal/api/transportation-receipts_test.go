package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransportationReceipts_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/transportation-receipts", r.URL.Path)
		_, _ = w.Write([]byte(`[{"id":"5","number":"APITT5","date":"2022-12-07","status":"open"}]`))
	})
	items, err := c.TransportationReceipts().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ID("5"), items[0].ID)
	assert.Equal(t, "open", items[0].Status)
}

func TestTransportationReceipts_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/transportation-receipts/5", r.URL.Path)
		_, _ = w.Write([]byte(`{
			"id":"5",
			"number":"APITT5",
			"date":"2022-12-07",
			"status":"open",
			"client":{"id":"3","name":"FLOR GARCIA ROJAS"},
			"items":[{"id":"2","name":"item 1","price":"1000.00","total":1000}]
		}`))
	})
	rec, err := c.TransportationReceipts().Get(context.Background(), "5")
	require.NoError(t, err)
	assert.Equal(t, ID("5"), rec.ID)
	assert.Equal(t, "APITT5", rec.Number)
	require.NotNil(t, rec.Client)
	assert.Equal(t, ID("3"), rec.Client.ID)
	require.Len(t, rec.Items, 1)
	assert.Equal(t, Money("1000"), rec.Items[0].Total)
}
