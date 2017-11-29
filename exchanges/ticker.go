package exchanges

type Ticker struct {
	LastTrade          *float64
	HighestBid         *float64
	LowestAsk          *float64
	TwentyFourHourHigh *float64
	TwentyFourHourLow  *float64
}

type TickerResponse struct {
	Data map[string]Ticker
	Err  error
}
