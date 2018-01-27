package logging

import (
	"os"
	"path/filepath"
	log "github.com/inconshreveable/log15"
	"time"
)


func Init(logPath string) {
	var basePath = filepath.Join("/var", "log", "tracr")

	if logPath != "" {
		basePath = logPath
	}

	log.Info("Using base log path", "module", "main", "path", basePath)

	var executorsBasePath = filepath.Join(basePath, "executors")
	var commandBasePath = filepath.Join(basePath, "command")
	var mainBasePath = filepath.Join(basePath, "main")
	var processorsBasePath = filepath.Join(basePath, "processors")
	var receiversBasePath = filepath.Join(basePath, "receivers")
	var storeBasePath = filepath.Join(basePath, "store")
	var brokerBasePath = filepath.Join(basePath, "broker")
	var collectorsBasePath = filepath.Join(basePath, "collectors")
	var exchangeCollectorsBasePath = filepath.Join(basePath, "exchangeCollectors")
	var exchangesBasePath = filepath.Join(basePath, "exchanges")
	var streamsBasePath = filepath.Join(basePath, "streams")

	var loggingPaths = []string{executorsBasePath, commandBasePath, mainBasePath, processorsBasePath, receiversBasePath, storeBasePath, brokerBasePath, collectorsBasePath, exchangeCollectorsBasePath, exchangesBasePath, streamsBasePath}

	// create base folder structure
	for _, path := range loggingPaths {
		err := os.MkdirAll(path, os.ModePerm)

		if err != nil {
			log.Error("Error creating log folder", "module", "main", "error", err, "path", path)
		}
	}

	// Set handlers for logs
	now := time.Now()
	formattedTime := now.Format("2-Jan-2006-15:04:05")
	log.Info("formatted time", "module", "logs", "time", formattedTime)
	handler := log.MultiHandler(
		// direct log output from modules to their given log file
		log.MatchFilterHandler("module", "executors", log.Must.FileHandler(filepath.Join(executorsBasePath, formattedTime + "-executor.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "broker", log.Must.FileHandler(filepath.Join(brokerBasePath, formattedTime + "-broker.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "collectors", log.Must.FileHandler(filepath.Join(collectorsBasePath, formattedTime + "-collectors.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "main", log.Must.FileHandler(filepath.Join(mainBasePath, formattedTime + "-main.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "processors", log.Must.FileHandler(filepath.Join(processorsBasePath, formattedTime + "-processors.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "receivers", log.Must.FileHandler(filepath.Join(receiversBasePath, formattedTime + "-receivers.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "store", log.Must.FileHandler(filepath.Join(storeBasePath, formattedTime + "-store.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "command", log.Must.FileHandler(filepath.Join(commandBasePath, formattedTime + "-command.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "streams", log.Must.FileHandler(filepath.Join(streamsBasePath, formattedTime + "-streams.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "exchangeCollectors", log.Must.FileHandler(filepath.Join(exchangeCollectorsBasePath, formattedTime + "-exchangeCollectors.txt"), log.JsonFormat())),
		log.MatchFilterHandler("module", "exchanges", log.Must.FileHandler(filepath.Join(exchangesBasePath, formattedTime + "-exchanges.txt"), log.JsonFormat())),
		// Also send log output to stdout
		log.LvlFilterHandler(log.LvlDebug, log.StderrHandler),
		log.LvlFilterHandler(log.LvlError, log.StderrHandler),
		log.LvlFilterHandler(log.LvlInfo, log.StderrHandler),
		log.LvlFilterHandler(log.LvlWarn, log.StderrHandler),
	)

	log.Root().SetHandler(handler)
}
