package poloniex

import (
	"encoding/json"
	"fmt"
	"time"
	"strconv"
	"exchange-client"
	log "github.com/inconshreveable/log15"
	"goku-bot/exchanges"
	"goku-bot/pairs"
)

func NewPoloniexClient(apiKey, apiSecret string) *Poloniex {
	client := exchange_client.NewClient(apiKey, apiSecret, "poloniex", "https://poloniex.com/tradingApi", "https://poloniex.com/public", exchanges.POLONIEX_THROTTLE)
	return &Poloniex{client}
}

type Poloniex struct {
	client *exchange_client.Client
}

func (self *Poloniex) Ticker() (resp exchanges.TickerResponse) {
	urlQueryArgs := make(map[string]string)
	urlQueryArgs["command"] = "returnTicker"

	clientRes, err := self.client.Do("GET", "", urlQueryArgs, nil, nil)

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
			log.Warn("error finding standard pair name skipping", "module", "exchanges", "exchangePair", pair, "exchange", "poloniex", "error", err)
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

func (self *Poloniex) Balances() (resp exchanges.BalancesResponse) {
	bodyArgs := make(map[string]string)
	bodyArgs["command"] = "returnCompleteBalances"
	res, err := self.client.Do("POST", "", nil, bodyArgs, nil)

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

	for pair, balance := range balancesResp {
		stdPairName, err := pairs.StandardPair(pair, "poloniex")
		if err != nil {
			log.Warn("error finding standard pair name skipping", "module", "exchanges", "exchangePair", pair, "exchange", "poloniex", "error", err)
			continue
		}

		available := balance.GetAvailable()
		resp.Data[stdPairName] = available
	}

	return
}

func (self *Poloniex) ChartData(currencyPair string, period int, start, end time.Time) (resp exchanges.ChartDataResponse) {
	urlQueryArgs := make(map[string]string)
	urlQueryArgs["currencyPair"], _ = pairs.ExchangePair(currencyPair, "poloniex")
	urlQueryArgs["period"] = fmt.Sprintf("%d", period)
	urlQueryArgs["start"] = fmt.Sprintf("%d", start.Unix())
	urlQueryArgs["end"] = fmt.Sprintf("%d", end.Unix())
	urlQueryArgs["command"] = "returnChartData"

	res, err := self.client.Do("GET", "", urlQueryArgs, nil, nil)

	if err != nil {
		log.Error("there was an error getting chart data", "module", "exchanges", "exchange", "poloniex", "error", err)
		resp.Data = nil
		resp.Err = err
		return
	}

	var poloniexResp PoloniexChartData
	err = json.Unmarshal(res, &poloniexResp)

	if err != nil {
		log.Error("there was an error un-marshalling chart data", "module", "exchanges", "exchange", "poloniex", "error", err)
		resp.Err = err
		resp.Data = nil
		return
	}

	for _, candle := range poloniexResp {
		c := exchanges.Candle{candle.Open, candle.High, candle.Low, candle.Close}
		resp.Data = append(resp.Data, c)
	}

	return
}

func (self *Poloniex) MyTradeHistory() (resp exchanges.TradeHistoryResponse) {
	body := make(map[string]string)
	body["currencyPair"] = "all"
	body["limit"] = "10000"
	body["start"] = "0"
	body["end"] = strconv.FormatInt(time.Now().Unix(), 10)
	body["command"] = "returnTradeHistory"

	data, err := self.client.Do("POST", "", nil, body, nil)

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
		log.Debug("std pair", "module", "exchanges", "pair", stdPair, "exchangePair", pair)
		resp.Data[stdPair] = pairTrades
	}

	return
}

func (self *Poloniex) DepositsWithdrawals() (resp exchanges.DepositsWithdrawalsResponse) {
	body := make(map[string]string)
	body["start"] = "0"
	body["end"] = strconv.FormatInt(time.Now().Unix(), 10)
	body["command"] = "returnDepositsWithdrawals"

	data, err := self.client.Do("POST", "", nil, body, nil)

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
			Status: d.Status,
			Timestamp: d.GetDatetime(),
			Currency: d.Currency,
			Amount: d.GetAmount(),
		}
		deposits = append(deposits, deposit)
	}

	resp.Data = &exchanges.DepositsWithdrawals{
		Deposits: deposits,
		Withdrawals: withdrawals,
	}
	resp.Err = nil

	return
}