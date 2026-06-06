package api

// Warehouse is an Alegra inventory warehouse (bodega).
type Warehouse struct {
	ID           ID                   `json:"id,omitempty"`
	Name         string               `json:"name,omitempty"`
	Observations string               `json:"observations,omitempty"`
	Address      string               `json:"address,omitempty"`
	IsDefault    bool                 `json:"isDefault,omitempty"`
	Status       string               `json:"status,omitempty"`
	CostCenter   *WarehouseCostCenter `json:"costCenter,omitempty"`
}

// WarehouseCostCenter is the cost center associated with a warehouse.
type WarehouseCostCenter struct {
	ID          ID     `json:"id,omitempty"`
	Code        string `json:"code,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

// Warehouses returns a typed handle to the /warehouses resource.
func (c *Client) Warehouses() *Resource[Warehouse] {
	return NewResource[Warehouse](c, "warehouses")
}
