package actions

import (
	"goku-bot/util"
)

// ActionIntents
const (
	OPEN_SHORT_POSITION = iota
	OPEN_LONG_POSITION  = iota
	CLOSE_POSITION      = iota
)

// OrderTypes
const (
	MARKET_ORDER = iota
	LIMIT_ORDER  = iota
)

// ActionConsumers
const (
	BOT      = iota
	EXECUTOR = iota
)

type ActionIntent int   // OPEN_SHORT_POSITION or OPEN_LONG_POSITION or CLOSE_POSITION or ...
type OrderType int      // MARKET_ORDER or LIMIT_ORDER or ...
type ActionConsumer int // BOT or EXECUTOR or ...

type ActionData map[string]interface{}

type Action struct {
	Intent   ActionIntent
	Id       string
	data     ActionData
	Consumer ActionConsumer
}

var ActionFunctions = make(map[string]func() *Action)

func (self *Action) SetVolume(volume float64) {
	self.data["volume"] = volume
}

func (self *Action) SetLeverage(leverage int) {
	self.data["leverage"] = leverage
}

func (self *Action) SetMargin(margin float64) {
	self.data["margin"] = margin
}

func (self *Action) SetOrderType(orderType OrderType) {
	self.data["orderType"] = orderType
}

func (self *Action) SetPair(pair string) {
	self.data["pair"] = pair
}

func (self *Action) SetExchange(exchange string) {
	self.data["exchange"] = exchange
}

func (self *Action) SetBotKey(botKey string) {
	self.data["botKey"] = botKey
}

func newAction(intent ActionIntent, data ActionData, consumer ActionConsumer) *Action {
	id := util.RandString(20)
	return &Action{intent, id, data, consumer}
}

func ShortPositionAction() *Action {
	var intent ActionIntent = OPEN_SHORT_POSITION
	data := make(ActionData)
	var consumer ActionConsumer = EXECUTOR

	return newAction(intent, data, consumer)
}
