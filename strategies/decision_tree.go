package strategies

import (
	log "github.com/sirupsen/logrus"
	"goku-bot/strategies/actions"
)

type DecisionTree struct {
	root     *Signal
	position Position
}

func newDecisionTree(rootSignal *Signal, position Position) *DecisionTree {
	return &DecisionTree{rootSignal, position}
}

func (self *DecisionTree) run(actionQueueChan chan<- *actions.ActionQueue) {
	log.WithFields(log.Fields{"module": "strategies"}).Debug("running tree")

	signalActionChan := make(chan *actions.Action)
	actionQueue := actions.NewActionQueue()

	go self.root.run(signalActionChan) // runs root signal of tree

	for action := range signalActionChan { // reads actions from signals thru channel
		actionQueue.Push(action)
	}

	actionQueueChan <- actionQueue // Sends queue of actions to Strategy
}

func BuildDecisionChain(position Position, signals ...*Signal) *DecisionTree {
	var rootSignal *Signal
	var refSignal *Signal

	for _, signal := range signals {
		if rootSignal == nil { // if root signal
			rootSignal = signal
			refSignal = signal
			continue
		}

		refSignal.addChild(signal)
		refSignal = signal
	}

	return newDecisionTree(rootSignal, position)
}
