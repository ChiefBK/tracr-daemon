package main

import (
	log "github.com/inconshreveable/log15"
	"os/exec"
	"errors"
	"tracr-store"
	"tracr-daemon/logging"
	"tracr-daemon/util"
	"tracr-cache"
)

func initialize(logPath string, cleanDb bool, onOsx bool) (err error) {
	log.Info("Initializing...", "module", "main")
	//single := flag.Bool("single", false, "")

	err = startMongoDb(onOsx)

	if err != nil {
		return
	}

	err = startRedis()

	if err != nil {
		return
	}

	store, err := tracr_store.NewStore()

	if err != nil {
		return
	}

	cache, err := tracr_cache.NewCacheClient()

	if err != nil {
		return
	}

	if cleanDb {
		log.Info("Dropping DB")
		err = store.DropDatabase()

		if err != nil {
			return
		}

		log.Info("Clearing cache")
		err = cache.ClearCache()

		if err != nil {
			return
		}

	}

	logging.Init(logPath)

	//processors.Init()
	//receivers.Init()

	return
}

func startMongoDb(onOsx bool) error {
	log.Info("Seeing if MongoDB is already running")
	testCmd := exec.Command("mongod", "--fork", "--logpath", "/Users/ian/mongod.log")

	returnCode := util.ExecCommandWithCode(testCmd)
	log.Debug("test return code", "module", "main", "code", returnCode)
	if returnCode != 0 { // if mongoDB is already running
		log.Info("MongoDB is already running", "module", "main")
		return nil
	}

	log.Info("Starting MongoDB", "module", "main", "onOsx", onOsx)

	var cmd *exec.Cmd
	if onOsx {
		cmd = exec.Command("mongod", "--fork", "--logpath", "/Users/ian/mongod.log")
	} else {
		cmd = exec.Command("sudo", "service", "mongod", "start")
	}

	err := util.ExecCommand(cmd)
	if err != nil {
		log.Error("Error starting mongodb", "module", "main", "error", err)
		return errors.New("error starting mongod service")
	}

	return nil
}

func startRedis() error {
	log.Info("Starting Redis", "module", "main")

	cmd := exec.Command("redis-server", "--daemonize", "yes")
	err := util.ExecCommand(cmd)
	if err != nil {
		log.Error("Error starting redis server", "module", "main", "error", err)
		return errors.New("error starting redis")
	}

	return nil
}
