package poloniex

import (
	"encoding/json"
	"fmt"
	"time"
	"strconv"
	log "github.com/inconshreveable/log15"
	"tracr-daemon/exchanges"
	"tracr-daemon/pairs"
	"tracr-client"
)

func NewPoloniexClient(apiKey, apiSecret string) *PoloniexClient {
	client := tracr_client.NewClient(apiKey, apiSecret, "poloniex", "https://poloniex.com/tradingApi", "https://poloniex.com/public", exchanges.POLONIEX_THROTTLE)
	return &PoloniexClient{client}
}

type PoloniexClient struct {
	apiClient *tracr_client.Client
}

func (self *PoloniexClient) OrderBook(stdPair string) (resp exchanges.OrderBookResponse) {
	panic("implement me")

	exchangePair, err := pairs.ExchangePair(stdPair, exchanges.POLONIEX)

	if err != nil {
		resp.Err = err
		return
	}

	urlQueryArgs := make(map[string]string)
	urlQueryArgs["command"] = "returnOrderBook"
	urlQueryArgs["pair"] = exchangePair
	urlQueryArgs["depth"] = "100"

	clientResponse, err := self.apiClient.Do("GET", "", urlQueryArgs, nil, nil)

	var exchangeResponse PoloniexOrderBook
	err = json.Unmarshal(clientResponse, &exchangeResponse)

	if err != nil {
		log.Error("there was an error un-marshalling the order book", "module", "exchanges", "exchange", "poloniex", "error", err)
		resp.Err = err
		return
	}

	orderBook := exchanges.OrderBook{Exchange: exchanges.POLONIEX, Pair: stdPair}
	asks := make(map[string]float64)
	bids := make(map[string]float64)

	for _, ask := range exchangeResponse.Asks {
		price := strconv.FormatFloat(ask[0], 'f', 8, 64)
		volume := ask[1]
		asks[price] = volume
	}

	for _, bid := range exchangeResponse.Bids {
		price := strconv.FormatFloat(bid[0], 'f', 8, 64)
		volume := bid[1]
		bids[price] = volume
	}

	orderBook.Bids = bids
	orderBook.Asks = asks

	resp.Data = orderBook
	resp.Err = nil

	return
}

func (self *PoloniexClient) Ticker() (resp exchanges.TickerResponse) {
	urlQueryArgs := make(map[string]string)
	urlQueryArgs["command"] = "returnTicker"

	clientRes, err := self.apiClient.Do("GET", "", urlQueryArgs, nil, nil)

	if err != nil {
		log.Error("there was an error getting the ticker", "module", "exchanges", "exchange", "poloniex", "error", err)
		resp.Err = err
		return
	}

	exchangeResp := make(map[string]PoloniexApiTicker) // map between each pair and ticker
	err = json.Unmarshal(clientRes, &exchangeResp)

	if err != nil {
		log.Error("there was an error un-marshalling the ticker", "module", "exchanges", "exchange", "poloniex", "error", err)
		resp.Err = err
		return
	}

	resp.Data = make(map[string]exchanges.Ticker)

	for pair, ticker := range exchangeResp {
		stdPairName, err := pairs.StandardPair(pair, "poloniex")
		if err != nil {
			log.Warn("error finding standard pair name skipping", "module", "exchanges", "exchangePair", pair, "exchange", "poloniex", "error", err, "retrieving", "ticker")
			continue
		}

		last := ticker.getLast()
		highest := ticker.getHighestBid()
		lowest := ticker.getLowestAsk()

		resp.Data[stdPairName] = exchanges.Ticker{
			LastTrade:          &last,
			HighestBid:         &highest,
			LowestAsk:          &lowest,
			TwentyFourHourLow:  nil,
			TwentyFourHourHigh: nil,
		}
	}

	return
}

func (self *PoloniexClient) Balances() (resp exchanges.BalancesResponse) {
	bodyArgs := make(map[string]string)
	bodyArgs["command"] = "returnCompleteBalances"
	res, err := self.apiClient.Do("POST", "", nil, bodyArgs, nil)

	if err != nil {
		log.Error("there was an error getting balances", "module", "exchanges", "exchange", "poloniex", "error", err)
		resp.Err = err
		resp.Data = nil
		return
	}

	balancesResp := make(PoloniexCompleteBalancesResponse)
	err = json.Unmarshal(res, &balancesResp)

	if err != nil {
		log.Error("there was an error un-marshalling balances", "module", "exchanges", "exchange", "poloniex", "error", err)
		resp.Err = err
		resp.Data = nil
		return
	}

	resp.Data = make(exchanges.Balances)

	for currency, balance := range balancesResp {
		available := balance.GetAvailable()
		resp.Data[currency] = available
	}

	return
}

func (self *PoloniexClient) ChartData(stdPair string, period time.Duration, start, end time.Time) (resp exchanges.ChartDataResponse) {
	urlQueryArgs := make(map[string]string)
	urlQueryArgs["currencyPair"], _ = pairs.ExchangePair(stdPair, "poloniex")
	urlQueryArgs["period"] = fmt.Sprintf("%d", int(period.Seconds()))
	urlQueryArgs["start"] = fmt.Sprintf("%d", start.Unix())
	urlQueryArgs["end"] = fmt.Sprintf("%d", end.Unix())
	urlQueryArgs["command"] = "returnChartData"

	res, err := self.apiClient.Do("GET", "", urlQueryArgs, nil, nil)

	if err != nil {
		log.Error("there was an error getting chart data", "module", "exchanges", "exchange", "poloniex", "error", err)
		resp.Data = nil
		resp.Err = err
		return
	}

	var poloniexResp []PoloniexCandle
	err = json.Unmarshal(res, &poloniexResp)

	if err != nil {
		log.Error("there was an error un-marshalling chart data", "module", "exchanges", "exchange", "poloniex", "error", err)
		resp.Err = err
		resp.Data = nil
		return
	}

	for _, candle := range poloniexResp {
		dateTime := time.Unix(int64(candle.Date), 0)
		c := exchanges.Candle{candle.Open, candle.High, candle.Low, candle.Close, dateTime}
		resp.Data = append(resp.Data, c)
	}

	return
}

func (self *PoloniexClient) MyTradeHistory() (resp exchanges.TradeHistoryResponse) {
	body := make(map[string]string)
	body["currencyPair"] = "all"
	body["limit"] = "10000"
	body["start"] = "0"
	body["end"] = strconv.FormatInt(time.Now().Unix(), 10)
	body["command"] = "returnTradeHistory"

	data, err := self.apiClient.Do("POST", "", nil, body, nil)

	if err != nil {
		log.Error("there was an error getting my trade history", "module", "exchanges", "exchange", "poloniex", "error", err)
		resp.Err = err
		resp.Data = nil
		return
	}

	var poloniexResp map[string][]PoloniexTrade // mapping of pair to list of trades
	err = json.Unmarshal(data, &poloniexResp)

	if err != nil {
		log.Error("there was an error un-marshalling my trade history", "module", "exchanges", "exchange", "poloniex", "error", err)
		resp.Err = err
		resp.Data = nil
		return
	}

	resp.Data = make(map[string][]exchanges.Trade)

	for pair, trades := range poloniexResp {
		var pairTrades []exchanges.Trade
		for _, trade := range trades {
			t := exchanges.Trade{
				ID:      trade.TradeID,
				Amount:  trade.GetAmount(),
				Rate:    trade.GetRate(),
				Date:    trade.GetDate(),
				Type:    trade.Type,
				Total:   trade.GetTotal(),
				Fee:     trade.GetFee(),
				OrderId: trade.OrderNumber,
			}
			pairTrades = append(pairTrades, t)
		}
		stdPair, _ := pairs.StandardPair(pair, exchanges.POLONIEX)
		resp.Data[stdPair] = pairTrades
	}

	return
}

func (self *PoloniexClient) DepositsWithdrawals() (resp exchanges.DepositsWithdrawalsResponse) {
	body := make(map[string]string)
	body["start"] = "0"
	body["end"] = strconv.FormatInt(time.Now().Unix(), 10)
	body["command"] = "returnDepositsWithdrawals"

	data, err := self.apiClient.Do("POST", "", nil, body, nil)

	if err != nil {
		resp.Data = nil
		resp.Err = err
		return
	}

	var poloniexResponse PoloniexDepositsWithdrawls

	json.Unmarshal(data, &poloniexResponse)

	var withdrawals []exchanges.Withdrawal
	var deposits []exchanges.Deposit

	for _, w := range poloniexResponse.Withdrawals {
		withdrawal := exchanges.Withdrawal{
			Amount:    w.GetAmount(),
			Timestamp: w.GetDatetime(),
			Currency:  w.Currency,
			Status:    w.Status,
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	for _, d := range poloniexResponse.Deposits {
		deposit := exchanges.Deposit{
			Status:    d.Status,
			Timestamp: d.GetDatetime(),
			Currency:  d.Currency,
			Amount:    d.GetAmount(),
		}
		deposits = append(deposits, deposit)
	}

	resp.Data = &exchanges.DepositsWithdrawals{
		Deposits:    deposits,
		Withdrawals: withdrawals,
	}
	resp.Err = nil

	return
}
