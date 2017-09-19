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
)

const (
	API_KEY    = "8650VFZA-2POLX348-D69ZFDTC-AKQ2NEFM"
	API_SECRET = "dc79063b5781e7926521fd1c9b87efa276189af5b298fb84574c777cf19816f45f42624e598f9d0add0297dee527f6294babf255cead7dec9b5df19f2f228562"
)

var (
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

	go monitor.SyncOhlc(&group)

	group.Wait()

	log.Println("Finished Monitor")

	analyze()
}

func analyze() {
	log.Println("Starting Analyze")

	bot1ActionQueueCh := make(chan goku_bot.ActionQueue)
	bot1ErrorCh := make(chan error)

	bot1 := goku_bot.NewBot("bot1", "poloniex", BTC_ETH_PAIR, strategies.Strategy1)
	go bot1.RunStrategy(bot1ActionQueueCh, bot1ErrorCh)

	<- bot1ActionQueueCh
	<- bot1ErrorCh

	log.Println("Finished Analyze")
}
