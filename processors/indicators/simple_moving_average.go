package indicators

import (
	"tracr-daemon/math"
	"errors"
)

func CalculateSimpleMovingAverage(periods int, candles []float64) (averages []float64, err error) {
	var window []float64
	for _, candle := range candles {
		window = append(window, candle)

		if len(window) > periods {
			window = window[1:]
		}

		if len(window) != periods {
			continue
		}

		avgValue := math.Avg(window)
		averages = append(averages, avgValue)
	}

	if len(averages) == 0 {
		return nil, errors.New("not enough data to calculate sma with provided parameters")
	}

	return
}
