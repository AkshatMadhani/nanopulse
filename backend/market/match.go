package market

import (
	"math/rand"
	"sync"
	"time"

	"github.com/AkshatMadhani/nanopulse/engine"
	"github.com/AkshatMadhani/nanopulse/logger"
)

type Bot struct {
	engine       *engine.MatchingEngine
	tradeChan    <-chan *engine.Trade
	logger       *logger.Logger
	profit       float64
	profitMu     sync.RWMutex
	totalOrders  int64
	activeOrders map[string]bool
	mu           sync.Mutex
}

func NewBot(eng *engine.MatchingEngine, tradeChan <-chan *engine.Trade, log *logger.Logger) *Bot {
	return &Bot{
		engine:       eng,
		tradeChan:    tradeChan,
		logger:       log,
		activeOrders: make(map[string]bool),
	}
}

func (b *Bot) Start() {
	b.logger.Info("Starting market maker bot")
	go b.trackTrades()
	go b.provideQuotes()
}

func (b *Bot) trackTrades() {
	for trade := range b.tradeChan {
		b.logger.Debug("Trade observed",
			"symbol", trade.Symbol,
			"price", trade.Price,
			"qty", trade.Qty,
		)
	}
}

func (b *Bot) provideQuotes() {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	symbols := []string{"RELIANCE", "TCS", "INFY"}

	for range ticker.C {
		symbol := symbols[rand.Intn(len(symbols))]
		b.makeMarket(symbol)
	}
}

func (b *Bot) makeMarket(symbol string) {
	book := b.engine.GetBook(symbol)
	if book == nil {
		b.createInitialQuotes(symbol, 2500.0)
		return
	}

	bestBid := book.GetBestBid()
	bestAsk := book.GetBestAsk()

	var fairValue float64
	if bestBid != nil && bestAsk != nil {
		fairValue = (*bestBid + *bestAsk) / 2
	} else if bestBid != nil {
		fairValue = *bestBid + 1.0
	} else if bestAsk != nil {
		fairValue = *bestAsk - 1.0
	} else {
		fairValue = 2500.0
	}

	spread := 2.0
	bidPrice := fairValue - spread/2
	askPrice := fairValue + spread/2

	b.placeQuote(symbol, engine.BUY, bidPrice, 10)
	b.placeQuote(symbol, engine.SELL, askPrice, 10)
}
func (b *Bot) createInitialQuotes(symbol string, basePrice float64) {
	spread := 2.0

	buyOrder := engine.NewOrder(symbol, engine.BUY, basePrice-spread/2, 10, "market-maker")
	sellOrder := engine.NewOrder(symbol, engine.SELL, basePrice+spread/2, 10, "market-maker")

	b.engine.GetOrderChan() <- buyOrder
	b.engine.GetOrderChan() <- sellOrder

	b.mu.Lock()
	b.totalOrders += 2
	b.mu.Unlock()

	b.logger.Info("Created initial quotes",
		"symbol", symbol,
		"bid", basePrice-spread/2,
		"ask", basePrice+spread/2,
	)
}
func (b *Bot) placeQuote(symbol string, side engine.Side, price float64, qty int) {
	order := engine.NewOrder(symbol, side, price, qty, "market-maker")

	b.engine.GetOrderChan() <- order

	b.mu.Lock()
	b.totalOrders++
	b.mu.Unlock()

	b.logger.Debug("Market maker quote",
		"symbol", symbol,
		"side", side,
		"price", price,
		"qty", qty,
	)
}

func (b *Bot) GetProfit() float64 {
	b.profitMu.RLock()
	defer b.profitMu.RUnlock()
	return b.profit
}

func (b *Bot) AddProfit(amount float64) {
	b.profitMu.Lock()
	b.profit += amount
	b.profitMu.Unlock()

	b.logger.Info("Market maker profit",
		"amount", amount,
		"total", b.profit,
	)
}

func (b *Bot) GetStats() Stats {
	b.mu.Lock()
	totalOrders := b.totalOrders
	b.mu.Unlock()

	return Stats{
		Profit:      b.GetProfit(),
		TotalOrders: totalOrders,
	}
}

type Stats struct {
	Profit      float64 `json:"profit"`
	TotalOrders int64   `json:"total_orders"`
}
