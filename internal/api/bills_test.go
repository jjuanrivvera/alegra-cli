package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBills_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bills", r.URL.Path)
		_, _ = w.Write([]byte(`[{"id":"1","date":"2019-02-21","status":"closed","total":2100,"balance":0,"provider":{"id":"1","name":"Coorporación Alegrate"}}]`))
	})
	items, err := c.Bills().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "closed", items[0].Status)
	assert.Equal(t, Money("2100"), items[0].Total)
	require.NotNil(t, items[0].Provider)
	assert.Equal(t, "Coorporación Alegrate", items[0].Provider.Name)
}

func TestBills_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bills/1", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":"1","date":"2019-02-21","dueDate":"2019-02-21","status":"open","total":12500.45,"totalPaid":500,"balance":12000.45,"provider":{"id":"1","name":"Coorporación Alegrate","identification":"159.549.847"}}`))
	})
	bill, err := c.Bills().Get(context.Background(), "1")
	require.NoError(t, err)
	require.NotNil(t, bill)
	assert.Equal(t, ID("1"), bill.ID)
	assert.Equal(t, "open", bill.Status)
	assert.Equal(t, Money("12500.45"), bill.Total)
	assert.Equal(t, Money("12000.45"), bill.Balance)
	require.NotNil(t, bill.Provider)
	assert.Equal(t, "159.549.847", bill.Provider.Identification)
}
