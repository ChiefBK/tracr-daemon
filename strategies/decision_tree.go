package strategies

import (
	"goku-bot"
	log "github.com/sirupsen/logrus"
)

type DecisionTree struct {
	root *Signal
}

func NewDecisionTree(rootSignal *Signal) *DecisionTree {
	return &DecisionTree{rootSignal}
}

func (self *DecisionTree) run(actionQueueChan chan<- *goku_bot.ActionQueue) {
	log.WithFields(log.Fields{"module": "strategies"}).Debug("running tree")

	signalActionChan := make(chan *goku_bot.Action)
	actionQueue := goku_bot.NewActionQueue()

	go self.root.run(signalActionChan) // runs root signal of tree

	for action := range signalActionChan { // reads actions from signals thru channel
		actionQueue.Push(action)
	}

	actionQueueChan <- actionQueue // Sends queue of actions to Strategy
}