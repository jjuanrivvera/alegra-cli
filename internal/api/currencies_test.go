package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrencies_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/currencies", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"name":"Peso colombiano","code":"COP","symbol":"$","status":"active","exchangeRate":0},
			{"name":"Dólar estadounidense","code":"USD","symbol":"$","status":"active","exchangeRate":4000}
		]`))
	})
	items, err := c.Currencies().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, "COP", items[0].Code)
	assert.Equal(t, "Peso colombiano", items[0].Name)
	assert.Equal(t, "USD", items[1].Code)
	assert.Equal(t, Money("4000"), items[1].ExchangeRate)
}

func TestCurrencies_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/currencies/USD", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"Dólar estadounidense","code":"USD","symbol":"$","status":"active","exchangeRate":4000}`))
	})
	cur, err := c.Currencies().Get(context.Background(), "USD")
	require.NoError(t, err)
	require.NotNil(t, cur)
	assert.Equal(t, "USD", cur.Code)
	assert.Equal(t, "Dólar estadounidense", cur.Name)
	assert.Equal(t, "$", cur.Symbol)
	assert.Equal(t, "active", cur.Status)
	assert.Equal(t, Money("4000"), cur.ExchangeRate)
}
