package goku_bot

type Bot struct {
	Pair     string
	Exchange string
	Interval string
	Strategy func()
}

func NewBot(func(pair, exchange, interval string)) {
	//bot := new(Bot)

}
