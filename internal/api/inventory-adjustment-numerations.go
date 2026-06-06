package api

// InventoryAdjustmentNumeration is an Alegra numbering scheme (numeración) used
// to number inventory adjustments. It defines the prefix, whether the sequence
// auto-increments, and which numeration is the default.
// See https://developer.alegra.com/reference/get_inventory-adjustments-numerations
type InventoryAdjustmentNumeration struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	// Prefix prepended to each generated number (may be empty).
	Prefix string `json:"prefix,omitempty"`
	// Status of the numeration: active or inactive.
	Status string `json:"status,omitempty"`
	// AutoIncrement reports whether the sequence increments automatically.
	AutoIncrement bool `json:"autoIncrement,omitempty"`
	// StartNumber is the initial number of the sequence; nil for
	// non-autoincremental numerations.
	StartNumber *Int `json:"startNumber,omitempty"`
	// IsDefault marks the preferred numeration for inventory adjustments.
	IsDefault bool `json:"isDefault,omitempty"`
	// NextNumber is the next number to be assigned. Only present for
	// autoincremental numerations when requested via fields=nextNumber.
	NextNumber *Int `json:"nextNumber,omitempty"`
}

// InventoryAdjustmentNumerations returns a typed handle to the
// inventory-adjustments/numerations resource.
func (c *Client) InventoryAdjustmentNumerations() *Resource[InventoryAdjustmentNumeration] {
	return NewResource[InventoryAdjustmentNumeration](c, "inventory-adjustments/numerations")
}
