package api

import "encoding/json"

// PurchaseOrder is an Alegra purchase order (orden de compra) issued to a provider.
// See https://developer.alegra.com/reference/get_purchase-orders
type PurchaseOrder struct {
	ID              ID                      `json:"id,omitempty"`
	Number          string                  `json:"number,omitempty"`
	FullNumber      string                  `json:"fullNumber,omitempty"`
	Date            string                  `json:"date,omitempty"`
	DeliveryDate    string                  `json:"deliveryDate,omitempty"`
	Status          string                  `json:"status,omitempty"`
	Observations    string                  `json:"observations,omitempty"`
	TermsConditions string                  `json:"termsConditions,omitempty"`
	Total           Money                   `json:"total,omitempty"`
	SubTotal        Money                   `json:"subtotal,omitempty"`
	Discount        Money                   `json:"discount,omitempty"`
	Tax             Money                   `json:"tax,omitempty"`
	Provider        *PurchaseOrderProvider  `json:"provider,omitempty"`
	Warehouse       *Ref                    `json:"warehouse,omitempty"`
	CostCenter      *Ref                    `json:"costCenter,omitempty"`
	Currency        *PurchaseOrderCurrency  `json:"currency,omitempty"`
	Purchases       *PurchaseOrderPurchases `json:"purchases,omitempty"`
	Bills           []Ref                   `json:"bills,omitempty"`
	// NumberTemplate is kept raw because Alegra serializes it as either an
	// object or an array of objects depending on the record.
	NumberTemplate json.RawMessage `json:"numberTemplate,omitempty"`
}

// PurchaseOrderProvider is the provider (supplier) block embedded in a purchase order.
type PurchaseOrderProvider struct {
	ID             ID     `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Identification string `json:"identification,omitempty"`
	Email          string `json:"email,omitempty"`
	PhonePrimary   string `json:"phonePrimary,omitempty"`
	Mobile         string `json:"mobile,omitempty"`
}

// PurchaseOrderCurrency is the currency block embedded in a purchase order.
type PurchaseOrderCurrency struct {
	Code         string `json:"code,omitempty"`
	Symbol       string `json:"symbol,omitempty"`
	ExchangeRate Money  `json:"exchangeRate,omitempty"`
}

// PurchaseOrderPurchases groups the items and categories on a purchase order.
type PurchaseOrderPurchases struct {
	Items      []PurchaseOrderItem     `json:"items,omitempty"`
	Categories []PurchaseOrderCategory `json:"categories,omitempty"`
}

// PurchaseOrderItem is a product line on a purchase order.
type PurchaseOrderItem struct {
	ID       ID     `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Price    Money  `json:"price,omitempty"`
	Discount Money  `json:"discount,omitempty"`
	Quantity Money  `json:"quantity,omitempty"`
	Subtotal Money  `json:"subtotal,omitempty"`
	Total    Money  `json:"total,omitempty"`
	Tax      []Ref  `json:"tax,omitempty"`
}

// PurchaseOrderCategory is an accounting category line on a purchase order.
type PurchaseOrderCategory struct {
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

// PurchaseOrders returns a typed handle to the /purchase-orders resource.
func (c *Client) PurchaseOrders() *Resource[PurchaseOrder] {
	return NewResource[PurchaseOrder](c, "purchase-orders")
}
