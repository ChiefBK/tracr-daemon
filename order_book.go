package goku_bot

import (
	"log"
	"sort"
	"sync"
)

type OrderBook struct {
	Asks *OrderBookAsks
	Bids *OrderBookBids
}

type OrderBookEntry struct {
	CurrencyPair string
	Sequence     int
	Price        float64
	Amount       float64
}


type OrderBookAsks struct {
	Orders map[float64]OrderBookEntry
	Lock   sync.Mutex
}

type OrderBookBids struct {
	Orders map[float64]OrderBookEntry
	Lock   sync.Mutex
}

func (self *OrderBook) DeepCopy() *OrderBook {
	self.Asks.Lock.Lock()
	self.Bids.Lock.Lock()
	defer self.Asks.Lock.Unlock()
	defer self.Bids.Lock.Unlock()

	askOrders := make(map[float64]OrderBookEntry)
	bidOrders := make(map[float64]OrderBookEntry)

	for key, order := range self.Asks.Orders {
		askOrders[key] = order
	}
	for key, order := range self.Bids.Orders {
		bidOrders[key] = order
	}

	asks := &OrderBookAsks{Orders: askOrders}
	bids := &OrderBookBids{Orders: bidOrders}

	orderBook := &OrderBook{asks, bids}

	return orderBook
}

func (self OrderBook) GetAsksAscending() (asks []OrderBookEntry) {
	var keys []float64

	log.Printf("There are %d Asks", len(self.Asks.Orders))

	for k := range self.Asks.Orders {
		keys = append(keys, k)
	}

	sort.Float64s(keys)
	for _, k := range keys {
		asks = append(asks, self.Asks.Orders[k])
	}

	return
}

func (self OrderBook) GetBidsDescending() (bids []OrderBookEntry) {
	var keys []float64

	log.Printf("There are %d Bids", len(self.Bids.Orders))

	for k := range self.Bids.Orders {
		keys = append(keys, k)
	}

	sort.Float64s(keys)
	for i := len(keys) - 1; i >= 0; i-- {
		bids = append(bids, self.Bids.Orders[keys[i]])
	}

	return
}

