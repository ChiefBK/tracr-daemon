package indicators

func CalculateAroon(periods int, candles []float64) (aroonUp []int, aroonDown []int) {
	var window []float64
	for _, candle := range candles {
		window = append(window, candle)

		if len(window) < periods {
			continue
		}

		if len(window) > periods {
			window = window[1:]
		}

		high, low := daysSinceLastHighLow(window)

		up := int(((float64(periods) - float64(high)) / float64(periods)) * 100)
		down := int(((float64(periods) - float64(low)) / float64(periods)) * 100)

		aroonUp = append(aroonUp, up)
		aroonDown = append(aroonDown, down)
	}

	return
}

func daysSinceLastHighLow(window []float64) (high, low int) {
	high = 0
	low = 0

	var highest float64
	var lowest float64

	for i := range window {
		index := len(window) - 1 - i // Start at end
		candle := window[index]

		if i == 0 {
			highest = candle
			lowest = candle
			continue
		}

		if candle > highest {
			highest = candle
			high = index
		}

		if candle < lowest {
			lowest = candle
			low = index
		}
	}

	return
}
