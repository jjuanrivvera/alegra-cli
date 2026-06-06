package api

// RecurringPayment is an Alegra recurring payment (pago recurrente): a payment
// that is automatically registered on a recurring schedule against a bank account.
// See https://developer.alegra.com/reference/get_recurring-payments
type RecurringPayment struct {
	ID               ID                        `json:"id,omitempty"`
	Number           string                    `json:"number,omitempty"`
	Status           string                    `json:"status,omitempty"`
	StartDate        string                    `json:"startDate,omitempty"`
	EndDate          string                    `json:"endDate,omitempty"`
	Amount           Money                     `json:"amount,omitempty"`
	RepeatEvery      ID                        `json:"repeatEvery,omitempty"`
	DatesEditable    bool                      `json:"datesEditable,omitempty"`
	DecimalPrecision ID                        `json:"decimalPrecision,omitempty"`
	Account          *RecurringPaymentAccount  `json:"account,omitempty"`
	Currency         *RecurringPaymentCurrency `json:"currency,omitempty"`
}

// RecurringPaymentAccount is the bank account a recurring payment is charged to.
type RecurringPaymentAccount struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// RecurringPaymentCurrency describes the currency (and exchange rate) of a recurring payment.
type RecurringPaymentCurrency struct {
	Code         string `json:"code,omitempty"`
	Symbol       string `json:"symbol,omitempty"`
	ExchangeRate Money  `json:"exchangeRate,omitempty"`
}

// RecurringPayments returns a typed handle to the /recurring-payments resource.
func (c *Client) RecurringPayments() *Resource[RecurringPayment] {
	return NewResource[RecurringPayment](c, "recurring-payments")
}
