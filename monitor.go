package goku_bot

import (
	"poloniex-go-api"
	"time"
	. "goku-bot/global"
	"log"
	"sync"
)

//TODO - Abstract 'Poloniex' to a list of 'Exchanges'

type Monitor struct {
	Poloniex *poloniex_go_api.Poloniex
	Store    *Store
}

func (m *Monitor) SyncOhlc(group *sync.WaitGroup) {
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
			go m.SyncOhlcPoloniex(pair, k, syncOhlcErrorsCh, &wg)
		}
	}

	wg.Wait()

	log.Println("Finished Syncing OHLC")
	group.Done()
}

func (m *Monitor) SyncOhlcPoloniex(pair string, interval int, err chan<- error, group *sync.WaitGroup) {
	end := time.Now()
	start := end.AddDate(0, 0, -1)

	resp := m.Poloniex.ReturnChartData(pair, interval, start, end)

	if resp.Err != nil {
		log.Println("error getting the chart data")
		err <- resp.Err
	}

	m.Store.SyncCandles(resp.Response, "poloniex", pair, POLONIEX_OHLC_INTERVALS[interval])

	group.Done()
}
