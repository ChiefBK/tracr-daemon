package command

import (
	log "github.com/inconshreveable/log15"
	"tracr-daemon/command/actions"
)

type Signal struct {
	condition func() bool
	children  []*Signal
	action    *actions.Action
	isRoot    bool
}

func NewSignal(condition func() bool, action *actions.Action, isRoot bool) *Signal {
	var children []*Signal
	return &Signal{condition, children, action, isRoot}
}

func (self *Signal) addChild(signal *Signal) {
	self.children = append(self.children, signal)
}

func (self *Signal) run(actionChan chan<- *actions.Action) {
	log.Debug("running signal", "module", "command", "children", len(self.children), "isRoot", self.isRoot)
	result := self.condition()

	if result { // if signal is true
		if len(self.children) == 0 && self.action != nil { // if leaf node
			log.Debug("sending action from signal", "module", "command", "action", self.action)
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
