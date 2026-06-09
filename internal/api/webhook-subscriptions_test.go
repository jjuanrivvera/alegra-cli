package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhookSubscriptions_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/webhooks/subscriptions", r.URL.Path)
		// Alegra wraps the list under "subscriptions" (not data/results/rows).
		_, _ = w.Write([]byte(`{"subscriptions":[{"id":"3","event":"new-invoice","url":"https://x.co/hook","status":"active"}]}`))
	})
	items, err := c.WebhookSubscriptions().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ID("3"), items[0].ID)
	assert.Equal(t, "new-invoice", items[0].Event)
	assert.Equal(t, "https://x.co/hook", items[0].URL)
}

func TestWebhookSubscriptions_Create(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/webhooks/subscriptions", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":"3","event":"new-invoice","url":"https://x.co/hook"}`))
	})
	out, err := c.WebhookSubscriptions().Create(context.Background(),
		map[string]any{"event": "new-invoice", "url": "https://x.co/hook"})
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, ID("3"), out.ID)
}
