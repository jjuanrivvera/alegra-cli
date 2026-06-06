package api

// User is an Alegra account user (a person with access to the company).
type User struct {
	ID          ID             `json:"id,omitempty"`
	Username    string         `json:"username,omitempty"`
	Name        string         `json:"name,omitempty"`
	LastName    string         `json:"lastName,omitempty"`
	Email       string         `json:"email,omitempty"`
	Phone       string         `json:"phone,omitempty"`
	PhoneCode   string         `json:"phoneCode,omitempty"`
	Role        string         `json:"role,omitempty"`
	Status      string         `json:"status,omitempty"`
	Language    string         `json:"language,omitempty"`
	Position    string         `json:"position,omitempty"`
	Permissions map[string]any `json:"permissions,omitempty"`
}

// Users returns a typed handle to the /users resource.
func (c *Client) Users() *Resource[User] {
	return NewResource[User](c, "users")
}
