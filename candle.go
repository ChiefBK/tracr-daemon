package goku_bot

import (
	"poloniex-go-api"
	"time"
)

func GetCandles(ohlc []*OhlcSchema) (results []poloniex_go_api.Candle) {
	for _, element := range ohlc {
		results = append(results, *element.Candle)
	}

	return
}

func GetDateValues(candles []poloniex_go_api.Candle) (results []DateValue) {
	for _, candle := range candles {
		results = append(results, DateValue{
			Date:  time.Unix(int64(candle.Date), 0),
			Value: candle.Close,
		})
	}

	return
}
