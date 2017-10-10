package goku_bot

import (
	"poloniex-go-api"
	"time"
	. "goku-bot/global"
	"log"
	"sync"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/websocket"
	"strconv"
	"net/url"
	"encoding/json"
	"bytes"
)

//TODO - Abstract 'Poloniex' to a list of 'Exchanges'

type OrderBookSteward struct {
	Pair     string
	Exchange string
	Store    *Store
}

type CandlestickSteward struct {
	Poloniex *poloniex_go_api.Poloniex
	Store    *Store
}

type OrderBookEventData struct {
	Type   string
	Rate   float64
	Amount float64
}

type OrderBookEvent struct {
	Data OrderBookEventData
	Type string
	Seq  int
}

type FullOrderBook struct {
	CurrencyPair string
	Asks         []OrderBookEntry
	Bids         []OrderBookEntry
	Sequence     float64
}

type OrderBookEntry struct {
	CurrencyPair string
	Sequence int
	Price  float64
	Amount float64
}

func (self *OrderBookSteward) ConnectPoloniexOrderBook(pair string, group *sync.WaitGroup) {
	address := "api2.poloniex.com:443"

	u := url.URL{Scheme: "wss", Host: address}

	connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Printf("error connecting: %s", err)
		return
	}

	defer connection.Close()

	cmdString := []byte("{\"command\" : \"subscribe\", \"channel\" : \"" + pair + "\"}")
	err = connection.WriteMessage(websocket.TextMessage, cmdString)
	if err != nil {
		log.Printf("there was an error writing command: %s", err)
		return
	}

	i := 0
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

		log.Printf("M: %s", m)

		//first := m[0]
		seq := int(m.([]interface{})[1].(float64))

		if i == 0 { //if full order book
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
		}
		i++
	}
}

func (self *CandlestickSteward) SyncCandles(group *sync.WaitGroup) {
	syncOhlcErrorsCh := make(chan error)
	defer close(syncOhlcErrorsCh)

	go func() {
		for {
			select {
			case err := <-syncOhlcErrorsCh:
				if err != nil {
					log.Println("There was an error syncing Poloniex candles")
					log.Println(err)
				}
			}

		}
	}()

	numWorkers := len(POLONIEX_PAIRS) * len(POLONIEX_OHLC_INTERVALS)

	log.Printf("There are %d workers", numWorkers)
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for _, pair := range POLONIEX_PAIRS {
		for k, _ := range POLONIEX_OHLC_INTERVALS {
			log.Printf("pair: %s, interval: %d", pair, k)
			go self.BuildCandlesPoloniex(pair, k, syncOhlcErrorsCh, &wg)
		}
	}

	wg.Wait()

	log.Println("Finished Syncing OHLC")
	group.Done()
}

func (self *CandlestickSteward) BuildCandlesPoloniex(pair string, interval int, err chan<- error, group *sync.WaitGroup) {
	end := time.Now()
	start := end.AddDate(0, 0, -1)

	resp := self.Poloniex.ReturnChartData(pair, interval, start, end)

	if resp.Err != nil {
		log.Println("error getting the chart data")
		err <- resp.Err
	}

	self.Store.SyncCandles(resp.Response, "poloniex", pair, POLONIEX_OHLC_INTERVALS[interval])
	self.CalculateIndicators("poloniex", pair, interval)

	group.Done()
}

func (self *CandlestickSteward) CalculateIndicators(exchange, pair string, interval int) {
	collectionName := BuildCandleSliceCollectionName(exchange, pair, POLONIEX_OHLC_INTERVALS[interval])
	allSlices := self.Store.RetrieveSlicesByQueue(exchange, pair, interval, -1, -1)

	CalculateExponentialMovingAverage(10, allSlices)
	CalculateMacd(12, 26, 9, allSlices)
	CalculateAroon(25, allSlices)

	for _, slice := range allSlices {
		self.Store.GetCollection(collectionName).Update(bson.M{"queue": slice.Queue}, slice)
	}
}
