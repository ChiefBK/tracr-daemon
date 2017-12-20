package processors

import (
	"time"
	"tracr-daemon/exchanges"
	log "github.com/inconshreveable/log15"
	"tracr-daemon/keys"
	"tracr-store"
)

type ChartDataProcessor struct {
	exchange string
	pair     string
	interval time.Duration
	store    tracr_store.Store
}

func NewChartDataProcessor(exchange, pair string, interval time.Duration) *ChartDataProcessor {
	store, _ := tracr_store.NewStore()

	return &ChartDataProcessor{exchange, pair, interval, store}
}

func (self *ChartDataProcessor) Process(input interface{}) {
	log.Debug("processing", "key", self.Key(), "module", "processors")
	candles := input.([]exchanges.Candle)

	//closes := exchanges.GetCloses(candles)
	//sma, err := indicators.CalculateSimpleMovingAverage(5, closes)
	//ema, err := indicators.CalculateExponentialMovingAverage(5, closes)
	//macd, signal, err := indicators.CalculateMacd(12, 26, 9, closes)

	//log.Debug("Candles", "key", self.Key(), "module", "processors", "candles", candles)
	//log.Debug("SMA", "key", self.Key(), "module", "processors", "sma5", sma)
	//log.Debug("EMA", "key", self.Key(), "module", "processors", "ema5", ema)

	//if err == nil {
	//	log.Debug("MACD", "key", self.Key(), "module", "processors", "macd12-26-9", macd, "signal", signal)
	//}

	self.store.ReplaceChartData(self.exchange, self.pair, self.interval, candles)
}

func (self *ChartDataProcessor) Key() string {
	return keys.BuildChartDataKey(self.exchange, self.pair, self.interval)
}
