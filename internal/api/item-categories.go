package api

// ItemCategory is an Alegra item category used to group products and services.
// See https://developer.alegra.com/reference/get_item-categories
type ItemCategory struct {
	ID          ID                 `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	Status      string             `json:"status,omitempty"` // active | inactive
	Image       *ItemCategoryImage `json:"image,omitempty"`
	Parent      *Ref               `json:"parent,omitempty"`   // parent category, for nested hierarchies
	Children    []Ref              `json:"children,omitempty"` // child categories, when nested
}

// ItemCategoryImage is the nested image associated with an item category.
type ItemCategoryImage struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// ItemCategories returns a typed handle to the /item-categories resource.
func (c *Client) ItemCategories() *Resource[ItemCategory] {
	return NewResource[ItemCategory](c, "item-categories")
}
