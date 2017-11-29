package goku_bot

//var TradeHistorySynced = make(chan bool)             // 'true' is sent to channel if trades have been updated
//var DepositWithdrawalHistorySynced = make(chan bool) // 'true' is sent to channel if trades have been updated
//
//type AccountSteward struct {
//	Poloniex *poloniex_go_api.Poloniex
//	Store    store.Store
//}
//
//func NewAccountSteward() (*AccountSteward, error) {
//	store, err := store.NewStore()
//
//	if err != nil {
//		return nil, errors.New("there was an error creating the store")
//	}
//
//	if PoloniexClient == nil {
//		return nil, errors.New("the poloniex client hasn't been initialized")
//	}
//
//	return &AccountSteward{PoloniexClient, store}, nil
//
//}
//
//func (self *AccountSteward) SyncBalances() {
//	response := self.Poloniex.ReturnCompleteBalances()
//
//	if response.Err != nil {
//		log.Println("there was an error getting the Poloniex balances - stopping balance sync")
//		return
//	}
//
//	data := response.Data
//
//	var balances CompleteBalances
//
//	balances.BTC = *data["BTC"]
//	balances.BCH = *data["BCH"]
//	balances.BCN = *data["BCN"]
//	balances.BTS = *data["BTS"]
//	balances.BURST = *data["BURST"]
//	balances.DASH = *data["DASH"]
//	balances.EMC2 = *data["EMC2"]
//	balances.ETH = *data["ETH"]
//	balances.EXP = *data["EXP"]
//	balances.FCT = *data["FCT"]
//	balances.LTC = *data["LTC"]
//	balances.PINK = *data["PINK"]
//	balances.VRC = *data["VRC"]
//	balances.XMR = *data["XMR"]
//	balances.ZEC = *data["ZEC"]
//	now := time.Now()
//	balances.Updated = now
//
//	select {
//	case PoloniexBalances <- balances:
//	case <-time.After(1 * time.Second):
//	}
//
//	log.Printf("Balances updated")
//}
//
//func (self *AccountSteward) SyncTradeHistory() {
//
//}
//
//func (self *AccountSteward) SyncDepositWithdrawlHistory() {
//	response := self.Poloniex.ReturnDepositsWithdrawals()
//
//	if response.Err != nil {
//		log.Printf("there was an error getting my Poloniex deposit-withdrawl history : %s", response.Err)
//		return
//	}
//
//	data := response.Data
//	deposits := data.Deposits
//	withdrawals := data.Withdrawals
//
//	self.Store.ReplaceDeposits("poloniex", deposits)
//	self.Store.ReplaceWithdrawals("poloniex", withdrawals)
//
//	select {
//	case DepositWithdrawalHistorySynced <- true:
//	case <-time.After(3 * time.Second):
//	}
//}
