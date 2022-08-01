package types

type AccountBalance struct {
	AccountID   int64  `json:"account_id" bun:"accountid"`
	AccountName string `json:"account_name" bun:"accountname"`
	Balance     string `json:"balance" bun:","`
}
