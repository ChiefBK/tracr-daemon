package strategies

import (
	"goku-bot"
	log "github.com/sirupsen/logrus"
)


var bots []*Bot

func Init() {
	log.WithFields(log.Fields{"module": "strategies"}).Info("Initializing strategies module")
	bot1 := NewBot("bot1", "poloniex", "USDT_BTC")

	rootSignal := NewSignal(func() bool {
		return true
	}, nil, true)
	leafSignal := NewSignal(func() bool {
		return true
	}, goku_bot.NewAction("leaf action", "Action stuffs"), false)
	rootSignal.addChild(leafSignal)

	tree := NewDecisionTree(rootSignal)
	var trees []*DecisionTree
	trees = append(trees, tree)

	strat := NewStategy(trees)
	bot1.addStrategy(CLOSED_POSITION, strat)

	bots = append(bots, bot1)
}

func Start() {
	log.WithFields(log.Fields{"module": "strategies"}).Info("Starting strategies module")
	for _, bot := range bots {
		go bot.start()
	}
}
