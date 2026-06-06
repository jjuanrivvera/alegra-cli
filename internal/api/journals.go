package api

// Journal is an Alegra accounting journal entry (comprobante contable): a
// double-entry voucher whose lines (entries) carry debit/credit amounts.
type Journal struct {
	ID           ID             `json:"id,omitempty"`
	Number       string         `json:"number,omitempty"`
	Date         string         `json:"date,omitempty"` // YYYY-MM-DD
	Observations string         `json:"observations,omitempty"`
	Reference    string         `json:"reference,omitempty"`
	Status       string         `json:"status,omitempty"` // open, void
	Total        Money          `json:"total,omitempty"`
	Client       *Ref           `json:"client,omitempty"`
	Employee     *Ref           `json:"employee,omitempty"`
	Entries      []JournalEntry `json:"entries,omitempty"`
}

// JournalEntry is a single debit/credit line of a journal voucher.
type JournalEntry struct {
	ID          ID     `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Account     *Ref   `json:"account,omitempty"` // ledger account this line hits
	Debit       Money  `json:"debit,omitempty"`
	Credit      Money  `json:"credit,omitempty"`
	Client      *Ref   `json:"client,omitempty"`
	CostCenter  *Ref   `json:"costCenter,omitempty"`
}

// Journals returns a typed handle to the /journals resource.
func (c *Client) Journals() *Resource[Journal] {
	return NewResource[Journal](c, "journals")
}
