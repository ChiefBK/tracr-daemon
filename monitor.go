package goku_bot

import (
	"poloniex-go-api"
	"time"
)

const (
	ETH_BTC_PAIR = "BTC_ETH"

	FIVE_MIN_PERIOD    = 300
	FIFTEEN_MIN_PERIOD = 900
	THIRTY_MIN_PERIOD  = 1800
	TWO_HOUR_PERIOD    = 7200
	FOUR_HOUR_PERIOD   = 14400
	ONE_DAY_PERIOD     = 86400 // number of seconds in one day
)

var poloniexCandlestickPeriods = map[int]string{
	FIVE_MIN_PERIOD: "5_minutes",
	FIFTEEN_MIN_PERIOD: "15_minutes",
	THIRTY_MIN_PERIOD: "30_minutes",
	TWO_HOUR_PERIOD: "2_hours",
	FOUR_HOUR_PERIOD: "4_hours",
	ONE_DAY_PERIOD: "1_day",
}

type Monitor struct {
	Poloniex *poloniex_go_api.Poloniex
	Store    *Store
}

// TODO - return error to main go-routine if error
func (m *Monitor) SyncOHLC(err chan error) {
	defer close(err)

	end := time.Now()
	start := end.AddDate(0, 0, -1)

	poloniexBtcEthFiveMin := make(chan *poloniex_go_api.ReturnChartDataResponse)
	poloniexBtcEthFifteenMin := make(chan *poloniex_go_api.ReturnChartDataResponse)

	go m.Poloniex.ReturnChartData(ETH_BTC_PAIR, FIVE_MIN_PERIOD, start, end, poloniexBtcEthFiveMin)
	go m.Poloniex.ReturnChartData(ETH_BTC_PAIR, FIFTEEN_MIN_PERIOD, start, end, poloniexBtcEthFifteenMin)

	poloniexBtcEthFiveMinResp := <-poloniexBtcEthFiveMin
	poloniexBtcEthFifteenMinResp := <-poloniexBtcEthFifteenMin

	if poloniexBtcEthFiveMinResp.Err != nil {
		err <- poloniexBtcEthFiveMinResp.Err
		return
	}

	if poloniexBtcEthFifteenMinResp.Err != nil {
		err <- poloniexBtcEthFifteenMinResp.Err
		return
	}

	m.Store.SyncCandles(poloniexBtcEthFiveMinResp.Response, "poloniex", ETH_BTC_PAIR, poloniexCandlestickPeriods[FIVE_MIN_PERIOD])
	m.Store.SyncCandles(poloniexBtcEthFifteenMinResp.Response, "poloniex", ETH_BTC_PAIR, poloniexCandlestickPeriods[FIFTEEN_MIN_PERIOD])
}
