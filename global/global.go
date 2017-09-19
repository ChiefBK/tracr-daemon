package global

const (
	DB_NAME = "goku-bot"
	DB_URI  = "localhost"

	// BTC Pairs
	BTC_ETH_PAIR = "BTC_ETH" // Ethereum
	BTC_LTC_PAIR = "BTC_LTC" // Litecoin
	BTC_XMR_PAIR = "BTC_XMR" // Monero
	BTC_XRP_PAIR = "BTC_XRP" // Ripple
	BTC_BCH_PAIR = "BTC_BCH" // Bitcoin Cash

	// USDT Pairs
	USDT_BTC_PAIR = "USDT_BTC" // Bitcoin
	USDT_ETH_PAIR = "USDT_ETH" // Ethereum
	USDT_LTC_PAIR = "USDT_LTC" // Litecoin
	USDT_XMR_PAIR = "USDT_XMR" // Monero
	USDT_BCH_PAIR = "USDT_BCH" // Bitcoin Cash
	USDT_XRP_PAIR = "USDT_XRP" // Ripple

	FIVE_MIN_INTERVAL    = 300
	FIFTEEN_MIN_INTERVAL = 900
	THIRTY_MIN_INTERVAL  = 1800
	TWO_HOUR_INTERVAL    = 7200
	FOUR_HOUR_INTERVAL   = 14400
	ONE_DAY_INTERVAL     = 86400 // number of seconds in one day
)

//var POLONIEX_PAIRS = []string{BTC_ETH_PAIR}
//
//var POLONIEX_OHLC_INTERVALS = map[int]string{
//	FIVE_MIN_INTERVAL:    "5_minutes",
//}

var POLONIEX_PAIRS = []string{BTC_ETH_PAIR, BTC_LTC_PAIR, BTC_XMR_PAIR, BTC_XRP_PAIR, BTC_BCH_PAIR, USDT_BTC_PAIR,
USDT_ETH_PAIR, USDT_LTC_PAIR, USDT_XMR_PAIR, USDT_BCH_PAIR, USDT_XRP_PAIR}

var POLONIEX_OHLC_INTERVALS = map[int]string{
	FIVE_MIN_INTERVAL:    "5_minutes",
	FIFTEEN_MIN_INTERVAL: "15_minutes",
	THIRTY_MIN_INTERVAL:  "30_minutes",
	TWO_HOUR_INTERVAL:    "2_hours",
	FOUR_HOUR_INTERVAL:   "4_hours",
	ONE_DAY_INTERVAL:     "1_day",
}
