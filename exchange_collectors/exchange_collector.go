package exchange_collectors

import (
	"time"
	log "github.com/inconshreveable/log15"
	"errors"
	"tracr-daemon/exchanges"
	"tracr-daemon/exchanges/poloniex"
	"tracr-daemon/exchanges/kraken"
	"os"
	"sort"
	"tracr-daemon/exchange_collectors/collectors"
	"tracr-daemon/pairs"
)

var (
	POLONIEX_API_KEY    = os.Getenv("POLONIEX_API_KEY")
	POLONIEX_API_SECRET = os.Getenv("POLONIEX_API_SECRET")

	KRAKEN_API_KEY = os.Getenv("KRAKEN_API_KEY")
	KRAKEN_API_SECRET = os.Getenv("KRAKEN_API_SECRET")
)

type ExchangeCollector struct {
	Exchange   string
	throttle   time.Duration
	collectors map[string]collectors.Collector
}

func NewExchangeCollector(exchange string, throttle time.Duration) *ExchangeCollector {
	ec := &ExchangeCollector{exchange, throttle, make(map[string]collectors.Collector)}

	err := ec.init()

	if err != nil {
		log.Error("exchange collector failed to init", "exchange", exchange)
		return nil
	}

	return ec
}

func (self *ExchangeCollector) init() error {
	switch self.Exchange {
	case exchanges.POLONIEX:
		if len(POLONIEX_API_KEY) == 0 || len(POLONIEX_API_SECRET) == 0 { // if environmental variables not specified
			return errors.New("api key and secret env vars not specified")
		}

		client := poloniex.NewPoloniexClient(POLONIEX_API_KEY, POLONIEX_API_SECRET)

		myTradeHistoryCollector := collectors.NewMyTradeHistoryCollector(exchanges.POLONIEX, client)
		self.collectors[myTradeHistoryCollector.Key()] = myTradeHistoryCollector

		for pair := range pairs.PoloniexStdPairs {
			for _, interval := range exchanges.POLONIEX_INTERVALS {
				chartDataCollector := collectors.NewChartDataCollector(exchanges.POLONIEX, pair, interval, client)
				self.collectors[chartDataCollector.Key()] = chartDataCollector
			}
		}

		balancesCollector := collectors.NewBalancesCollector(exchanges.POLONIEX, client)
		self.collectors[balancesCollector.Key()] = balancesCollector

	case exchanges.KRAKEN:
		if len(KRAKEN_API_KEY) == 0 || len(KRAKEN_API_SECRET) == 0 { // if environmental variables not specified
			return errors.New("api key and secret env vars not specified")
		}

		apiClient := kraken.NewKrakenClient(KRAKEN_API_KEY, KRAKEN_API_SECRET)

		balancesCollector := collectors.NewBalancesCollector(exchanges.KRAKEN, apiClient)
		self.collectors[balancesCollector.Key()] = balancesCollector

		myTradeHistoryCollector := collectors.NewMyTradeHistoryCollector(exchanges.KRAKEN, apiClient)
		self.collectors[myTradeHistoryCollector.Key()] = myTradeHistoryCollector

		for stdPair := range pairs.KrakenStdPairs {
			for _, interval := range exchanges.KRAKEN_INTERVALS {
				chartDataCollector := collectors.NewChartDataCollector(exchanges.KRAKEN, stdPair, interval, apiClient)
				self.collectors[chartDataCollector.Key()] = chartDataCollector
			}
		}
	}

	for key := range self.collectors {
		log.Debug("Initialized collector", "module", "exchangeCollectors", "key", key)
	}

	return nil
}

func (self *ExchangeCollector) Start() error {
	log.Info("Starting exchange collector", "module", "exchangeCollectors", "exchange", self.Exchange)

	if len(self.collectors) == 0 {
		return errors.New("failed to start Exchange Collector. It contains no collectors")
	}

	sortedCollectorKeys := sortKeys(self.collectors)

	for {
		for _, key := range sortedCollectorKeys {
			go self.runCollector(key)
			<-time.After(self.throttle) // wait throttle time between running each collector
		}
	}

	return nil
}

func (self *ExchangeCollector) runCollector(key string) {
	log.Debug("running collector", "module", "exchangeCollectors", "key", key)
	collector := self.collectors[key]

	if collector == nil {
		log.Error("attempted to run nil collector", "module", "exchangeCollectors", "key", key)
		return
	}

	collector.Collect()
}

type SortedCollectorKeys []string

func (a SortedCollectorKeys) Len() int           { return len(a) }
func (a SortedCollectorKeys) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortedCollectorKeys) Less(i, j int) bool { return a[i] < a[j] }

func sortKeys(collectors map[string]collectors.Collector) (keys SortedCollectorKeys) {
	log.Debug("sorting keys", "module", "exchangeCollectors")
	for k := range collectors {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return
}