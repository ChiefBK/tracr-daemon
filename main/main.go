package main

import (
	"time"
	log "github.com/inconshreveable/log15"
	"os"
	"tracr-daemon/collectors"
	"flag"
)

var (
	API_KEY    = os.Getenv("POLONIEX_API_KEY")
	API_SECRET = os.Getenv("POLONIEX_API_SECRET")
)

/*
	Program usage

	tracrd start [--clean] [--log-path <path>]
	tracrd stop
	tracrd monitor <exchange name>
	tracrd monitor <exchange name> <indicator>
	tracrd monitor <indicator>


	Options

	--clean							wipe database and cache before start
	--help -h						show help
	--log-path <path>, -l <path>	specify log path
	--clear-logs, -c				delete logs before starting


	see http://docopt.org/ for docs on program usage syntax

 */

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		// show help
		return
	}

	logPath1 := flag.String("log-path", "", "The base logging path")
	logPath2 := flag.String("l", "", "The base logging path")
	clean := flag.Bool("clean", false, "Clean DB on start")
	onOsx := flag.Bool("osx", false, "Is running on Mac OSX?")
	flag.Parse()

	var logPath string

	if logPath1 != nil {
		logPath = *logPath1
	} else if logPath2 != nil {
		logPath = *logPath2
	} else {
		logPath = ""
	}

	log.Debug("log path cmd line args", "module", "main", "logpath1", *logPath1, "logpath2", *logPath2)
	action := args[len(args)-1]

	switch action {
	case "start":
		start(logPath, *clean, *onOsx)
	case "stop":
		stop()
	case "monitor":
		monitor()
	default:
		// error
		log.Error("action not defined - exiting")
		return
	}

	// TODO - add signal for catching when user wants to terminate process. Add graceful shutdown
	timer := time.NewTimer(time.Minute * 3)
	<-timer.C
}

func start(logPath string, cleanDb bool, onOsx bool) {
	err := initialize(logPath, cleanDb, onOsx)

	if err != nil {
		log.Error("Initialization error", "module", "main", "error", err)
		return
	}

	log.Info("Initialization Complete", "module", "main")

	go collectors.Start()
	//go processors.StartProcessingCollectors()
	//go processors.StartProcessingReceivers()
	//go receivers.Start()
}

func stop() {

}

func monitor() {

}

//func startGatheringAccountInfo() {
//	log.Println("Starting Account")
//
//	accountSteward, err := goku_bot.NewAccountSteward()
//
//	if err != nil {
//		log.Printf("There was an error creating the Account Steward: %s", err)
//		return
//	}
//
//	//go accountSteward.SyncBalances()
//	go repeatFunction(accountSteward.SyncBalances, time.Second*5)
//	go repeatFunction(accountSteward.SyncTradeHistory, time.Second*10)
//	go repeatFunction(accountSteward.SyncDepositWithdrawlHistory, time.Second*10)
//}
//
//func runCandles() {
//	log.Println("Starting Candles")
//
//	if firstMonitor.IsZero() {
//		firstMonitor = time.Now()
//	}
//
//	var group sync.WaitGroup
//	group.Add(1)
//
//	candlestickSteward, err := goku_bot.NewCandleStickSteward()
//
//	if err != nil {
//		log.Printf("There was an error creating the candlestick steward: %s", err)
//		return
//	}
//
//	go candlestickSteward.SyncCandles(&group)
//
//	group.Wait()
//
//	log.Println("Finished Monitor")
//
//	runAnalyze()
//}
//
//func runAnalyze() {
//	log.Println("Starting Analyze")
//
//	bot1ActionQueueCh := make(chan actions.ActionQueue)
//	bot1ErrorCh := make(chan error)
//
//	//bot1 := goku_bot.NewBot("bot1", "poloniex", BTC_ETH_PAIR, command.Strategy1)
//	//go bot1.RunStrategy(bot1ActionQueueCh, bot1ErrorCh)
//
//	<-bot1ActionQueueCh
//	<-bot1ErrorCh
//
//	log.Println("Finished Analyze")
//}

//func repeatFunction(f func(), every time.Duration) {
//	for {
//		f()
//		<-time.After(every)
//	}
//}
