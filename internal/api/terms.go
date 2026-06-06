package api

// Term is an Alegra payment term (término de pago): the agreed number of days
// a client has to pay an invoice.
type Term struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Days Int    `json:"days,omitempty"`
	// Status is not part of every API response but is kept for convenient
	// table rendering when present.
	Status string `json:"status,omitempty"`
	// Deletable is returned only when the request asks for fields=deletable.
	Deletable *bool `json:"deletable,omitempty"`
}

// Terms returns a typed handle to the /terms resource (payment terms).
func (c *Client) Terms() *Resource[Term] {
	return NewResource[Term](c, "terms")
}
