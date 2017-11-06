package collectors

import (
	"poloniex-go-api"
	log "github.com/sirupsen/logrus"
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
	log.WithFields(log.Fields{"key": self.Key(), "module": "collectors"}).Debug("Collecting")
	response := self.Poloniex.ReturnMyTradeHistory()

	if response.Err != nil {
		log.WithFields(log.Fields{"key": self.Key(), "module": "collectors", "error": response.Err}).Warn("Error collecting")
		return
	}

	data := response.Data

	sendToProcessor(self.Key(), data)
}
