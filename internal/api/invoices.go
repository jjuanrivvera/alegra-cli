package api

import "encoding/json"

// Invoice is an Alegra sales invoice (factura de venta).
// See https://developer.alegra.com/reference/get_invoices
type Invoice struct {
	ID               ID               `json:"id,omitempty"`
	Date             string           `json:"date,omitempty"`
	DueDate          string           `json:"dueDate,omitempty"`
	Datetime         string           `json:"datetime,omitempty"`
	Status           string           `json:"status,omitempty"`
	Observations     string           `json:"observations,omitempty"`
	Anotation        string           `json:"anotation,omitempty"`
	TermsConditions  string           `json:"termsConditions,omitempty"`
	Total            Money            `json:"total,omitempty"`
	TotalPaid        Money            `json:"totalPaid,omitempty"`
	Balance          Money            `json:"balance,omitempty"`
	DecimalPrecision Int              `json:"decimalPrecision,omitempty"`
	Client           *Ref             `json:"client,omitempty"`
	Seller           *Ref             `json:"seller,omitempty"`
	PriceList        *Ref             `json:"priceList,omitempty"`
	Currency         *InvoiceCurrency `json:"currency,omitempty"`
	Retentions       []Ref            `json:"retentions,omitempty"`
	Items            []InvoiceItem    `json:"items,omitempty"`
	// NumberTemplate is kept raw because Alegra serializes it as either an
	// object or an array of objects depending on the record.
	NumberTemplate json.RawMessage `json:"numberTemplate,omitempty"`
}

// InvoiceCurrency is the currency block embedded in an invoice.
type InvoiceCurrency struct {
	Code         string `json:"code,omitempty"`
	Symbol       string `json:"symbol,omitempty"`
	ExchangeRate Money  `json:"exchangeRate,omitempty"`
}

// InvoiceItem is a line item on an invoice.
type InvoiceItem struct {
	ID          ID     `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Reference   string `json:"reference,omitempty"`
	Price       Money  `json:"price,omitempty"`
	Quantity    Money  `json:"quantity,omitempty"`
	Tax         []Ref  `json:"tax,omitempty"`
}

// Invoices returns a typed handle to the /invoices resource.
func (c *Client) Invoices() *Resource[Invoice] {
	return NewResource[Invoice](c, "invoices")
}
