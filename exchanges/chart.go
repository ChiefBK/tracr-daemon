package exchanges

import "time"

type Candle struct {
	Open float64
	High float64
	Low float64
	Close float64
	DateTime time.Time
}

type ChartDataResponse struct {
	Data []Candle
	Err error
}
