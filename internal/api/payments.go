package api

import "encoding/json"

// Payment is an Alegra payment (a registered income "in" or expense "out").
// See https://developer.alegra.com/reference/get_payments
type Payment struct {
	ID            ID               `json:"id,omitempty"`
	Date          string           `json:"date,omitempty"`
	Number        string           `json:"number,omitempty"`
	Amount        Money            `json:"amount,omitempty"`
	Type          string           `json:"type,omitempty"`
	Status        string           `json:"status,omitempty"`
	PaymentMethod string           `json:"paymentMethod,omitempty"`
	Observations  string           `json:"observations,omitempty"`
	Anotation     string           `json:"anotation,omitempty"`
	Client        *Ref             `json:"client,omitempty"`
	BankAccount   *PaymentAccount  `json:"bankAccount,omitempty"`
	Currency      *PaymentCurrency `json:"currency,omitempty"`
	Invoices      []PaymentInvoice `json:"invoices,omitempty"`
	Bills         []PaymentInvoice `json:"bills,omitempty"`
	Categories    []Ref            `json:"categories,omitempty"`
	// NumberTemplate is kept raw because Alegra serializes it as either an
	// object or an array of objects depending on the record.
	NumberTemplate json.RawMessage `json:"numberTemplate,omitempty"`
}

// PaymentAccount is the bank or cash account a payment moves money through.
type PaymentAccount struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// PaymentCurrency is the currency block embedded in a payment.
type PaymentCurrency struct {
	Code         string `json:"code,omitempty"`
	Symbol       string `json:"symbol,omitempty"`
	ExchangeRate Money  `json:"exchangeRate,omitempty"`
}

// PaymentInvoice is an invoice or bill settled by a payment, with the paid amount.
type PaymentInvoice struct {
	ID        ID     `json:"id,omitempty"`
	Number    string `json:"number,omitempty"`
	Date      string `json:"date,omitempty"`
	Total     Money  `json:"total,omitempty"`
	TotalPaid Money  `json:"totalPaid,omitempty"`
	Balance   Money  `json:"balance,omitempty"`
	Amount    Money  `json:"amount,omitempty"`
}

// Payments returns a typed handle to the /payments resource.
func (c *Client) Payments() *Resource[Payment] {
	return NewResource[Payment](c, "payments")
}
