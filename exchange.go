package goku_bot

import (
	"poloniex-go-api"
	"time"
)

type Exchange interface {
	ReturnTicker() *poloniex_go_api.ReturnTickerResponse
	ReturnChartData(currencyPair string, period int, start, end time.Time, respCh chan *poloniex_go_api.ReturnChartDataResponse)
}
