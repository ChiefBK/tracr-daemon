package indicators

import "errors"

func CalculateExponentialMovingAverage(periods int, candles []float64) (exponentialAverages []float64, err error) {
	multiplier := float64(2) / (float64(periods) + float64(1)) // Multiplier: (2 / (periods + 1) )

	simpleAverages, err := CalculateSimpleMovingAverage(periods, candles)

	if err != nil {
		return nil, errors.New("not enough data to calculate ema with provided parameters")
	}

	for i, candle := range candles {
		if i < periods-1 { // skip first <period> candles
			continue
		}

		if i == periods-1 { // the first result of the EMA will be the same as the first result of SMA
			exponentialAverages = append(exponentialAverages, simpleAverages[0])
			continue
		}

		lastExponentialMovingAvg := exponentialAverages[len(exponentialAverages)-1]

		nextEma := calculateNextEma(candle, lastExponentialMovingAvg, multiplier)
		exponentialAverages = append(exponentialAverages, nextEma)
	}

	if len(exponentialAverages) == 0 {
		return nil, errors.New("not enough data to calculate ema with provided parameters")
	}

	return
}

// EMA: (Close - EMA(previous day)) x multiplier + EMA(previous day).
func calculateNextEma(currentCandle, lastExponentialMovingAvg, multiplier float64) float64 {
	return (currentCandle-lastExponentialMovingAvg)*multiplier + lastExponentialMovingAvg
}
