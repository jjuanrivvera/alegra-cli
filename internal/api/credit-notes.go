package api

// CreditNote is an Alegra credit note (nota de crédito): a document that
// reduces or cancels a previously issued sales invoice.
// See https://developer.alegra.com/reference/get_credit-notes-id
type CreditNote struct {
	ID              ID                  `json:"id,omitempty"`
	Date            string              `json:"date,omitempty"`
	Status          string              `json:"status,omitempty"`
	Observations    string              `json:"observations,omitempty"`
	Anotation       string              `json:"anotation,omitempty"`
	TermsConditions string              `json:"termsConditions,omitempty"`
	Type            string              `json:"type,omitempty"`
	Total           Money               `json:"total,omitempty"`
	Balance         Money               `json:"balance,omitempty"`
	TotalApplied    Money               `json:"totalApplied,omitempty"`
	Client          *Ref                `json:"client,omitempty"`
	Warehouse       *Ref                `json:"warehouse,omitempty"`
	CostCenter      *Ref                `json:"costCenter,omitempty"`
	NumberTemplate  *CreditNoteNumber   `json:"numberTemplate,omitempty"`
	Items           []CreditNoteItem    `json:"items,omitempty"`
	Invoices        []CreditNoteInvoice `json:"invoices,omitempty"`
	Refunds         []CreditNoteRefund  `json:"refunds,omitempty"`
}

// CreditNoteNumber is the numbering template metadata of a credit note.
type CreditNoteNumber struct {
	Number       string `json:"number,omitempty"`
	DocumentType string `json:"documentType,omitempty"`
}

// CreditNoteItem is a line item included in a credit note.
type CreditNoteItem struct {
	ID          ID     `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Reference   string `json:"reference,omitempty"`
	Description string `json:"description,omitempty"`
	Price       Money  `json:"price,omitempty"`
	Discount    Money  `json:"discount,omitempty"`
	Quantity    Money  `json:"quantity,omitempty"`
	Total       Money  `json:"total,omitempty"`
}

// CreditNoteInvoice is an invoice that a credit note is applied against.
type CreditNoteInvoice struct {
	ID      ID     `json:"id,omitempty"`
	Prefix  string `json:"prefix,omitempty"`
	Number  string `json:"number,omitempty"`
	Date    string `json:"date,omitempty"`
	DueDate string `json:"dueDate,omitempty"`
	Amount  Money  `json:"amount,omitempty"`
	Total   Money  `json:"total,omitempty"`
	Balance Money  `json:"balance,omitempty"`
}

// CreditNoteRefund is a refund associated with a credit note.
type CreditNoteRefund struct {
	ID           ID     `json:"id,omitempty"`
	Number       string `json:"number,omitempty"`
	Date         string `json:"date,omitempty"`
	Amount       Money  `json:"amount,omitempty"`
	Account      *Ref   `json:"account,omitempty"`
	Observations string `json:"observations,omitempty"`
}

// CreditNotes returns a typed handle to the /credit-notes resource.
func (c *Client) CreditNotes() *Resource[CreditNote] {
	return NewResource[CreditNote](c, "credit-notes")
}
