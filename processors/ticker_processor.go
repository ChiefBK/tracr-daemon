package processors

import (
	log "github.com/inconshreveable/log15"
	"tracr-daemon/exchanges"
	"tracr-daemon/keys"
	"tracr-cache"
)

type TickerProcessor struct {
	exchange string
	pair     string
}

func NewTickerProcessor(exchange, pair string) *TickerProcessor {
	return &TickerProcessor{exchange, pair}
}

func (self *TickerProcessor) Process(input interface{}) {
	log.Debug("processing", "key", self.Key(), "module", "processors")
	ticker := input.(exchanges.Ticker)

	tracr_cache.PutTicker(self.Key(), ticker)
}

func (self *TickerProcessor) Key() string {
	return keys.BuildTickerKey(self.exchange, self.pair)
}
