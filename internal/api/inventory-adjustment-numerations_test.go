package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventoryAdjustmentNumerations_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/inventory-adjustments/numerations", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"01HVVJWB1SJJNCG1W4KM7JSZ3A","name":"Ajuste de inventario","status":"active","autoIncrement":true,"prefix":null,"startNumber":1,"isDefault":true},
			{"id":"01HVXJWB1SJH2CG1X4KG7ZSLM1","name":"No autoincrement","status":"active","autoIncrement":false,"prefix":"PR","startNumber":null,"isDefault":false}
		]`))
	})

	nums, err := c.InventoryAdjustmentNumerations().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, nums, 2)
	assert.Equal(t, ID("01HVVJWB1SJJNCG1W4KM7JSZ3A"), nums[0].ID)
	assert.Equal(t, "Ajuste de inventario", nums[0].Name)
	assert.True(t, nums[0].AutoIncrement)
	assert.True(t, nums[0].IsDefault)
	require.NotNil(t, nums[0].StartNumber)
	assert.Equal(t, Int(1), *nums[0].StartNumber)

	assert.Equal(t, "PR", nums[1].Prefix)
	assert.False(t, nums[1].AutoIncrement)
	assert.Nil(t, nums[1].StartNumber)
}

func TestInventoryAdjustmentNumerations_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/inventory-adjustments/numerations/01HVXJWB1SJH2CG1X4KG7ZSLM1", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":"01HVXJWB1SJH2CG1X4KG7ZSLM1","name":"No autoincrement","status":"active","autoIncrement":false,"prefix":"PR","startNumber":null,"isDefault":false}`))
	})
	num, err := c.InventoryAdjustmentNumerations().Get(context.Background(), "01HVXJWB1SJH2CG1X4KG7ZSLM1")
	require.NoError(t, err)
	assert.Equal(t, ID("01HVXJWB1SJH2CG1X4KG7ZSLM1"), num.ID)
	assert.Equal(t, "No autoincrement", num.Name)
	assert.Equal(t, "active", num.Status)
	assert.Equal(t, "PR", num.Prefix)
}
