package goku_bot

type ActionQueue struct {
}

type Bot struct {
	Name     string
	Exchange string
	Pair     string
	Strategy func(exchange, pair string, store Store) (actionQueue ActionQueue, err error)
}

func NewBot(name, exchange, pair string, strategy func(exchange, pair string, store Store) (actionQueue ActionQueue, err error)) (bot *Bot) {
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

	store, err := NewStore()

	if err != nil {
		errCh <- err
		return
	}

	queue, err := b.Strategy(b.Exchange, b.Pair, *store)

	queueCh <- queue
	errCh <- err
}
