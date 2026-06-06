package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReconciliations_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/reconciliations", r.URL.Path)
		_, _ = w.Write([]byte(`[{"id":"1","account":{"id":"2","name":"Banco 1","type":"bank"},"realBalance":-1999,"date":"2020-07-01","status":"open","transactions":[]}]`))
	})
	items, err := c.Reconciliations().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "open", items[0].Status)
	assert.Equal(t, Money(-1999), items[0].RealBalance)
	require.NotNil(t, items[0].Account)
	assert.Equal(t, "Banco 1", items[0].Account.Name)
	assert.Equal(t, "bank", items[0].Account.Type)
}

func TestReconciliations_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/reconciliations/74", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":"74","account":{"id":"1","name":"Caja general","type":"cash"},"realBalance":100,"bankExpenses":10,"bankInterests":0,"date":"2020-07-02","status":"open","transactions":[{"id":"4","date":"2020-07-02","amount":100,"type":"in","status":"open","associations":[{"id":"5073","name":"Descuentos financieros","price":100,"quantity":"1.00","total":100,"type":"category"}]}]}`))
	})
	rec, err := c.Reconciliations().Get(context.Background(), "74")
	require.NoError(t, err)
	require.NotNil(t, rec)
	assert.Equal(t, ID("74"), rec.ID)
	assert.Equal(t, "open", rec.Status)
	assert.Equal(t, Money(100), rec.RealBalance)
	assert.Equal(t, Money(10), rec.BankExpenses)
	require.NotNil(t, rec.Account)
	assert.Equal(t, "Caja general", rec.Account.Name)
	require.Len(t, rec.Transactions, 1)
	assert.Equal(t, ID("4"), rec.Transactions[0].ID)
	assert.Equal(t, "in", rec.Transactions[0].Type)
	require.Len(t, rec.Transactions[0].Associations, 1)
	assert.Equal(t, "Descuentos financieros", rec.Transactions[0].Associations[0].Name)
}
