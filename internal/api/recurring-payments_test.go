package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecurringPayments_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/recurring-payments", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","startDate":"2019-09-21","endDate":null,"amount":"100.0000000000","repeatEvery":"1","datesEditable":true,"decimalPrecision":"2"},
			{"id":2,"startDate":"2019-09-24","amount":200,"repeatEvery":1}
		]`))
	})

	payments, err := c.RecurringPayments().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, payments, 2)
	assert.Equal(t, ID("1"), payments[0].ID)
	assert.Equal(t, "2019-09-21", payments[0].StartDate)
	assert.Equal(t, Money(100), payments[0].Amount)
	// Numeric id normalizes to a string and numeric money parses.
	assert.Equal(t, ID("2"), payments[1].ID)
	assert.Equal(t, Money(200), payments[1].Amount)
}

func TestRecurringPayments_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/recurring-payments/1", r.URL.Path)
		_, _ = w.Write([]byte(`{
			"id":1,
			"startDate":"2019-09-21",
			"endDate":null,
			"amount":"100.0000000000",
			"account":{"id":"5","name":"Banco en dolares","type":"bank"},
			"repeatEvery":"1",
			"datesEditable":true,
			"decimalPrecision":"2",
			"currency":{"code":"AUD","symbol":"$","exchangeRate":2950}
		}`))
	})
	payment, err := c.RecurringPayments().Get(context.Background(), "1")
	require.NoError(t, err)
	assert.Equal(t, ID("1"), payment.ID)
	assert.Equal(t, Money(100), payment.Amount)
	require.NotNil(t, payment.Account)
	assert.Equal(t, ID("5"), payment.Account.ID)
	assert.Equal(t, "Banco en dolares", payment.Account.Name)
	assert.Equal(t, "bank", payment.Account.Type)
	require.NotNil(t, payment.Currency)
	assert.Equal(t, "AUD", payment.Currency.Code)
	assert.Equal(t, Money(2950), payment.Currency.ExchangeRate)
}
