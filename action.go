package goku_bot

type Action struct {
	Action string
	Id     string
}

func NewAction(id, action string) *Action {
	return &Action{action, id}
}

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

func (aq *ActionQueue) Pop() *Action {
	if len(aq.Queue) < 1 {
		return nil
	}

	action := aq.Queue[0]
	aq.Queue = aq.Queue[1:]

	return action
}
