package processors

import (
	log "github.com/inconshreveable/log15"
	"tracr-daemon/keys"
	"tracr-daemon/exchanges"
	"tracr-cache"
)

type OrderBookProcessor struct {
	exchange string
	pair     string
}

func NewOrderBookProcessor(exchange, pair string) *OrderBookProcessor {
	return &OrderBookProcessor{exchange, pair}
}

func (self *OrderBookProcessor) Process(input interface{}) {
	log.Debug("processing", "key", self.Key(), "module", "processors")
	orderBook := input.(exchanges.OrderBook)
	tracr_cache.PutOrderBook(self.Key(), orderBook)
}

func (self *OrderBookProcessor) Key() string {
	return keys.BuildOrderBookKey(self.exchange, self.pair)
}
