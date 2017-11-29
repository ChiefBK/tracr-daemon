package exchanges

type Candle struct {
	Open float64
	High float64
	Low float64
	Close float64
}

type ChartData []Candle

type ChartDataResponse struct {
	Data ChartData
	Err error
}
