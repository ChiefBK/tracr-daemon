package goku_bot

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
	"sort"
	"sync"
)

type OrderBookSteward struct {
	Pair      string
	Exchange  string
	orderBook *OrderBook
	channel   chan<- OrderBook
}

type OrderBookEntry struct {
	CurrencyPair string
	Sequence     int
	Price        float64
	Amount       float64
}

var OrderBookChannels = make(map[string]chan OrderBook) // Maps all order book Pairs

func NewOrderBookSteward(pair, exchange string) *OrderBookSteward {
	var asks = &OrderBookAsks{orders: make(map[float64]OrderBookEntry)}
	var bids = &OrderBookBids{orders: make(map[float64]OrderBookEntry)}
	orderBook := &OrderBook{asks: asks, bids: bids}

	OrderBookChannels[pair] = make(chan OrderBook)

	return &OrderBookSteward{pair, exchange, orderBook, OrderBookChannels[pair]}
}

type OrderBook struct {
	asks     *OrderBookAsks
	bids     *OrderBookBids
}

type OrderBookAsks struct {
	orders map[float64]OrderBookEntry
	lock sync.Mutex
}

type OrderBookBids struct {
	orders map[float64]OrderBookEntry
	lock sync.Mutex
}

func (self *OrderBook) deepCopy() *OrderBook {
	self.asks.lock.Lock()
	self.bids.lock.Lock()
	defer self.asks.lock.Unlock()
	defer self.bids.lock.Unlock()

	askOrders := make(map[float64]OrderBookEntry)
	bidOrders := make(map[float64]OrderBookEntry)

	for key, order := range self.asks.orders {
		askOrders[key] = order
	}
	for key, order := range self.bids.orders {
		bidOrders[key] = order
	}

	asks := &OrderBookAsks{orders: askOrders}
	bids := &OrderBookBids{orders: bidOrders}

	orderBook := &OrderBook{asks, bids}

	return orderBook
}

func (self OrderBook) GetAsksAscending() (asks []OrderBookEntry) {
	var keys []float64

	log.Printf("There are %d asks", len(self.asks.orders))

	for k := range self.asks.orders {
		keys = append(keys, k)
	}

	sort.Float64s(keys)
	for _, k := range keys {
		asks = append(asks, self.asks.orders[k])
	}

	return
}

func (self OrderBook) GetBidsDescending() (bids []OrderBookEntry) {
	var keys []float64

	log.Printf("There are %d bids", len(self.bids.orders))

	for k := range self.bids.orders {
		keys = append(keys, k)
	}

	sort.Float64s(keys)
	for i := len(keys) - 1; i >= 0; i-- {
		bids = append(bids, self.bids.orders[keys[i]])
	}

	return
}

func (self *OrderBookSteward) broadcastOrderBook() {
	self.channel <- *self.orderBook.deepCopy()
}

func (self *OrderBookSteward) getBidsDecending() (bids []OrderBookEntry) {
	self.orderBook.bids.lock.Lock()
	defer self.orderBook.bids.lock.Unlock()
	return self.orderBook.GetBidsDescending()
}

func (self *OrderBookSteward) getAsksAscending() (asks []OrderBookEntry) {
	self.orderBook.asks.lock.Lock()
	defer self.orderBook.asks.lock.Unlock()
	return self.orderBook.GetAsksAscending()
}

func (self *OrderBookSteward) syncAsks(orders []OrderBookEntry) {
	self.orderBook.asks.lock.Lock()
	defer self.orderBook.asks.lock.Unlock()
	for _, order := range orders {
		if order.Amount == 0 {
			delete(self.orderBook.asks.orders, order.Amount)
			continue
		}
		self.orderBook.asks.orders[order.Price] = order
	}
}

func (self *OrderBookSteward) syncBids(orders []OrderBookEntry) {
	self.orderBook.bids.lock.Lock()
	defer self.orderBook.bids.lock.Unlock()
	for _, order := range orders {
		if order.Amount == 0 {
			delete(self.orderBook.bids.orders, order.Amount)
			continue
		}
		self.orderBook.bids.orders[order.Price] = order
	}
}

func (self *OrderBookSteward) ConnectPoloniexOrderBook(pair string) {
	address := "api2.poloniex.com:443"

	connection, err := websocketConnect(address, 5)

	if err != nil {
		log.Printf("Could not connect to Poloniex Order book: %s", err)
		return
	}

	defer connection.Close()

	cmdString := []byte("{\"command\" : \"subscribe\", \"channel\" : \"" + pair + "\"}")
	err = connection.WriteMessage(websocket.TextMessage, cmdString)
	if err != nil {
		log.Printf("there was an error writing command: %s", err)
		return
	}

	isFirst := true
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		dec := json.NewDecoder(bytes.NewReader(message))

		var m interface{}

		// decode an array value (Message)
		err = dec.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}

		if len(m.([]interface{})) <= 1 {
			continue
		}

		seq := int(m.([]interface{})[1].(float64))

		if isFirst { //if full order book
			ob := m.([]interface{})[2].([]interface{})[0].([]interface{})[1].(map[string]interface{})

			// Store asks and bids
			var asks []OrderBookEntry
			for k, v := range ob["orderBook"].([]interface{})[0].(map[string]interface{}) {
				price, _ := strconv.ParseFloat(k, 64)
				value, _ := strconv.ParseFloat(v.(string), 64)
				entry := OrderBookEntry{pair, seq, price, value}

				asks = append(asks, entry)
			}
			self.syncAsks(asks)

			var bids []OrderBookEntry
			for k, v := range ob["orderBook"].([]interface{})[1].(map[string]interface{}) {
				price, _ := strconv.ParseFloat(k, 64)
				value, _ := strconv.ParseFloat(v.(string), 64)
				entry := OrderBookEntry{pair, seq, price, value}

				bids = append(bids, entry)
			}
			self.syncBids(bids)
			isFirst = false
			log.Printf("Initial sync of full order book complete at %s", time.Now())
		} else { // if a set of changes
			changes := m.([]interface{})[2].([]interface{})

			var asks []OrderBookEntry
			var bids []OrderBookEntry

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

				entry := OrderBookEntry{pair, seq, price, amount}

				if isAsk {
					asks = append(asks, entry)
				} else {
					bids = append(bids, entry)
				}
			}

			self.syncBids(bids)
			self.syncAsks(asks)
			log.Printf("Recieved Orderbook update")
		}

		self.broadcastOrderBook()
	}
}
