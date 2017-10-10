package goku_bot

func Avg(arr []*CandleSlice) float64 {
	var sum float64 = 0

	for _, element := range arr {
		val := element.Candle.Close
		sum += val
	}

	return sum / float64(len(arr))
}
