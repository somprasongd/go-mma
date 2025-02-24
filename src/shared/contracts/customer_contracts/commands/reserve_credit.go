package commands

type ReserveCreditCommand struct {
	ID           int `json:"id"`
	CreditAmount int `json:"credit_amount"`
}
