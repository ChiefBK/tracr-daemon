package goku_bot

//type CompleteBalances struct {
//	BTC     poloniex_go_api.Balance
//	LTC     poloniex_go_api.Balance
//	EXP     poloniex_go_api.Balance
//	EMC2    poloniex_go_api.Balance
//	PINK    poloniex_go_api.Balance
//	BCN     poloniex_go_api.Balance
//	FCT     poloniex_go_api.Balance
//	BTS     poloniex_go_api.Balance
//	VRC     poloniex_go_api.Balance
//	BURST   poloniex_go_api.Balance
//	ETH     poloniex_go_api.Balance
//	BCH     poloniex_go_api.Balance
//	ZEC     poloniex_go_api.Balance
//	DASH    poloniex_go_api.Balance
//	XMR     poloniex_go_api.Balance
//	Updated time.Time
//}
//
//func (self CompleteBalances) SumBalances() (sum float64) {
//	sum = 0
//	v := reflect.ValueOf(self)
//
//	for i := 0; i < v.NumField(); i++ {
//		prop := v.Field(i).Interface()
//
//		bal, ok := prop.(poloniex_go_api.Balance)
//
//		if ok {
//			sum = sum + bal.GetBtcValue()
//		}
//	}
//
//	return sum
//}
//
//var PoloniexBalances = make(chan CompleteBalances)
