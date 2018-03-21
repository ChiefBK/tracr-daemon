package exchanges

import "time"

const (
	POLONIEX          = "poloniex"
	KRAKEN            = "kraken"
	POLONIEX_THROTTLE = 200 * time.Millisecond
	KRAKEN_THROTTLE = 1 * time.Second
)

var POLONIEX_INTERVALS = []time.Duration{5 * time.Minute, 15 * time.Minute, 30 * time.Minute, 2 * time.Hour, 4 * time.Hour, 24 * time.Hour}

// 1 minute, 5 minutes, 15 minutes, 30 minutes, 1 hour, 4 hours, 1 day, 7 days, 15 days
var KRAKEN_INTERVALS = []time.Duration{1 * time.Minute, 5 * time.Minute, 15 * time.Minute, 30 * time.Minute, 1 * time.Hour, 4 * time.Hour, 24 * time.Hour, 168 * time.Hour, 360 * time.Hour}

type ExchangeClient interface {
	Ticker() TickerResponse
	Balances() BalancesResponse
	ChartData(stdPair string, period time.Duration, start, end time.Time) ChartDataResponse
	MyTradeHistory() TradeHistoryResponse
	DepositsWithdrawals() DepositsWithdrawalsResponse
	OrderBook(stdPair string) OrderBookResponse
}
