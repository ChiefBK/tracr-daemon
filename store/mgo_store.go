package store

import (
	. "goku-bot/global"
	"gopkg.in/mgo.v2"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"poloniex-go-api"
)

type MgoStore struct {
	session  *mgo.Session
	database *mgo.Database
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

func (self *MgoStore) GetTrades(exchange, pair string, sort *string) (trades []poloniex_go_api.Trade) {
	name := buildMyTradeHistoryCollectionName(exchange, pair)

	sortVal := ""
	if sort != nil {
		sortVal = *sort
	}

	self.getCollection(name).Find(bson.M{}).Sort(sortVal).All(&trades)

	return
}

func (self *MgoStore) InsertTrades(exchange, pair string, trades []poloniex_go_api.Trade) {
	collectionName := buildMyTradeHistoryCollectionName(exchange, pair)
	for _, trade := range trades {
		self.getCollection(collectionName).Insert(trade)
	}
}

func (self *MgoStore) ReplaceTrades(exchange, pair string, trades []poloniex_go_api.Trade) {
	collectionName := buildMyTradeHistoryCollectionName(exchange, pair)
	self.EmptyCollection(collectionName)
	self.InsertTrades(exchange, pair, trades)
}

func (self *MgoStore) InsertDeposits(exchange string, deposits []poloniex_go_api.Deposit) {
	collectionName := buildDepositHistoryCollectionName(exchange)
	for _, deposit := range deposits {
		self.getCollection(collectionName).Insert(deposit)
	}
}

func (self *MgoStore) GetDeposits(exchange string, sort *string) (deposits []poloniex_go_api.Deposit) {
	name := buildDepositHistoryCollectionName(exchange)
	self.get(name, nil, nil, &deposits)
	return
}

func (self *MgoStore) ReplaceDeposits(exchange string, deposits []poloniex_go_api.Deposit) {
	name := buildDepositHistoryCollectionName(exchange)
	self.EmptyCollection(name)
	self.InsertDeposits(exchange, deposits)
}

func (self *MgoStore) InsertWithdrawals(exchange string, withdrawals []poloniex_go_api.Withdrawal) {
	collectionName := buildWithdrawalHistoryCollectionName(exchange)
	for _, deposit := range withdrawals {
		self.getCollection(collectionName).Insert(deposit)
	}
}

func (self *MgoStore) GetWithdrawals(exchange string, sort *string) (withdrawals []poloniex_go_api.Withdrawal) {
	name := buildDepositHistoryCollectionName(exchange)
	self.get(name, nil, nil, &withdrawals)
	return
}

func (self *MgoStore) ReplaceWithdrawals(exchange string, withdrawals []poloniex_go_api.Withdrawal) {
	name := buildWithdrawalHistoryCollectionName(exchange)
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