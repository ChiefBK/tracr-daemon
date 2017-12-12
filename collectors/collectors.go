/**

The collectors module

 */

package collectors

import (
	"os"
	"sort"
	"time"
	"tracr-daemon/channels"
	log "github.com/inconshreveable/log15"
	"tracr-daemon/exchanges"
)

var (
	API_KEY    = os.Getenv("POLONIEX_API_KEY")
	API_SECRET = os.Getenv("POLONIEX_API_SECRET")
)

type SortedCollectorKeys []string

func (a SortedCollectorKeys) Len() int           { return len(a) }
func (a SortedCollectorKeys) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortedCollectorKeys) Less(i, j int) bool { return a[i] < a[j] }

func sortKeys(collectors map[string]Collector) (keys SortedCollectorKeys) {
	log.Debug("sorting keys", "module", "exchangeCollectors")
	for k := range collectors {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return
}

var exchangeCollectors []*ExchangeCollector

// initialize all exchange exchangeCollectors
func Init() {
	poloniexCollector := NewExchangeCollector(exchanges.POLONIEX, 200*time.Millisecond)
	poloniexCollector.Init()

	exchangeCollectors = append(exchangeCollectors, poloniexCollector)

	for _, ec := range exchangeCollectors {
		log.Debug("Initialized exchange collector", "module", "exchangeCollectors", "exchange", ec.exchange)
	}
	log.Info("Finished initialization of Collectors module", "module", "exchangeCollectors")
}

func Start() {
	for _, exchangeCollector := range exchangeCollectors {
		go exchangeCollector.Start()
	}
}

func sendToProcessor(key string, output interface{}) {
	log.Debug("sending to processor module", "module", "exchangeCollectors", "key", key)
	collectorOutput := channels.CollectorOutputProcessorInput{output, key}

	channels.CollectorProcessorChan <- collectorOutput
}
