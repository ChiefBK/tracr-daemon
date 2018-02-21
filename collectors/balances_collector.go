package collectors

import (
	"tracr-daemon/exchanges"
	log "github.com/inconshreveable/log15"
	"tracr-daemon/keys"
	"tracr-cache"
)

type BalancesCollector struct {
	exchange       string
	exchangeClient exchanges.ExchangeClient
	cacheClient    *tracr_cache.CacheClient
}

func (self *BalancesCollector) Collect() {
	log.Debug("Collecting", "module", "exchangeCollectors", "key", self.Key())

	response := self.exchangeClient.Balances()

	log.Debug("balances response", "module", "exchangeCollectors", "key", self.Key(), "response", response)

	if response.Err != nil {
		log.Error("Error collecting", "module", "exchangeCollectors", "key", self.Key(), "error", response.Err)
		return
	}

	data := response.Data

	self.cacheClient.PutBalances(self.Key(), data)
}

func (self *BalancesCollector) Key() string {
	return keys.BuildBalancesKey(self.exchange)
}

func NewBalancesCollector(exchange string, exchangeClient exchanges.ExchangeClient) *BalancesCollector {
	cacheClient, err := tracr_cache.NewCacheClient()

	if err != nil {
		log.Error("error creating balances collector", "module", "exchangeCollectors", "exchange", exchange, "error", err)
		return nil
	}

	return &BalancesCollector{exchange, exchangeClient, cacheClient}
}
