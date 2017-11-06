package goku_bot

import (
	"errors"
	. "goku-bot/global"
	"log"
	"poloniex-go-api"
	"sync"
	"time"
	"goku-bot/store"
)

type CandlestickSteward struct {
	Poloniex *poloniex_go_api.Poloniex
	Store    store.Store
}

func NewCandleStickSteward() (*CandlestickSteward, error) {
	store, err := store.NewStore()

	if err != nil {
		return nil, errors.New("there was an error creating the store")
	}

	if PoloniexClient == nil {
		return nil, errors.New("the poloniex client hasn't been initialized")
	}

	return &CandlestickSteward{PoloniexClient, store}, nil
}

func (self *CandlestickSteward) SyncCandles(group *sync.WaitGroup) {
	syncOhlcErrorsCh := make(chan error)
	defer close(syncOhlcErrorsCh)

	go func() {
		for {
			select {
			case err := <-syncOhlcErrorsCh:
				if err != nil {
					log.Println("There was an error syncing Poloniex candles")
					log.Println(err)
				}
			}

		}
	}()

	numWorkers := len(POLONIEX_PAIRS) * len(POLONIEX_OHLC_INTERVALS)

	log.Printf("There are %d workers", numWorkers)
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for _, pair := range POLONIEX_PAIRS {
		for k := range POLONIEX_OHLC_INTERVALS {
			log.Printf("pair: %s, interval: %d", pair, k)
			go self.BuildCandlesPoloniex(pair, k, syncOhlcErrorsCh, &wg)
		}
	}

	wg.Wait()

	log.Println("Finished Syncing OHLC")
	group.Done()
}

func (self *CandlestickSteward) BuildCandlesPoloniex(pair string, interval int, err chan<- error, group *sync.WaitGroup) {
	end := time.Now()
	start := end.AddDate(0, 0, -1)

	resp := self.Poloniex.ReturnChartData(pair, interval, start, end)

	if resp.Err != nil {
		log.Println("error getting the chart data")
		err <- resp.Err
	}

	self.Store.SyncCandles(resp.Data, "poloniex", pair, POLONIEX_OHLC_INTERVALS[interval])
	self.CalculateIndicators("poloniex", pair, interval)

	group.Done()
}

func (self *CandlestickSteward) CalculateIndicators(exchange, pair string, interval int) {
	//collectionName := buildCandleSliceCollectionName(exchange, pair, POLONIEX_OHLC_INTERVALS[interval])
	//allSlices := self.Store.RetrieveSlicesByQueue(exchange, pair, interval, -1, -1)
	//
	//CalculateExponentialMovingAverage(10, allSlices)
	//CalculateMacd(12, 26, 9, allSlices)
	//CalculateAroon(25, allSlices)
	//
	//for _, slice := range allSlices {
	//	self.Store.getCollection(collectionName).Update(bson.M{"queue": slice.Queue}, slice)
	//}
}
