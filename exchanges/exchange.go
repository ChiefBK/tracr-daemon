package exchanges

import "time"

const (
	POLONIEX          = "poloniex"
	KRAKEN            = "kraken"
	POLONIEX_THROTTLE = 200 * time.Millisecond
)

var POLONIEX_INTERVALS = []time.Duration{5 * time.Minute, 15 * time.Minute, 30 * time.Minute, 2 * time.Hour, 4 * time.Hour, 24 * time.Hour}

type Exchange interface {
	Ticker() TickerResponse
	Balances() BalancesResponse
	ChartData(stdPair string, period time.Duration, start, end time.Time) ChartDataResponse
	MyTradeHistory() TradeHistoryResponse
	DepositsWithdrawals() DepositsWithdrawalsResponse
	OrderBook(stdPair string) OrderBookResponse
}
