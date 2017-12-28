package exchanges

import (
	"sync"
)

type OrderBook struct {
	Exchange string
	Pair     string
	Asks     map[string]float64 // key is price; value is volume
	Bids     map[string]float64
	asksLock sync.Mutex
	bidsLock sync.Mutex
}

func NewOrderBook(exchange, pair string) *OrderBook {
	return &OrderBook{exchange, pair, make(map[string]float64), make(map[string]float64), sync.Mutex{}, sync.Mutex{}}
}

type OrderBookResponse struct {
	Data OrderBook
	Err  error
}

func (self *OrderBook) SyncAsks(orders map[string]float64) {
	self.asksLock.Lock()
	defer self.asksLock.Unlock()
	for price, volume := range orders {
		if volume == 0 {
			delete(self.Asks, price)
			continue
		}
		self.Asks[price] = volume
	}
}

func (self *OrderBook) SyncBids(orders map[string]float64) {
	self.bidsLock.Lock()
	defer self.bidsLock.Unlock()
	for price, volume := range orders {
		if volume == 0 {
			delete(self.Bids, price)
			continue
		}
		self.Bids[price] = volume
	}
}

func (self *OrderBook) DeepCopy() *OrderBook {
	self.asksLock.Lock()
	self.bidsLock.Lock()
	defer self.asksLock.Unlock()
	defer self.bidsLock.Unlock()

	askOrders := make(map[string]float64)
	bidOrders := make(map[string]float64)

	for price, volume := range self.Asks {
		askOrders[price] = volume
	}
	for price, volume := range self.Bids {
		bidOrders[price] = volume
	}

	orderBook := &OrderBook{self.Exchange, self.Pair, askOrders, bidOrders, sync.Mutex{}, sync.Mutex{}}

	return orderBook
}

// TODO
func (self OrderBook) GetAsksAscending()  {
	//var keys []float64
	//
	//log.Printf("There are %d Asks", len(self.Asks.Orders))
	//
	//for k := range self.Asks.Orders {
	//	keys = append(keys, k)
	//}
	//
	//sort.Float64s(keys)
	//for _, k := range keys {
	//	asks = append(asks, self.Asks.Orders[k])
	//}

	return
}

// TODO
func (self OrderBook) GetBidsDescending() {
	//var keys []float64
	//
	//log.Printf("There are %d Bids", len(self.Bids.Orders))
	//
	//for k := range self.Bids.Orders {
	//	keys = append(keys, k)
	//}
	//
	//sort.Float64s(keys)
	//for i := len(keys) - 1; i >= 0; i-- {
	//	bids = append(bids, self.Bids.Orders[keys[i]])
	//}

	return
}
