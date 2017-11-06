package receivers

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"bytes"
	"strconv"
	"sync"
	"fmt"
	"goku-bot"
	log "github.com/sirupsen/logrus"
)

type OrderBookReceiver struct {
	Exchange  string
	Pair      string
	orderBook *goku_bot.OrderBook
}

func NewOrderBookReceiver(exchange, pair string) *OrderBookReceiver {
	var asks = &goku_bot.OrderBookAsks{make(map[float64]goku_bot.OrderBookEntry), sync.Mutex{}}
	var bids = &goku_bot.OrderBookBids{make(map[float64]goku_bot.OrderBookEntry), sync.Mutex{}}
	orderBook := &goku_bot.OrderBook{asks, bids}

	return &OrderBookReceiver{exchange, pair, orderBook}
}

func (self *OrderBookReceiver) Key() string {
	return fmt.Sprintf("%s-OrderBook-%s", self.Exchange, self.Pair)
}

func (self *OrderBookReceiver) Start() {
	address := "api2.poloniex.com:443"

	connection, err := websocketConnect(address, 5)

	if err != nil {
		log.WithFields(log.Fields{"key": self.Key(), "module": "receivers", "error": err}).Error("could not connect to Poloniex Order book")
		return
	}

	defer connection.Close()

	cmdString := []byte("{\"command\" : \"subscribe\", \"channel\" : \"" + self.Pair + "\"}")
	err = connection.WriteMessage(websocket.TextMessage, cmdString)
	if err != nil {
		log.WithFields(log.Fields{"key": self.Key(), "module": "receivers", "error": err}).Error("there was an error writing command")
		return
	}

	isFirst := true
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.WithFields(log.Fields{"key": self.Key(), "module": "receivers", "error": err}).Warn("there was an error reading message")
			return
		}
		dec := json.NewDecoder(bytes.NewReader(message))

		var m interface{}

		// decode an array value (Message)
		err = dec.Decode(&m)
		if err != nil {
			log.WithFields(log.Fields{"key": self.Key(), "module": "receivers", "error": err}).Warn("error decoding poloniex order book message")
		}

		if len(m.([]interface{})) <= 1 {
			continue
		}

		seq := int(m.([]interface{})[1].(float64))

		if isFirst { //if full order book
			ob := m.([]interface{})[2].([]interface{})[0].([]interface{})[1].(map[string]interface{})

			// Store asks and bids
			var asks []goku_bot.OrderBookEntry
			for k, v := range ob["orderBook"].([]interface{})[0].(map[string]interface{}) {
				price, _ := strconv.ParseFloat(k, 64)
				value, _ := strconv.ParseFloat(v.(string), 64)
				entry := goku_bot.OrderBookEntry{self.Pair, seq, price, value}

				asks = append(asks, entry)
			}
			self.syncAsks(asks)

			var bids []goku_bot.OrderBookEntry
			for k, v := range ob["orderBook"].([]interface{})[1].(map[string]interface{}) {
				price, _ := strconv.ParseFloat(k, 64)
				value, _ := strconv.ParseFloat(v.(string), 64)
				entry := goku_bot.OrderBookEntry{self.Pair, seq, price, value}

				bids = append(bids, entry)
			}
			self.syncBids(bids)
			isFirst = false
			log.WithFields(log.Fields{"key": self.Key(), "module": "receivers"}).Info("initial sync of full order book complete")
		} else { // if a set of changes
			changes := m.([]interface{})[2].([]interface{})

			var asks []goku_bot.OrderBookEntry
			var bids []goku_bot.OrderBookEntry

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

				entry := goku_bot.OrderBookEntry{self.Pair, seq, price, amount}

				if isAsk {
					asks = append(asks, entry)
				} else {
					bids = append(bids, entry)
				}
			}

			self.syncBids(bids)
			self.syncAsks(asks)
			log.WithFields(log.Fields{"key": self.Key(), "module": "receivers"}).Debug("recieved Orderbook update")
		}

		self.broadcastOrderBook()
	}
}

func (self *OrderBookReceiver) broadcastOrderBook() {
	sendToProcessor(self.Key(), self.orderBook.DeepCopy())
}

func (self *OrderBookReceiver) syncAsks(orders []goku_bot.OrderBookEntry) {
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

func (self *OrderBookReceiver) syncBids(orders []goku_bot.OrderBookEntry) {
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
