package api

// Report is a single row of an Alegra sales report. Report endpoints are
// read-only and aggregate sales documents over a date range, so a Report row
// is a sparse superset of the fields returned by the various report subpaths
// (sales-documents, sales-totals, sales-by-client, sales-by-seller). Unknown
// JSON fields are ignored, so each subpath only populates the columns it owns.
type Report struct {
	// sales-documents rows.
	ID             ID     `json:"id,omitempty"`
	DocumentNumber string `json:"documentNumber,omitempty"`
	DocumentType   string `json:"documentType,omitempty"`
	Date           string `json:"date,omitempty"`
	Status         string `json:"status,omitempty"`

	// sales-totals rows (grouped by day/month/year).
	Period      string `json:"period,omitempty"`
	BeforeTaxes Money  `json:"beforeTaxes,omitempty"`
	Tax         Money  `json:"tax,omitempty"`
	Discount    Money  `json:"discount,omitempty"`
	CreditNote  Money  `json:"creditNote,omitempty"`

	// sales-by-client / sales-by-seller rows.
	ClientName     string `json:"clientName,omitempty"`
	SellerName     string `json:"sellerName,omitempty"`
	TotalDocuments Int    `json:"totalDocuments,omitempty"`
	SubTotal       Money  `json:"subTotal,omitempty"`
	TotalPayed     Money  `json:"totalPayed,omitempty"`

	// Shared monetary total across all report rows.
	Total Money `json:"total,omitempty"`
}

// ReportEnvelope wraps the array returned by Alegra report endpoints, which nest
// rows under a "data" key.
type ReportEnvelope struct {
	Data []Report `json:"data"`
}

// Reports returns a typed handle to the /reports resource. Report endpoints are
// read-only; use the bespoke `reports` command (or client.GetInto) to fetch the
// individual report subpaths.
func (c *Client) Reports() *Resource[Report] {
	return NewResource[Report](c, "reports")
}
