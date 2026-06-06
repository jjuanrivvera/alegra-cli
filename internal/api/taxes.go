package api

// Tax is an Alegra tax definition (e.g. IVA) applied to items and documents.
// See https://developer.alegra.com/reference/get_taxes
type Tax struct {
	ID          ID     `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Percentage  Money  `json:"percentage,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"` // active | inactive
	Type        string `json:"type,omitempty"`   // e.g. IVA
	Rate        string `json:"rate,omitempty"`

	// CategoryFavorable is the accounting category where the tax is recorded in favor.
	CategoryFavorable *TaxCategory `json:"categoryFavorable,omitempty"`
	// CategoryToBePaid is the accounting category where the tax is recorded as payable.
	CategoryToBePaid *TaxCategory `json:"categoryToBePaid,omitempty"`
}

// TaxCategory is the accounting category associated with a tax (favorable or payable).
type TaxCategory struct {
	ID           ID               `json:"id,omitempty"`
	IDParent     ID               `json:"idParent,omitempty"`
	Name         string           `json:"name,omitempty"`
	Text         string           `json:"text,omitempty"`
	Code         string           `json:"code,omitempty"`
	Description  string           `json:"description,omitempty"`
	Type         string           `json:"type,omitempty"`   // asset | liability
	Nature       string           `json:"nature,omitempty"` // debit | credit
	Status       string           `json:"status,omitempty"`
	CategoryRule *TaxCategoryRule `json:"categoryRule,omitempty"`
}

// TaxCategoryRule is the rule applied to a tax accounting category.
type TaxCategoryRule struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

// Taxes returns a typed handle to the /taxes resource.
func (c *Client) Taxes() *Resource[Tax] {
	return NewResource[Tax](c, "taxes")
}
