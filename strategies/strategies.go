package strategies

import (
	log "github.com/sirupsen/logrus"
	"goku-bot/strategies/conditions"
	"goku-bot/strategies/actions"
)

var bots []*Bot

func Init() {
	log.WithFields(log.Fields{"module": "strategies"}).Info("Initializing strategies module")
	log.WithFields(log.Fields{"module": "strategies"}).Debug("Creating bots")

	rootSignal := NewSignal(conditions.TrueFunction(), nil, true)
	leafSignal := NewSignal(conditions.TrueFunction(), actions.ShortPositionAction(), false)
	tree := BuildDecisionChain(CLOSED_POSITION, rootSignal, leafSignal)
	addBot("bot1", "poloniex", "USDT_BTC", nil, tree)

}

func Start() {
	log.WithFields(log.Fields{"module": "strategies"}).Info("Starting strategies module")
	for _, bot := range bots {
		go bot.start()
	}
}
