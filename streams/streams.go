package streams

import (
	"goku-bot"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type AllStreams map[string]chan interface{}

var streams AllStreams

func Init() {
	streams = make(AllStreams)

	streams["poloniex-OrderBook-USDT_BTC"] = make(chan interface{})
}

func BroadcastStream(key string, value interface{}) {
	streams[key] <- value
}

func BroadcastOrderBook(key string, value goku_bot.OrderBook) {
	log.WithFields(log.Fields{"key": key, "module": "streams"}).Debug("broadcasting order book")
	streams[key] <- value
}

func ReadStream(key string) interface{} {
	value := <-streams[key]

	return value
}

func ReadOrderBookStream(exchange, pair string) goku_bot.OrderBook{
	log.WithFields(log.Fields{"exchange": exchange, "pair": pair, "module": "streams"}).Debug("reading order book")
	key := fmt.Sprintf("%s-OrderBook-%s", exchange, pair)

	streamOutput := <- streams[key]
	orderBook := streamOutput.(goku_bot.OrderBook)

	return orderBook
}

// Todo add functions for reading specific streams
