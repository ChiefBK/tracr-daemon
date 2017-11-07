package processors

import (
	"fmt"
	"goku-bot/streams"
	"goku-bot"
	log "github.com/sirupsen/logrus"
)

type OrderBookProcessor struct {
	exchange string
	pair     string
}

func NewOrderBookProcessor(exchange, pair string) *OrderBookProcessor {
	return &OrderBookProcessor{exchange, pair}
}

func (self *OrderBookProcessor) Process(input interface{}) {
	log.WithFields(log.Fields{"key": self.Key(), "module": "processors"}).Debug("processing")
	orderBook := input.(*goku_bot.OrderBook)
	streams.BroadcastOrderBook(self.Key(), *orderBook)
}

func (self *OrderBookProcessor) Key() string {
	return fmt.Sprintf("%s-OrderBook-%s", self.exchange, self.pair)
}
