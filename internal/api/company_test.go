package api

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompany_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/company", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name": "Empresa Ejemplar",
			"identification": "1234",
			"email": "empresa.ejemplar@example.com",
			"regime": "Responsable de IVA",
			"applicationVersion": "colombia",
			"address": {"city": "Bogotá", "description": "Calle 10 # 10-10"},
			"currency": {"code": "COP", "symbol": "$"},
			"identificationObject": {"type": "NIT", "number": "1234", "dv": "1"}
		}`))
	})

	var company Company
	err := c.GetInto(context.Background(), "company", nil, &company)
	require.NoError(t, err)
	assert.Equal(t, "Empresa Ejemplar", company.Name)
	assert.Equal(t, "colombia", company.ApplicationVersion)
	require.NotNil(t, company.Address)
	assert.Equal(t, "Bogotá", company.Address.City)
	require.NotNil(t, company.Currency)
	assert.Equal(t, "COP", company.Currency.Code)
	require.NotNil(t, company.IdentificationObject)
	assert.Equal(t, "NIT", company.IdentificationObject.Type)
}

func TestCompany_Update(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/company", r.URL.Path)
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "Nueva Empresa", body["name"])
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"Nueva Empresa","email":"new@example.com"}`))
	})

	var company Company
	err := c.PutInto(context.Background(), "company", map[string]any{"name": "Nueva Empresa"}, &company)
	require.NoError(t, err)
	assert.Equal(t, "Nueva Empresa", company.Name)
	assert.Equal(t, "new@example.com", company.Email)
}

func TestCompany_Accessor(t *testing.T) {
	c := newTestClient(t, func(http.ResponseWriter, *http.Request) {})
	assert.Equal(t, "company", c.Company().Path())
}
