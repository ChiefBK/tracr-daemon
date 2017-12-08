package logging

import (
	"os"
	"path/filepath"
	log "github.com/inconshreveable/log15"
	"time"
)

var basePath = "logs"
var executorsBasePath = filepath.Join(basePath, "executors")
var strategiesBasePath = filepath.Join(basePath, "command")
var mainBasePath = filepath.Join(basePath, "main")
var processorsBasePath = filepath.Join(basePath, "processors")
var receiversBasePath = filepath.Join(basePath, "receivers")
var storeBasePath = filepath.Join(basePath, "store")
var brokerBasePath = filepath.Join(basePath, "broker")
var collectorsBasePath = filepath.Join(basePath, "collectors")
var exchangeCollectorsBasePath = filepath.Join(basePath, "exchangeCollectors")
var exchangesBasePath = filepath.Join(basePath, "exchanges")
var streamsBasePath = filepath.Join(basePath, "streams")

var basePaths = []string{executorsBasePath, strategiesBasePath, mainBasePath, processorsBasePath, receiversBasePath, storeBasePath, brokerBasePath, collectorsBasePath, exchangeCollectorsBasePath, exchangesBasePath, streamsBasePath}

func Init() {
	// create base folder structure
	for _, path := range basePaths {
		os.MkdirAll(path, os.ModePerm)
	}

	// Set handlers for logs
	now := time.Now()
	formattedTime := now.Format("2-Jan-2006-15:04:05")
	handler := log.MultiHandler(
		// direct log output from modules to their given log file
		log.MatchFilterHandler("module", "executors", log.Must.FileHandler(filepath.Join(executorsBasePath, formattedTime + "-executor.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "broker", log.Must.FileHandler(filepath.Join(brokerBasePath, formattedTime + "-broker.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "collectors", log.Must.FileHandler(filepath.Join(collectorsBasePath, formattedTime + "-collectors.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "main", log.Must.FileHandler(filepath.Join(mainBasePath, formattedTime + "-main.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "processors", log.Must.FileHandler(filepath.Join(processorsBasePath, formattedTime + "-processors.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "receivers", log.Must.FileHandler(filepath.Join(receiversBasePath, formattedTime + "-receivers.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "store", log.Must.FileHandler(filepath.Join(storeBasePath, formattedTime + "-store.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "command", log.Must.FileHandler(filepath.Join(strategiesBasePath, formattedTime + "-command.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "streams", log.Must.FileHandler(filepath.Join(streamsBasePath, formattedTime + "-streams.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "exchangeCollectors", log.Must.FileHandler(filepath.Join(exchangeCollectorsBasePath, formattedTime + "-exchangeCollectors.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "exchanges", log.Must.FileHandler(filepath.Join(exchangesBasePath, formattedTime + "-exchanges.txt"), log.JsonFormat())),
		// Also send log output to stdout
		log.LvlFilterHandler(log.LvlError, log.StderrHandler),
		log.LvlFilterHandler(log.LvlInfo, log.StderrHandler),
		log.LvlFilterHandler(log.LvlWarn, log.StderrHandler),
		log.LvlFilterHandler(log.LvlDebug, log.StderrHandler),
	)

	log.Root().SetHandler(handler)
}
