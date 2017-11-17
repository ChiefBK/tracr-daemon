package strategies

import (
	log "github.com/inconshreveable/log15"
	"goku-bot/strategies/conditions"
	"goku-bot/strategies/actions"
)

var bots []*Bot

func Init() {
	log.Info("Initializing strategies module", "module", "strategies")
	log.Debug("Creating bots", "module", "strategies")

	rootSignal := NewSignal(conditions.TrueFunction(), nil, true)
	leafSignal := NewSignal(conditions.TrueFunction(), actions.ShortPositionAction(), false)
	tree := BuildDecisionChain(CLOSED_POSITION, rootSignal, leafSignal)
	addBot("bot1", "poloniex", "USDT_BTC", nil, tree)

}

func Start() {
	log.Info("Starting strategies module", "module", "strategies")
	for _, bot := range bots {
		go bot.start()
	}
}
