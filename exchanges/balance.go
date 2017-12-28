package exchanges

type Balances map[string]float64 // mapping between currency and available balance

type BalancesResponse struct {
	Data Balances
	Err  error
}
