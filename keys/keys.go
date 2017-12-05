package keys

import (
	"strings"
	"time"
	"fmt"
)

func buildKey(params ...string) string {
	return strings.Join(params, "-")
}

func BuildChartDataKey(exchange, pair string, interval time.Duration) string {
	intervalMins := int64(interval / time.Minute)
	return buildKey("ChartData", exchange, pair, fmt.Sprintf("%d", intervalMins))
}

func BuildBalancesKey(exchange string) string {
	return buildKey("Balances", exchange)
}

func BuildMyTradeHistoryKey(exchange, pair string) string {
	return buildKey("MyTradeHistory", exchange, pair)
}

func BuildDepositHistoryKey(exchange string) string {
	return buildKey("MyDepositHistory", exchange)
}

func BuildWithdrawalHistoryKey(exchange string) string {
	return buildKey("MyWithdrawlHistory", exchange)
}

func BuildOrderBookKey(exchange, pair string) string {
	return buildKey("OrderBook", exchange, pair)
}
