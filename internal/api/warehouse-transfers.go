package api

// WarehouseTransfer represents a movement of inventory items between two
// warehouses (origin -> destination) on a given date.
type WarehouseTransfer struct {
	ID           ID                      `json:"id,omitempty"`
	Date         string                  `json:"date,omitempty"`
	Observations string                  `json:"observations,omitempty"`
	Origin       *WarehouseTransferStore `json:"origin,omitempty"`
	Destination  *WarehouseTransferStore `json:"destination,omitempty"`
	Items        []WarehouseTransferItem `json:"items,omitempty"`
}

// WarehouseTransferStore is the origin or destination warehouse embedded in a
// transfer, including its cost center.
type WarehouseTransferStore struct {
	ID           ID     `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Observations string `json:"observations,omitempty"`
	Address      string `json:"address,omitempty"`
	Status       string `json:"status,omitempty"`
	IsDefault    bool   `json:"isDefault,omitempty"`
	CostCenter   *Ref   `json:"costCenter,omitempty"`
}

// WarehouseTransferItem is a single product/service line moved by a transfer.
type WarehouseTransferItem struct {
	ID       ID     `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Quantity Money  `json:"quantity,omitempty"`
}

// WarehouseTransfers returns a typed handle to the /warehouse-transfers resource.
func (c *Client) WarehouseTransfers() *Resource[WarehouseTransfer] {
	return NewResource[WarehouseTransfer](c, "warehouse-transfers")
}
