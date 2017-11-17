package collectors

import (
	"poloniex-go-api"
	log "github.com/inconshreveable/log15"
)

type MyTradeHistoryCollector struct {
	Poloniex *poloniex_go_api.Poloniex
}

func NewMyTradeHistoryCollector() *MyTradeHistoryCollector {
	p := poloniex_go_api.New(API_KEY, API_SECRET)

	return &MyTradeHistoryCollector{p}
}

func (self *MyTradeHistoryCollector) Key() string {
	return "MyTradeHistory"
}

func (self *MyTradeHistoryCollector) Collect() {
	log.Debug("Collecting", "module", "collectors", "key", self.Key())
	response := self.Poloniex.ReturnMyTradeHistory()

	if response.Err != nil {
		log.Warn("Error collecting", "module", "collectors", "key", self.Key(), "error", response.Err)
		return
	}

	data := response.Data

	sendToProcessor(self.Key(), data)
}
