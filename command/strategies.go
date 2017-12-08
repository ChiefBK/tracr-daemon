package command

import (
	log "github.com/inconshreveable/log15"
	"goku-bot/command/conditions"
	"goku-bot/command/actions"
)

var bots []*Bot

func Init() {
	log.Info("Initializing command module", "module", "command")
	log.Debug("Creating bots", "module", "command")

	rootSignal := NewSignal(conditions.TrueFunction(), nil, true)
	leafSignal := NewSignal(conditions.TrueFunction(), actions.ShortPositionAction(), false)
	tree := BuildDecisionChain(CLOSED_POSITION, rootSignal, leafSignal)
	addBot("bot1", "poloniex", "USDT_BTC", nil, tree)

}

func Start() {
	log.Info("Starting command module", "module", "command")
	for _, bot := range bots {
		go bot.start()
	}
}
