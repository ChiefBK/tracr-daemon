package main

import (
	"errors"
	"flag"
	"goku-bot"
	"log"
	"poloniex-go-api"
	"time"

	. "goku-bot/global"
	"goku-bot/strategies"
	"os"
	"sync"
)

var (
	API_KEY    = os.Getenv("POLONIEX_API_KEY")
	API_SECRET = os.Getenv("POLONIEX_API_SECRET")

	firstMonitor time.Time
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
		Store:    store,
		Exchange: "poloniex",
		Pair:     "USDT_BTC",
	}

	tickerSteward, err := goku_bot.NewTickerSteward()

	if err != nil {
		log.Println(err)
		return
	}

	// TODO - create new order book steward for each pair
	var websocketConnections sync.WaitGroup
	websocketConnections.Add(1)

	go orderBookSteward.ConnectPoloniexOrderBook("USDT_BTC", &websocketConnections)
	go tickerSteward.ConnectPoloniexTicker()

	websocketConnections.Wait()

	log.Printf("Websocket connections established and receiving")

	startGatheringAccountInfo()

	isTradesSynced := <-goku_bot.TradesSynced // make sure my trades have been synced
	log.Printf("trades received: %s", isTradesSynced)
	balances := <-goku_bot.PoloniexBalances // make sure balances have received data
	log.Printf("Balances received: %s", balances)
	ticker := <-goku_bot.TickerUsdtBtc // make sure usdt_btc ticker has received data
	log.Printf("ticker received: %s", ticker)

	log.Printf("Seeing how things go for 3 min")
	timer := time.NewTimer(time.Minute * 3)
	<-timer.C

	//runCandles()

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

	store, err := goku_bot.NewStore()

	if err != nil {
		err = errors.New("error creating store")
		return
	}

	log.Println("Initializing global variables")
	goku_bot.PoloniexClient = poloniex_go_api.New(API_KEY, API_SECRET)
	log.Println("Initialized Poloniex Client")

	if *clean {
		log.Println("Dropping DB")
		err = store.DropDatabase()
	}

	return
}

func startGatheringAccountInfo() {
	log.Println("Starting Account")

	accountSteward, err := goku_bot.NewAccountSteward()

	if err != nil {
		log.Printf("There was an error creating the Account Steward: %s", err)
		return
	}

	//go accountSteward.SyncBalances()
	go repeatFunction(accountSteward.SyncBalances, time.Second*5)
	go repeatFunction(accountSteward.SyncTrades, time.Second*10)
}

func runCandles() {
	log.Println("Starting Candles")

	if firstMonitor.IsZero() {
		firstMonitor = time.Now()
	}

	var group sync.WaitGroup
	group.Add(1)

	candlestickSteward, err := goku_bot.NewCandleStickSteward()

	if err != nil {
		log.Printf("There was an error creating the candlestick steward: %s", err)
		return
	}

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

func repeatFunction(f func(), every time.Duration) {
	for {
		f()
		<-time.After(every)
	}
}
