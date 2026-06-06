package api

// GlobalInvoice is an Alegra global invoice (factura global): a single CFDI that
// groups one or more sale tickets issued to the general public.
// See https://developer.alegra.com/reference/get_global-invoices
type GlobalInvoice struct {
	ID             ID                       `json:"id,omitempty"`
	Date           string                   `json:"date,omitempty"`
	Status         string                   `json:"status,omitempty"`
	Total          Money                    `json:"total,omitempty"`
	NumberTemplate *GlobalInvoiceNumberTmpl `json:"numberTemplate,omitempty"`
	SaleTickets    []GlobalInvoiceSaleTkt   `json:"saleTickets,omitempty"`
	Stamp          *GlobalInvoiceStamp      `json:"stamp,omitempty"`
}

// GlobalInvoiceNumberTmpl is the numbering template (prefix/number) of a global invoice.
type GlobalInvoiceNumberTmpl struct {
	ID     ID     `json:"id,omitempty"`
	Prefix string `json:"prefix,omitempty"`
	Number string `json:"number,omitempty"`
	Text   string `json:"text,omitempty"`
}

// GlobalInvoiceSaleTkt is a sale ticket grouped into a global invoice.
type GlobalInvoiceSaleTkt struct {
	ID      ID     `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	DueDate string `json:"dueDate,omitempty"`
	Status  string `json:"status,omitempty"`
	Client  *Ref   `json:"client,omitempty"`
	Total   Money  `json:"total,omitempty"`
	Balance Money  `json:"balance,omitempty"`
}

// GlobalInvoiceStamp holds the SAT electronic stamping (CFDI) data of a global invoice.
type GlobalInvoiceStamp struct {
	UUID                 string `json:"uuid,omitempty"`
	StampDate            string `json:"stampDate,omitempty"`
	SATCertificateNumber string `json:"satCertificateNumber,omitempty"`
	CertificateNumber    string `json:"certificateNumber,omitempty"`
	PaymentMethod        string `json:"paymentMethod,omitempty"`
}

// GlobalInvoices returns a typed handle to the /global-invoices resource.
func (c *Client) GlobalInvoices() *Resource[GlobalInvoice] {
	return NewResource[GlobalInvoice](c, "global-invoices")
}
