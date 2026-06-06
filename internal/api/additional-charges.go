package api

// AdditionalCharge is an Alegra additional charge (e.g. a tip or parafiscal
// contribution) that can be applied to sales and/or purchase documents.
type AdditionalCharge struct {
	ID                  ID     `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	Code                string `json:"code,omitempty"`
	Description         string `json:"description,omitempty"`
	Percentage          Money  `json:"percentage,omitempty"`
	CategoryPurchasesID ID     `json:"categoryPurchasesId,omitempty"`
	CategorySalesID     ID     `json:"categorySalesId,omitempty"`
	Status              string `json:"status,omitempty"` // active, inactive
	CompanyID           ID     `json:"companyId,omitempty"`
}

// AdditionalCharges returns a typed handle to the /additional-charges resource.
func (c *Client) AdditionalCharges() *Resource[AdditionalCharge] {
	return NewResource[AdditionalCharge](c, "additional-charges")
}
