package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return New(
		WithBaseURL(srv.URL),
		WithBasicAuth("test@example.com", "token"),
		WithRequestsPerSecond(1000),
	)
}

func TestContacts_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/contacts", r.URL.Path)
		assert.Equal(t, "30", r.URL.Query().Get("limit"))
		// Basic auth header is present.
		user, pass, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "test@example.com", user)
		assert.Equal(t, "token", pass)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","name":"Coorporación Alegrate","email":"a@x.com","type":["client"]},
			{"id":2,"name":"Proveedor","type":["provider"]}
		]`))
	})

	contacts, err := c.Contacts().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, contacts, 2)
	assert.Equal(t, ID("1"), contacts[0].ID)
	assert.Equal(t, "Coorporación Alegrate", contacts[0].Name)
	// Numeric id is normalized to a string.
	assert.Equal(t, ID("2"), contacts[1].ID)
	assert.Equal(t, StringOrSlice{"provider"}, contacts[1].Type)
}

func TestContacts_TypeStringOrArray(t *testing.T) {
	// Alegra returns "type" as a bare string in simple mode and an array in
	// advanced mode; both must decode.
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"1","type":"client"},{"id":"2","type":["client","provider"]}]`))
	})
	contacts, err := c.Contacts().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, contacts, 2)
	assert.Equal(t, StringOrSlice{"client"}, contacts[0].Type)
	assert.Equal(t, StringOrSlice{"client", "provider"}, contacts[1].Type)
}

func TestContacts_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/contacts/5", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":5,"name":"Acme","address":{"city":"Cali"}}`))
	})
	contact, err := c.Contacts().Get(context.Background(), "5")
	require.NoError(t, err)
	assert.Equal(t, ID("5"), contact.ID)
	require.NotNil(t, contact.Address)
	assert.Equal(t, "Cali", contact.Address.City)
}

func TestContacts_Create(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/contacts", r.URL.Path)
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "New Co", body["name"])
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":99,"name":"New Co"}`))
	})
	created, err := c.Contacts().Create(context.Background(), map[string]any{"name": "New Co"})
	require.NoError(t, err)
	assert.Equal(t, ID("99"), created.ID)
}

func TestContacts_Delete(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/contacts/7", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	})
	require.NoError(t, c.Contacts().Delete(context.Background(), "7"))
}

func TestAPIError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"message":"Validation failed","code":"1001"}`))
	})
	_, err := c.Contacts().Get(context.Background(), "1")
	require.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 422, apiErr.StatusCode)
	assert.Equal(t, "Validation failed", apiErr.Message)
}
