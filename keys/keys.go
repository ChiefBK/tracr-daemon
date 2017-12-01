package keys

import (
	"strings"
	"time"
	"fmt"
)

func buildCollectionKey(params ...string) string {
	return strings.Join(params, "-")
}

func BuildChartDataKey(exchange, pair string, interval time.Duration) string {
	intervalMins := int64(interval / time.Minute)
	return buildCollectionKey("ChartData", exchange, pair, fmt.Sprintf("%d", intervalMins))
}

func BuildBalancesKey(exchange string) string {
	return buildCollectionKey("Balances", exchange)
}

func BuildMyTradeHistoryKey(exchange, pair string) string {
	return buildCollectionKey("MyTradeHistory", exchange, pair)
}

func BuildDepositHistoryKey(exchange string) string {
	return buildCollectionKey("MyDepositHistory", exchange)
}

func BuildWithdrawalHistoryKey(exchange string) string {
	return buildCollectionKey("MyWithdrawlHistory", exchange)
}
