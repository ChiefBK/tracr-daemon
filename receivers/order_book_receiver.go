package receivers

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"bytes"
	"strconv"
	"sync"
	log "github.com/inconshreveable/log15"
	"goku-bot/keys"
	"goku-bot/exchanges"
	"goku-bot/pairs"
	"goku-bot/exchanges/poloniex"
)

type PoloniexOrderBookReceiver struct {
	exchangePair string
	orderBook    *poloniex.OrderBook
}

func NewPoloniexOrderBookReceiver(exchangePair string) *PoloniexOrderBookReceiver {
	var asks = &poloniex.OrderBookAsks{make(map[float64]poloniex.OrderBookEntry), sync.Mutex{}}
	var bids = &poloniex.OrderBookBids{make(map[float64]poloniex.OrderBookEntry), sync.Mutex{}}
	orderBook := &poloniex.OrderBook{asks, bids}

	return &PoloniexOrderBookReceiver{exchangePair, orderBook}
}

func (self *PoloniexOrderBookReceiver) Key() string {
	stdPair, err := pairs.StandardPair(self.exchangePair, exchanges.POLONIEX)

	if err != nil {
		log.Error("could not find standard pair. Are you sure the exchange pair is correct?", "module", "receivers", "exchangePair", self.exchangePair)
		return ""
	}

	return keys.BuildOrderBookKey(exchanges.POLONIEX, stdPair)
}

func (self *PoloniexOrderBookReceiver) Start() {
	address := "api2.poloniex.com"

	connection, err := websocketConnect(address, 5)

	if err != nil {
		log.Error("could not connect to Poloniex Order book", "key", self.Key(), "module", "receivers", "error", err)
		return
	}

	defer connection.Close()

	cmdString := []byte("{\"command\" : \"subscribe\", \"channel\" : \"" + self.exchangePair + "\"}")
	err = connection.WriteMessage(websocket.TextMessage, cmdString)
	if err != nil {
		log.Error("there was an error writing command", "key", self.Key(), "module", "receivers", "error", err)

		return
	}

	isFirst := true
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Warn("there was an error reading message", "key", self.Key(), "module", "receivers", "error", err)
			return
		}
		dec := json.NewDecoder(bytes.NewReader(message))

		var m interface{}

		// decode an array value (Message)
		err = dec.Decode(&m)
		if err != nil {
			log.Warn("error decoding poloniex order book message", "key", self.Key(), "module", "receivers", "error", err)
		}

		if len(m.([]interface{})) <= 1 {
			continue
		}

		seq := int(m.([]interface{})[1].(float64))

		if isFirst { //if full order book
			ob := m.([]interface{})[2].([]interface{})[0].([]interface{})[1].(map[string]interface{})

			// Store asks and bids
			var asks []poloniex.OrderBookEntry
			for k, v := range ob["orderBook"].([]interface{})[0].(map[string]interface{}) {
				price, _ := strconv.ParseFloat(k, 64)
				value, _ := strconv.ParseFloat(v.(string), 64)
				entry := poloniex.OrderBookEntry{self.exchangePair, seq, price, value}

				asks = append(asks, entry)
			}
			self.syncAsks(asks)

			var bids []poloniex.OrderBookEntry
			for k, v := range ob["orderBook"].([]interface{})[1].(map[string]interface{}) {
				price, _ := strconv.ParseFloat(k, 64)
				value, _ := strconv.ParseFloat(v.(string), 64)
				entry := poloniex.OrderBookEntry{self.exchangePair, seq, price, value}

				bids = append(bids, entry)
			}
			self.syncBids(bids)
			isFirst = false
			log.Info("initial sync of full order book complete", "key", self.Key(), "module", "receivers")
		} else { // if a set of changes
			changes := m.([]interface{})[2].([]interface{})

			var asks []poloniex.OrderBookEntry
			var bids []poloniex.OrderBookEntry

			for i := range changes {
				change := changes[i].([]interface{})
				// The starting index changes based on the first element. If it's "t" there's an extra element
				// and the index is bumped up one. TODO - figure out what the "t" events are for
				first := change[0].(string)
				index := 1
				if first == "t" {
					index++
				}

				var isAsk bool = change[index].(float64) == 1
				index++
				price, _ := strconv.ParseFloat(change[index].(string), 64)
				index++
				amount, _ := strconv.ParseFloat(change[index].(string), 64)

				entry := poloniex.OrderBookEntry{self.exchangePair, seq, price, amount}

				if isAsk {
					asks = append(asks, entry)
				} else {
					bids = append(bids, entry)
				}
			}

			self.syncBids(bids)
			self.syncAsks(asks)
			log.Debug("received Orderbook update", "key", self.Key(), "module", "receivers")
		}

		self.broadcastOrderBook()
	}
}

func (self *PoloniexOrderBookReceiver) broadcastOrderBook() {
	sendToProcessor(self.Key(), *self.orderBook.DeepCopy())
}

func (self *PoloniexOrderBookReceiver) syncAsks(orders []poloniex.OrderBookEntry) {
	self.orderBook.Asks.Lock.Lock()
	defer self.orderBook.Asks.Lock.Unlock()
	for _, order := range orders {
		if order.Amount == 0 {
			delete(self.orderBook.Asks.Orders, order.Amount)
			continue
		}
		self.orderBook.Asks.Orders[order.Price] = order
	}
}

func (self *PoloniexOrderBookReceiver) syncBids(orders []poloniex.OrderBookEntry) {
	self.orderBook.Bids.Lock.Lock()
	defer self.orderBook.Bids.Lock.Unlock()
	for _, order := range orders {
		if order.Amount == 0 {
			delete(self.orderBook.Bids.Orders, order.Amount)
			continue
		}
		self.orderBook.Bids.Orders[order.Price] = order
	}
}
