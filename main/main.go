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

	candlestickSteward *goku_bot.CandlestickSteward
)

func main() {
	err := initialize()

	if err != nil {
		log.Println("Initialization Error")
		log.Println(err)
		return
	}

	log.Println("Initialization Complete")

	store, err := goku_bot.NewStore()

	if err != nil {
		log.Printf("Error connecting to store: %s", err)
	}

	orderBookSteward := &goku_bot.OrderBookSteward{
		Store: store,
		Exchange: "poloniex",
		Pair: "USDT_BTC",
	}

	// TODO - create new order book steward for each pair
	var orderBookConnections sync.WaitGroup
	orderBookConnections.Add(1)
	go orderBookSteward.ConnectPoloniexOrderBook("USDT_BTC", &orderBookConnections)

	orderBookConnections.Wait()
	log.Printf("Done Waiting - orderbook stored")
	log.Printf("Starting timer")
	timer := time.NewTimer(time.Minute * 3)
	<-timer.C

	//runMonitor()

	//log.Println("Starting Cron job")
	//c := cron.New()
	//c.AddFunc("0 */1 * * * *", runMonitor)
	//c.Run()
}

func initialize() (err error) {
	log.Println("Initializing...")
	clean := flag.Bool("clean", false, "Clean DB on start")
	//single := flag.Bool("single", false, "")
	flag.Parse()

	startDateTime = time.Now()

	store, err := goku_bot.NewStore()
	log.Println("Created store for candlestick steward")

	if err != nil {
		err = errors.New("error creating store")
		return
	}

	poloniex := poloniex_go_api.New(API_KEY, API_SECRET)
	log.Println("Created Poloniex API interface for candlestick steward")

	if *clean {
		log.Println("Dropping DB")
		err = store.DropDatabase()
	}

	candlestickSteward = &goku_bot.CandlestickSteward{poloniex, store}
	log.Println("Initialized candlestick steward")

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

	go candlestickSteward.SyncCandles(&group)

	group.Wait()

	log.Println("Finished Monitor")

	runAnalyze()
}

func runAnalyze() {
	log.Println("Starting Analyze")

	bot1ActionQueueCh := make(chan goku_bot.ActionQueue)
	bot1ErrorCh := make(chan error)

	bot1 := goku_bot.NewBot("bot1", "poloniex", BTC_ETH_PAIR, strategies.Strategy1)
	go bot1.RunStrategy(bot1ActionQueueCh, bot1ErrorCh)

	<-bot1ActionQueueCh
	<-bot1ErrorCh

	log.Println("Finished Analyze")
}
