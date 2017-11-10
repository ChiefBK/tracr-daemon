package strategies

import (
	"goku-bot"
	log "github.com/sirupsen/logrus"
)

type Strategy struct {
	decisionTrees []*DecisionTree
}

func NewStategy(trees []*DecisionTree) *Strategy {
	return &Strategy{trees}
}

func (self *Strategy) run(botActionChan chan<- *goku_bot.ActionQueue) {
	log.WithFields(log.Fields{"module": "strategies"}).Debug("running strategy")
	botActionQueue := goku_bot.NewActionQueue() // the queue that will be sent back to the bot

	for _, tree := range self.decisionTrees {
		treeActionChan := make(chan *goku_bot.ActionQueue)
		go tree.run(treeActionChan)

		treeActionQueue := <- treeActionChan // gets actions from tree

		for _, action := range treeActionQueue.Queue { // add actions from tree action queue to bot action queue
			botActionQueue.Push(action)
		}
	}

	botActionChan <- botActionQueue
}
