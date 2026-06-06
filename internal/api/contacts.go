package api

// Contact is an Alegra contact (a client and/or provider).
// See https://developer.alegra.com/reference/listcontacts-1
type Contact struct {
	ID             ID              `json:"id,omitempty"`
	Name           string          `json:"name,omitempty"`
	Identification string          `json:"identification,omitempty"`
	Email          string          `json:"email,omitempty"`
	PhonePrimary   string          `json:"phonePrimary,omitempty"`
	PhoneSecondary string          `json:"phoneSecondary,omitempty"`
	Mobile         string          `json:"mobile,omitempty"`
	Fax            string          `json:"fax,omitempty"`
	Observations   string          `json:"observations,omitempty"`
	Status         string          `json:"status,omitempty"`
	Type           []string        `json:"type,omitempty"` // e.g. ["client"], ["provider"]
	KindOfPerson   string          `json:"kindOfPerson,omitempty"`
	Regime         string          `json:"regime,omitempty"`
	Address        *ContactAddress `json:"address,omitempty"`
	Seller         *Ref            `json:"seller,omitempty"`
	Term           *Ref            `json:"term,omitempty"`
	PriceList      *Ref            `json:"priceList,omitempty"`
}

// ContactAddress is the nested address of a contact.
type ContactAddress struct {
	Address    string `json:"address,omitempty"`
	City       string `json:"city,omitempty"`
	Department string `json:"department,omitempty"`
	Country    string `json:"country,omitempty"`
	ZipCode    string `json:"zipCode,omitempty"`
}

// Contacts returns a typed handle to the /contacts resource.
func (c *Client) Contacts() *Resource[Contact] {
	return NewResource[Contact](c, "contacts")
}
