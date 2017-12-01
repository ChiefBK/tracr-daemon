package processors

import (
	"time"
	"goku-bot/exchanges"
	log "github.com/inconshreveable/log15"
	"goku-bot/store"
)

type ChartDataProcessor struct {
	exchange string
	pair     string
	interval time.Duration
	store    store.Store
}

func NewChartDataProcessor(exchange, pair string, interval time.Duration) *ChartDataProcessor {
	store, _ := store.NewStore()

	return &ChartDataProcessor{exchange, pair, interval, store}
}

func (self *ChartDataProcessor) Process(input interface{}) {
	log.Debug("processing", "key", self.Key(), "module", "processors")
	candles := input.([]exchanges.Candle)

	self.store.ReplaceChartData(self.exchange, self.pair, self.interval, candles)
}

func (self *ChartDataProcessor) Key() string {
	return store.BuildChartDataCollectionName(self.exchange, self.pair, self.interval)
}
