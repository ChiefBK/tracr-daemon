package executors

import "goku-bot"

type ActionChannel chan goku_bot.Action

func (self ActionChannel) start() {
	for {
		action := <-self
		processAction(action)
	}
}

var actionChannels = make(map[string]ActionChannel) // Channel which receives actions from bots

func AddActionChannel(botKey string) {
	actionChannels[botKey] = make(ActionChannel)
	actionChannels[botKey].start()
}
