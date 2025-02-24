package commands

type ReleaseCreditCommand struct {
	ID           int `json:"id"`
	CreditAmount int `json:"credit_amount"`
}
