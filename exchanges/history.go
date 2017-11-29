package exchanges

import "time"

type Trade struct {
	ID       string    `json:"id"`
	Date     time.Time `json:"date"`
	Rate     float64   `json:"rate"`
	Amount   float64   `json:"amount"`
	Total    float64   `json:"total"`
	Fee      float64   `json:"fee"`
	OrderId  string    `json:"orderId"`
	Type     string    `json:"type"`
}

type TradeHistoryResponse struct {
	Data map[string][]Trade // mapping of pairs to list of trades
	Err  error
}
