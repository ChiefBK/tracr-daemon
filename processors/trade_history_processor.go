package processors

import (
	"goku-bot/store"
	"poloniex-go-api"
	log "github.com/inconshreveable/log15"
)

type MyTradeHistoryProcessor struct {
	Store store.Store
}

func NewMyTradeHistoryProcessor() *MyTradeHistoryProcessor {
	s, _ := store.NewStore()

	return &MyTradeHistoryProcessor{s}
}

func (self *MyTradeHistoryProcessor) Process(input interface{}) {
	log.Debug("processing", "key", self.Key(), "module", "processors")
	data := input.(map[string][]poloniex_go_api.Trade)

	for pair, trades := range data {
		self.Store.ReplaceTrades("poloniex", pair, trades)
	}
}

func (self *MyTradeHistoryProcessor) Key() string {
	return "MyTradeHistory"
}
