package store

import (
	. "tracr-daemon/global"
	"gopkg.in/mgo.v2"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"tracr-daemon/exchanges"
	"time"
	"tracr-daemon/keys"
)

type MgoStore struct {
	session  *mgo.Session
	database *mgo.Database
}

func (self *MgoStore) GetChartData(exchange, pair string, interval time.Duration, sort *string) (candles []exchanges.Candle) {
	name := keys.BuildChartDataKey(exchange, pair, interval)

	sortVal := ""
	if sort != nil {
		sortVal = *sort
	}

	self.getCollection(name).Find(bson.M{}).Sort(sortVal).All(&candles)

	return
}

func (self *MgoStore) InsertChartData(exchange, pair string, interal time.Duration, candles []exchanges.Candle) {
	name := keys.BuildChartDataKey(exchange, pair, interal)

	for _, candle := range candles {
		self.getCollection(name).Insert(candle)
	}
}

func (self *MgoStore) ReplaceChartData(exchange, pair string, interval time.Duration, candles []exchanges.Candle) {
	name := keys.BuildChartDataKey(exchange, pair, interval)
	self.EmptyCollection(name)
	self.InsertChartData(exchange, pair, interval, candles)
}

func (self *MgoStore) RetrieveSlicesByQueue(exchange, pair string, interval, start, end int) (slices []*CandleSlice) {
	panic("implement me")
}

func (self *MgoStore) SyncCandles(candles []exchanges.Candle, exchange, pair, interval string) {
	panic("implement me")
}

func newMgoStore() (store *MgoStore, err error) {
	session, dbErr := mgo.Dial(DB_URI)
	if dbErr != nil {
		err = errors.New("could not connect to store")
		return
	}

	db := session.DB(DB_NAME)

	store = new(MgoStore)

	store.session = session
	store.database = db

	return
}

func (self *MgoStore) CloseStore() {
	self.session.Close()
}

func (self *MgoStore) getCollection(name string) *mgo.Collection {
	return self.database.C(name)
}

func (self *MgoStore) EmptyCollection(name string) error {
	_, err := self.getCollection(name).RemoveAll(bson.M{})
	return err
}

func (self *MgoStore) GetTrades(exchange, pair string, sort *string) (trades []exchanges.Trade) {
	name := keys.BuildMyTradeHistoryKey(exchange, pair)

	sortVal := ""
	if sort != nil {
		sortVal = *sort
	}

	self.getCollection(name).Find(bson.M{}).Sort(sortVal).All(&trades)

	return
}

func (self *MgoStore) InsertTrades(exchange, pair string, trades []exchanges.Trade) {
	collectionName := keys.BuildMyTradeHistoryKey(exchange, pair)
	for _, trade := range trades {
		self.getCollection(collectionName).Insert(trade)
	}
}

func (self *MgoStore) ReplaceTrades(exchange, pair string, trades []exchanges.Trade) {
	collectionName := keys.BuildMyTradeHistoryKey(exchange, pair)
	self.EmptyCollection(collectionName)
	self.InsertTrades(exchange, pair, trades)
}

func (self *MgoStore) InsertDeposits(exchange string, deposits []exchanges.Deposit) {
	collectionName := keys.BuildDepositHistoryKey(exchange)
	for _, deposit := range deposits {
		self.getCollection(collectionName).Insert(deposit)
	}
}

func (self *MgoStore) GetDeposits(exchange string, sort *string) (deposits []exchanges.Deposit) {
	name := keys.BuildDepositHistoryKey(exchange)
	self.get(name, nil, nil, &deposits)
	return
}

func (self *MgoStore) ReplaceDeposits(exchange string, deposits []exchanges.Deposit) {
	name := keys.BuildDepositHistoryKey(exchange)
	self.EmptyCollection(name)
	self.InsertDeposits(exchange, deposits)
}

func (self *MgoStore) InsertWithdrawals(exchange string, withdrawals []exchanges.Withdrawal) {
	collectionName := keys.BuildWithdrawalHistoryKey(exchange)
	for _, deposit := range withdrawals {
		self.getCollection(collectionName).Insert(deposit)
	}
}

func (self *MgoStore) GetWithdrawals(exchange string, sort *string) (withdrawals []exchanges.Withdrawal) {
	name := keys.BuildDepositHistoryKey(exchange)
	self.get(name, nil, nil, &withdrawals)
	return
}

func (self *MgoStore) ReplaceWithdrawals(exchange string, withdrawals []exchanges.Withdrawal) {
	name := keys.BuildWithdrawalHistoryKey(exchange)
	self.EmptyCollection(name)
	self.InsertWithdrawals(exchange, withdrawals)
}

func (self *MgoStore) get(collectionName string, find *bson.M, sort *string, result interface{}) {
	sortVal := ""
	if sort != nil {
		sortVal = *sort
	}

	findVal := bson.M{}
	if find != nil {
		findVal = *find
	}

	self.getCollection(collectionName).Find(findVal).Sort(sortVal).All(&result)
}

func (self *MgoStore) DropDatabase() error {
	return self.database.DropDatabase()
}
