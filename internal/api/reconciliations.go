package api

import "encoding/json"

// Reconciliation is an Alegra bank reconciliation, matching the recorded
// transactions of a bank/cash account against its real balance for a period.
// See https://developer.alegra.com/reference/reconciliationdetails
type Reconciliation struct {
	ID            ID                     `json:"id,omitempty"`
	Account       *ReconciliationAccount `json:"account,omitempty"`
	RealBalance   Money                  `json:"realBalance,omitempty"`
	BankExpenses  Money                  `json:"bankExpenses,omitempty"`
	BankTaxes     Money                  `json:"bankTaxes,omitempty"`
	BankInterests Money                  `json:"bankInterests,omitempty"`
	Date          string                 `json:"date,omitempty"`
	Status        string                 `json:"status,omitempty"`
	Transactions  []ReconciliationTxn    `json:"transactions,omitempty"`
}

// ReconciliationAccount is the bank or cash account a reconciliation belongs to.
type ReconciliationAccount struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// ReconciliationTxn is a single transaction included in a reconciliation.
type ReconciliationTxn struct {
	ID            ID                    `json:"id,omitempty"`
	Client        *Ref                  `json:"client,omitempty"`
	Date          string                `json:"date,omitempty"`
	Amount        Money                 `json:"amount,omitempty"`
	PaymentMethod string                `json:"paymentMethod,omitempty"`
	Observations  string                `json:"observations,omitempty"`
	Anotation     string                `json:"anotation,omitempty"`
	Status        string                `json:"status,omitempty"`
	Type          string                `json:"type,omitempty"`
	Associations  []ReconciliationAssoc `json:"associations,omitempty"`
	// NumberTemplate is kept raw because Alegra serializes it as either an
	// object or an array of objects depending on the record.
	NumberTemplate json.RawMessage `json:"numberTemplate,omitempty"`
}

// ReconciliationAssoc is an item or category associated with a transaction.
type ReconciliationAssoc struct {
	ID           ID     `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Price        Money  `json:"price,omitempty"`
	Quantity     string `json:"quantity,omitempty"`
	Observations string `json:"observations,omitempty"`
	Total        Money  `json:"total,omitempty"`
	Type         string `json:"type,omitempty"`
}

// Reconciliations returns a typed handle to the bank-reconciliation resource.
//
// The official Alegra API documents this resource at /conciliations (no "re-"
// prefix); the CLI keeps the "reconciliations" command name for UX but talks to
// the documented path.
//
// Note: create performs a create-or-update (upsert); the API has no PUT endpoint.
func (c *Client) Reconciliations() *Resource[Reconciliation] {
	return NewResource[Reconciliation](c, "conciliations")
}
