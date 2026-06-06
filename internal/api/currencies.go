package api

// Currency is an Alegra currency (moneda) registered in the account.
//
// IMPORTANT: currencies are identified by their ISO Code (e.g. "USD"), not by a
// numeric id. Get and Update therefore take the code string as the identifier.
type Currency struct {
	// Code is the ISO currency code (e.g. "USD", "COP") and is the resource identifier.
	Code string `json:"code,omitempty"`
	// Name is the human-readable currency name (e.g. "Dólar estadounidense").
	Name string `json:"name,omitempty"`
	// Symbol is the currency symbol (e.g. "$").
	Symbol string `json:"symbol,omitempty"`
	// Status is "active" or "inactive".
	Status string `json:"status,omitempty"`
	// ExchangeRate is the exchange rate against the account's main currency.
	ExchangeRate Money `json:"exchangeRate,omitempty"`
	// Deletable is returned only when the request asks for fields=deletable.
	Deletable *bool `json:"deletable,omitempty"`
	// CanBeInactive is returned only when the request asks for fields=canBeInactive.
	CanBeInactive *bool `json:"canBeInactive,omitempty"`
	// AutoUpdate is returned only when the request asks for fields=autoUpdate.
	AutoUpdate *bool `json:"autoUpdate,omitempty"`
}

// Currencies returns a typed handle to the /currencies resource.
//
// Note: currencies are identified by their ISO Code (e.g. "USD"), so Get and
// Update expect the code string rather than a numeric id.
func (c *Client) Currencies() *Resource[Currency] {
	return NewResource[Currency](c, "currencies")
}
