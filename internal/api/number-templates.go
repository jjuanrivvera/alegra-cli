package api

// NumberTemplate is an Alegra document numbering (numeración de facturación):
// the prefix, ranges and resolution used to number invoices, estimates,
// credit/debit notes and similar documents.
// See https://developer.alegra.com/reference/get_number-templates-1
type NumberTemplate struct {
	ID     ID     `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Prefix string `json:"prefix,omitempty"`
	Status string `json:"status,omitempty"` // active | inactive

	// DocumentType is the document this numbering applies to, e.g.
	// invoice, estimate, creditNote, debitNote, incomeDebitNote,
	// transactionIn or transactionOut.
	DocumentType string `json:"documentType,omitempty"`
	// SubDocumentType is the country-specific document subtype (e.g. INVOICE_A in AR).
	SubDocumentType string `json:"subDocumentType,omitempty"`
	InvoiceText     string `json:"invoiceText,omitempty"`

	// IsDefault marks the preferred numbering for its document type.
	IsDefault bool `json:"isDefault,omitempty"`
	// Autoincrement indicates whether the numbering increments automatically.
	Autoincrement bool `json:"autoincrement,omitempty"`
	// IsElectronic indicates whether the numbering is used for electronic invoicing.
	IsElectronic bool `json:"isElectronic,omitempty"`
	// Deletable indicates whether the numbering can be deleted.
	Deletable bool `json:"deletable,omitempty"`

	NextInvoiceNumber Int `json:"nextInvoiceNumber,omitempty"`
	MinInvoiceNumber  Int `json:"minInvoiceNumber,omitempty"`
	MaxInvoiceNumber  Int `json:"maxInvoiceNumber,omitempty"`

	ResolutionNumber string `json:"resolutionNumber,omitempty"`
	StartDate        string `json:"startDate,omitempty"`
	EndDate          string `json:"endDate,omitempty"`
	BranchOffice     string `json:"branchOffice,omitempty"`
}

// NumberTemplates returns a typed handle to the /number-templates resource.
func (c *Client) NumberTemplates() *Resource[NumberTemplate] {
	return NewResource[NumberTemplate](c, "number-templates")
}
