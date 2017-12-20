package main

import (
	log "github.com/inconshreveable/log15"
	"os/exec"
	"errors"
	"flag"
	"tracr-store"
	"tracr-daemon/logging"
	"tracr-daemon/collectors"
	"tracr-daemon/processors"
	"tracr-daemon/receivers"
	"tracr-daemon/streams"
)

func initialize() (err error) {
	log.Info("Initializing...", "module", "main")
	clean := flag.Bool("clean", false, "Clean DB on start")
	//single := flag.Bool("single", false, "")
	flag.Parse()

	err = startMongoDb()

	if err != nil {
		return
	}

	store, err := tracr_store.NewStore()

	if err != nil {
		return
	}

	if *clean {
		log.Info("Dropping DB")
		err = store.DropDatabase()
	}

	logging.Init()
	err = collectors.Init()

	if err != nil {
		return
	}

	processors.Init()
	receivers.Init()
	streams.Init()

	return
}

func startMongoDb() error {
	log.Info("Starting MongoDB", "module", "main")

	cmd := exec.Command("sudo", "service", "mongod", "start")
	startErr := cmd.Start()
	if startErr != nil {
		log.Error("Error starting mongodb", "module", "main", "error", startErr)
		return errors.New("error starting mongod service (start)")
	}

	waitErr := cmd.Wait()
	if waitErr != nil {
		log.Error("Error starting mongodb", "module", "main", "error", waitErr)
		return errors.New("error starting mongod service (wait)")
	}

	return nil
}
