package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPurchaseOrders_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/purchase-orders", r.URL.Path)
		_, _ = w.Write([]byte(`[{"id":"1","date":"2023-05-10","status":"open","total":1500.5,"provider":{"id":"7","name":"Coorporación Alegrate"}}]`))
	})
	items, err := c.PurchaseOrders().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "open", items[0].Status)
	assert.Equal(t, Money("1500.5"), items[0].Total)
	require.NotNil(t, items[0].Provider)
	assert.Equal(t, "Coorporación Alegrate", items[0].Provider.Name)
}

func TestPurchaseOrders_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/purchase-orders/1", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":"1","date":"2023-05-10","deliveryDate":"2023-05-20","status":"open","total":1500.5,"subtotal":1300,"provider":{"id":"7","name":"Coorporación Alegrate","identification":"159.549.847"},"purchases":{"items":[{"id":"99","name":"Cable","price":100,"quantity":13}]}}`))
	})
	po, err := c.PurchaseOrders().Get(context.Background(), "1")
	require.NoError(t, err)
	require.NotNil(t, po)
	assert.Equal(t, ID("1"), po.ID)
	assert.Equal(t, "open", po.Status)
	assert.Equal(t, "2023-05-20", po.DeliveryDate)
	assert.Equal(t, Money("1500.5"), po.Total)
	require.NotNil(t, po.Provider)
	assert.Equal(t, "159.549.847", po.Provider.Identification)
	require.NotNil(t, po.Purchases)
	require.Len(t, po.Purchases.Items, 1)
	assert.Equal(t, ID("99"), po.Purchases.Items[0].ID)
}
