package poloniex

import "strconv"

type PoloniexApiTicker struct {
	Last          string `json:"last"`
	LowestAsk     string `json:"lowestAsk"`
	HighestBid    string `json:"highestBid"`
	PercentChange string `json:"percentChange"`
	BaseVolume    string `json:"baseVolume"`
	QuoteVolume   string `json:"quoteVolume"`
}

func (self PoloniexApiTicker) getLast() float64 {
	ret, _ := strconv.ParseFloat(self.Last, 64)

	return ret
}

func (self PoloniexApiTicker) getLowestAsk() float64 {
	ret, _ := strconv.ParseFloat(self.LowestAsk, 64)

	return ret
}

func (self PoloniexApiTicker) getHighestBid() float64 {
	ret, _ := strconv.ParseFloat(self.HighestBid, 64)

	return ret
}