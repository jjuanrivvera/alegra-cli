package api

// CostCenter is an Alegra cost center, used to group and track accounting
// movements (income and expenses) by business unit, project, or department.
type CostCenter struct {
	ID          ID     `json:"id,omitempty"`
	Code        string `json:"code,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"` // active, inactive
}

// CostCenters returns a typed handle to the /cost-centers resource.
func (c *Client) CostCenters() *Resource[CostCenter] {
	return NewResource[CostCenter](c, "cost-centers")
}
