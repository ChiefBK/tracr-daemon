package broker

import (
	log "github.com/inconshreveable/log15"
	"tracr/actions"
	"tracr/executors/responses"
)

type ActionReceiverChannel chan actions.ActionQueue

var BotResponseChannels = make(map[string]chan responses.ExecutorResponse) // Channel which receives response from executor
var ActionReceiverChannels = make(map[string]ActionReceiverChannel)        // Channel which receives actions from bots

var logger = log.New("module", "broker")

// Used by bot to send actions to executors module
func SendToExecutor(botKey string, queue actions.ActionQueue) {
	log.Debug("sending action to receiver", "module", "broker", "botKey", botKey)
	if _, ok := ActionReceiverChannels[botKey]; ok { // if action receivers channels contains bot key
		ActionReceiverChannels[botKey] <- queue
	} else {
		log.Warn("Could not find channel with bot key. Intent channel hasn't been added to action receiver", "module", "broker", "botKey", botKey)
	}
}

// Adds a channel to executors module to handle Actions from bots
func AddActionReceiverChannel(botKey string) {
	log.Debug("adding action channel for bot", "module", "broker", "botKey", botKey)
	ActionReceiverChannels[botKey] = make(ActionReceiverChannel)
}

func SendResponseToBot(botKey string, response responses.ExecutorResponse) {
	log.Debug("sending response to bot", "module", "broker", "botKey", botKey)
	channel := GetBotResponseChannel(botKey)
	channel <- response
}

func GetBotResponseChannel(botKey string) chan responses.ExecutorResponse {
	log.Debug("getting bot response channel", "module", "broker", "botKey", botKey)
	return BotResponseChannels[botKey]
}
