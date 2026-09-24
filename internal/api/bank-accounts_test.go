package api

import (
	"context"
	"encoding/json"
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

func TestBankTransfer_JSONRoundTrip(t *testing.T) {
	in := BankTransfer{
		IDDestination:         "2",
		Amount:                Money("150.00"),
		Date:                  "2026-09-12",
		Observations:          "move funds",
		ExchangeRate:          Money("1"),
		CostCenterOrigin:      "10",
		CostCenterDestination: "20",
		IDResolution:          "3",
		IDResolutionOut:       "4",
	}
	raw, err := json.Marshal(in)
	require.NoError(t, err)
	assert.Contains(t, string(raw), `"costCenterOrigin":"10"`)
	assert.Contains(t, string(raw), `"costCenterDestination":"20"`)
	assert.Contains(t, string(raw), `"idResolution":"3"`)
	assert.Contains(t, string(raw), `"idResolutionOut":"4"`)

	var out BankTransfer
	require.NoError(t, json.Unmarshal(raw, &out))
	assert.Equal(t, in, out)

	omitted, err := json.Marshal(BankTransfer{IDDestination: "2", Amount: Money("1")})
	require.NoError(t, err)
	assert.NotContains(t, string(omitted), "costCenterOrigin")
	assert.NotContains(t, string(omitted), "idResolution")
}
