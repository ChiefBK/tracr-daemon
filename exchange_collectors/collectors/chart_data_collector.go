package collectors

import (
	"tracr-daemon/exchanges"
	log "github.com/inconshreveable/log15"
	"time"
	"tracr-daemon/keys"
	"tracr-store"
)

type ChartDataCollector struct {
	exchange string
	pair     string
	interval time.Duration
	client   exchanges.ExchangeClient
	store    tracr_store.Store
}

func NewChartDataCollector(exchange, pair string, interval time.Duration, client exchanges.ExchangeClient) *ChartDataCollector {
	store, err := tracr_store.NewStore()

	if err != nil {
		log.Error("error creating store", "module", "collectors", "collector", "chartData")
		return nil
	}

	return &ChartDataCollector{exchange, pair, interval, client, store}
}

func (self *ChartDataCollector) Collect() {
	log.Debug("Collecting", "module", "exchangeCollectors", "key", self.Key())
	end := time.Now()
	start := end.Add(-10 * 24 * time.Hour) // get last 10 days
	response := self.client.ChartData(self.pair, self.interval, start, end)

	if response.Err != nil {
		log.Error("Error collecting", "module", "exchangeCollectors", "key", self.Key(), "error", response.Err)
		return
	}

	candles := response.Data

	self.store.ReplaceChartData(self.exchange, self.pair, self.interval, candles)

}

func (self *ChartDataCollector) Key() string {
	return keys.BuildChartDataKey(self.exchange, self.pair, self.interval)
}
