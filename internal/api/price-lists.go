package api

// PriceList is an Alegra price list. A price list either sets a specific amount
// (type "amount") or discounts a percentage off the regular price
// (type "percentage"), in which case Percentage is populated.
type PriceList struct {
	ID          ID     `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Percentage  Money  `json:"percentage,omitempty"`
	Status      string `json:"status,omitempty"`
}

// PriceLists returns a typed handle to the /price-lists resource.
func (c *Client) PriceLists() *Resource[PriceList] {
	return NewResource[PriceList](c, "price-lists")
}
