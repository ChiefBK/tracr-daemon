package receivers

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/websocket"
	"encoding/json"
	"bytes"
	"strconv"
	"time"
	"goku-bot"
)

type TickerReceiver struct {
	exchange string
	pair     string
}

func NewTickerReceiver(exchange, pair string) *TickerReceiver {
	return &TickerReceiver{exchange, pair}
}

func (self *TickerReceiver) Start() {
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
		newTicker := new(goku_bot.Ticker)
		newTicker.CurrencyPair = "USDT_BTC"
		newTicker.Last = last
		newTicker.LowestAsk = lowestAsk
		newTicker.HighestBid = highestBid
		newTicker.PercentChange = percentChange
		newTicker.BaseVolume = baseVolume
		newTicker.QuoteVolume = quoteVolume
		newTicker.IsFrozen = isFrozen
		newTicker.TwentyFourHourHigh = twentyFourHourHigh
		newTicker.TwentyFourHourLow = twentyFourHourLow
		newTicker.Updated = now

		log.WithFields(log.Fields{"key": self.Key(), "module": "receivers"}).Debug("recieved ticker update")

		sendToProcessor(self.Key(), newTicker)
	}
}

func (self *TickerReceiver) Key() string {
	return fmt.Sprintf("%s-Ticker-%s", self.exchange, self.pair)
}
