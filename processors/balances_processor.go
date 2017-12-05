package processors

import (
	"goku-bot/store"
	log "github.com/inconshreveable/log15"
	"goku-bot/keys"
	"goku-bot/exchanges"
	"goku-bot/streams"
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
