package poloniex

import (
	"time"
	"strconv"
)

type PoloniexDepositsWithdrawls struct {
	Deposits    []PoloniexDeposit
	Withdrawals []PoloniexWithdrawal
}

type PoloniexDeposit struct {
	Currency      string `json:"currency"`
	Address       string `json:"address"`
	Amount        string `json:"amount"`
	Confirmations int    `json:"confirmations"`
	Txid          string `json:"txid"`
	Timestamp     int    `json:"timestamp"`
	Status        string `json:"status"`
}

func (self PoloniexDeposit) GetDatetime() time.Time {
	return time.Unix(int64(self.Timestamp), 0)
}

func (self PoloniexDeposit) GetAmount() float64 {
	amount, _ := strconv.ParseFloat(self.Amount, 64)
	return amount
}


type PoloniexWithdrawal struct {
	WithdrawalNumber int    `json:"withdrawalNumber"`
	Currency         string `json:"currency"`
	Address          string `json:"address"`
	Amount           string `json:"amount"`
	Timestamp        int    `json:"timestamp"`
	Status           string `json:"status"`
	IPAddress        string `json:"ipAddress"`
}

func (self PoloniexWithdrawal) GetDatetime() time.Time {
	return time.Unix(int64(self.Timestamp), 0)
}

func (self PoloniexWithdrawal) GetAmount() float64 {
	amount, _ := strconv.ParseFloat(self.Amount, 64)
	return amount
}

