package goku_bot

import "goku-bot/store"

func Avg(arr []*store.CandleSlice) float64 {
	var sum float64 = 0

	for _, element := range arr {
		val := element.Candle.Close
		sum += val
	}

	return sum / float64(len(arr))
}
