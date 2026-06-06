package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemissions_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/remissions", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","number":"1","date":"2024-08-30","status":"open","total":110000,"client":{"id":"1","name":"Usuario Alegra"}},
			{"id":2,"number":"2","date":"2019-08-30","status":"open","total":110000}
		]`))
	})

	remissions, err := c.Remissions().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, remissions, 2)
	assert.Equal(t, ID("1"), remissions[0].ID)
	assert.Equal(t, "open", remissions[0].Status)
	require.NotNil(t, remissions[0].Client)
	assert.Equal(t, "Usuario Alegra", remissions[0].Client.Name)
	// Numeric id is normalized to a string.
	assert.Equal(t, ID("2"), remissions[1].ID)
}

func TestRemissions_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/remissions/1", r.URL.Path)
		_, _ = w.Write([]byte(`{
			"id":1,
			"number":"1",
			"date":"2024-08-30",
			"status":"open",
			"total":102000,
			"warehouse":{"id":"1","name":"Principal"},
			"items":[
				{"id":"1","name":"Item 1","price":100000,"quantity":"1.00","total":100000},
				{"id":2,"name":"Item 2","price":2000,"quantity":"1.00","total":2000}
			]
		}`))
	})

	remission, err := c.Remissions().Get(context.Background(), "1")
	require.NoError(t, err)
	assert.Equal(t, ID("1"), remission.ID)
	assert.Equal(t, "open", remission.Status)
	require.NotNil(t, remission.Warehouse)
	assert.Equal(t, "Principal", remission.Warehouse.Name)
	require.Len(t, remission.Items, 2)
	assert.Equal(t, ID("1"), remission.Items[0].ID)
	assert.Equal(t, "Item 1", remission.Items[0].Name)
	assert.Equal(t, ID("2"), remission.Items[1].ID)
}
