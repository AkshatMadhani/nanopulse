package monitor

import (
	"time"

	"github.com/AkshatMadhani/nanopulse/engine"
	"github.com/AkshatMadhani/nanopulse/logger"
)

type SelfHealer struct {
	monitor      *Monitor
	engine       *engine.MatchingEngine
	logger       *logger.Logger
	injectionLog []LiquidityInjection
}

type LiquidityInjection struct {
	Symbol    string
	Side      engine.Side
	Price     float64
	Qty       int
	Reason    string
	Timestamp int64
}

func NewSelfHealer(mon *Monitor, eng *engine.MatchingEngine, log *logger.Logger) *SelfHealer {
	return &SelfHealer{
		monitor:      mon,
		engine:       eng,
		logger:       log,
		injectionLog: make([]LiquidityInjection, 0),
	}
}

func (sh *SelfHealer) Start() {
	sh.logger.Info("Starting self-healing system")
	go sh.monitorLiquidity()
}

func (sh *SelfHealer) monitorLiquidity() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for range ticker.C {
		book := sh.engine.GetBook("RELIANCE")
		if book != nil {
			sh.checkAndInjectLiquidity(book)
		}
	}
}

func (sh *SelfHealer) checkAndInjectLiquidity(book *engine.OrderBook) {
	bestBid := book.GetBestBid()
	bestAsk := book.GetBestAsk()

	if bestBid != nil && bestAsk == nil {
		sh.injectSellOrder(book, *bestBid)
		return
	}

	if bestBid == nil && bestAsk != nil {
		sh.injectBuyOrder(book, *bestAsk)
		return
	}

	if bestBid != nil && bestAsk != nil {
		spread := *bestAsk - *bestBid
		midPrice := (*bestBid + *bestAsk) / 2
		spreadPct := (spread / midPrice) * 100

		if spreadPct > 1.0 {
			sh.logger.Warn("Wide spread detected",
				"symbol", book.Symbol,
				"spread_pct", spreadPct,
			)
		}
	}
}

func (sh *SelfHealer) injectSellOrder(book *engine.OrderBook, basedOnBid float64) {
	price := basedOnBid + 2.0
	qty := 5

	order := engine.NewOrder(book.Symbol, engine.SELL, price, qty, "self-healer")

	sh.logger.Info("Injecting liquidity - SELL",
		"symbol", book.Symbol,
		"price", price,
		"qty", qty,
	)

	sh.engine.GetOrderChan() <- order

	sh.injectionLog = append(sh.injectionLog, LiquidityInjection{
		Symbol:    book.Symbol,
		Side:      engine.SELL,
		Price:     price,
		Qty:       qty,
		Reason:    "one-sided book (only bids)",
		Timestamp: time.Now().UnixNano(),
	})
}

func (sh *SelfHealer) injectBuyOrder(book *engine.OrderBook, basedOnAsk float64) {
	price := basedOnAsk - 2.0
	qty := 5

	order := engine.NewOrder(book.Symbol, engine.BUY, price, qty, "self-healer")

	sh.logger.Info("Injecting liquidity - BUY",
		"symbol", book.Symbol,
		"price", price,
		"qty", qty,
	)

	sh.engine.GetOrderChan() <- order
	sh.injectionLog = append(sh.injectionLog, LiquidityInjection{
		Symbol:    book.Symbol,
		Side:      engine.BUY,
		Price:     price,
		Qty:       qty,
		Reason:    "one-sided book (only asks)",
		Timestamp: time.Now().UnixNano(),
	})
}

func (sh *SelfHealer) GetInjectionCount() int {
	return len(sh.injectionLog)
}
