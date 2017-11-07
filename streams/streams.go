package streams

import (
	"goku-bot"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

type AllStreams map[string]chan interface{}
type AllValues map[string]interface{}

var streams AllStreams
var values AllValues

func Init() {
	streams = make(AllStreams)
	values = make(AllValues)

	streams["poloniex-OrderBook-USDT_BTC"] = make(chan interface{})
	streams["poloniex-Ticker-USDT_BTC"] = make(chan interface{})
}

func Start() {
	for {
		broadcastStreams()
		<-time.After(1 * time.Second)
	}
}

func BroadcastStream(key string, value interface{}) {
	streams[key] <- value
}

func BroadcastOrderBook(key string, value goku_bot.OrderBook) {
	log.WithFields(log.Fields{"key": key, "module": "streams"}).Debug("broadcasting order book")
	values[key] = value
}

func BroadcastTicker(key string, value goku_bot.Ticker) {
	log.WithFields(log.Fields{"key": key, "module": "streams"}).Debug("broadcasting ticker")
	values[key] = value
}

func broadcastStreams() {
	log.WithFields(log.Fields{"module": "streams", "numOfValues": len(values)}).Debug("broadcasting streams")
	for key, value := range values {
		select {
		case streams[key] <- value:
		case <-time.After(1 * time.Nanosecond):
		}
	}
}

func ReadStream(key string) interface{} {
	value := <-streams[key]

	return value
}

func ReadOrderBookStream(exchange, pair string) goku_bot.OrderBook {
	log.WithFields(log.Fields{"exchange": exchange, "pair": pair, "module": "streams"}).Debug("reading order book")
	key := fmt.Sprintf("%s-OrderBook-%s", exchange, pair)

	streamOutput := <-streams[key]
	orderBook := streamOutput.(goku_bot.OrderBook)

	return orderBook
}

func ReadTickerStream(exchange, pair string) goku_bot.Ticker {
	log.WithFields(log.Fields{"exchange": exchange, "pair": pair, "module": "streams"}).Debug("reading ticker stream")
	key := fmt.Sprintf("%s-Ticker-%s", exchange, pair)

	streamOutput := <-streams[key]
	ticker := streamOutput.(goku_bot.Ticker)

	return ticker
}

// Todo add functions for reading specific streams
