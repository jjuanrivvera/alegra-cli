package api

// DebitNote is an Alegra debit note (nota de débito): a document that
// increases the amount a client owes against previously issued documents.
// See https://developer.alegra.com/reference/get_debit-notes-id
type DebitNote struct {
	ID               ID                `json:"id,omitempty"`
	Date             string            `json:"date,omitempty"`
	Status           string            `json:"status,omitempty"`
	Observations     string            `json:"observations,omitempty"`
	TermsConditions  string            `json:"termsConditions,omitempty"`
	DecimalPrecision Int               `json:"decimalPrecision,omitempty"`
	Total            Money             `json:"total,omitempty"`
	Balance          Money             `json:"balance,omitempty"`
	TotalApplied     Money             `json:"totalApplied,omitempty"`
	Client           *Ref              `json:"client,omitempty"`
	Warehouse        *Ref              `json:"warehouse,omitempty"`
	NumberTemplate   *DebitNoteNumber  `json:"numberTemplate,omitempty"`
	Items            []DebitNoteItem   `json:"items,omitempty"`
	Refunds          []DebitNoteRefund `json:"refunds,omitempty"`
}

// DebitNoteNumber is the numbering template metadata of a debit note.
type DebitNoteNumber struct {
	Number       string `json:"number,omitempty"`
	DocumentType string `json:"documentType,omitempty"`
}

// DebitNoteItem is a line item included in a debit note.
type DebitNoteItem struct {
	ID           ID     `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Reference    string `json:"reference,omitempty"`
	Description  string `json:"description,omitempty"`
	Price        Money  `json:"price,omitempty"`
	Quantity     Money  `json:"quantity,omitempty"`
	Observations string `json:"observations,omitempty"`
	Subtotal     Money  `json:"subtotal,omitempty"`
	Total        Money  `json:"total,omitempty"`
}

// DebitNoteRefund is a refund associated with a debit note.
type DebitNoteRefund struct {
	ID           ID     `json:"id,omitempty"`
	Number       string `json:"number,omitempty"`
	Date         string `json:"date,omitempty"`
	Amount       Money  `json:"amount,omitempty"`
	Type         string `json:"type,omitempty"`
	Status       string `json:"status,omitempty"`
	Observations string `json:"observations,omitempty"`
	BankAccount  *Ref   `json:"bankAccount,omitempty"`
}

// DebitNotes returns a typed handle to the /debit-notes resource.
func (c *Client) DebitNotes() *Resource[DebitNote] {
	return NewResource[DebitNote](c, "debit-notes")
}
