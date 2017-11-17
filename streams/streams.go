package streams

import (
	"goku-bot"
	"fmt"
	log "github.com/inconshreveable/log15"
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

func broadcastStreams() {
	log.Debug("broadcasting streams", "module", "streams", "lenValues", len(values))
	for key, value := range values {
		select {
		case streams[key] <- value:
		case <-time.After(1 * time.Nanosecond):
		}
	}
}

func BroadcastOrderBook(key string, value goku_bot.OrderBook) {
	log.Debug("broadcasting order book", "module", "streams", "key", key)
	values[key] = value
}

func BroadcastTicker(key string, value goku_bot.Ticker) {
	log.Debug("broadcasting ticker", "module", "streams", "key", key)

	values[key] = value
}

func ReadOrderBook(exchange, pair string) goku_bot.OrderBook {
	log.Debug("reading order book", "module", "streams", "exchange", exchange, "pair", pair)
	key := fmt.Sprintf("%s-OrderBook-%s", exchange, pair)

	streamOutput := <-streams[key]
	orderBook := streamOutput.(goku_bot.OrderBook)

	return orderBook
}

func ReadTicker(exchange, pair string) goku_bot.Ticker {
	log.Debug("reading ticker", "module", "streams", "exchange", exchange, "pair", pair)
	key := fmt.Sprintf("%s-Ticker-%s", exchange, pair)

	streamOutput := <-streams[key]
	ticker := streamOutput.(goku_bot.Ticker)

	return ticker
}

// Todo add functions for reading specific streams
