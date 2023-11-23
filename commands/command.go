package commands

type OpenAccountCommand struct {
	AccountHolder  string
	AccountType    int
	OpeningBalance float64
}

type DepositFundCommand struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

type WithdrawFundCommand struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

type CloseAccountCommand struct {
	ID string `json:"id"`
}
