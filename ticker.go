package goku_bot

import "time"

type Ticker struct {
	CurrencyPair       string
	Last               float64
	lowestAsk          float64
	HighestBid         float64
	PercentChange      float64
	BaseVolume         float64
	QuoteVolume        float64
	IsFrozen           bool
	TwentyFourHourHigh float64
	TwentyFourHourLow  float64
	Updated            time.Time
}

var TickerUsdtBtc *Ticker
