package goku_bot

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
	"github.com/gorilla/websocket"
	"encoding/json"
	"bytes"
)

type Ticker struct {
	CurrencyPair       string
	Last               float64
	lowestAsk          float64
	HighestBid         float64
	PercentChange      float64
	BaseVolume         float64
	QuoteVolume        float64
	IsFrozen           bool
	TwentyFourHourHigh float64
	TwentyFourHourLow  float64
	Updated            time.Time
}

var TickerUsdtBtc = make(chan Ticker)

type TickerSteward struct {
	Store *Store
}

func NewTickerSteward() (*TickerSteward, error) {
	store, err := NewStore()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error creating Ticker Steward: %s", err))
	}

	return &TickerSteward{store}, nil
}

func (self *TickerSteward) ConnectPoloniexTicker() {
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
		newTicker := new(Ticker)
		newTicker.CurrencyPair = "USDT_BTC"
		newTicker.Last = last
		newTicker.lowestAsk = lowestAsk
		newTicker.HighestBid = highestBid
		newTicker.PercentChange = percentChange
		newTicker.BaseVolume = baseVolume
		newTicker.QuoteVolume = quoteVolume
		newTicker.IsFrozen = isFrozen
		newTicker.TwentyFourHourHigh = twentyFourHourHigh
		newTicker.TwentyFourHourLow = twentyFourHourLow
		newTicker.Updated = now

		select {
		case TickerUsdtBtc <- *newTicker:
		case <-time.After(1 * time.Second):
		}

		log.Printf("Ticker USDT_BTC updated")
	}
}
