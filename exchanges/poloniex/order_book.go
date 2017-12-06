package poloniex

type PoloniexOrderBook struct {
	Asks     [][]float64
	Bids     [][]float64
	IsFrozen int
	Seq      int
}
