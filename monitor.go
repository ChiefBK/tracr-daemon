package goku_bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	. "goku-bot/global"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/url"
	"poloniex-go-api"
	"strconv"
	"sync"
	"time"
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

type TickerSteward struct {
	Store *Store
}

type AccountSteward struct {
	Poloniex *poloniex_go_api.Poloniex
	Store    *Store
}

func NewAccountSteward() (*AccountSteward, error) {
	store, err := NewStore()

	if err != nil {
		return nil, errors.New("there was an error creating the store")
	}

	if PoloniexClient == nil {
		return nil, errors.New("the poloniex client hasn't been initialized")
	}

	return &AccountSteward{PoloniexClient, store}, nil

}

func NewCandleStickSteward() (*CandlestickSteward, error) {
	store, err := NewStore()

	if err != nil {
		return nil, errors.New("there was an error creating the store")
	}

	if PoloniexClient == nil {
		return nil, errors.New("the poloniex client hasn't been initialized")
	}

	return &CandlestickSteward{PoloniexClient, store}, nil
}

func NewTickerSteward() (*TickerSteward, error) {
	store, err := NewStore()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error creating Ticker Steward: %s", err))
	}

	return &TickerSteward{store}, nil
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
	Sequence     int
	Price        float64
	Amount       float64
}

func websocketConnect(address string, retries int) (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: address}

	var connection *websocket.Conn
	var err error
	retriesLeft := retries
	for retriesLeft > 0 {
		connection, _, err = websocket.DefaultDialer.Dial(u.String(), nil)

		if err != nil {
			log.Printf("error connecting: %s", err)
			log.Println("Retrying connection after 5 seconds")
			retriesLeft--

			timer := time.NewTimer(time.Second * 5)
			<-timer.C
		} else {
			break
		}
	}

	if retriesLeft == 0 {
		log.Println()
		return connection, errors.New(fmt.Sprintf("Could not connect after %d attemps", retries))
	}

	return connection, nil
}

func (self *AccountSteward) SyncBalances() {
	response := self.Poloniex.ReturnCompleteBalances()

	if response.Err != nil {
		log.Println("there was an error getting the Poloniex balances - stopping balance sync")
		return
	}

	balances := response.Data

	PoloniexBalances.BTC = *balances["BTC"]
	PoloniexBalances.BCH = *balances["BCH"]
	PoloniexBalances.BCN = *balances["BCN"]
	PoloniexBalances.BTS = *balances["BTS"]
	PoloniexBalances.BURST = *balances["BURST"]
	PoloniexBalances.DASH = *balances["DASH"]
	PoloniexBalances.EMC2 = *balances["EMC2"]
	PoloniexBalances.ETH = *balances["ETH"]
	PoloniexBalances.EXP = *balances["EXP"]
	PoloniexBalances.FCT = *balances["FCT"]
	PoloniexBalances.LTC = *balances["LTC"]
	PoloniexBalances.PINK = *balances["PINK"]
	PoloniexBalances.VRC = *balances["VRC"]
	PoloniexBalances.XMR = *balances["XMR"]
	PoloniexBalances.ZEC = *balances["ZEC"]
	now := time.Now()
	PoloniexBalances.Updated = now

	log.Printf("Balances updated at %s", now)
}

func (self *TickerSteward) ConnectPoloniexTicker(group *sync.WaitGroup) {
	address := "api2.poloniex.com:443"

	connection, err := websocketConnect(address, 5)

	if err != nil {
		log.Printf("Could not connect to Poloniex Ticker: %s", err)
	}

	defer connection.Close()

	cmdString := []byte("{\"command\" : \"subscribe\", \"channel\" : 1002}")
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

		if len(m.([]interface{})) < 3 {
			continue
		}

		ticker := m.([]interface{})[2].([]interface{})
		pairCode := ticker[0].(float64)

		if pairCode != 121 { // if NOT usdt_btc ticker
			continue
		}

		// Extract information from decoded message
		last, _ := strconv.ParseFloat(ticker[1].(string), 64)
		lowestAsk, _ := strconv.ParseFloat(ticker[2].(string), 64)
		highestBid, _ := strconv.ParseFloat(ticker[3].(string), 64)
		percentChange, _ := strconv.ParseFloat(ticker[4].(string), 64)
		baseVolume, _ := strconv.ParseFloat(ticker[5].(string), 64)
		quoteVolume, _ := strconv.ParseFloat(ticker[6].(string), 64)
		isFrozen := ticker[7].(float64) == 1
		twentyFourHourHigh, _ := strconv.ParseFloat(ticker[8].(string), 64)
		twentyFourHourLow, _ := strconv.ParseFloat(ticker[9].(string), 64)

		// Update Ticker
		now := time.Now()
		TickerUsdtBtc.CurrencyPair = "USDT_BTC"
		TickerUsdtBtc.Last = last
		TickerUsdtBtc.lowestAsk = lowestAsk
		TickerUsdtBtc.HighestBid = highestBid
		TickerUsdtBtc.PercentChange = percentChange
		TickerUsdtBtc.BaseVolume = baseVolume
		TickerUsdtBtc.QuoteVolume = quoteVolume
		TickerUsdtBtc.IsFrozen = isFrozen
		TickerUsdtBtc.TwentyFourHourHigh = twentyFourHourHigh
		TickerUsdtBtc.TwentyFourHourLow = twentyFourHourLow
		TickerUsdtBtc.Updated = now

		log.Printf("Ticker USDT_BTC updated at %s", now)

		if isFirst {
			group.Done()
		}

		isFirst = false
	}
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
			log.Printf("Recieved Orderbook update at %s", time.Now())
		}
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
		for k := range POLONIEX_OHLC_INTERVALS {
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
