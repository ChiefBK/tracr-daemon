package collectors

import (
	"tracr-daemon/exchanges"
	log "github.com/inconshreveable/log15"
	"time"
	"tracr-daemon/keys"
)

type ChartDataCollector struct {
	exchange string
	pair     string
	interval time.Duration
	client   exchanges.ExchangeClient
}

func NewChartDataCollector(exchange, pair string, interval time.Duration, client exchanges.ExchangeClient) *ChartDataCollector {
	return &ChartDataCollector{exchange, pair, interval, client}
}

func (self *ChartDataCollector) Collect() {
	log.Debug("Collecting", "module", "exchangeCollectors", "key", self.Key())
	end := time.Now()
	start := end.Add(-10 * 24 * time.Hour) // get last 10 days
	response := self.client.ChartData(self.pair, self.interval, start, end)

	if response.Err != nil {
		log.Warn("Error collecting", "module", "exchangeCollectors", "key", self.Key(), "error", response.Err)
		return
	}

	//data := response.Data

}

func (self *ChartDataCollector) Key() string {
	return keys.BuildChartDataKey(self.exchange, self.pair, self.interval)
}
