package api

// Seller is an Alegra salesperson (vendedor) that can be assigned to sales
// documents and rolled up in reports (see `alegra reports sales-by-seller`).
type Seller struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	// Identification is the seller's document number. Alegra returns it as a
	// JSON number (not a string), so use the flexible ID type.
	Identification ID     `json:"identification,omitempty"`
	Observations   string `json:"observations,omitempty"`
	Status         string `json:"status,omitempty"`
}

// Sellers returns a typed handle to the /sellers resource.
func (c *Client) Sellers() *Resource[Seller] {
	return NewResource[Seller](c, "sellers")
}
