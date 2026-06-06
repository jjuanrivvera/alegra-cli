package api

// Category is an Alegra chart-of-accounts account (cuenta contable). These are
// accounting ledger accounts, not item/product categories.
// See https://developer.alegra.com/reference/categoriesview
type Category struct {
	ID          ID            `json:"id,omitempty"`
	IDParent    ID            `json:"idParent,omitempty"` // parent account, for nested hierarchies
	Name        string        `json:"name,omitempty"`
	Code        string        `json:"code,omitempty"`
	Description string        `json:"description,omitempty"`
	Type        string        `json:"type,omitempty"`   // income | expense | asset | liability | equity | cost | productionCost | order
	Nature      string        `json:"nature,omitempty"` // debit | credit
	Status      string        `json:"status,omitempty"` // active | inactive | deleted
	Blocked     string        `json:"blocked,omitempty"`
	ReadOnly    bool          `json:"readOnly,omitempty"`
	Rule        *CategoryRule `json:"categoryRule,omitempty"`
	Children    []Category    `json:"children,omitempty"` // child accounts, when requested in tree format
}

// CategoryRule is the accounting rule associated with a chart-of-accounts account.
type CategoryRule struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

// Categories returns a typed handle to the /categories resource
// (chart-of-accounts accounts / cuentas contables).
func (c *Client) Categories() *Resource[Category] {
	return NewResource[Category](c, "categories")
}
