package poloniex

import (
	"strconv"
	"time"
)

type PoloniexTrade struct {
	GlobalTradeID int    `json:"globalTradeID"`
	TradeID       string `json:"tradeID"`
	Date          string `json:"date"`
	Rate          string `json:"rate"`
	Amount        string `json:"amount"`
	Total         string `json:"total"`
	Fee           string `json:"fee"`
	OrderNumber   string `json:"orderNumber"`
	Type          string `json:"type"`
	Category      string `json:"category"`
}

func (self PoloniexTrade) GetTotal() float64 {
	total, _ := strconv.ParseFloat(self.Total, 64)

	return total
}

func (self PoloniexTrade) GetAmount() float64 {
	amount, _ := strconv.ParseFloat(self.Amount, 64)

	return amount
}

func (self PoloniexTrade) GetOrderNumber() int64 {
	number, _ := strconv.ParseInt(self.OrderNumber, 10, 64)

	return number
}

func (self PoloniexTrade) GetRate() float64 {
	rate, _ := strconv.ParseFloat(self.Rate, 64)

	return rate
}

func (self PoloniexTrade) GetFee() float64 {
	fee, _ := strconv.ParseFloat(self.Fee, 64)

	return fee
}


func (self PoloniexTrade) GetDate() time.Time {
	const dateForm = "2006-01-02 15:04:05"
	t, _ := time.Parse(dateForm, self.Date)
	return t
}