package processors

import (
	log "github.com/inconshreveable/log15"
	"goku-bot/channels"
)

type Processor interface {
	Process(input interface{})
	Key() string
}

var processors = make(map[string]Processor)

func Init() {
	thp := NewMyTradeHistoryProcessor()
	processors[thp.Key()] = thp

	obr := NewOrderBookProcessor("poloniex", "USDT_BTC")
	processors[obr.Key()] = obr

	tp := NewTickerProcessor("poloniex", "USDT_BTC")
	processors[tp.Key()] = tp
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