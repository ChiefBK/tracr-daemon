package goku_bot

import (
	"poloniex-go-api"
	"time"
)

type CompleteBalances struct {
	BTC     poloniex_go_api.Balance
	LTC     poloniex_go_api.Balance
	EXP     poloniex_go_api.Balance
	EMC2    poloniex_go_api.Balance
	PINK    poloniex_go_api.Balance
	BCN     poloniex_go_api.Balance
	FCT     poloniex_go_api.Balance
	BTS     poloniex_go_api.Balance
	VRC     poloniex_go_api.Balance
	BURST   poloniex_go_api.Balance
	ETH     poloniex_go_api.Balance
	BCH     poloniex_go_api.Balance
	ZEC     poloniex_go_api.Balance
	DASH    poloniex_go_api.Balance
	XMR     poloniex_go_api.Balance
	Updated time.Time
}

var PoloniexClient *poloniex_go_api.Poloniex
var PoloniexBalances *CompleteBalances
