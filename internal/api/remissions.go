package api

// Remission is an Alegra remission (remisión / delivery note).
// See https://developer.alegra.com/reference/get_remissions
type Remission struct {
	ID           ID              `json:"id,omitempty"`
	Number       string          `json:"number,omitempty"`
	DocumentName string          `json:"documentName,omitempty"`
	Date         string          `json:"date,omitempty"`
	DueDate      string          `json:"dueDate,omitempty"`
	Status       string          `json:"status,omitempty"`
	Observations string          `json:"observations,omitempty"`
	Anotation    string          `json:"anotation,omitempty"`
	Total        Money           `json:"total,omitempty"`
	Client       *Ref            `json:"client,omitempty"`
	Seller       *Ref            `json:"seller,omitempty"`
	PriceList    *Ref            `json:"priceList,omitempty"`
	Warehouse    *Ref            `json:"warehouse,omitempty"`
	Items        []RemissionItem `json:"items,omitempty"`
}

// RemissionItem is a line item on a remission.
type RemissionItem struct {
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

// Remissions returns a typed handle to the /remissions resource.
func (c *Client) Remissions() *Resource[Remission] {
	return NewResource[Remission](c, "remissions")
}
