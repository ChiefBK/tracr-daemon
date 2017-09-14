package goku_bot

import (
	"gopkg.in/mgo.v2"
	"time"
	"poloniex-go-api"
	"log"
	"strings"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

type Store struct {
	Session  *mgo.Session
	Database *mgo.Database
}

type BtcBalanceStore struct {
	Balance *poloniex_go_api.Balance
	Created time.Time
}

type LoanOffersStore struct {
	LoanOffers []*poloniex_go_api.Order
	Created    time.Time
}

type ActiveLoansStore struct {
	ActiveLoans []*poloniex_go_api.Loan
	Created     time.Time
}

type OhlcPoloniexEthBtcFiveMinSchema struct {
	Candle *poloniex_go_api.Candle
}

const (
	OHLC_MAX_CANDLES = 200
)

//type ChartOHLCSchema struct {
//	//ExchangeOHLCSchema `bson:"exchange,inline"`
//	//Exchanges map[string]map[string]map[string]
//}
//
//type ExchangeOHLCSchema struct {
//	//PairOHLCSchema `bson:"pair,inline"`
//}
//
//type PairOHLCSchema struct {
//	//IntervalOHLCSchema `bson:"interval,inline"`
//}
//
//type IntervalOHLCSchema struct {
//	//Candles []*Poloniex_Go_Api.Candle
//}

// TODO - use Exchange abstraction to get info from multiple exchanges
type Exchange struct {
	Poloniex *poloniex_go_api.Poloniex
}

//func (s *Store) StoreBtcBalances(wg *sync.WaitGroup) {
//	balanceCh := make(chan *Poloniex_Go_Api.Balance)
//	go s.PoloniexApi.ReturnCompleteBalancesBtc(balanceCh)
//	balance := <-balanceCh
//
//	s.Database.C("BtcBalances").Insert(&BtcBalanceStore{balance, time.Now()})
//	wg.Done()
//}
//
//func (s *Store) StoreLoanOffers(wg *sync.WaitGroup) {
//	loanOffersCh := make(chan []*Poloniex_Go_Api.Order)
//	go s.PoloniexApi.ReturnLoanOffers(loanOffersCh)
//	loanOffers := <-loanOffersCh
//
//	s.Database.C("LoanOffers").Insert(&LoanOffersStore{loanOffers, time.Now()})
//	wg.Done()
//}
//
//func (s *Store) StoreActiveLoans(wg *sync.WaitGroup) {
//	activeLoansCh := make(chan *Poloniex_Go_Api.ReturnActiveLoansResponse)
//	go s.PoloniexApi.ReturnActiveLoans(activeLoansCh)
//	activeLoans := <-activeLoansCh
//
//	loans := activeLoans.Response["provided"]
//
//	s.Database.C("ActiveLoans").Insert(&ActiveLoansStore{loans, time.Now()})
//	wg.Done()
//}

func (s *Store) SyncCandles(candles []*poloniex_go_api.Candle, exchange, pair, interval string) {
	log.Println("Syncing Candles")

	if len(candles) == 0 {
		return
	}

	collectionName := BuildCollectionName(exchange, pair, interval)

	startWindow := candles[0].Date

	var dbCandles []*poloniex_go_api.Candle
	s.Database.C(collectionName).Find(bson.M{"date": bson.M{"$gte": startWindow}}).All(&dbCandles)

	if len(dbCandles) == 0 {
		fmt.Println("No existing candles in db. Storing all candles")
		s.storeCandles(candles, collectionName)
		return
	}

	lastDbCandle := dbCandles[len(dbCandles)-1]

	var candlesToAdd []*poloniex_go_api.Candle

	for i := range candles {
		if candles[i].Date > lastDbCandle.Date {
			candlesToAdd = append(candlesToAdd, candles[i])
		}
	}

	s.storeCandles(candlesToAdd, collectionName)
}

func (s *Store) storeCandles(candles []*poloniex_go_api.Candle, collectionName string) {
	log.Printf("Storing %d candles", len(candles))

	for i := range candles {
		s.Database.C(collectionName).Insert(candles[i])
	}
}

//TODO - remove old candles
func (s *Store) trimCandles(collectionName string) {

}

func BuildCollectionName(exchange, pair, interval string) string{
	return strings.Join([]string{"OHLC", exchange, pair, interval}, "-")
}
