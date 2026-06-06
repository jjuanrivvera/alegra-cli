package api

// VariantAttribute is an Alegra item variant attribute (e.g. "Color" or "Talla")
// together with its selectable options.
type VariantAttribute struct {
	ID      ID              `json:"id,omitempty"`
	Name    string          `json:"name,omitempty"`
	Status  string          `json:"status,omitempty"`
	Options []VariantOption `json:"options,omitempty"`
}

// VariantOption is a single selectable value for a variant attribute
// (e.g. "Rojo" for the "Color" attribute).
type VariantOption struct {
	ID    ID     `json:"id,omitempty"`
	Value string `json:"value,omitempty"`
}

// VariantAttributes returns a typed handle to the /variant-attributes resource.
func (c *Client) VariantAttributes() *Resource[VariantAttribute] {
	return NewResource[VariantAttribute](c, "variant-attributes")
}
