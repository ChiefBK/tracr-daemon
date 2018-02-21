package collectors

import (
	log "github.com/inconshreveable/log15"
	"tracr-daemon/exchanges"
	"fmt"
)

type MyTradeHistoryCollector struct {
	exchange string
	client   exchanges.ExchangeClient
}

func NewMyTradeHistoryCollector(exchange string, client exchanges.ExchangeClient) *MyTradeHistoryCollector {
	return &MyTradeHistoryCollector{exchange, client}
}

func (self *MyTradeHistoryCollector) Key() string {
	return fmt.Sprintf("MyTradeHistory-%s", self.exchange)
}

func (self *MyTradeHistoryCollector) Collect() {
	log.Debug("Collecting", "module", "exchangeCollectors", "key", self.Key())
	response := self.client.MyTradeHistory()

	if response.Err != nil {
		log.Warn("Error collecting", "module", "exchangeCollectors", "key", self.Key(), "error", response.Err)
		return
	}

	//data := response.Data

	//sendToProcessor(self.Key(), data)
}
