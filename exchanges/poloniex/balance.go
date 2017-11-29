package poloniex

import "strconv"

type PoloniexCompleteBalance struct {
	Available string `json:"available"`
	OnOrders  string `json:"onOrders"`
	BtcValue  string `json:"btcValue"`
}

type PoloniexCompleteBalancesResponse map[string]PoloniexCompleteBalance

func (self PoloniexCompleteBalance) GetBtcValue() float64 {
	floatVal, _ := strconv.ParseFloat(self.BtcValue, 64)

	return floatVal
}

func (self PoloniexCompleteBalance) GetOnOrders() float64 {
	floatVal, _ := strconv.ParseFloat(self.OnOrders, 64)

	return floatVal
}

func (self PoloniexCompleteBalance) GetAvailable() float64 {
	floatVal, _ := strconv.ParseFloat(self.Available, 64)

	return floatVal
}