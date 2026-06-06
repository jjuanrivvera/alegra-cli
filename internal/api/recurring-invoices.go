package api

// RecurringInvoice is an Alegra recurring invoice (factura recurrente): a
// template that automatically generates sales invoices on a schedule.
type RecurringInvoice struct {
	ID               ID                        `json:"id,omitempty"`
	StartDate        string                    `json:"startDate,omitempty"`
	EndDate          string                    `json:"endDate,omitempty"`
	RepeatEvery      Int                       `json:"repeatEvery,omitempty"`
	DueDate          string                    `json:"dueDate,omitempty"`
	Term             Int                       `json:"term,omitempty"`
	Observations     string                    `json:"observations,omitempty"`
	TermsConditions  string                    `json:"termsConditions,omitempty"`
	Anotation        string                    `json:"anotation,omitempty"`
	DecimalPrecision Int                       `json:"decimalPrecision,omitempty"`
	CalculationScale Int                       `json:"calculationScale,omitempty"`
	Status           string                    `json:"status,omitempty"`
	Warehouse        *Ref                      `json:"warehouse,omitempty"`
	NumberTemplate   *RecurringNumberTemplate  `json:"numberTemplate,omitempty"`
	PriceList        *Ref                      `json:"priceList,omitempty"`
	CostCenter       *Ref                      `json:"costCenter,omitempty"`
	Currency         *RecurringInvoiceCurrency `json:"currency,omitempty"`
	Client           *Ref                      `json:"client,omitempty"`
	Total            Money                     `json:"total,omitempty"`
	NextCreation     string                    `json:"nextCreation,omitempty"`
	LastCreation     string                    `json:"lastCreation,omitempty"`
	Items            []RecurringInvoiceItem    `json:"items,omitempty"`
}

// RecurringNumberTemplate is the numbering template a recurring invoice uses,
// including the document type it emits.
type RecurringNumberTemplate struct {
	ID           ID     `json:"id,omitempty"`
	DocumentType string `json:"documentType,omitempty"`
}

// RecurringInvoiceCurrency is the currency in which a recurring invoice is
// denominated.
type RecurringInvoiceCurrency struct {
	Code   string `json:"code,omitempty"`
	Symbol string `json:"symbol,omitempty"`
}

// RecurringInvoiceItem is a line item on a recurring invoice.
type RecurringInvoiceItem struct {
	ID          ID                    `json:"id,omitempty"`
	ItemID      ID                    `json:"itemId,omitempty"`
	Name        string                `json:"name,omitempty"`
	Description string                `json:"description,omitempty"`
	Reference   string                `json:"reference,omitempty"`
	Price       Money                 `json:"price,omitempty"`
	Quantity    Money                 `json:"quantity,omitempty"`
	Discount    Money                 `json:"discount,omitempty"`
	Total       Money                 `json:"total,omitempty"`
	Taxes       []RecurringInvoiceTax `json:"taxes,omitempty"`
}

// RecurringInvoiceTax is a tax applied to a recurring invoice line item.
type RecurringInvoiceTax struct {
	ID         ID     `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Percentage Money  `json:"percentage,omitempty"`
	Amount     Money  `json:"amount,omitempty"`
}

// RecurringInvoices returns a typed handle to the /recurring-invoices resource.
func (c *Client) RecurringInvoices() *Resource[RecurringInvoice] {
	return NewResource[RecurringInvoice](c, "recurring-invoices")
}
