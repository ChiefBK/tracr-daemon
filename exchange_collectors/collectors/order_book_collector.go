package collectors

import (
	"tracr-daemon/exchanges"
	"tracr-daemon/keys"
	log "github.com/inconshreveable/log15"
)

type OrderBookCollector struct {
	exchange string
	pair string
	exchangeClient exchanges.ExchangeClient
}

func NewOrderBookCollector(exchange, pair string, exchangeClient exchanges.ExchangeClient) *OrderBookCollector {
	return &OrderBookCollector{exchange, pair, exchangeClient}
}

func (self *OrderBookCollector) Collect() {
	response := self.exchangeClient.OrderBook(self.pair)

	if response.Err != nil {
		log.Warn("There was an error collecting the order book", "error", response.Err, "module", "exchangeCollectors")
		return
	}

	//orderBook := response.Data

	//sendToProcessor(self.Key(), orderBook)
}

func (self *OrderBookCollector) Key() string {
	return keys.BuildOrderBookKey(self.exchange, self.pair)
}
