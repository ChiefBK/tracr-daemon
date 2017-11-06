package strategies

import (
	"goku-bot"
	"log"
	"goku-bot/store"
)

func Strategy1(pair, exchange string, indicator *goku_bot.Indicator, store store.Store) (actionQueue goku_bot.ActionQueue, err error) {
	log.Println("Running Strategy 1")

	return
}
