package main

import (
	"poloniex-go-api"
	"time"
	"goku-bot"
	"gopkg.in/mgo.v2"
	"log"
	"flag"
	"errors"
	//"github.com/robfig/cron"
)

const (
	API_KEY                  = "8650VFZA-2POLX348-D69ZFDTC-AKQ2NEFM"
	API_SECRET               = "dc79063b5781e7926521fd1c9b87efa276189af5b298fb84574c777cf19816f45f42624e598f9d0add0297dee527f6294babf255cead7dec9b5df19f2f228562"
	DB_URI                   = "localhost"
	DB_NAME                  = "goku-bot"
	NUMBER_MONITOR_PROCESSES = 3
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

	var dbErr error
	session, dbErr := mgo.Dial(DB_URI)
	if dbErr != nil {
		err = errors.New("Could not connect to db")
		return
	}

	db := session.DB(DB_NAME)

	if *clean {
		log.Println("Dropping DB")
		err = db.DropDatabase()
	}

	store := &goku_bot.Store{Database: db, Session: session}
	polo := poloniex_go_api.New(API_KEY, API_SECRET)

	monitor = &goku_bot.Monitor{polo, store}

	return
}

func runMonitor() {
	log.Println("Starting Monitor")

	if firstMonitor.IsZero() {
		firstMonitor = time.Now()
	}

	lastMonitor = time.Now()

	defer analyze()

	//var wg sync.WaitGroup
	//wg.Add(NUMBER_MONITOR_PROCESSES)
	//
	//go store.StoreBtcBalances(&wg)
	//go store.StoreLoanOffers(&wg)
	//go store.StoreActiveLoans(&wg)

	syncOHLCch := make(chan error)
	go monitor.SyncOHLC(syncOHLCch)

	syncOHLCRes := <-syncOHLCch

	if syncOHLCRes != nil {
		log.Println("Error syncing OHLC")
	}
	//wg.Wait()

	log.Println("Finished Monitor")
}

func analyze() {
	log.Println("Starting Analyze")

}

//func poller(){
//	done := make(chan bool)
//
//	monitor(done)
//}
