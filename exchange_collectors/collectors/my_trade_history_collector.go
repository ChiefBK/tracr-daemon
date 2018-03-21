package collectors

import (
	log "github.com/inconshreveable/log15"
	"tracr-daemon/exchanges"
	"fmt"
	"tracr-store"
)

type MyTradeHistoryCollector struct {
	exchange string
	client   exchanges.ExchangeClient
	store    tracr_store.Store
}

func NewMyTradeHistoryCollector(exchange string, client exchanges.ExchangeClient) *MyTradeHistoryCollector {
	store, err := tracr_store.NewStore()

	if err != nil {
		log.Error("error creating store", "module", "collectors", "exchange", exchange, "collector", "myTradeHistory")
		return nil
	}

	return &MyTradeHistoryCollector{exchange, client, store}
}

func (self *MyTradeHistoryCollector) Key() string {
	return fmt.Sprintf("MyTradeHistory-%s", self.exchange)
}

func (self *MyTradeHistoryCollector) Collect() {
	log.Debug("Collecting", "module", "collectors", "key", self.Key())
	response := self.client.MyTradeHistory()

	if response.Err != nil {
		log.Error("Error collecting", "module", "collectors", "key", self.Key(), "error", response.Err)
		return
	}

	pairTradesMap := response.Data

	for pair, trades := range pairTradesMap {
		log.Debug("replacing trades for pair", "module", "collectors", "pair", pair, "numOfTrades", len(trades))
		self.store.ReplaceTrades(self.exchange, pair, trades)
	}

}
