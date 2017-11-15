package broker

import (
	log "github.com/sirupsen/logrus"
	"goku-bot/strategies/actions"
	"goku-bot/executors/responses"
)

type ActionReceiverChannel chan actions.ActionQueue

var BotResponseChannels = make(map[string]chan responses.ExecutorResponse) // Channel which receives response from executor
var ActionReceiverChannels = make(map[string]ActionReceiverChannel)        // Channel which receives actions from bots

// Used by bot to send actions to executors module
func SendToExecutor(botKey string, queue actions.ActionQueue) {
	log.WithFields(log.Fields{"module": "executors", "botKey": botKey}).Debug("sending action to receiver")
	if _, ok := ActionReceiverChannels[botKey]; ok { // if action receivers channels contains bot key
		ActionReceiverChannels[botKey] <- queue
	} else {
		log.WithFields(log.Fields{"module": "executors", "botKey": botKey}).Warn("Could not find channel with bot key. Intent channel hasn't been added to action receiver")
	}
}

// Adds a channel to executors module to handle Actions from bots
func AddActionReceiverChannel(botKey string) {
	log.WithFields(log.Fields{"module": "executors", "botKey": botKey}).Debug("adding action channel for bot")
	ActionReceiverChannels[botKey] = make(ActionReceiverChannel)
}

func SendResponseToBot(botKey string, response responses.ExecutorResponse) {
	channel := GetBotResponseChannel(botKey)
	channel <- response
}

func GetBotResponseChannel(botKey string) chan responses.ExecutorResponse {
	return BotResponseChannels[botKey]
}
