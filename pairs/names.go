package pairs

import (
	"errors"
	"goku-bot/exchanges"
)

/*

	There are many pairs we are monitoring. The first is what is being traded. The second is the unit of trading.
	Not all exchanges support every pair listed

	They include:

	BTC-USD
	BCH-USD
	LTC-USD
	ETH-USD
	XMR-USD
	XRP-USD

	BCH-BTC
	LTC-BTC
	ETH-BTC
	XMR-BTC
	STR-BTC
	ZEC-BTC
	EMC2-BTC
	FCT-BTC
	XRP-BTC
	BTS-BTC
	BURST-BTC
	PINK-BTC
	VRC-BTC
	BCN-BTC


	BTC: 		Bitcoin
	BCH: 		Bitcoin Cash
	LTC: 		Litecoin
	ETH: 		Ethereum
	XMR:		Monero
	STR:		Stellar
	ZEC:		Zcash
	EMC2:		Einsteinium
	FCT:		Factom
	XRP: 		Ripple
	BTS: 		BitShares
	BURST: 		Burst
	USD: 		United States Dollars
	PINK:		Pinkcoin
	BCN: 		Bytecoin
	VRC:		Vericoin
 */

const (
	BTC_USD   = "BTC-USD"
	BCH_USD   = "BCH-USD"
	LTC_USD   = "LTC-USD"
	ETH_USD   = "ETH-USD"
	XMR_USD   = "XMR-USD"
	BCH_BTC   = "BCH-BTC"
	LTC_BTC   = "LTC-BTC"
	ETH_BTC   = "ETH-BTC"
	XMR_BTC   = "XMR-BTC"
	STR_BTC   = "STR-BTC"
	ZEC_BTC   = "ZEC-BTC"
	EMC2_BTC  = "EMC2-BTC"
	FCT_BTC   = "FCT-BTC"
	XRP_BTC   = "XRP-BTC"
	BTS_BTC   = "BTS-BTC"
	BURST_BTC = "BURST-BTC"
	PINK_BTC  = "PINK-BTC"
	BCN_BTC   = "BCN-BTC"
	VRC_BTC   = "VRC-BTC"

	BTC_USD_POLONIEX   = "USDT_BTC"
	BCH_USD_POLONIEX   = "USDT_BCH"
	LTC_USD_POLONIEX   = "USDT_LTC"
	ETH_USD_POLONIEX   = "USDT_ETH"
	XMR_USD_POLONIEX   = "USDT_XMR"
	BCH_BTC_POLONIEX   = "BTC_BCH"
	LTC_BTC_POLONIEX   = "BTC_LTC"
	ETH_BTC_POLONIEX   = "BTC_ETH"
	XMR_BTC_POLONIEX   = "BTC_XMR"
	STR_BTC_POLONIEX   = "BTC_STR"
	ZEC_BTC_POLONIEX   = "BTC_ZEC"
	EMC2_BTC_POLONIEX  = "BTC_EMC2"
	FCT_BTC_POLONIEX   = "BTC_FCT"
	XRP_BTC_POLONIEX   = "BTC_XRP"
	BTS_BTC_POLONIEX   = "BTC_BTS"
	BURST_BTC_POLONIEX = "BTC_BURST"
	PINK_BTC_POLONIEX  = "BTC_PINK"
	VRC_BTC_POLONIEX   = "BTC_VRC"
	BCN_BTC_POLONIEX   = "BTC_BCN"

	BTC_USD_KRAKEN = "XXBTZUSD"
	BCH_USD_KRAKEN = "BCHUSD"
	LTC_USD_KRAKEN = "XLTCZUSD"
	ETH_USD_KRAKEN = "XETCZUSD"
	XMR_USD_KRAKEN = "XXMRZUSD"
	BCH_BTC_KRAKEN = "BCHXBT"
	LTC_BTC_KRAKEN = "XLTCXXBT"
	ETH_BTC_KRAKEN = "XETHXXBT"
	XMR_BTC_KRAKEN = "XXMRXXBT"
	ZEC_BTC_KRAKEN = "XZECXXBT"
	XRP_BTC_KRAKEN = "XXRPXXBT"
)

var poloniexStdPairs = map[string]string{
	"BTC-USD": BTC_USD_POLONIEX,
	"BCH-USD": BCH_USD_POLONIEX,
	"LTC-USD": LTC_USD_POLONIEX,
	"ETH-USD": ETH_USD_POLONIEX,
	"XMR-USD": XMR_USD_POLONIEX,

	"BCH-BTC":   BCH_BTC_POLONIEX,
	"LTC-BTC":   LTC_BTC_POLONIEX,
	"ETH-BTC":   ETH_BTC_POLONIEX,
	"XMR-BTC":   XMR_BTC_POLONIEX,
	"STR-BTC":   STR_BTC_POLONIEX,
	"ZEC-BTC":   ZEC_BTC_POLONIEX,
	"EMC2-BTC":  EMC2_BTC_POLONIEX,
	"FCT-BTC":   FCT_BTC_POLONIEX,
	"XRP-BTC":   XRP_BTC_POLONIEX,
	"BTS-BTC":   BTS_BTC_POLONIEX,
	"BURST-BTC": BURST_BTC_POLONIEX,
	PINK_BTC:    PINK_BTC_POLONIEX,
	VRC_BTC:     VRC_BTC_POLONIEX,
	BCN_BTC:     BCN_BTC_POLONIEX,
}

var poloniexExchPairs = map[string]string{
	BTC_USD_POLONIEX: "BTC-USD",
	BCH_USD_POLONIEX: "BCH-USD",
	LTC_USD_POLONIEX: "LTC-USD",
	ETH_USD_POLONIEX: "ETH-USD",
	XMR_USD_POLONIEX: "XMR-USD",

	BCH_BTC_POLONIEX:   "BCH-BTC",
	LTC_BTC_POLONIEX:   "LTC-BTC",
	ETH_BTC_POLONIEX:   "ETH-BTC",
	XMR_BTC_POLONIEX:   "XMR-BTC",
	STR_BTC_POLONIEX:   "STR-BTC",
	ZEC_BTC_POLONIEX:   "ZEC-BTC",
	EMC2_BTC_POLONIEX:  "EMC2-BTC",
	FCT_BTC_POLONIEX:   "FCT-BTC",
	XRP_BTC_POLONIEX:   "XRP-BTC",
	BTS_BTC_POLONIEX:   "BTS-BTC",
	BURST_BTC_POLONIEX: "BURST-BTC",
	PINK_BTC_POLONIEX:  PINK_BTC,
	VRC_BTC_POLONIEX:   VRC_BTC,
	BCN_BTC_POLONIEX:   BCN_BTC,
}

var krakenStdPairs = map[string]string{
	"BTC-USD": "XXBTZUSD",
	"BCH-USD": "BCHUSD",
	"LTC-USD": "XLTCZUSD",
	"ETH-USD": "XETCZUSD",
	"XMR-USD": "XXMRZUSD",

	"BCH-BTC": "BCHXBT",
	"LTC-BTC": "XLTCXXBT",
	"ETH-BTC": "XETHXXBT",
	"XMR-BTC": "XXMRXXBT",
	//"STR-BTC":   "STR_BTC",
	"ZEC-BTC": "XZECXXBT",
	//"EMC2-BTC":  "EMC2_BTC",
	//"FCT-BTC":   "FCT_BTC",
	"XRP-BTC": "XXRPXXBT",
	//"BTS-BTC":   "BTS_BTC",
	//"BURST-BTC": "BURST_BTC",
}

var krakenExchPairs = map[string]string{
	"XXBTZUSD": "BTC-USD",
	"BCHUSD":   "BCH-USD",
	"XLTCZUSD": "LTC-USD",
	"XETCZUSD": "ETH-USD",
	"XXMRZUSD": "XMR-USD",

	"BCHXBT":   "BCH-BTC",
	"XLTCXXBT": "LTC-BTC",
	"XETHXXBT": "ETH-BTC",
	"XXMRXXBT": "XMR-BTC",
	//"STR-BTC":   "STR_BTC",
	"XZECXXBT": "ZEC-BTC",
	//"EMC2-BTC":  "EMC2_BTC",
	//"FCT-BTC":   "FCT_BTC",
	"XXRPXXBT": "XRP-BTC",
	//"BTS-BTC":   "BTS_BTC",
	//"BURST-BTC": "BURST_BTC",
}

func ExchangePair(stdName, exchange string) (string, error) {
	switch exchange {
	case exchanges.POLONIEX:
		return getPoloniexExchPair(stdName)
	case exchanges.KRAKEN:
		return getKrakenExchPair(stdName)
	default:
		return "", errors.New("exchange specified not listed")
	}
}

func StandardPair(exchangePairName, exchange string) (string, error) {
	switch exchange {
	case exchanges.POLONIEX:
		return getPoloniexStdPair(exchangePairName)
	case exchanges.KRAKEN:
		return getKrakenStdPair(exchangePairName)
	default:
		return "", errors.New("exchange specified not listed")
	}
}

func getPoloniexExchPair(pair string) (string, error) {
	if pair, ok := poloniexStdPairs[pair]; ok {
		return pair, nil
	} else {
		return "", errors.New("could not find pair specified")
	}
}

func getKrakenExchPair(pair string) (string, error) {
	if pair, ok := krakenStdPairs[pair]; ok {
		return pair, nil
	} else {
		return "", errors.New("could not find pair specified")
	}
}

func getPoloniexStdPair(pair string) (string, error) {
	if pair, ok := poloniexExchPairs[pair]; ok {
		return pair, nil
	} else {
		return "", errors.New("could not find pair specified")
	}
}

func getKrakenStdPair(pair string) (string, error) {
	if pair, ok := krakenExchPairs[pair]; ok {
		return pair, nil
	} else {
		return "", errors.New("could not find pair specified")
	}
}
