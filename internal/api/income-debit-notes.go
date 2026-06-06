package api

// IncomeDebitNote is an Alegra customer debit note (nota débito cliente): a
// document that increases the amount a client owes against a previously issued
// sales invoice.
// See https://developer.alegra.com/reference/get_income-debit-notes-id
type IncomeDebitNote struct {
	ID               ID                     `json:"id,omitempty"`
	Date             string                 `json:"date,omitempty"`
	Status           string                 `json:"status,omitempty"`
	Type             string                 `json:"type,omitempty"`
	Subtotal         Money                  `json:"subtotal,omitempty"`
	Discount         Money                  `json:"discount,omitempty"`
	Tax              Money                  `json:"tax,omitempty"`
	Total            Money                  `json:"total,omitempty"`
	Balance          Money                  `json:"balance,omitempty"`
	DecimalPrecision Int                    `json:"decimalPrecision,omitempty"`
	PaymentMethod    string                 `json:"paymentMethod,omitempty"`
	PaymentType      string                 `json:"paymentType,omitempty"`
	Client           *Ref                   `json:"client,omitempty"`
	Warehouse        *Ref                   `json:"warehouse,omitempty"`
	PriceList        *Ref                   `json:"priceList,omitempty"`
	NumberTemplate   *IncomeDebitNoteNumber `json:"numberTemplate,omitempty"`
	Items            []IncomeDebitNoteItem  `json:"items,omitempty"`
	CostCenter       []Ref                  `json:"costCenter,omitempty"`
}

// IncomeDebitNoteNumber is the numbering template metadata of a debit note.
type IncomeDebitNoteNumber struct {
	ID                     ID     `json:"id,omitempty"`
	Prefix                 string `json:"prefix,omitempty"`
	Number                 Int    `json:"number,omitempty"`
	IsElectronicResolution bool   `json:"isElectronicResolution,omitempty"`
}

// IncomeDebitNoteItem is a line item included in a customer debit note.
type IncomeDebitNoteItem struct {
	ID          ID                   `json:"id,omitempty"`
	Name        string               `json:"name,omitempty"`
	Reference   string               `json:"reference,omitempty"`
	Description string               `json:"description,omitempty"`
	Price       Money                `json:"price,omitempty"`
	Discount    Money                `json:"discount,omitempty"`
	Quantity    Money                `json:"quantity,omitempty"`
	Total       Money                `json:"total,omitempty"`
	Tax         []IncomeDebitNoteTax `json:"tax,omitempty"`
}

// IncomeDebitNoteTax is a tax applied to a debit note line item.
type IncomeDebitNoteTax struct {
	ID          ID     `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Percentage  Money  `json:"percentage,omitempty"`
}

// IncomeDebitNotes returns a typed handle to the /income-debit-notes resource.
func (c *Client) IncomeDebitNotes() *Resource[IncomeDebitNote] {
	return NewResource[IncomeDebitNote](c, "income-debit-notes")
}
