package collectors

import (
	"os"
	"sort"
	"errors"
	"time"
	"poloniex-go-api"
	"goku-bot/channels"
	log "github.com/sirupsen/logrus"
)

type Collector interface {
	Collect()
	Key() string
}

var (
	API_KEY    = os.Getenv("POLONIEX_API_KEY")
	API_SECRET = os.Getenv("POLONIEX_API_SECRET")
)

type AllCollectors map[string]Collector

type SortedCollectorKeys []string

func (a SortedCollectorKeys) Len() int           { return len(a) }
func (a SortedCollectorKeys) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortedCollectorKeys) Less(i, j int) bool { return a[i] < a[j] }

func sortKeys(collectors AllCollectors) (keys SortedCollectorKeys) {
	log.WithFields(log.Fields{"module": "collectors"}).Debug("sorting keys")
	for k := range collectors {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return
}

var collectors AllCollectors = make(AllCollectors)
var poloniex *poloniex_go_api.Poloniex

func Init() {
	thc := NewMyTradeHistoryCollector()
	collectors[thc.Key()] = thc

	poloniex = poloniex_go_api.New(API_KEY, API_SECRET)
	log.WithFields(log.Fields{"module": "collectors"}).Info("Finished initialization of Collectors module")
}

func Start() error {
	log.WithFields(log.Fields{"module": "collectors"}).Info("Starting Collectors module")
	if len(collectors) == 0 {
		return errors.New("failed to start collectors module. Collectors have not been initialized")
	}

	sortedCollectorKeys := sortKeys(collectors)

	for {
		for _, key := range sortedCollectorKeys {
			go runCollector(key)
			<-time.After(250 * time.Millisecond) // run collector every quarter second
		}
	}

	return nil
}

func runCollector(key string) {
	log.WithFields(log.Fields{"key": key, "module": "collectors"}).Debug("running collector")
	collectors[key].Collect()
}

func sendToProcessor(key string, output interface{}) {
	log.WithFields(log.Fields{"key": key, "module": "collectors"}).Debug("sending to processor module")
	collectorOutput := channels.CollectorOutputProcessorInput{output, key}

	channels.CollectorProcessorChan <- collectorOutput
}
