package exchanges

import "time"

const (
	POLONIEX = "poloniex"
	KRAKEN   = "kraken"
	POLONIEX_THROTTLE = 200 * time.Millisecond
)

type Exchange interface {
	Ticker() TickerResponse
	Balances() BalancesResponse
	ChartData(currencyPair string, period int, start, end time.Time) ChartDataResponse
	MyTradeHistory() TradeHistoryResponse
	DepositsWithdrawals() DepositsWithdrawalsResponse
}
