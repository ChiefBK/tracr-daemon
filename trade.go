package tracr_daemon

import (
	"errors"
	"fmt"
	"time"
	"tracr-store"
)

type TradeSteward struct {
	Store tracr_store.Store
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
	store, err := tracr_store.NewStore()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("There was an error connecting to the store: %s", err))
	}

	return &TradeSteward{store}, nil
}

//func (self *TradeSteward) CalculateTradeNetProfit(exchange, pair string) (net float64) {
//	sort := "date"
//	trades := self.Store.GetTrades(exchange, pair, &sort)
//
//	net = 0
//
//	for _, trade := range trades {
//		if trade.Type == "buy" {
//			net -= trade.GetTotal()
//		} else { // "sell"
//			net += trade.GetTotal()
//		}
//	}
//
//	return
//}

//type PositionsById []Position
//
//func (self PositionsById) Len() int           { return len(self) }
//func (self PositionsById) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
//func (self PositionsById) Less(i, j int) bool { return self[i].Id < self[j].Id }
//
//func (self *TradeSteward) GetPositions(exchange, pair string) (positions []Position) {
//	s := "date"
//	trades := self.Store.GetTrades(exchange, pair, &s)
//
//	tradesAggregation := make(map[int64][]poloniex_go_api.Trade)
//	for _, trade := range trades { // aggregate all trades with same order number
//		tradesAggregation[trade.GetOrderNumber()] = append(tradesAggregation[trade.GetOrderNumber()], trade)
//	}
//
//	for id, trades := range tradesAggregation {
//		var totalAmount float64 = 0
//		for _, trade := range trades { // sum up amount of all trades in same order
//			totalAmount += trade.GetAmount()
//		}
//
//		var effectiveRate float64 = 0
//		for _, trade := range trades { // get effective rate of all trades in order
//			effectiveRate += (trade.GetAmount() / totalAmount) * trade.GetRate() // effective rate = sum(percentage of total amount * trade rate)
//		}
//
//		position := Position{id, totalAmount, effectiveRate, trades[0].Type, len(trades), trades[len(trades)-1].GetDate()}
//
//		positions = append(positions, position)
//	}
//
//	sort.Sort(PositionsById(positions))
//
//	return
//}
//
//type PositionResult struct {
//	Open          Position
//	Close         Position
//	Amount        float64
//	NetPercentage float64
//	NetUsd        float64
//}
//
//type PositionStack struct {
//	Positions []*Position
//}
//
//func (self *PositionStack) push(position Position) {
//	self.Positions = append(self.Positions, &position)
//}
//
//func (self *PositionStack) pop() {
//	self.Positions = self.Positions[:len(self.Positions)-1]
//}
//
//func (self *PositionStack) peek() *Position {
//	l := len(self.Positions)
//	return self.Positions[l-1]
//}
//
//func (self *PositionStack) len() int {
//	return len(self.Positions)
//}
//
//func net(open, close Position) (netUsd, netPercentage float64) {
//	var buyRate float64
//	var sellRate float64
//	if open.Direction == "buy" {
//		buyRate = open.Rate
//		sellRate = close.Rate
//	} else {
//		buyRate = close.Rate
//		sellRate = open.Rate
//	}
//	var amount float64 = math.Min(open.Amount, close.Amount)
//
//	netPercentage = 100 * ((sellRate / buyRate) - 1)
//	netUsd = (sellRate - buyRate) * amount
//
//	return
//}
//
//func (self *TradeSteward) CalculatePositionNetProfits(exchange, pair string) (closedPositions []PositionResult) {
//	positions := self.GetPositions(exchange, pair)
//
//	var sellPositionStack *PositionStack = new(PositionStack)
//	var buyPositionStack *PositionStack = new(PositionStack)
//
//	for _, position := range positions { // Add positions to respective stacks
//		if position.Direction == "buy" {
//			buyPositionStack.push(position)
//		} else {
//			sellPositionStack.push(position)
//		}
//	}
//
//	for buyPositionStack.len() > 0 && sellPositionStack.len() > 0 { // while both stacks still have one or more elements
//		buyPosition := buyPositionStack.peek()
//		sellPosition := sellPositionStack.peek()
//		//log.Printf("sellPosition : %s", sellPosition)
//		//log.Printf("buyPosition : %s", buyPosition)
//
//		var closePositionStack *PositionStack
//		var openPositionStack *PositionStack
//		if buyPosition.Date.Before(sellPosition.Date) { // figure out which stack is closing and opening
//			closePositionStack = sellPositionStack
//			openPositionStack = buyPositionStack
//		} else {
//			closePositionStack = buyPositionStack
//			openPositionStack = sellPositionStack
//		}
//		closePosition := closePositionStack.peek()
//		openPosition := openPositionStack.peek()
//
//		var amountToClose float64
//		var amountLeft float64
//		var popOpenPosition bool = true
//
//		amountToClose = openPosition.Amount
//		amountLeft = closePosition.Amount - amountToClose
//
//		if closePosition.Amount < openPosition.Amount {
//			amountToClose = closePosition.Amount
//			amountLeft = openPosition.Amount - amountToClose
//			popOpenPosition = false
//		}
//
//		if amountLeft == 0 { // if both open and close positions are of the same amount
//			netUsd, netPercentage := net(*openPosition, *closePosition)
//			closedPositions = append(closedPositions,
//				PositionResult{
//					Amount:        openPosition.Amount,
//					Open:          *openPosition,
//					Close:         *closePosition,
//					NetPercentage: netPercentage,
//					NetUsd:        netUsd,
//				})
//			closePositionStack.pop()
//			openPositionStack.pop()
//		} else if popOpenPosition {
//			netUsd, netPercentage := net(*openPosition, *closePosition)
//			closedPositions = append(closedPositions,
//				PositionResult{
//					Amount:        openPosition.Amount,
//					Open:          *openPosition,
//					Close:         *closePosition,
//					NetPercentage: netPercentage,
//					NetUsd:        netUsd,
//				})
//			closePosition.Amount = amountLeft
//			openPositionStack.pop()
//		} else {
//			netUsd, netPercentage := net(*openPosition, *closePosition)
//			closedPositions = append(closedPositions,
//				PositionResult{
//					Amount:        closePosition.Amount,
//					Open:          *openPosition,
//					Close:         *closePosition,
//					NetPercentage: netPercentage,
//					NetUsd:        netUsd,
//				})
//			openPosition.Amount = amountLeft
//			closePositionStack.pop()
//		}
//	}
//
//	return
//}
