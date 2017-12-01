package collectors

import (
	"goku-bot/exchanges"
	log "github.com/inconshreveable/log15"
	"time"
	"goku-bot/keys"
)

type ChartDataCollector struct {
	exchange string
	pair     string
	interval time.Duration
	client   exchanges.Exchange
}

func NewChartDataCollector(exchange, pair string, interval time.Duration, client exchanges.Exchange) *ChartDataCollector {
	return &ChartDataCollector{exchange, pair, interval, client}
}

func (self *ChartDataCollector) Collect() {
	log.Debug("Collecting", "module", "exchangeCollectors", "key", self.Key())
	end := time.Now()
	start := end.Add(-10 * 24 * time.Hour)
	response := self.client.ChartData(self.pair, self.interval, start, end)

	if response.Err != nil {
		log.Warn("Error collecting", "module", "exchangeCollectors", "key", self.Key(), "error", response.Err)
		return
	}

	data := response.Data
	log.Debug("about to process trades", "module", "exchangeCollectors", "key", self.Key(), "trades", data)

	sendToProcessor(self.Key(), data)
}

func (self *ChartDataCollector) Key() string {
	return keys.BuildChartDataKey(self.exchange, self.pair, self.interval)
}
