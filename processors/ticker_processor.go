package processors

import (
	"goku-bot"
	"goku-bot/streams"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type TickerProcessor struct {
	exchange string
	pair     string
}

func NewTickerProcessor(exchange, pair string) *TickerProcessor {
	return &TickerProcessor{exchange, pair}
}

func (self *TickerProcessor) Process(input interface{}) {
	log.WithFields(log.Fields{"key": self.Key(), "module": "processors"}).Debug("processing")
	ticker := input.(*goku_bot.Ticker)

	streams.BroadcastTicker(self.Key(), *ticker)
}

func (self *TickerProcessor) Key() string {
	return fmt.Sprintf("%s-Ticker-%s", self.exchange, self.pair)
}
