package collectors

import (
	"tracr-daemon/exchanges"
	log "github.com/inconshreveable/log15"
	"tracr-daemon/keys"
)

type BalancesCollector struct {
	exchange       string
	exchangeClient exchanges.Exchange
}

func (self *BalancesCollector) Collect() {
	log.Debug("Collecting", "module", "exchangeCollectors", "key", self.Key())

	response := self.exchangeClient.Balances()

	log.Debug("balances response", "module", "exchangeCollectors", "key", self.Key(), "response", response)

	if response.Err != nil {
		log.Warn("Error collecting", "module", "exchangeCollectors", "key", self.Key(), "error", response.Err)
		return
	}

	data := response.Data

	sendToProcessor(self.Key(), data)
}

func (self *BalancesCollector) Key() string {
	return keys.BuildBalancesKey(self.exchange)
}

func NewBalancesCollector(exchange string, exchangeClient exchanges.Exchange) *BalancesCollector {
	return &BalancesCollector{exchange, exchangeClient}
}
