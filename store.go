package goku_bot

import (
	"errors"
	. "goku-bot/global"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"poloniex-go-api"
	"strings"
	"time"
)

type Store struct {
	Session  *mgo.Session
	database *mgo.Database
}

func (self *Store) closeStore() {
	self.Session.Close()
}

func (self *Store) GetCollection(name string) *mgo.Collection {
	return self.database.C(name)
}

func (self *Store) EmptyCollection(name string) {
	self.GetCollection(name).RemoveAll(bson.M{})
}

func (self *Store) InsertAll(name string, objects []interface{}) {
	for _, ob := range objects {
		self.GetCollection(name).Insert(ob)
	}
}

func (self *Store) DropDatabase() error {
	return self.database.DropDatabase()
}

type OhlcSchema struct {
	Candle *poloniex_go_api.Candle
	Step   int
	Queue  int
}

type MacdSchema struct {
	Macd  MacdValue
	Step  int
	Queue int
}

type MacdValue struct {
	Macd      *float64
	Signal    *float64
	Histogram *float64
}

type AroonValue struct {
	Up   *int
	Down *int
}

type CandleSlice struct {
	Date   time.Time
	Step   int
	Queue  int
	Sma    map[string]*float64
	Ema    map[string]*float64
	Candle poloniex_go_api.Candle
	Macd   map[string]*MacdValue
	Volume float64
	Rsi    map[string]*float64
	Aroon  map[string]*AroonValue
}

const (
	OHLC_MAX_CANDLES = 200
)

func NewStore() (store *Store, err error) {
	session, dbErr := mgo.Dial(DB_URI)
	if dbErr != nil {
		err = errors.New("could not connect to store")
		return
	}

	db := session.DB(DB_NAME)

	store = new(Store)

	store.Session = session
	store.database = db

	return
}

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

//type Exchange struct {
//	Poloniex *poloniex_go_api.Poloniex
//}

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
	if candles == nil {
		return
	}

	collectionName := BuildCandleSliceCollectionName(exchange, pair, interval)
	log.Printf("Syncing Candles for collection %s", collectionName)

	startWindow := candles[0].Date

	// TODO - put in initialization
	index := mgo.Index{
		Key:    []string{"-step"},
		Unique: true,
	}
	s.GetCollection(collectionName).EnsureIndex(index)

	var ohlc []*CandleSlice
	s.GetCollection(collectionName).Find(bson.M{"date": bson.M{"$gte": startWindow}}).All(&ohlc)

	if len(ohlc) == 0 {
		log.Println("No existing candles in db. Storing all candles")
		s.storeCandles(candles, collectionName, 0)
		return
	}

	lastOhlc := ohlc[len(ohlc)-1]

	var candlesToAdd []*poloniex_go_api.Candle

	for i := range candles {
		if candles[i].Date > lastOhlc.Candle.Date {
			candlesToAdd = append(candlesToAdd, candles[i])
		}
	}

	s.storeCandles(candlesToAdd, collectionName, lastOhlc.Queue+1)
}

func (s *Store) storeCandles(candles []*poloniex_go_api.Candle, collectionName string, startingStep int) {
	log.Printf("Storing %d candles", len(candles))

	step := startingStep
	queue := len(candles) - 1

	s.GetCollection(collectionName).UpdateAll(bson.M{}, bson.M{"$inc": bson.M{"queue": len(candles)}})

	for _, candle := range candles {
		ohlc := &CandleSlice{Candle: *candle, Date: time.Unix(int64(candle.Date), 0), Step: step, Queue: queue, Volume: candle.Volume}
		s.GetCollection(collectionName).Insert(ohlc)
		step++
		queue--
	}
}

func (s *Store) SyncAsks(exchange, currencyPair string, orders []OrderBookEntry) {
	collectionName := BuildOrderBookCollectionName(exchange, currencyPair, true)

	for _, order := range orders {
		if order.Amount == 0 {
			s.GetCollection(collectionName).Remove(bson.M{"price": order.Price})
			continue
		}
		s.GetCollection(collectionName).Upsert(bson.M{"price": order.Price}, order)
	}
}

func (s *Store) SyncBids(exchange, currencyPair string, orders []OrderBookEntry) {
	collectionName := BuildOrderBookCollectionName(exchange, currencyPair, false)

	for _, order := range orders {
		if order.Amount == 0 {
			s.GetCollection(collectionName).Remove(bson.M{"price": order.Price})
			continue
		}
		s.GetCollection(collectionName).Upsert(bson.M{"price": order.Price}, order)
	}
}

func (s *Store) RetrieveCandlesByDate(exchange, pair, interval string, start, end time.Time) (candles []*OhlcSchema) {
	collectionName := BuildCandleSliceCollectionName(exchange, pair, interval)
	err := s.GetCollection(collectionName).Find(bson.M{"candle.date": bson.M{"$gte": start.Unix(), "$lte": end.Unix()}}).All(&candles)

	if err != nil {
		log.Println("Error retrieving candles by date")
		log.Println(err)
	}

	return
}

func (s *Store) RetrieveSlicesByQueue(exchange, pair string, interval, start, end int) (slices []*CandleSlice) {
	collectionName := BuildCandleSliceCollectionName(exchange, pair, POLONIEX_OHLC_INTERVALS[interval])
	log.Printf("Getting candles from collection %s within queue (%d, %d)", collectionName, start, end)

	var err error
	if start == -1 || end == -1 {
		err = s.GetCollection(collectionName).Find(bson.M{}).All(&slices)
	} else {
		err = s.GetCollection(collectionName).Find(bson.M{"queue": bson.M{"$lte": start, "$gte": end}}).All(&slices)
	}

	if err != nil {
		log.Println("Error retrieving candles by queue")
		log.Println(err)
	}

	return
}

func (s *Store) RetrieveMacdByQueue(exchange, pair string, interval int, macdParams []int, start, end int) (results []MacdSchema) {
	//stringParams := strings.Split(fmt.Sprint(macdParams), " ")

	//collectionName := BuildCollectionName("Indicator", exchange, pair, POLONIEX_OHLC_INTERVALS[interval], strings.Join(stringParams, "-"))
	return
}

//TODO - remove old candles
func (s *Store) trimCandles(collectionName string) {

}

func BuildCollectionName(params ...string) string {
	return strings.Join(params, "-")
}

func BuildOrderBookCollectionName(exchange, pair string, isAsk bool) string {
	if isAsk {
		return BuildCollectionName("OrderBook", exchange, pair, "ask")
	} else {
		return BuildCollectionName("OrderBook", exchange, pair, "bid")
	}
}

func BuildCandleSliceCollectionName(exchange, pair, interval string) string {
	return BuildCollectionName("CandleSlice", exchange, pair, interval)
}

func BuildMyTradeHistoryCollectionName(exchange, pair string) string {
	return BuildCollectionName("History", exchange, pair)
}
