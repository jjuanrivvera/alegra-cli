package api

// CustomField is an Alegra custom field (campo adicional) attached to a
// resource such as an item, letting companies store extra typed data.
type CustomField struct {
	ID           ID                   `json:"id,omitempty"`
	Name         string               `json:"name,omitempty"`
	Description  string               `json:"description,omitempty"`
	DefaultValue string               `json:"defaultValue,omitempty"`
	ResourceType string               `json:"resourceType,omitempty"`
	Status       string               `json:"status,omitempty"`
	Key          string               `json:"key,omitempty"`
	Type         string               `json:"type,omitempty"`
	Options      []string             `json:"options,omitempty"`
	Settings     *CustomFieldSettings `json:"settings,omitempty"`
}

// CustomFieldSettings holds the behavioral configuration of a custom field.
type CustomFieldSettings struct {
	IsRequired         bool `json:"isRequired,omitempty"`
	PrintOnInvoices    bool `json:"printOnInvoices,omitempty"`
	ShowInItemVariants bool `json:"showInItemVariants,omitempty"`
}

// CustomFields returns a typed handle to the /custom-fields resource.
func (c *Client) CustomFields() *Resource[CustomField] {
	return NewResource[CustomField](c, "custom-fields")
}
