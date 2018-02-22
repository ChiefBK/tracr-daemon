package exchanges

import "time"

type Trade struct {
	ID        string    `json:"id"` // the transaction ID
	Date      time.Time `json:"date"`
	Price     float64   `json:"rate"`
	Volume    float64   `json:"amount"`
	TotalCost float64   `json:"total"`
	Fee       float64   `json:"fee"`
	OrderId   string    `json:"orderId"` // the order ID
	Type      string    `json:"type"`
}

type TradeHistoryResponse struct {
	Data map[string][]Trade // mapping of pairs to list of trades
	Err  error
}
