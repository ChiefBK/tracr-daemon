package processors

import (
	"goku-bot/streams"
	log "github.com/inconshreveable/log15"
	"goku-bot/exchanges/poloniex"
	"goku-bot/keys"
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
	orderBook := input.(poloniex.OrderBook)
	streams.PutValue(self.Key(), orderBook)
}

func (self *OrderBookProcessor) Key() string {
	return keys.BuildOrderBookKey(self.exchange, self.pair)
}
