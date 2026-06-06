package api

// Item is an Alegra item: a product or service that can be sold or purchased.
// An item may be inventariable (carries an Inventory object) or a service, and
// can be a simple product, a kit, or a variant.
type Item struct {
	ID          ID             `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Reference   string         `json:"reference,omitempty"`
	Status      string         `json:"status,omitempty"`
	Type        string         `json:"type,omitempty"`
	Price       []ItemPrice    `json:"price,omitempty"`
	Tax         []ItemTax      `json:"tax,omitempty"`
	Category    *Ref           `json:"category,omitempty"`
	Inventory   *ItemInventory `json:"inventory,omitempty"`
}

// ItemPrice is one entry in an item's price list array. A bare item price is
// represented by a single entry carrying only Price.
type ItemPrice struct {
	IDPriceList ID     `json:"idPriceList,omitempty"`
	Name        string `json:"name,omitempty"`
	Price       Money  `json:"price,omitempty"`
}

// ItemTax is a tax associated with an item.
type ItemTax struct {
	ID          ID     `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Percentage  Money  `json:"percentage,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

// ItemInventory holds inventory data for inventariable items. Its presence
// indicates the item is tracked in stock; absence implies a service.
type ItemInventory struct {
	Unit              string `json:"unit,omitempty"`
	AvailableQuantity Money  `json:"availableQuantity,omitempty"`
	UnitCost          Money  `json:"unitCost,omitempty"`
	InitialQuantity   Money  `json:"initialQuantity,omitempty"`
}

// Items returns a typed handle to the /items resource.
func (c *Client) Items() *Resource[Item] {
	return NewResource[Item](c, "items")
}
