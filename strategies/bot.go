package strategies

import (
	"goku-bot"
	log "github.com/sirupsen/logrus"
)

const (
	CLOSED_POSITION = iota
	LONG_POSITION   = iota
	SHORT_POSITION  = iota
)

type Position int // equal to one of the POSITION constants above

type Bot struct {
	name       string // must be unique amongst bots
	exchange   string
	pair       string
	position   Position
	strategies map[Position]*Strategy
}

func NewBot(name, exchange, pair string) (bot *Bot) {
	bot = new(Bot)
	bot.strategies = make(map[Position]*Strategy)
	bot.exchange = exchange
	bot.pair = pair
	bot.name = name
	bot.position = CLOSED_POSITION

	return
}

func (self *Bot) start() {
	log.WithFields(log.Fields{"bot": self.name, "module": "strategies"}).Info("Starting bot")
	var actionChan = make(chan *goku_bot.ActionQueue)
	self.runStrategy(actionChan)

	actionQueue := <- actionChan

	//run actions
	for _, action := range actionQueue.Queue {
		log.WithFields(log.Fields{"action": action.Action}).Debug("Running action")
	}

}

func (self *Bot) addStrategy(pos Position, strategy *Strategy) {
	self.strategies[pos] = strategy
}

func (self *Bot) runStrategy(actionChan chan<- *goku_bot.ActionQueue) {
	go self.strategies[self.position].run(actionChan)
}
