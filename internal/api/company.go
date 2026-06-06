package api

// Company is the singleton company (empresa) registered on the Alegra account.
// The Alegra API exposes it only at GET /company and PUT /company; it has no id
// or list endpoint. Fields vary by country (applicationVersion); unknown fields
// are ignored at decode time, so this struct only models the common ones.
type Company struct {
	Name                 string                 `json:"name,omitempty"`
	Identification       string                 `json:"identification,omitempty"`
	Phone                string                 `json:"phone,omitempty"`
	Website              string                 `json:"website,omitempty"`
	Email                string                 `json:"email,omitempty"`
	Regime               string                 `json:"regime,omitempty"`
	ApplicationVersion   string                 `json:"applicationVersion,omitempty"`
	Timezone             string                 `json:"timezone,omitempty"`
	DecimalPrecision     string                 `json:"decimalPrecision,omitempty"`
	DecimalSeparator     string                 `json:"decimalSeparator,omitempty"`
	Logo                 string                 `json:"logo,omitempty"`
	KindOfPerson         string                 `json:"kindOfPerson,omitempty"`
	Address              *CompanyAddress        `json:"address,omitempty"`
	Currency             *CompanyCurrency       `json:"currency,omitempty"`
	IdentificationObject *CompanyIdentification `json:"identificationObject,omitempty"`
	InvoicePreferences   *CompanyInvoicePrefs   `json:"invoicePreferences,omitempty"`
}

// CompanyAddress is the company's postal address. Field set varies by country
// (e.g. department/zipCode in CO, province/postalCode in AR, street/state in MX).
type CompanyAddress struct {
	Description    string `json:"description,omitempty"`
	City           string `json:"city,omitempty"`
	Department     string `json:"department,omitempty"`
	Province       string `json:"province,omitempty"`
	State          string `json:"state,omitempty"`
	Region         string `json:"region,omitempty"`
	Street         string `json:"street,omitempty"`
	ExteriorNumber string `json:"exteriorNumber,omitempty"`
	InteriorNumber string `json:"interiorNumber,omitempty"`
	Neighborhood   string `json:"neighborhood,omitempty"`
	ZipCode        string `json:"zipCode,omitempty"`
	PostalCode     string `json:"postalCode,omitempty"`
}

// CompanyCurrency is the company's default currency.
type CompanyCurrency struct {
	Code   string `json:"code,omitempty"`
	Symbol string `json:"symbol,omitempty"`
}

// CompanyIdentification is the structured identification document of the company.
type CompanyIdentification struct {
	Type                    string `json:"type,omitempty"`
	Number                  string `json:"number,omitempty"`
	DV                      string `json:"dv,omitempty"`
	NationalityKindOfPerson string `json:"nationalityKindOfPerson,omitempty"`
}

// CompanyInvoicePrefs holds the company's default invoicing preferences.
type CompanyInvoicePrefs struct {
	DefaultAnotation          string `json:"defaultAnotation,omitempty"`
	DefaultTermsAndConditions string `json:"defaultTermsAndConditions,omitempty"`
}

// Company returns a typed handle to the /company singleton resource. The generic
// Resource list/get-by-id helpers do not apply (there is no collection or id);
// use Client.GetInto / Client.PutInto against the "company" path directly.
func (c *Client) Company() *Resource[Company] {
	return NewResource[Company](c, "company")
}
