package api

// InventoryAdjustment represents a manual correction of stock levels in a
// warehouse, increasing ("in") or decreasing ("out") item quantities on a
// given date.
type InventoryAdjustment struct {
	ID           ID                        `json:"id,omitempty"`
	Date         string                    `json:"date,omitempty"`
	Observations string                    `json:"observations,omitempty"`
	Warehouse    *InventoryAdjustmentStore `json:"warehouse,omitempty"`
	Items        []InventoryAdjustmentItem `json:"items,omitempty"`
	CostCenter   *Ref                      `json:"costCenter,omitempty"`
	IDResolution ID                        `json:"idResolution,omitempty"`
	Prefix       string                    `json:"prefix,omitempty"`
	Number       string                    `json:"number,omitempty"`
}

// InventoryAdjustmentStore is the warehouse embedded in an adjustment.
type InventoryAdjustmentStore struct {
	ID           ID     `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Observations string `json:"observations,omitempty"`
	Address      string `json:"address,omitempty"`
	Status       string `json:"status,omitempty"`
	IsDefault    bool   `json:"isDefault,omitempty"`
	CostCenter   *Ref   `json:"costCenter,omitempty"`
}

// InventoryAdjustmentItem is a single product/service line of an adjustment.
// Type is "in" (increase) or "out" (decrease).
type InventoryAdjustmentItem struct {
	ID        ID     `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Quantity  Money  `json:"quantity,omitempty"`
	Type      string `json:"type,omitempty"`
	UnitCost  Money  `json:"unitCost,omitempty"`
	Reference string `json:"reference,omitempty"`
}

// InventoryAdjustments returns a typed handle to the /inventory-adjustments resource.
func (c *Client) InventoryAdjustments() *Resource[InventoryAdjustment] {
	return NewResource[InventoryAdjustment](c, "inventory-adjustments")
}
