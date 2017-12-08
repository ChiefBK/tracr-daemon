package math

func Avg(arr []float64) float64 {
	var sum float64 = 0

	for _, candle := range arr {
		val := candle
		sum += val
	}

	return sum / float64(len(arr))
}
