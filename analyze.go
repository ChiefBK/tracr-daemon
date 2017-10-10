package goku_bot

import (
	"time"
	"fmt"
	"log"
	"strconv"
)

type position string

const (
	openShort     = "openShort"
	openLong      = "openLong"
	closePosition = "closePosition"
)

type Action struct {
	Position string
}

type ActionQueue struct {
	queue []*Action
}

func (aq *ActionQueue) push(action *Action) {
	aq.queue = append(aq.queue, action)
}

func (aq *ActionQueue) pop() *Action {
	if len(aq.queue) < 1 {
		return nil
	}

	action := aq.queue[0]
	aq.queue = aq.queue[1:]

	return action
}

type Indicator struct {
	Store    *Store
	Exchange string
	Pair     string
	Interval int
}

type MovingAverageResult struct {
	Value float64
	Date  time.Time
}

type DateValue struct {
	Date  time.Time
	Value interface{}
}

type DateFloat struct {
	Date  time.Time
	Value float64
}

func NewIndicator() (indicator *Indicator) {
	indicator = new(Indicator)

	return
}

func CalculateSimpleMovingAverage(periods int, slices []*CandleSlice) {
	p := strconv.Itoa(periods)

	var window []*CandleSlice
	for _, slice := range slices {
		window = append(window, slice)

		if len(window) > periods {
			window = window[1:]
		}

		if len(window) != periods {
			continue
		}

		avgValue := Avg(window)
		lastSlice := window[len(window)-1]
		lastSlice.Sma[p] = &avgValue
	}
}

func CalculateExponentialMovingAverage(periods int, slices []*CandleSlice) {
	p := strconv.Itoa(periods)
	log.Printf("Calculating EMA - %s", p)
	multiplier := float64(2) / (float64(periods) + float64(1))

	CalculateSimpleMovingAverage(periods, slices)

	for i, current := range slices {
		if i < periods-1 {
			continue
		}

		if i == periods-1 { // the first result of the EMA will be the same as the first result of SMA
			if current.Ema[p] == nil { // Make sure that this EMA hasn't been calculated before
				current.Ema[p] = current.Sma[p]
			}
			continue
		}

		lastSlice := slices[i-1]

		next := CalculateNextEma(current.Candle.Close, *lastSlice.Ema[p], multiplier)
		current.Ema[p] = &next
	}
}

func CalculateNextEma(current, last, multiplier float64) float64 {
	return (current-last)*multiplier + last
}

func CalculateMacdLine(fastEmaPeriods, slowEmaPeriods, signalPeriods int, slices []*CandleSlice) {
	macdParams := fmt.Sprintf("%d-%d-%d", fastEmaPeriods, slowEmaPeriods, signalPeriods)
	fastP := strconv.Itoa(fastEmaPeriods)
	slowP := strconv.Itoa(slowEmaPeriods)

	CalculateExponentialMovingAverage(fastEmaPeriods, slices)
	CalculateExponentialMovingAverage(slowEmaPeriods, slices)

	for _, slice := range slices {
		fastValue := slice.Ema[fastP]
		slowValue := slice.Ema[slowP]

		if fastValue == nil || slowValue == nil {
			continue
		}

		macd := *fastValue - *slowValue

		slice.Macd[macdParams] = &MacdValue{
			Macd: &macd,
		}
	}
}

func CalculateMacdSignalLine(fastEmaPeriods, slowEmaPeriods, signalPeriods int, slices []*CandleSlice) {
	macdParams := fmt.Sprintf("%d-%d-%d", fastEmaPeriods, slowEmaPeriods, signalPeriods)
	multiplier := float64(2) / (float64(signalPeriods) + float64(1))

	var window []*CandleSlice
	for _, slice := range slices {
		macd := slice.Macd[macdParams]

		if macd == nil {
			continue
		}

		window = append(window, slice)

		if len(window) < signalPeriods {
			continue
		}

		if len(window) > signalPeriods {
			window = window[1:]
		}

		if len(window) == 1 { // if signal period is 1 than signal line is same as macd line
			slice.Macd[macdParams].Signal = macd.Macd
			continue
		}

		last := window[len(window)-2]
		current := window[len(window)-1]

		lastSignal := last.Macd[macdParams].Signal
		currentMacd := current.Macd[macdParams].Macd

		if lastSignal == nil { // if first signal calculated than calculate using avg of macd's in window
			var sum float64 = 0
			for _, s := range window {
				sum += *s.Macd[macdParams].Macd
			}
			avg := sum / float64(len(window))

			current.Macd[macdParams].Signal = &avg
			continue
		}

		nextEma := CalculateNextEma(*currentMacd, *lastSignal, multiplier)

		current.Macd[macdParams].Signal = &nextEma
		histogram := *currentMacd - nextEma
		current.Macd[macdParams].Histogram = &histogram
	}
}

func CalculateMacd(fastEmaPeriods, slowEmaPeriods, signalPeriods int, slices []*CandleSlice) {
	log.Printf("Calculating MACD - %d, %d, %d", fastEmaPeriods, slowEmaPeriods, signalPeriods)
	CalculateMacdLine(fastEmaPeriods, slowEmaPeriods, signalPeriods, slices)
	CalculateMacdSignalLine(fastEmaPeriods, slowEmaPeriods, signalPeriods, slices)
}

func CalculateAroon(periods int, slices []*CandleSlice) {
	log.Printf("Calculating AROON - %d", periods)
	p := strconv.Itoa(periods)

	var window []*CandleSlice
	for _, slice := range slices {
		window = append(window, slice)

		if len(window) < periods {
			continue
		}

		if len(window) > periods {
			window = window[1:]
		}

		high, low := DaysSinceLastHighLow(window)

		up := int(((float64(periods) - float64(high)) / float64(periods)) * 100)
		down := int(((float64(periods) - float64(low)) / float64(periods)) * 100)

		slice.Aroon[p] = &AroonValue{
			Up: &up,
			Down: &down,
		}
	}
}

func DaysSinceLastHighLow(window []*CandleSlice) (high, low int){
	high = 0
	low = 0

	var highest float64
	var lowest float64

	for i := range window {
		index := len(window) - 1 - i
		last := window[index]
		candleClose := last.Candle.Close

		if i == 0 {
			highest = candleClose
			lowest = candleClose
			continue
		}

		if candleClose > highest {
			highest = candleClose
			high = index
		}

		if candleClose < lowest {
			lowest = candleClose
			low = index
		}
	}

	return
}
