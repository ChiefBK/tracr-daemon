package kraken

type KrakenTradeResult struct {
	Trades map[string]KrakenTrade // mapping of "postxid" to trade object
	Count  int
}

type KrakenTrade struct {
	Ordertxid string
	Pair      string
	Time      float64
	Type      string
	Ordertype string
	Price     string
	Cost      string
	Fee       string
	Vol       string
	Margin    string
	Misc      string
}
