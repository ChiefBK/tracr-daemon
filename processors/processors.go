package processors

import (
	log "github.com/inconshreveable/log15"
	"goku-bot/channels"
	"goku-bot/exchanges"
	"goku-bot/pairs"
)

type Processor interface {
	Process(input interface{})
	Key() string
}

var processors = make(map[string]Processor)

func Init() {
	thp := NewMyTradeHistoryProcessor(exchanges.POLONIEX)
	processors[thp.Key()] = thp

	for pair := range pairs.PoloniexStdPairs {
		for _, interval := range exchanges.POLONIEX_INTERVALS {
			cdp := NewChartDataProcessor(exchanges.POLONIEX, pair, interval)
			processors[cdp.Key()] = cdp
		}
	}

	balancesProcessor := NewBalanceProcessor(exchanges.POLONIEX)
	processors[balancesProcessor.Key()] = balancesProcessor

	obr := NewOrderBookProcessor(exchanges.POLONIEX, pairs.BTC_USD)
	processors[obr.Key()] = obr

	tp := NewTickerProcessor(exchanges.POLONIEX, pairs.BTC_USD)
	processors[tp.Key()] = tp

	for pair := range pairs.PoloniexStdPairs {
		orderBookProcessor := NewOrderBookProcessor(exchanges.POLONIEX, pair)
		processors[orderBookProcessor.Key()] = orderBookProcessor
	}

	for key := range processors {
		log.Debug("Initialized processor", "module", "processors", "key", key)
	}
}

func StartProcessingCollectors() {
	log.Info("ready to process collector data", "module", "processors")
	for {
		input := <-channels.CollectorProcessorChan
		log.Debug("sending to processor", "module", "processors", "key", input.Key)
		go processors[input.Key].Process(input.Output)
	}
}

func StartProcessingReceivers() {
	log.Info("ready to process receiver data", "module", "processors")

	for {
		input := <-channels.ReceiverProcessorChan
		log.Debug("sending to processor", "module", "processors", "key", input.Key)
		go processors[input.Key].Process(input.Output)
	}
}