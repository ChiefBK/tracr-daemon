package indicators

import (
	"errors"
)

func CalculateMacd(fastEmaPeriods, slowEmaPeriods, signalPeriods int, candles []float64) (macdLine []float64, signalLine []float64, err error) {
	// fast ema periods must be less than slow
	if fastEmaPeriods >= slowEmaPeriods {
		return nil, nil, errors.New("fast ema periods must be less than slow ema periods")
	}

	fastEma, err := CalculateExponentialMovingAverage(fastEmaPeriods, candles)

	if err != nil {
		return
	}

	slowEma, err := CalculateExponentialMovingAverage(slowEmaPeriods, candles)

	if err != nil {
		return
	}

	// If there's not enough data to calculate ema's than return error
	if len(fastEma) == 0 || len(slowEma) == 0 {
		return nil, nil, errors.New("there's not enough data to compute the MACD with provided params")
	}

	// Start at end of both ema's and calculate their differences to create MACD line
	slowIndex := len(slowEma) - 1
	fastIndex := len(fastEma) - 1
	for slowIndex >= 0 {
		diff := fastEma[fastIndex] - slowEma[slowIndex]
		macdLine = append([]float64{diff}, macdLine...) // prepend data
		fastIndex--
		slowIndex--
	}

	signalLine, err = CalculateExponentialMovingAverage(signalPeriods, macdLine)

	if err != nil {
		return
	}

	return
}
