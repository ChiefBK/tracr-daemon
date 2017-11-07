package processors

import (
	log "github.com/sirupsen/logrus"
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
	log.WithFields(log.Fields{"module": "processors"}).Info("ready to process collector data")
	for {
		input := <-channels.CollectorProcessorChan
		log.WithFields(log.Fields{"key": input.Key, "module": "processors"}).Debug("sending to processor")
		go processors[input.Key].Process(input.Output)
	}
}

func StartProcessingReceivers() {
	log.WithFields(log.Fields{"module": "processors"}).Info("ready to process receiver data")
	for {
		input := <-channels.ReceiverProcessorChan
		log.WithFields(log.Fields{"key": input.Key, "module": "processors"}).Debug("sending to processor")
		go processors[input.Key].Process(input.Output)
	}
}