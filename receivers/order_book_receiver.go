package receivers

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"bytes"
	"strconv"
	log "github.com/inconshreveable/log15"
	"goku-bot/keys"
	"goku-bot/exchanges"
	"goku-bot/pairs"
)

type PoloniexOrderBookReceiver struct {
	pair      string
	orderBook *exchanges.OrderBook
}

func NewPoloniexOrderBookReceiver(pair string) *PoloniexOrderBookReceiver {
	orderBook := exchanges.NewOrderBook(exchanges.POLONIEX, pair)
	return &PoloniexOrderBookReceiver{pair, orderBook}
}

func (self *PoloniexOrderBookReceiver) Key() string {
	return keys.BuildOrderBookKey(exchanges.POLONIEX, self.pair)
}

func (self *PoloniexOrderBookReceiver) Start() {
	exchangePair, _ := pairs.ExchangePair(self.pair, exchanges.POLONIEX)
	address := "api2.poloniex.com"

	connection, err := websocketConnect(address, 5)

	if err != nil {
		log.Error("could not connect to Poloniex Order book", "key", self.Key(), "module", "receivers", "error", err)
		return
	}

	defer connection.Close()

	cmdString := []byte("{\"command\" : \"subscribe\", \"channel\" : \"" + exchangePair + "\"}")
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

		//seq := int(m.([]interface{})[1].(float64))

		if isFirst { //if full order book
			ob := m.([]interface{})[2].([]interface{})[0].([]interface{})[1].(map[string]interface{})

			// Store asks and bids
			asks := make(map[float64]float64)
			for k, v := range ob["orderBook"].([]interface{})[0].(map[string]interface{}) {
				price, _ := strconv.ParseFloat(k, 64)
				value, _ := strconv.ParseFloat(v.(string), 64)
				asks[price] = value
			}
			self.orderBook.SyncAsks(asks)

			bids := make(map[float64]float64)
			for k, v := range ob["orderBook"].([]interface{})[1].(map[string]interface{}) {
				price, _ := strconv.ParseFloat(k, 64)
				value, _ := strconv.ParseFloat(v.(string), 64)
				bids[price] = value
			}
			self.orderBook.SyncBids(bids)
			isFirst = false
			log.Info("initial sync of full order book complete", "key", self.Key(), "module", "receivers")
		} else { // if a set of changes
			changes := m.([]interface{})[2].([]interface{})

			asks := make(map[float64]float64)
			bids := make(map[float64]float64)

			for i := range changes {
				change := changes[i].([]interface{})
				// The starting index changes based on the first element. If it's "t" there's an extra element
				// and the index is bumped up one. TODO - figure out what the "t" events are for
				first := change[0].(string)
				index := 1
				if first == "t" {
					index++
				}

				var isAsk = change[index].(float64) == 1
				index++
				price, _ := strconv.ParseFloat(change[index].(string), 64)
				index++
				amount, _ := strconv.ParseFloat(change[index].(string), 64)

				if isAsk {
					asks[price] = amount
				} else {
					bids[price] = amount
				}
			}

			self.orderBook.SyncBids(bids)
			self.orderBook.SyncAsks(asks)
			log.Debug("received Orderbook update", "key", self.Key(), "module", "receivers")
		}

		self.broadcastOrderBook()
	}
}

func (self *PoloniexOrderBookReceiver) broadcastOrderBook() {
	sendToProcessor(self.Key(), *self.orderBook.DeepCopy())
}
