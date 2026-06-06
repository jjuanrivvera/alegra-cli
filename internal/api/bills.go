package api

import "encoding/json"

// Bill is an Alegra provider bill (factura de proveedor / gasto).
// See https://developer.alegra.com/reference/get_bills
type Bill struct {
	ID               ID             `json:"id,omitempty"`
	Date             string         `json:"date,omitempty"`
	DueDate          string         `json:"dueDate,omitempty"`
	Status           string         `json:"status,omitempty"`
	Observations     string         `json:"observations,omitempty"`
	TermsConditions  string         `json:"termsConditions,omitempty"`
	Total            Money          `json:"total,omitempty"`
	TotalPaid        Money          `json:"totalPaid,omitempty"`
	Balance          Money          `json:"balance,omitempty"`
	DecimalPrecision Int            `json:"decimalPrecision,omitempty"`
	Provider         *BillProvider  `json:"provider,omitempty"`
	Warehouse        *Ref           `json:"warehouse,omitempty"`
	Payments         []BillPayment  `json:"payments,omitempty"`
	Purchases        *BillPurchases `json:"purchases,omitempty"`
	Retentions       []Ref          `json:"retentions,omitempty"`
	// NumberTemplate is kept raw because Alegra serializes it as either an
	// object or an array of objects depending on the record.
	NumberTemplate json.RawMessage `json:"numberTemplate,omitempty"`
}

// BillProvider is the provider (supplier) block embedded in a bill.
type BillProvider struct {
	ID             ID     `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Identification string `json:"identification,omitempty"`
	Email          string `json:"email,omitempty"`
	PhonePrimary   string `json:"phonePrimary,omitempty"`
	PhoneSecondary string `json:"phoneSecondary,omitempty"`
	Mobile         string `json:"mobile,omitempty"`
	Fax            string `json:"fax,omitempty"`
}

// BillPayment is a payment applied to a bill.
type BillPayment struct {
	ID            ID     `json:"id,omitempty"`
	Prefix        string `json:"prefix,omitempty"`
	Number        string `json:"number,omitempty"`
	Date          string `json:"date,omitempty"`
	Amount        Money  `json:"amount,omitempty"`
	PaymentMethod string `json:"paymentMethod,omitempty"`
	Observations  string `json:"observations,omitempty"`
	Anotation     string `json:"anotation,omitempty"`
	Status        string `json:"status,omitempty"`
}

// BillPurchases groups the items and categories purchased on a bill.
type BillPurchases struct {
	Items      []BillItem     `json:"items,omitempty"`
	Categories []BillCategory `json:"categories,omitempty"`
}

// BillItem is a product line on a bill.
type BillItem struct {
	ID       ID     `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Price    Money  `json:"price,omitempty"`
	Discount Money  `json:"discount,omitempty"`
	Quantity Money  `json:"quantity,omitempty"`
	Subtotal Money  `json:"subtotal,omitempty"`
	Total    Money  `json:"total,omitempty"`
	Tax      []Ref  `json:"tax,omitempty"`
}

// BillCategory is an accounting category line on a bill.
type BillCategory struct {
	ID           ID     `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Price        Money  `json:"price,omitempty"`
	Discount     Money  `json:"discount,omitempty"`
	Quantity     Money  `json:"quantity,omitempty"`
	Subtotal     Money  `json:"subtotal,omitempty"`
	Total        Money  `json:"total,omitempty"`
	Observations string `json:"observations,omitempty"`
	Tax          []Ref  `json:"tax,omitempty"`
}

// Bills returns a typed handle to the /bills resource.
func (c *Client) Bills() *Resource[Bill] {
	return NewResource[Bill](c, "bills")
}
