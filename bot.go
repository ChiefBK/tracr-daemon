package goku_bot

//type Position struct {
//	Type string // 'LONG' or 'SHORT'
//	Door string // 'OPEN' or 'CLOSED'
//}

//func (self Position) isOpen() bool {
//	return self.Door == "OPEN"
//}
//
//func (self Position) isLong() bool {
//	return self.Type == "LONG"
//}

type Bot struct {
	Name     string
	Exchange string
	Pair     string
	Position *Position
	Strategy func(exchange, pair string, indicator *Indicator, store *MgoStore) (actionQueue ActionQueue, err error)
}

func NewBot(name, exchange, pair string, strategy func(exchange, pair string, indicator *Indicator, store *MgoStore) (actionQueue ActionQueue, err error)) (bot *Bot) {
	bot = new(Bot)
	bot.Strategy = strategy
	bot.Exchange = exchange
	bot.Pair = pair
	bot.Name = name

	return
}

func (b *Bot) RunStrategy(queueCh chan ActionQueue, errCh chan error) {
	defer close(queueCh)
	defer close(errCh)

	store, err := NewMgoStore()

	if err != nil {
		errCh <- err
		return
	}

	indicator := NewIndicator()

	queue, err := b.Strategy(b.Exchange, b.Pair, indicator, store)

	queueCh <- queue
	errCh <- err
}
