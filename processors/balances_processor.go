package processors

import (
	"tracr-daemon/store"
	log "github.com/inconshreveable/log15"
	"tracr-daemon/keys"
	"tracr-daemon/exchanges"
	"tracr-daemon/streams"
)

type BalanceProcessor struct {
	exchange string
	store    store.Store
}

func (self *BalanceProcessor) Process(input interface{}) {
	log.Debug("processing", "key", self.Key(), "module", "processors")
	balances := input.(exchanges.Balances)

	streams.PutValue(self.Key(), balances)
}

func (self *BalanceProcessor) Key() string {
	return keys.BuildBalancesKey(self.exchange)
}

func NewBalanceProcessor(exchange string) *BalanceProcessor {
	store, _ := store.NewStore()

	return &BalanceProcessor{exchange, store}
}
