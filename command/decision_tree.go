package command

import (
	log "github.com/inconshreveable/log15"
	"tracr-daemon/command/actions"
)

type DecisionTree struct {
	root     *Signal
}

func newDecisionTree(rootSignal *Signal) *DecisionTree {
	return &DecisionTree{rootSignal}
}

func (self *DecisionTree) run(actionQueueChan chan<- *actions.ActionQueue) {
	log.Debug("running tree", "module", "command")

	signalActionChan := make(chan *actions.Action)
	actionQueue := actions.NewActionQueue()

	go self.root.run(signalActionChan) // runs root signal of tree

	for action := range signalActionChan { // reads actions from signals thru channel
		actionQueue.Push(action)
	}

	actionQueueChan <- actionQueue // Sends queue of actions to Strategy
}

func BuildDecisionChain(position string, signals ...*Signal) *DecisionTree {
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

	return newDecisionTree(rootSignal)
}
