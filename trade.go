package goku_bot

import (
	"errors"
	"fmt"
	"poloniex-go-api"
	"sort"
	"time"
)

var TradesSynced = make(chan bool) // 'true' is sent to channel if trades have been updated

type TradeSteward struct {
	Store *MgoStore
}

type Position struct {
	Id          int64
	Amount      float64
	Rate        float64
	Direction   string
	NumOfTrades int
	Date        time.Time
}

func NewTradeStewared() (*TradeSteward, error) {
	store, err := NewMgoStore()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("There was an error connecting to the store: %s", err))
	}

	return &TradeSteward{store}, nil
}

func (self *TradeSteward) CalculateTradeNetProfit(exchange, pair string) (net float64) {
	sort := "date"
	trades := self.Store.GetTrades(exchange, pair, &sort)

	net = 0

	for _, trade := range trades {
		if trade.Type == "buy" {
			net -= trade.GetTotal()
		} else { // "sell"
			net += trade.GetTotal()
		}
	}

	return
}

type PositionsById []Position

func (self PositionsById) Len() int           { return len(self) }
func (self PositionsById) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self PositionsById) Less(i, j int) bool { return self[i].Id < self[j].Id }

func (self *TradeSteward) GetPositions(exchange, pair string) (positions []Position) {
	s := "date"
	trades := self.Store.GetTrades(exchange, pair, &s)

	tradesAggregation := make(map[int64][]poloniex_go_api.Trade)
	for _, trade := range trades { // aggregate all trades with same order number
		tradesAggregation[trade.GetOrderNumber()] = append(tradesAggregation[trade.GetOrderNumber()], trade)
	}

	for id, trades := range tradesAggregation {
		var totalAmount float64 = 0
		for _, trade := range trades { // sum up amount of all trades in same order
			totalAmount += trade.GetAmount()
		}

		var effectiveRate float64 = 0
		for _, trade := range trades { // get effective rate of all trades in order
			effectiveRate += (trade.GetAmount() / totalAmount) * trade.GetRate() // effective rate = sum(percentage of total amount * trade rate)
		}

		position := Position{id, totalAmount, effectiveRate, trades[0].Type, len(trades), trades[len(trades)-1].GetDate()}

		positions = append(positions, position)
	}

	sort.Sort(PositionsById(positions))

	return
}

//func (self *TradeSteward)
