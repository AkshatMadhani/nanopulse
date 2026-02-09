package engine

import (
	"container/heap"
	"sync"
	"time"

	"github.com/AkshatMadhani/nanopulse/logger"
)

type MatchingEngine struct {
	books       map[string]*OrderBook
	orderChan   chan *Order
	tradeChan   chan *Trade
	metricsChan chan Metric
	mu          sync.RWMutex
	logger      *logger.Logger
}

type Metric struct {
	Type      string
	Value     float64
	Timestamp int64
}

func NewMatchingEngine(orderBufferSize int, log *logger.Logger) *MatchingEngine {
	return &MatchingEngine{
		books:       make(map[string]*OrderBook),
		orderChan:   make(chan *Order, orderBufferSize),
		tradeChan:   make(chan *Trade, 1000),
		metricsChan: make(chan Metric, 1000),
		logger:      log,
	}
}

func (me *MatchingEngine) GetOrderChan() chan<- *Order {
	return me.orderChan
}

func (me *MatchingEngine) GetTradeChan() <-chan *Trade {
	return me.tradeChan
}

func (me *MatchingEngine) GetMetricsChan() <-chan Metric {
	return me.metricsChan
}

func (me *MatchingEngine) GetOrCreateBook(symbol string) *OrderBook {
	me.mu.Lock()
	defer me.mu.Unlock()

	if book, exists := me.books[symbol]; exists {
		return book
	}

	book := NewOrderBook(symbol)
	me.books[symbol] = book
	me.logger.Info("Created order book", "symbol", symbol)
	return book
}

func (me *MatchingEngine) GetBook(symbol string) *OrderBook {
	me.mu.RLock()
	defer me.mu.RUnlock()
	return me.books[symbol]
}

func (me *MatchingEngine) Start() {
	me.logger.Info("Starting matching engine")
	go me.processOrders()
}

func (me *MatchingEngine) processOrders() {
	for order := range me.orderChan {
		startTime := time.Now()

		me.matchOrder(order)

		latency := time.Since(startTime).Microseconds()
		me.metricsChan <- Metric{
			Type:      "latency",
			Value:     float64(latency),
			Timestamp: time.Now().UnixNano(),
		}
	}
}

func (me *MatchingEngine) matchOrder(order *Order) {
	book := me.GetOrCreateBook(order.Symbol)

	if order.Side == BUY {
		me.matchBuyOrder(book, order)
	} else {
		me.matchSellOrder(book, order)
	}
}

func (me *MatchingEngine) matchBuyOrder(book *OrderBook, buyOrder *Order) {
	book.mu.Lock()
	defer book.mu.Unlock()

	for buyOrder.Qty > 0 && book.SellHeap.Len() > 0 {
		bestSell := book.SellHeap.Peek()

		if buyOrder.Price < bestSell.Price {
			break
		}
		tradeQty := min(buyOrder.Qty, bestSell.Qty)
		tradePrice := bestSell.Price
		trade := NewTrade(
			book.Symbol,
			buyOrder.ID,
			bestSell.ID,
			tradePrice,
			tradeQty,
			buyOrder.Side,
		)
		me.tradeChan <- trade

		buyOrder.Qty -= tradeQty
		bestSell.Qty -= tradeQty

		if bestSell.Qty == 0 {
			heap.Pop(book.SellHeap)
			me.logger.Debug("Sell order fully filled", "order_id", bestSell.ID)
		}

		me.logger.Info("Trade executed",
			"symbol", book.Symbol,
			"price", tradePrice,
			"qty", tradeQty,
			"trade_id", trade.ID,
		)
	}
	if buyOrder.Qty > 0 {
		heap.Push(book.BuyHeap, buyOrder)
		me.logger.Debug("Buy order added to book",
			"order_id", buyOrder.ID,
			"remaining_qty", buyOrder.Qty,
		)
	}
}
func (me *MatchingEngine) matchSellOrder(book *OrderBook, sellOrder *Order) {
	book.mu.Lock()
	defer book.mu.Unlock()
	for sellOrder.Qty > 0 && book.BuyHeap.Len() > 0 {
		bestBuy := book.BuyHeap.Peek()
		if sellOrder.Price > bestBuy.Price {
			break
		}
		tradeQty := min(sellOrder.Qty, bestBuy.Qty)
		tradePrice := bestBuy.Price

		trade := NewTrade(
			book.Symbol,
			bestBuy.ID,
			sellOrder.ID,
			tradePrice,
			tradeQty,
			sellOrder.Side,
		)
		me.tradeChan <- trade

		sellOrder.Qty -= tradeQty
		bestBuy.Qty -= tradeQty

		if bestBuy.Qty == 0 {
			heap.Pop(book.BuyHeap)
			me.logger.Debug("Buy order fully filled", "order_id", bestBuy.ID)
		}

		me.logger.Info("Trade executed",
			"symbol", book.Symbol,
			"price", tradePrice,
			"qty", tradeQty,
			"trade_id", trade.ID,
		)
	}
	if sellOrder.Qty > 0 {
		heap.Push(book.SellHeap, sellOrder)
		me.logger.Debug("Sell order added to book",
			"order_id", sellOrder.ID,
			"remaining_qty", sellOrder.Qty,
		)
	}
}

func (me *MatchingEngine) GetQueueDepth() int {
	return len(me.orderChan)
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
