package strategies

import (
	"goku-bot"
	log "github.com/sirupsen/logrus"
)

type Signal struct {
	condition func() bool
	children  []*Signal
	action    *goku_bot.Action
	isRoot    bool
}

func NewSignal(condition func() bool, action *goku_bot.Action, isRoot bool) *Signal {
	var children []*Signal
	return &Signal{condition, children, action, isRoot}
}

func (self *Signal) addChild(signal *Signal) {
	self.children = append(self.children, signal)
}

func (self *Signal) run(actionChan chan<- *goku_bot.Action) {
	log.WithFields(log.Fields{"module": "strategies", "children": len(self.children), "isRoot": self.isRoot}).Debug("running signal")
	result := self.condition()

	if result { // if signal is true
		if len(self.children) == 0 && self.action != nil { // if leaf node
			log.WithFields(log.Fields{"module": "strategies"}).Debug("sending action from signal")
			actionChan <- self.action // send action to tree
		}

		for _, child := range self.children {
			child.run(actionChan)
		}
	}

	if self.isRoot { // when all children of root have run then close action channel
		close(actionChan)
	}
}
