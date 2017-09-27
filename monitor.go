package goku_bot

import (
	"poloniex-go-api"
	"time"
	. "goku-bot/global"
	"log"
	"sync"
	"gopkg.in/mgo.v2/bson"
)

//TODO - Abstract 'Poloniex' to a list of 'Exchanges'

type Monitor struct {
	Poloniex *poloniex_go_api.Poloniex
	Store    *Store
}

func (m *Monitor) SyncMonitor(group *sync.WaitGroup) {
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
		for k, _ := range POLONIEX_OHLC_INTERVALS {
			log.Printf("pair: %s, interval: %d", pair, k)
			go m.BuildMonitorPoloniex(pair, k, syncOhlcErrorsCh, &wg)
		}
	}

	wg.Wait()

	log.Println("Finished Syncing OHLC")
	group.Done()
}

func (m *Monitor) BuildMonitorPoloniex(pair string, interval int, err chan<- error, group *sync.WaitGroup) {
	end := time.Now()
	start := end.AddDate(0, 0, -1)

	resp := m.Poloniex.ReturnChartData(pair, interval, start, end)

	if resp.Err != nil {
		log.Println("error getting the chart data")
		err <- resp.Err
	}

	m.Store.SyncCandles(resp.Response, "poloniex", pair, POLONIEX_OHLC_INTERVALS[interval])
	m.CalculateIndicators("poloniex", pair, interval)

	group.Done()
}

func (m *Monitor) CalculateIndicators(exchange, pair string, interval int) {
	collectionName := BuildTimeSliceCollectionName(exchange, pair, POLONIEX_OHLC_INTERVALS[interval])
	allSlices := m.Store.RetrieveSlicesByQueue(exchange, pair, interval, -1, -1)

	CalculateExponentialMovingAverage(10, allSlices)
	CalculateMacd(12, 26, 9, allSlices)
	CalculateAroon(25, allSlices)

	for _, slice := range allSlices {
		m.Store.Database.C(collectionName).Update(bson.M{"queue": slice.Queue}, slice)
	}
}
