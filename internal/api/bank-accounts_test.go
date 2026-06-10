package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBankAccounts_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bank-accounts", r.URL.Path)
		_, _ = w.Write([]byte(`[{"id":"3","name":"Banco 1","number":"100294","type":"bank","status":"active","initialBalance":200}]`))
	})
	items, err := c.BankAccounts().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ID("3"), items[0].ID)
	assert.Equal(t, "bank", items[0].Type)
	assert.Equal(t, Money("200"), items[0].InitialBalance)
}

func TestBankAccounts_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bank-accounts/12", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":"12","name":"Banco X","number":"1234 56789","type":"bank","status":"active","initialBalance":200,"initialBalanceDate":"2022-10-03"}`))
	})
	acc, err := c.BankAccounts().Get(context.Background(), "12")
	require.NoError(t, err)
	require.NotNil(t, acc)
	assert.Equal(t, ID("12"), acc.ID)
	assert.Equal(t, "Banco X", acc.Name)
	assert.Equal(t, "2022-10-03", acc.InitialBalanceDate)
}
