package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayments_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/payments", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":1,"date":"2015-12-10","type":"in","status":"open","amount":150,"paymentMethod":"cash","client":{"id":"20","name":"Juan Carlos"},"bankAccount":{"id":2,"name":"Bancolombia","type":"bank"}},
			{"id":"2","date":"2015-12-15","type":"out","amount":"500.50","paymentMethod":"transfer"}
		]`))
	})

	payments, err := c.Payments().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, payments, 2)
	assert.Equal(t, ID("1"), payments[0].ID)
	assert.Equal(t, "in", payments[0].Type)
	assert.Equal(t, "open", payments[0].Status)
	assert.Equal(t, Money(150), payments[0].Amount)
	require.NotNil(t, payments[0].Client)
	assert.Equal(t, "Juan Carlos", payments[0].Client.Name)
	require.NotNil(t, payments[0].BankAccount)
	assert.Equal(t, "bank", payments[0].BankAccount.Type)
	// Numeric-string id and amount are normalized.
	assert.Equal(t, ID("2"), payments[1].ID)
	assert.Equal(t, Money(500.50), payments[1].Amount)
}

func TestPayments_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/payments/5", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":5,"type":"in","status":"open","amount":100,"invoices":[{"id":"6","number":"AL-12","total":150,"balance":50,"amount":100}]}`))
	})
	payment, err := c.Payments().Get(context.Background(), "5")
	require.NoError(t, err)
	assert.Equal(t, ID("5"), payment.ID)
	assert.Equal(t, "open", payment.Status)
	require.Len(t, payment.Invoices, 1)
	assert.Equal(t, "AL-12", payment.Invoices[0].Number)
	assert.Equal(t, Money(100), payment.Invoices[0].Amount)
}

func TestPayments_Void(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/payments/7/void", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":7,"status":"void"}`))
	})
	var out Payment
	require.NoError(t, c.Payments().Action(context.Background(), "7", "void", nil, &out))
	assert.Equal(t, "void", out.Status)
}
