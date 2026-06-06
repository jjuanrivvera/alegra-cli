package api

// Estimate is an Alegra estimate (cotización), a sales quote.
// See https://developer.alegra.com/reference/get_estimates
type Estimate struct {
	ID           ID                `json:"id,omitempty"`
	Number       string            `json:"number,omitempty"`
	Date         string            `json:"date,omitempty"`
	DueDate      string            `json:"dueDate,omitempty"`
	Status       string            `json:"status,omitempty"`
	Observations string            `json:"observations,omitempty"`
	Anotation    string            `json:"anotation,omitempty"`
	Total        Money             `json:"total,omitempty"`
	Client       *Ref              `json:"client,omitempty"`
	Seller       *Ref              `json:"seller,omitempty"`
	PriceList    *Ref              `json:"priceList,omitempty"`
	Warehouse    *Ref              `json:"warehouse,omitempty"`
	Currency     *EstimateCurrency `json:"currency,omitempty"`
	Items        []EstimateItem    `json:"items,omitempty"`
}

// EstimateCurrency is the currency block embedded in an estimate.
type EstimateCurrency struct {
	Code         string `json:"code,omitempty"`
	Symbol       string `json:"symbol,omitempty"`
	ExchangeRate Money  `json:"exchangeRate,omitempty"`
}

// EstimateItem is a line item on an estimate.
type EstimateItem struct {
	ID          ID     `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Reference   string `json:"reference,omitempty"`
	Price       Money  `json:"price,omitempty"`
	Discount    Money  `json:"discount,omitempty"`
	Quantity    Money  `json:"quantity,omitempty"`
	Total       Money  `json:"total,omitempty"`
	Tax         []Ref  `json:"tax,omitempty"`
}

// Estimates returns a typed handle to the /estimates resource.
func (c *Client) Estimates() *Resource[Estimate] {
	return NewResource[Estimate](c, "estimates")
}
