package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsers_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users", r.URL.Path)
		_, _ = w.Write([]byte(`[{"id":"1","name":"Pedro","email":"pedro.perez@example.com","role":"admin","status":"active"}]`))
	})
	items, err := c.Users().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "Pedro", items[0].Name)
	assert.Equal(t, "admin", items[0].Role)
}

func TestUsers_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users/1", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":"1","username":"pedro.perez","name":"Pedro","lastName":"Perez","email":"pedro.perez@example.com","status":"active","position":"owner"}`))
	})
	item, err := c.Users().Get(context.Background(), "1")
	require.NoError(t, err)
	assert.Equal(t, ID("1"), item.ID)
	assert.Equal(t, "pedro.perez", item.Username)
	assert.Equal(t, "Perez", item.LastName)
	assert.Equal(t, "owner", item.Position)
}
