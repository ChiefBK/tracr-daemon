package streams

import (
	"goku-bot"
	"fmt"
	log "github.com/inconshreveable/log15"
	"time"
	"goku-bot/keys"
	"goku-bot/exchanges"
)

var streams map[string]chan interface{}
var values map[string]interface{}

func Init() {
	streams = make(map[string]chan interface{})
	values = make(map[string]interface{})
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

func PutValue(key string, value interface{}) {
	log.Debug("putting value", "module", "streams", "key", key, "value", value)
	if _, ok := values[key]; !ok { // if values does not contain key than create a stream channel
		log.Debug("adding key to stream channels", "module", "streams", "key", key)
		streams[key] = make(chan interface{})
	}
	values[key] = value
}

func ReadBalance(exchange, currency string) float64 {
	key := keys.BuildBalancesKey(exchange)
	waitForChannelInitialization(key)

	log.Debug("reading stream", "module", "streams", "key", key)

	streamOutput := <-streams[key]
	balances := streamOutput.(exchanges.Balances)
	return balances[currency]
}

func ReadOrderBook(exchange, pair string) exchanges.OrderBook {
	log.Debug("reading order book", "module", "streams", "exchange", exchange, "pair", pair)
	key := keys.BuildOrderBookKey(exchange, pair)
	waitForChannelInitialization(key)

	streamOutput := <-streams[key]
	orderBook := streamOutput.(exchanges.OrderBook)

	return orderBook
}

func ReadTicker(exchange, pair string) goku_bot.Ticker {
	log.Debug("reading ticker", "module", "streams", "exchange", exchange, "pair", pair)
	key := fmt.Sprintf("%s-Ticker-%s", exchange, pair)

	streamOutput := <-streams[key]
	ticker := streamOutput.(goku_bot.Ticker)

	return ticker
}

func waitForChannelInitialization(key string) {
	if _, ok := streams[key]; !ok { // if stream channel hasn't been initialized yet; wait until it is
		for !ok {
			<-time.After(1 * time.Second)
			_, ok = streams[key]
		}
	}
}
