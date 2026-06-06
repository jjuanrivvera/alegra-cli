package api

// TransportationReceipt is an Alegra transportation receipt (documento de
// traslado), a shipping/transfer document used for moving goods.
type TransportationReceipt struct {
	ID              ID                          `json:"id,omitempty"`
	Number          string                      `json:"number,omitempty"`
	Date            string                      `json:"date,omitempty"`
	DateShipping    string                      `json:"dateShipping,omitempty"`
	Observations    string                      `json:"observations,omitempty"`
	Anotation       string                      `json:"anotation,omitempty"`
	TermsConditions string                      `json:"termsConditions,omitempty"`
	Status          string                      `json:"status,omitempty"`
	Client          *Ref                        `json:"client,omitempty"`
	NumberTemplate  *Ref                        `json:"numberTemplate,omitempty"`
	Invoice         *Ref                        `json:"invoice,omitempty"`
	Seller          *Ref                        `json:"seller,omitempty"`
	CfdiUse         string                      `json:"cfdiUse,omitempty"`
	Regime          string                      `json:"regime,omitempty"`
	RegimeClient    string                      `json:"regimeClient,omitempty"`
	PriceList       *Ref                        `json:"priceList,omitempty"`
	Warehouse       *Ref                        `json:"warehouse,omitempty"`
	Items           []TransportationReceiptItem `json:"items,omitempty"`
}

// TransportationReceiptItem is a single line item of a transportation receipt.
type TransportationReceiptItem struct {
	ID          ID     `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Unit        string `json:"unit,omitempty"`
	Description string `json:"description,omitempty"`
	Reference   string `json:"reference,omitempty"`
	Quantity    string `json:"quantity,omitempty"`
	ProductKey  string `json:"productKey,omitempty"`
	Price       Money  `json:"price,omitempty"`
	Total       Money  `json:"total,omitempty"`
	Subtotal    Money  `json:"subtotal,omitempty"`
	Weight      string `json:"weight,omitempty"`
}

// TransportationReceipts returns a typed handle to the /transportation-receipts resource.
func (c *Client) TransportationReceipts() *Resource[TransportationReceipt] {
	return NewResource[TransportationReceipt](c, "transportation-receipts")
}
