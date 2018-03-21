package kraken

import (
	"tracr-client"
	"tracr-daemon/exchanges"
	"time"
	"tracr-daemon/pairs"
	"encoding/json"
	log "github.com/inconshreveable/log15"
	"strconv"
	"strings"
	"fmt"
	"errors"
)

const BASE_URL = "https://api.kraken.com"

type KrakenClient struct {
	apiClient tracr_client.BaseApiClient
}

func NewKrakenClient(apiKey, apiSecret string) *KrakenClient {
	client := tracr_client.NewApiClient(apiKey, apiSecret, exchanges.KRAKEN, BASE_URL, BASE_URL, exchanges.KRAKEN_THROTTLE)
	return &KrakenClient{client}
}

func handleKrakenErrors(response KrakenResponse) error {
	if len(response.Error) > 0 {
		errString := ""
		for i, err := range response.Error {
			errString += fmt.Sprintf("(Error %d) - %s ", i, err)
		}

		return errors.New(errString)
	}

	return nil
}

func (self *KrakenClient) Ticker() (resp exchanges.TickerResponse) {
	bodyArgs := make(map[string]string)
	var krakenPairs []string

	for krakenPair := range pairs.KrakenExchPairs {
		krakenPairs = append(krakenPairs, krakenPair)
	}

	bodyArgs["pairs"] = strings.Join(krakenPairs, ", ")

	exchangeResp, err := self.apiClient.Do("POST", "/0/public/Ticker", nil, nil, nil)

	if err != nil {
		log.Error("error getting making request", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "ticker", "error", err)
		resp.Err = err
		return
	}

	var krakenResponse KrakenResponse
	krakenResponse.Result = new(map[string]KrakenTickerResult)

	err = json.Unmarshal(exchangeResp, &krakenResponse)

	if err != nil {
		log.Error("error un-marshalling exchange response", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "ticker", "error", err)
		resp.Err = err
		return
	}

	err = handleKrakenErrors(krakenResponse)

	if err != nil {
		log.Error("kraken responded successfully with errors", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "ticker", "error", err)
		resp.Err = err
		return
	}

	data := make(map[string]exchanges.Ticker)

	for pair, result := range *krakenResponse.Result.(*map[string]KrakenTickerResult) {
		lastTrade, _ := strconv.ParseFloat(result.C[0], 64)
		highestBid, _ := strconv.ParseFloat(result.B[0], 64)
		lowestAsk, _ := strconv.ParseFloat(result.A[0], 64)
		twentyFourHourHigh, _ := strconv.ParseFloat(result.H[1], 64)
		twentyFourHourLow, _ := strconv.ParseFloat(result.L[1], 64)

		ticker := exchanges.Ticker{
			&lastTrade,
			&highestBid,
			&lowestAsk,
			&twentyFourHourHigh,
			&twentyFourHourLow,
		}

		stdPair, err := pairs.StandardPair(pair, exchanges.KRAKEN)

		if err != nil {
			log.Error("error retrieving standard pair", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "ticker", "exchangePair", pair, "error", err)
		}

		data[stdPair] = ticker
	}

	resp.Data = data
	resp.Err = nil
	return
}

func (self *KrakenClient) Balances() (resp exchanges.BalancesResponse) {
	exchangeResp, err := self.apiClient.Do("POST", "/0/private/Balance", nil, nil, nil)

	if err != nil {
		log.Error("error getting making request", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "balances", "error", err)
		resp.Err = err
		return
	}

	var krakenResponse KrakenResponse
	krakenResponse.Result = new(map[string]string)

	err = json.Unmarshal(exchangeResp, &krakenResponse)

	if err != nil {
		log.Error("error un-marshalling exchange response", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "balances", "error", err)
		resp.Err = err
		return
	}

	err = handleKrakenErrors(krakenResponse)

	if err != nil {
		log.Error("kraken responded successfully with errors", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "balances", "error", err)
		resp.Err = err
		return
	}

	results := *krakenResponse.Result.(*map[string]string)
	data := make(map[string]float64)

	for currency, amount := range results {
		amountFloat, _ := strconv.ParseFloat(amount, 64)
		data[currency] = amountFloat
	}

	resp.Data = data
	resp.Err = nil
	return
}

func (self *KrakenClient) ChartData(stdPair string, period time.Duration, start, end time.Time) (resp exchanges.ChartDataResponse) {
	krakenPair, err := pairs.ExchangePair(stdPair, exchanges.KRAKEN)

	if err != nil {
		log.Error("std pair argument not valid", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "chartData", "stdPair", stdPair, "error", err)
		resp.Err = err
		return
	}

	bodyArgs := make(map[string]string)
	bodyArgs["pair"] = krakenPair
	bodyArgs["interval"] = strconv.FormatInt(int64(period.Minutes()), 10)
	bodyArgs["since"] = strconv.FormatInt(start.Unix(), 10)

	exchangeResp, err := self.apiClient.Do("POST", "/0/public/OHLC", nil, bodyArgs, nil)

	if err != nil {
		log.Error("error getting making request", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "chartData", "stdPair", stdPair, "error", err)
		resp.Err = err
		return
	}

	var krakenResponse KrakenResponse
	krakenResponse.Result = new(map[string]interface{})

	err = json.Unmarshal(exchangeResp, &krakenResponse)

	if err != nil {
		log.Error("error un-marshalling exchange response", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "chartData", "stdPair", stdPair, "error", err)
		resp.Err = err
		return
	}

	err = handleKrakenErrors(krakenResponse)

	if err != nil {
		log.Error("kraken responded successfully with errors", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "chartData", "error", err)
		resp.Err = err
		return
	}

	resultPtr := krakenResponse.Result.(*map[string]interface{})
	result := *resultPtr

	ohlcData := result[krakenPair].([]interface{})

	var candles []exchanges.Candle

	for _, c := range ohlcData {
		data := c.([]interface{})
		timestamp := time.Unix(int64(data[0].(float64)), 0)
		open, _ := strconv.ParseFloat(data[1].(string), 64)
		high, _ := strconv.ParseFloat(data[2].(string), 64)
		low, _ := strconv.ParseFloat(data[3].(string), 64)
		close, _ := strconv.ParseFloat(data[4].(string), 64)

		candle := exchanges.Candle{
			Open:     open,
			High:     high,
			Low:      low,
			Close:    close,
			DateTime: timestamp,
		}

		candles = append(candles, candle)
	}

	resp.Data = candles
	resp.Err = nil
	return
}

func (self *KrakenClient) MyTradeHistory() (resp exchanges.TradeHistoryResponse) {
	bodyArgs := make(map[string]string)
	bodyArgs["type"] = "all"
	response, err := self.apiClient.Do("POST", "/0/private/TradesHistory", nil, bodyArgs, nil)

	if err != nil {
		log.Error("error getting making request", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "myTradeHistory", "error", err)
		resp.Err = err
		return
	}

	var krakenResponse KrakenResponse
	krakenResponse.Result = new(KrakenTradeResult)

	err = json.Unmarshal(response, &krakenResponse)

	if err != nil {
		log.Error("error un-marshalling exchange response", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "myTradeHistory", "error", err)
		resp.Err = err
		return
	}

	err = handleKrakenErrors(krakenResponse)

	if err != nil {
		log.Error("kraken responded successfully with errors", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "myTradeHistory", "error", err)
		resp.Err = err
		return
	}

	results := *krakenResponse.Result.(*KrakenTradeResult)

	data := make(map[string][]exchanges.Trade)

	for id, exchangeTrade := range results.Trades {
		date := time.Unix(int64(exchangeTrade.Time), 0)
		price, _ := strconv.ParseFloat(exchangeTrade.Price, 64)
		volume, _ := strconv.ParseFloat(exchangeTrade.Vol, 64)
		totalCost, _ := strconv.ParseFloat(exchangeTrade.Cost, 64)
		orderId := exchangeTrade.Ordertxid
		fee, _ := strconv.ParseFloat(exchangeTrade.Fee, 64)
		type2 := exchangeTrade.Type
		exchangePair := exchangeTrade.Pair

		stdPair, err := pairs.StandardPair(exchangePair, exchanges.KRAKEN)

		if err != nil {
			log.Error("error converting exchange pair", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "myTradeHistory", "exchangePair", exchangePair, "error", err)
			continue
		}

		trade := exchanges.Trade{
			Price:     price,
			Volume:    volume,
			Date:      date,
			Fee:       fee,
			ID:        id,
			OrderId:   orderId,
			TotalCost: totalCost,
			Type:      type2,
		}

		if _, hasPair := data[stdPair]; !hasPair {
			data[stdPair] = make([]exchanges.Trade, 0)
		}

		data[stdPair] = append(data[stdPair], trade)
	}

	resp.Data = data
	resp.Err = nil
	return
}

func (*KrakenClient) DepositsWithdrawals() exchanges.DepositsWithdrawalsResponse {
	panic("implement me")
}

func (self *KrakenClient) OrderBook(stdPair string) (resp exchanges.OrderBookResponse) {
	krakenPair, err := pairs.ExchangePair(stdPair, exchanges.KRAKEN)

	if err != nil {
		resp.Err = err
		return
	}

	bodyArgs := make(map[string]string)
	bodyArgs["pair"] = krakenPair

	clientResponse, err := self.apiClient.Do("POST", "/0/public/Depth", nil, bodyArgs, nil)

	var krakenResponse KrakenResponse
	krakenResponse.Result = &KrakenOrderBookResponse{}

	err = json.Unmarshal(clientResponse, &krakenResponse)

	if err != nil {
		log.Error("there was an error un-marshalling the order book", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "orderBook", "error", err)
		resp.Err = err
		return
	}

	err = handleKrakenErrors(krakenResponse)

	if err != nil {
		log.Error("kraken responded successfully with errors", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "chartData", "error", err)
		resp.Err = err
		return
	}

	orderBook := exchanges.OrderBook{Exchange: exchanges.KRAKEN, Pair: stdPair}
	asks := make(map[string]float64)
	bids := make(map[string]float64)

	krakenOrderBookRespPtr := krakenResponse.Result.(*KrakenOrderBookResponse)
	krakenOrderBookResp := *krakenOrderBookRespPtr

	for _, ask := range krakenOrderBookResp[krakenPair].Asks {
		price := ask.Price
		volume, err := strconv.ParseFloat(ask.Amount, 64)

		if err != nil {
			log.Error("error converting kraken amount to float", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "orderBook", "error", err)
			continue
		}

		asks[price] = volume
	}

	for _, bid := range krakenOrderBookResp[krakenPair].Bids {
		price := bid.Price
		volume, err := strconv.ParseFloat(bid.Amount, 64)

		if err != nil {
			log.Error("error converting kraken amount to float", "module", "exchanges", "exchange", exchanges.KRAKEN, "method", "orderBook", "error", err)
			continue
		}

		bids[price] = volume
	}

	orderBook.Bids = bids
	orderBook.Asks = asks

	resp.Data = orderBook
	resp.Err = nil

	return
}
