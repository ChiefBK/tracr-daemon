package store

import (
	. "goku-bot/global"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"poloniex-go-api"
	"strings"
	"time"
)

type Store interface {
	CloseStore()
	EmptyCollection(name string) error
	DropDatabase() error

	getCollection(name string) *mgo.Collection
	get(collectionName string, find *bson.M, sort *string, result interface{})

	GetTrades(exchange, pair string, sort *string) (trades []poloniex_go_api.Trade)
	InsertTrades(exchange, pair string, trades []poloniex_go_api.Trade)
	ReplaceTrades(exchange, pair string, trades []poloniex_go_api.Trade)

	InsertDeposits(exchange string, deposits []poloniex_go_api.Deposit)
	GetDeposits(exchange string, sort *string) (deposits []poloniex_go_api.Deposit)
	ReplaceDeposits(exchange string, deposits []poloniex_go_api.Deposit)

	InsertWithdrawals(exchange string, withdrawals []poloniex_go_api.Withdrawal)
	GetWithdrawals(exchange string, sort *string) (withdrawals []poloniex_go_api.Withdrawal)
	ReplaceWithdrawals(exchange string, withdrawals []poloniex_go_api.Withdrawal)

	SyncCandles(candles []*poloniex_go_api.Candle, exchange, pair, interval string)
	RetrieveSlicesByQueue(exchange, pair string, interval, start, end int) (slices []*CandleSlice)
}

func NewStore() (store Store, err error) {
	store, err = newMgoStore()

	return
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

//func (s *MgoStore) StoreBtcBalances(wg *sync.WaitGroup) {
//	balanceCh := make(chan *Poloniex_Go_Api.Balance)
//	go s.PoloniexApi.ReturnCompleteBalancesBtc(balanceCh)
//	balance := <-balanceCh
//
//	s.Database.C("BtcBalances").Insert(&BtcBalanceStore{balance, time.Now()})
//	wg.Done()
//}
//
//func (s *MgoStore) StoreLoanOffers(wg *sync.WaitGroup) {
//	loanOffersCh := make(chan []*Poloniex_Go_Api.Order)
//	go s.PoloniexApi.ReturnLoanOffers(loanOffersCh)
//	loanOffers := <-loanOffersCh
//
//	s.Database.C("LoanOffers").Insert(&LoanOffersStore{loanOffers, time.Now()})
//	wg.Done()
//}
//
//func (s *MgoStore) StoreActiveLoans(wg *sync.WaitGroup) {
//	activeLoansCh := make(chan *Poloniex_Go_Api.ReturnActiveLoansResponse)
//	go s.PoloniexApi.ReturnActiveLoans(activeLoansCh)
//	activeLoans := <-activeLoansCh
//
//	loans := activeLoans.Response["provided"]
//
//	s.Database.C("ActiveLoans").Insert(&ActiveLoansStore{loans, time.Now()})
//	wg.Done()
//}

func (s *MgoStore) SyncCandles(candles []*poloniex_go_api.Candle, exchange, pair, interval string) {
	if candles == nil {
		return
	}

	collectionName := buildCandleSliceCollectionName(exchange, pair, interval)
	log.Printf("Syncing Candles for collection %s", collectionName)

	startWindow := candles[0].Date

	// TODO - put in initialization
	index := mgo.Index{
		Key:    []string{"-step"},
		Unique: true,
	}
	s.getCollection(collectionName).EnsureIndex(index)

	var ohlc []*CandleSlice
	s.getCollection(collectionName).Find(bson.M{"date": bson.M{"$gte": startWindow}}).All(&ohlc)

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

func (s *MgoStore) storeCandles(candles []*poloniex_go_api.Candle, collectionName string, startingStep int) {
	log.Printf("Storing %d candles", len(candles))

	step := startingStep
	queue := len(candles) - 1

	s.getCollection(collectionName).UpdateAll(bson.M{}, bson.M{"$inc": bson.M{"queue": len(candles)}})

	for _, candle := range candles {
		ohlc := &CandleSlice{Candle: *candle, Date: time.Unix(int64(candle.Date), 0), Step: step, Queue: queue, Volume: candle.Volume}
		s.getCollection(collectionName).Insert(ohlc)
		step++
		queue--
	}
}

func (s *MgoStore) RetrieveCandlesByDate(exchange, pair, interval string, start, end time.Time) (candles []*OhlcSchema) {
	collectionName := buildCandleSliceCollectionName(exchange, pair, interval)
	err := s.getCollection(collectionName).Find(bson.M{"candle.date": bson.M{"$gte": start.Unix(), "$lte": end.Unix()}}).All(&candles)

	if err != nil {
		log.Println("Error retrieving candles by date")
		log.Println(err)
	}

	return
}

func (s *MgoStore) RetrieveSlicesByQueue(exchange, pair string, interval, start, end int) (slices []*CandleSlice) {
	collectionName := buildCandleSliceCollectionName(exchange, pair, POLONIEX_OHLC_INTERVALS[interval])
	log.Printf("Getting candles from collection %s within queue (%d, %d)", collectionName, start, end)

	var err error
	if start == -1 || end == -1 {
		err = s.getCollection(collectionName).Find(bson.M{}).All(&slices)
	} else {
		err = s.getCollection(collectionName).Find(bson.M{"queue": bson.M{"$lte": start, "$gte": end}}).All(&slices)
	}

	if err != nil {
		log.Println("Error retrieving candles by queue")
		log.Println(err)
	}

	return
}

func (s *MgoStore) RetrieveMacdByQueue(exchange, pair string, interval int, macdParams []int, start, end int) (results []MacdSchema) {
	//stringParams := strings.Split(fmt.Sprint(macdParams), " ")

	//collectionName := buildCollectionName("Indicator", exchange, pair, POLONIEX_OHLC_INTERVALS[interval], strings.Join(stringParams, "-"))
	return
}

//TODO - remove old candles
func (s *MgoStore) trimCandles(collectionName string) {

}

func buildCollectionName(params ...string) string {
	return strings.Join(params, "-")
}

func buildCandleSliceCollectionName(exchange, pair, interval string) string {
	return buildCollectionName("CandleSlice", exchange, pair, interval)
}

func buildMyTradeHistoryCollectionName(exchange, pair string) string {
	return buildCollectionName("MyTradeHistory", exchange, pair)
}

func buildDepositHistoryCollectionName(exchange string) string {
	return buildCollectionName("MyDepositHistory", exchange)
}

func buildWithdrawalHistoryCollectionName(exchange string) string {
	return buildCollectionName("MyWithdrawlHistory", exchange)
}
