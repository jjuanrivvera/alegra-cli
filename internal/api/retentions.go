package api

// Retention is an Alegra retention (withholding) configuration.
// See https://developer.alegra.com/reference/get_retentions
type Retention struct {
	ID                   ID     `json:"id,omitempty"`
	IDRetentionReference ID     `json:"idRetentionReference,omitempty"`
	Name                 string `json:"name,omitempty"`
	CalculatedBy         string `json:"calculatedBy,omitempty"` // "percentage" or "fixedAmount"
	Percentage           Money  `json:"percentage,omitempty"`
	Description          string `json:"description,omitempty"`
	Type                 string `json:"type,omitempty"`   // e.g. "FUENTE"
	Status               string `json:"status,omitempty"` // "active" or "inactive"
}

// Retentions returns a typed handle to the /retentions resource.
func (c *Client) Retentions() *Resource[Retention] {
	return NewResource[Retention](c, "retentions")
}
