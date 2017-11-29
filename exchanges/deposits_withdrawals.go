package exchanges

import "time"

type Deposit struct {
	Currency      string
	Amount        float64
	Confirmations int
	Timestamp     time.Time
	Status        string
}

type Withdrawal struct {
	Currency  string
	Amount    float64
	Timestamp time.Time
	Status    string
}

type DepositsWithdrawals struct {
	Deposits    []Deposit
	Withdrawals []Withdrawal
}

type DepositsWithdrawalsResponse struct {
	Data *DepositsWithdrawals
	Err  error
}
