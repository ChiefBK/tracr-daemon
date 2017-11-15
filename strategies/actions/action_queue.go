package actions

import log "github.com/sirupsen/logrus"

type ActionQueue struct {
	Queue []*Action
}

func NewActionQueue() *ActionQueue {
	var queue []*Action
	return &ActionQueue{queue}
}

func (aq *ActionQueue) Push(action *Action) {
	aq.Queue = append(aq.Queue, action)
}

func (aq *ActionQueue) Dequeue() *Action {
	log.WithFields(log.Fields{"module": "actions"}).Debug("dequeuing")

	if len(aq.Queue) < 1 {
		return nil
	}

	action := aq.Queue[0]
	aq.Queue = aq.Queue[1:]

	log.WithFields(log.Fields{"module": "actions", "len": len(aq.Queue)}).Debug("length of queue")

	return action
}

func (self ActionQueue) Length() int {
	return len(self.Queue)
}
