package processors

import (
	log "github.com/inconshreveable/log15"
	"tracr-daemon/keys"
	"tracr-store"
)

type BalanceProcessor struct {
	exchange string
	store    tracr_store.Store
}

func (self *BalanceProcessor) Process(input interface{}) {
	log.Debug("processing", "key", self.Key(), "module", "processors")
	//balances := input.(exchanges.Balances)

	//tracr_cache.PutBalances(self.Key(), balances)
}

func (self *BalanceProcessor) Key() string {
	return keys.BuildBalancesKey(self.exchange)
}

func NewBalanceProcessor(exchange string) *BalanceProcessor {
	store, _ := tracr_store.NewStore()

	return &BalanceProcessor{exchange, store}
}
