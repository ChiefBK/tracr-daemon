package collectors

import (
	"time"
	log "github.com/inconshreveable/log15"
	"errors"
	"goku-bot/exchanges"
	"goku-bot/exchanges/poloniex"
	"os"
	"goku-bot/pairs"
)

type ExchangeCollector struct {
	exchange   string
	throttle   time.Duration
	collectors map[string]Collector
}

func NewExchangeCollector(exchange string, throttle time.Duration) *ExchangeCollector {
	return &ExchangeCollector{exchange, throttle, make(map[string]Collector)}
}

func (self *ExchangeCollector) Init() {
	switch self.exchange {
	case exchanges.POLONIEX:
		client := poloniex.NewPoloniexClient(os.Getenv("POLONIEX_API_KEY"), os.Getenv("POLONIEX_API_SECRET"))

		myTradeHistoryCollector := NewMyTradeHistoryCollector(exchanges.POLONIEX, client)
		self.collectors[myTradeHistoryCollector.Key()] = myTradeHistoryCollector

		for pair := range pairs.PoloniexStdPairs {
			for _, interval := range exchanges.POLONIEX_INTERVALS {
				chartDataCollector := NewChartDataCollector(exchanges.POLONIEX, pair, interval, client)
				self.collectors[chartDataCollector.Key()] = chartDataCollector
			}
		}

		balancesCollector := NewBalancesCollector(exchanges.POLONIEX, client)
		self.collectors[balancesCollector.Key()] = balancesCollector



	case exchanges.KRAKEN:

	}

	for key := range self.collectors {
		log.Debug("Initialized collector", "module", "exchangeCollectors", "key", key)
	}
}

func (self *ExchangeCollector) Start() error {
	log.Info("Starting exchange collector", "module", "exchangeCollectors", "exchange", self.exchange)

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
	self.collectors[key].Collect()
}
