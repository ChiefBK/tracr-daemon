package processors

import (
	"fmt"
	"goku-bot/streams"
	"goku-bot"
	log "github.com/inconshreveable/log15"
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
	orderBook := input.(*goku_bot.OrderBook)
	streams.BroadcastOrderBook(self.Key(), *orderBook)
}

func (self *OrderBookProcessor) Key() string {
	return fmt.Sprintf("%s-OrderBook-%s", self.exchange, self.pair)
}
