package api

// BankAccount is an Alegra bank account (a bank, credit card, or cash account)
// used to register and reconcile money movements.
type BankAccount struct {
	ID                 ID     `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	Number             string `json:"number,omitempty"`
	Description        string `json:"description,omitempty"`
	Type               string `json:"type,omitempty"`   // bank, credit-card, cash
	Status             string `json:"status,omitempty"` // active, inactive
	InitialBalance     Money  `json:"initialBalance,omitempty"`
	InitialBalanceDate string `json:"initialBalanceDate,omitempty"`
	Balance            Money  `json:"balance,omitempty"`
	DecimalPrecision   string `json:"decimalPrecision,omitempty"`
	CalculationScale   string `json:"calculationScale,omitempty"`
	Category           *Ref   `json:"category,omitempty"`
}

// BankTransfer is the body of a transfer from one bank account to another
// (POST /bank-accounts/{id}/transfer).
type BankTransfer struct {
	IDDestination string `json:"idDestination,omitempty"`
	Amount        Money  `json:"amount,omitempty"`
	Date          string `json:"date,omitempty"`
	Observations  string `json:"observations,omitempty"`
	ExchangeRate  Money  `json:"exchangeRate,omitempty"`
}

// BankAccounts returns a typed handle to the /bank-accounts resource.
func (c *Client) BankAccounts() *Resource[BankAccount] {
	return NewResource[BankAccount](c, "bank-accounts")
}
