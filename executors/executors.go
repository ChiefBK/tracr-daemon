package executors

import (
	"goku-bot/strategies/actions"
	. "goku-bot/strategies/actions"
	"fmt"
	log "github.com/sirupsen/logrus"
	"goku-bot/executors/responses"
	"goku-bot/broker"
	"time"
)

// TODO - create third action receiver / sender module to link executors module with strategies module
func processActions(botKey string, actions actions.ActionQueue) {
	log.WithFields(log.Fields{"module": "executors", "botKey": botKey, "actionLen": actions.Length()}).Debug("executor processing actions")

	action := actions.Dequeue()
	for action != nil { // While actions is still contains actions
		switch action.Intent {
		case OPEN_SHORT_POSITION:
			response := responses.ExecutorResponse{botKey, "short was opened", action.Id}
			broker.SendResponseToBot(botKey, response)
		case CLOSE_POSITION:
			fmt.Println("closing position")
			response := responses.ExecutorResponse{botKey, "position closed", action.Id}
			broker.SendResponseToBot(botKey, response)
		case OPEN_LONG_POSITION:
			fmt.Println("opening long")
			response := responses.ExecutorResponse{botKey, "long was opened", action.Id}
			broker.SendResponseToBot(botKey, response)
		}

		action = actions.Dequeue()
	}
}

// Sends market buy
func openLongPosition() {

}

var actionChannelsMonitored []string

func inActionChannelsList(botKey string) bool {
	for _, key := range actionChannelsMonitored {
		if botKey == key {
			return true
		}
	}
	return false
}

func Start() {
	for {
		for botKey, channel := range broker.ActionReceiverChannels { // monitor action receiver channels
			if !inActionChannelsList(botKey) { // if not currently monitoring channel
				actionChannelsMonitored = append(actionChannelsMonitored, botKey)
				go monitorActionChannel(botKey, channel)
			}
		}
		<-time.After(5 * time.Second)
	}
}

func monitorActionChannel(botKey string, channel broker.ActionReceiverChannel) {
	log.WithFields(log.Fields{"module": "executors", "botKey": botKey}).Debug("starting action channel for bot")
	for {
		actionQueue := <-channel
		processActions(botKey, actionQueue)
	}
}
