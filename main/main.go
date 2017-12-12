package main

import (
	"errors"
	"flag"
	"time"

	log "github.com/inconshreveable/log15"
	"os"
	store2 "tracr-daemon/store"
	"tracr-daemon/collectors"
	"tracr-daemon/processors"
	"tracr-daemon/streams"
	"tracr-daemon/receivers"
	"tracr-daemon/logging"
)

var (
	API_KEY    = os.Getenv("POLONIEX_API_KEY")
	API_SECRET = os.Getenv("POLONIEX_API_SECRET")

	firstMonitor time.Time
)

func main() {
	err := initialize()

	if err != nil {
		log.Error("Initialization error", "module", "main")
		return
	}

	log.Info("Initialization Complete", "module", "main")

	//go collectors.Start()
	//go processors.StartProcessingCollectors()
	//go processors.StartProcessingReceivers()
	//go receivers.Start()
	//go streams.Start()
	//go executors.Start()

	//orderBook := streams.ReadOrderBook("poloniex", "USDT_BTC")
	//ticker := streams.ReadTicker("poloniex", "USDT_BTC")

	//if err != nil {
	//	log.Warn("There was an error Marshalling orderbook", "module", "main")
	//}

	//log.Printf("OrderBook: %s", ob)
	//log.Printf("OrderBook2: %s", orderBook)
	//log.Printf("ticker: %s", ticker)

	//go command.Start()

	//orderBookSteward := goku_bot.NewOrderBookSteward("USDT_BTC", "poloniex")
	//tickerSteward := goku_bot.NewTickerSteward()
	//
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//go orderBookSteward.ConnectPoloniexOrderBook("USDT_BTC")
	//go tickerSteward.ConnectPoloniexTicker()
	//
	//log.Printf("Websocket connections established and receiving")
	//
	//startGatheringAccountInfo()
	//
	//isTradesSynced := <-goku_bot.TradeHistorySynced // make sure my trades have been synced
	//log.Printf("trade history received: %s", isTradesSynced)
	//isDepositWithdrawalHistorySynced := <-goku_bot.DepositWithdrawalHistorySynced // make sure my deposit-withdrawal history have been synced
	//log.Printf("deposit-withdrawal history received: %s", isDepositWithdrawalHistorySynced)
	//balances := <-goku_bot.PoloniexBalances // make sure balances have received data
	//log.Printf("Balances received: %s", balances)
	//ticker := <-goku_bot.TickerUsdtBtc // make sure usdt_btc ticker has received data
	//log.Printf("ticker received: %s", ticker)
	//orderBookUsdtBtc := <-goku_bot.OrderBookChannels["USDT_BTC"]
	//log.Printf("orderBook bids : %s", orderBookUsdtBtc.GetBidsDescending())
	//log.Printf("orderBook asks : %s", orderBookUsdtBtc.GetAsksAscending())
	//
	//log.Printf("Seeing how things go for 3 min")
	//
	//tradeSteward, _ := goku_bot.NewTradeStewared()
	//
	//netProfit := tradeSteward.CalculateTradeNetProfit("poloniex", "USDT_BTC")
	//positions := tradeSteward.GetPositions("poloniex", "USDT_BTC")
	//positionResults := tradeSteward.CalculatePositionNetProfits("poloniex", "USDT_BTC")
	//
	//var netUsdSum float64 = 0
	//for _, result := range positionResults {
	//	netUsdSum += result.NetUsd
	//}
	//
	//log.Printf("Net USD: %f", netUsdSum)

	//btcBalance := streams.ReadBalance(exchanges.POLONIEX, pairs.BTC_POLONIEX)
	//log.Info("BTC balance", "module", "main", "balance", btcBalance)
	//
	//btcUsdOrderBook := streams.ReadOrderBook(exchanges.POLONIEX, pairs.BTC_USD)
	//log.Info("orderbook", "module", "main", "value", len(btcUsdOrderBook.Asks))

	timer := time.NewTimer(time.Minute * 3)
	<-timer.C

	//runCandles()

	//log.Println("Starting Cron job")
	//c := cron.New()
	//c.AddFunc("0 */1 * * * *", runMonitor)
	//c.Run()
}

func initialize() (err error) {
	log.Info("Initializing...", "module", "main")
	clean := flag.Bool("clean", false, "Clean DB on start")
	//single := flag.Bool("single", false, "")
	flag.Parse()

	store, err := store2.NewStore()

	if err != nil {
		err = errors.New("error creating connection to store")
		return
	}

	if *clean {
		log.Info("Dropping DB")
		err = store.DropDatabase()
	}

	logging.Init()
	collectors.Init()
	processors.Init()
	receivers.Init()
	streams.Init()

	return
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
