package goku_bot

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync"
	"time"
)

type OrderBookSteward struct {
	Pair     string
	Exchange string
	Store    *Store
}

type OrderBookEntry struct {
	CurrencyPair string
	Sequence     int
	Price        float64
	Amount       float64
}

func (self *OrderBookSteward) ConnectPoloniexOrderBook(pair string, group *sync.WaitGroup) {
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
			self.Store.SyncAsks("poloniex", pair, asks)

			var bids []OrderBookEntry
			for k, v := range ob["orderBook"].([]interface{})[1].(map[string]interface{}) {
				price, _ := strconv.ParseFloat(k, 64)
				value, _ := strconv.ParseFloat(v.(string), 64)
				entry := OrderBookEntry{pair, seq, price, value}

				bids = append(bids, entry)
			}
			self.Store.SyncBids("poloniex", pair, bids)
			isFirst = false
			log.Printf("Initial sync of full order book complete at %s", time.Now())
			group.Done() // After storing full order book - mark as done
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

			self.Store.SyncBids("poloniex", pair, bids)
			self.Store.SyncAsks("poloniex", pair, asks)
			log.Printf("Recieved Orderbook update")
		}
	}
}
