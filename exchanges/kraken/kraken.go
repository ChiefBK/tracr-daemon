package kraken

import (
	"tracr-client"
	"tracr-daemon/exchanges"
	"time"
	"tracr-daemon/pairs"
	"encoding/json"
	log "github.com/inconshreveable/log15"
	"strconv"
)

const BASE_URL = "https://api.kraken.com/0"

type KrakenClient struct {
	apiClient tracr_client.BaseClient
}

func NewKrakenClient(apiKey, apiSecret string) *KrakenClient {
	client := tracr_client.NewClient(apiKey, apiSecret, exchanges.KRAKEN, BASE_URL, BASE_URL, exchanges.KRAKEN_THROTTLE)
	return &KrakenClient{client}
}

func (self *KrakenClient) Ticker() exchanges.TickerResponse {
	panic("implement me")
}

func (*KrakenClient) Balances() exchanges.BalancesResponse {
	panic("implement me")
}

func (self *KrakenClient) ChartData(stdPair string, period time.Duration, start, end time.Time) (resp exchanges.ChartDataResponse) {
	krakenPair, err := pairs.ExchangePair(stdPair, exchanges.KRAKEN)

	if err != nil {
		log.Error("std pair argument not valid", "module", "exchanges", "exchange", exchanges.KRAKEN, "stdPair", stdPair, "error", err)
		resp.Err = err
		return
	}

	bodyArgs := make(map[string]string)
	bodyArgs["pair"] = krakenPair
	bodyArgs["interval"] = strconv.FormatInt(int64(period.Minutes()), 10)
	bodyArgs["since"] = strconv.FormatInt(start.Unix(), 10)

	exchangeResp, err := self.apiClient.Do("POST", "/public/OHLC", nil, bodyArgs, nil)

	if err != nil {
		log.Error("error getting making request", "module", "exchanges", "exchange", exchanges.KRAKEN, "stdPair", stdPair, "error", err)
		resp.Err = err
		return
	}

	var krakenResponse KrakenResponse
	krakenResponse.Result = new(map[string]interface{})

	err = json.Unmarshal(exchangeResp, &krakenResponse)

	if err != nil {
		log.Error("error un-marshalling exchange response", "module", "exchanges", "exchange", exchanges.KRAKEN, "stdPair", stdPair, "error", err)
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

func (*KrakenClient) MyTradeHistory() exchanges.TradeHistoryResponse {
	panic("implement me")
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

	clientResponse, err := self.apiClient.Do("POST", "/public/Depth", nil, bodyArgs, nil)

	var exchangeResponse KrakenResponse
	exchangeResponse.Result = &KrakenOrderBookResponse{}

	err = json.Unmarshal(clientResponse, &exchangeResponse)

	if err != nil {
		log.Error("there was an error un-marshalling the order book", "module", "exchanges", "exchange", exchanges.KRAKEN, "error", err)
		resp.Err = err
		return
	}

	orderBook := exchanges.OrderBook{Exchange: exchanges.KRAKEN, Pair: stdPair}
	asks := make(map[string]float64)
	bids := make(map[string]float64)

	krakenOrderBookRespPtr := exchangeResponse.Result.(*KrakenOrderBookResponse)
	krakenOrderBookResp := *krakenOrderBookRespPtr

	for _, ask := range krakenOrderBookResp[krakenPair].Asks {
		price := ask.Price
		volume, err := strconv.ParseFloat(ask.Amount, 64)

		if err != nil {
			log.Error("error converting kraken amount to float", "module", "exchanges", "exchange", exchanges.KRAKEN, "error", err)
			continue
		}

		asks[price] = volume
	}

	for _, bid := range krakenOrderBookResp[krakenPair].Bids {
		price := bid.Price
		volume, err := strconv.ParseFloat(bid.Amount, 64)

		if err != nil {
			log.Error("error converting kraken amount to float", "module", "exchanges", "exchange", exchanges.KRAKEN, "error", err)
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
