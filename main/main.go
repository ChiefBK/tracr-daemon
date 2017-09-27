package main

import (
	"poloniex-go-api"
	"time"
	"goku-bot"
	"log"
	"flag"
	"errors"

	"sync"
	"goku-bot/strategies"
	. "goku-bot/global"
	"os"
)

var (
	API_KEY    = os.Getenv("POLONIEX_API_KEY")
	API_SECRET = os.Getenv("POLONIEX_API_SECRET")

	startDateTime time.Time
	firstMonitor  time.Time
	lastMonitor   time.Time

	monitor *goku_bot.Monitor
)

func main() {
	err := initialize()

	if err != nil {
		log.Println("Initialization Error")
		log.Println(err)
	}

	log.Println("Initialization Complete")

	runMonitor()

	//log.Println("Starting Cron job")
	//c := cron.New()
	//c.AddFunc("0 */1 * * * *", runMonitor)
	//c.Run()
}

func initialize() (err error) {
	clean := flag.Bool("clean", false, "Clean DB on start")
	//single := flag.Bool("single", false, "")
	flag.Parse()

	startDateTime = time.Now()

	store, err := goku_bot.NewStore()

	if err != nil {
		err = errors.New("error creating store")
		return
	}

	poloniex := poloniex_go_api.New(API_KEY, API_SECRET)

	if *clean {
		log.Println("Dropping DB")
		err = store.Database.DropDatabase()
	}

	monitor = &goku_bot.Monitor{poloniex, store}

	return
}

func runMonitor() {
	log.Println("Starting Monitor")

	if firstMonitor.IsZero() {
		firstMonitor = time.Now()
	}

	lastMonitor = time.Now()

	var group sync.WaitGroup
	group.Add(1)

	go monitor.SyncMonitor(&group)

	group.Wait()

	log.Println("Finished Monitor")

	//runAnalyze()
}

//func calculateTechnicalIndicators() error {
//	store, err := goku_bot.NewStore()
//
//	if err != nil {
//		return errors.New("error creating store - can not calculate technical indicators")
//	}
//
//	ohlc := store.RetrieveSlicesByQueue(EXCHANGE_POLONIEX, USDT_BTC_PAIR, FIVE_MIN_INTERVAL, -1, -1) // retrieve all candles
//
//	if len(ohlc) == 0 { // If table is empty or doesn't exist
//
//	}
//
//	candles := goku_bot.GetCandles(ohlc)
//	dateValues := goku_bot.GetDateValues(candles)
//	log.Println("DATE VALUES:")
//	log.Println(dateValues)
//	log.Println(len(dateValues))
//
//	sma := goku_bot.CalculateSimpleMovingAverage(3, dateValues)
//	log.Println("SMA:")
//	log.Println(sma)
//	log.Println(len(sma))
//
//	ema := goku_bot.CalculateExponentialMovingAverage(3, dateValues)
//	log.Println("EMA:")
//	log.Println(ema)
//	log.Println(len(ema))
//
//	macd := goku_bot.Macd(12, 26, 9, dateValues)
//
//	log.Println("MACD:")
//	log.Println(macd)
//	log.Println(len(macd))
//
//	return nil
//}

func runAnalyze() {
	log.Println("Starting Analyze")

	//err := calculateTechnicalIndicators()

	//if err != nil {
	//	log.Println("Error during Analyze - Aborting")
	//	log.Println(err)
	//	return
	//}

	bot1ActionQueueCh := make(chan goku_bot.ActionQueue)
	bot1ErrorCh := make(chan error)

	bot1 := goku_bot.NewBot("bot1", "poloniex", BTC_ETH_PAIR, strategies.Strategy1)
	go bot1.RunStrategy(bot1ActionQueueCh, bot1ErrorCh)

	<-bot1ActionQueueCh
	<-bot1ErrorCh

	log.Println("Finished Analyze")
}
