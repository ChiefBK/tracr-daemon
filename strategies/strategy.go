package strategies

import (
	"goku-bot"
	"log"
)

func Strategy1(pair, exchange string, indicator *goku_bot.Indicator, store *goku_bot.MgoStore) (actionQueue goku_bot.ActionQueue, err error) {
	log.Println("Running Strategy 1")

	return
}
