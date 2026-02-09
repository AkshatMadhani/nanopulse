package engine_test

import (
	"testing"
	"time"

	"github.com/AkshatMadhani/nanopulse/engine"
	"github.com/AkshatMadhani/nanopulse/logger"
)

func TestOrderCreation(t *testing.T) {
	order := engine.NewOrder("RELIANCE", engine.BUY, 2500.0, 10, "test-user")

	if order.Symbol != "RELIANCE" {
		t.Errorf("Expected symbol RELIANCE, got %s", order.Symbol)
	}

	if order.Side != engine.BUY {
		t.Errorf("Expected BUY side, got %v", order.Side)
	}

	if order.Price != 2500.0 {
		t.Errorf("Expected price 2500.0, got %f", order.Price)
	}

	if order.Qty != 10 {
		t.Errorf("Expected qty 10, got %d", order.Qty)
	}
}

func TestBuyHeapOrdering(t *testing.T) {
	heap := engine.NewBuyHeap()

	o1 := engine.NewOrder("TEST", engine.BUY, 2500.0, 10, "user1")
	o2 := engine.NewOrder("TEST", engine.BUY, 2505.0, 5, "user2")
	o3 := engine.NewOrder("TEST", engine.BUY, 2495.0, 15, "user3")

	heap.Push(o1)
	heap.Push(o2)
	heap.Push(o3)

	top := heap.Peek()
	if top.Price != 2505.0 {
		t.Errorf("Expected top price 2505.0, got %f", top.Price)
	}
}

func TestSellHeapOrdering(t *testing.T) {
	heap := engine.NewSellHeap()

	o1 := engine.NewOrder("TEST", engine.SELL, 2500.0, 10, "user1")
	o2 := engine.NewOrder("TEST", engine.SELL, 2505.0, 5, "user2")
	o3 := engine.NewOrder("TEST", engine.SELL, 2495.0, 15, "user3")

	heap.Push(o1)
	heap.Push(o2)
	heap.Push(o3)

	top := heap.Peek()
	if top.Price != 2495.0 {
		t.Errorf("Expected top price 2495.0, got %f", top.Price)
	}
}

func TestOrderBookCreation(t *testing.T) {
	book := engine.NewOrderBook("RELIANCE")

	if book.Symbol != "RELIANCE" {
		t.Errorf("Expected symbol RELIANCE, got %s", book.Symbol)
	}

	if book.GetBestBid() != nil {
		t.Error("Expected empty buy heap")
	}

	if book.GetBestAsk() != nil {
		t.Error("Expected empty sell heap")
	}
}

func TestOrderBookBestPrices(t *testing.T) {
	book := engine.NewOrderBook("TEST")

	buyOrder := engine.NewOrder("TEST", engine.BUY, 2500.0, 10, "user1")
	book.AddOrder(buyOrder)

	bestBid := book.GetBestBid()
	if bestBid == nil || *bestBid != 2500.0 {
		t.Errorf("Expected best bid 2500.0, got %v", bestBid)
	}

	sellOrder := engine.NewOrder("TEST", engine.SELL, 2505.0, 5, "user2")
	book.AddOrder(sellOrder)

	bestAsk := book.GetBestAsk()
	if bestAsk == nil || *bestAsk != 2505.0 {
		t.Errorf("Expected best ask 2505.0, got %v", bestAsk)
	}

	spread := book.GetSpread()
	if spread == nil || *spread != 5.0 {
		t.Errorf("Expected spread 5.0, got %v", spread)
	}
}

func TestMatchingEngineBasicMatch(t *testing.T) {
	log := logger.New(logger.ERROR)
	me := engine.NewMatchingEngine(100, log)
	me.Start()

	buyOrder := engine.NewOrder("TEST", engine.BUY, 2500.0, 10, "buyer")
	me.GetOrderChan() <- buyOrder

	time.Sleep(time.Millisecond * 10)

	sellOrder := engine.NewOrder("TEST", engine.SELL, 2500.0, 10, "seller")
	me.GetOrderChan() <- sellOrder

	time.Sleep(time.Millisecond * 10)

	select {
	case trade := <-me.GetTradeChan():
		if trade.Price != 2500.0 {
			t.Errorf("Expected trade price 2500.0, got %f", trade.Price)
		}
		if trade.Qty != 10 {
			t.Errorf("Expected trade qty 10, got %d", trade.Qty)
		}
	case <-time.After(time.Second):
		t.Error("Expected trade, but none received")
	}
}

func TestMatchingEnginePartialFill(t *testing.T) {
	log := logger.New(logger.ERROR)
	me := engine.NewMatchingEngine(100, log)
	me.Start()

	buyOrder := engine.NewOrder("TEST", engine.BUY, 2500.0, 20, "buyer")
	me.GetOrderChan() <- buyOrder

	time.Sleep(time.Millisecond * 10)

	sellOrder := engine.NewOrder("TEST", engine.SELL, 2500.0, 10, "seller")
	me.GetOrderChan() <- sellOrder

	time.Sleep(time.Millisecond * 10)

	select {
	case trade := <-me.GetTradeChan():
		if trade.Qty != 10 {
			t.Errorf("Expected partial fill of 10, got %d", trade.Qty)
		}
	case <-time.After(time.Second):
		t.Error("Expected trade, but none received")
	}

	book := me.GetBook("TEST")
	if book == nil {
		t.Fatal("Expected order book to exist")
	}

	bestBid := book.GetBestBid()
	if bestBid == nil || *bestBid != 2500.0 {
		t.Error("Expected remaining buy order in book")
	}
}

func TestMatchingEngineNoMatch(t *testing.T) {
	log := logger.New(logger.ERROR)
	me := engine.NewMatchingEngine(100, log)
	me.Start()

	buyOrder := engine.NewOrder("TEST", engine.BUY, 2500.0, 10, "buyer")
	me.GetOrderChan() <- buyOrder

	time.Sleep(time.Millisecond * 10)

	sellOrder := engine.NewOrder("TEST", engine.SELL, 2510.0, 10, "seller")
	me.GetOrderChan() <- sellOrder

	time.Sleep(time.Millisecond * 10)

	select {
	case <-me.GetTradeChan():
		t.Error("Did not expect a trade")
	case <-time.After(time.Millisecond * 100):
	}

	book := me.GetBook("TEST")
	if book.GetBestBid() == nil || book.GetBestAsk() == nil {
		t.Error("Expected both orders in book")
	}
}

func BenchmarkOrderCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = engine.NewOrder("TEST", engine.BUY, 2500.0, 10, "user")
	}
}

func BenchmarkMatchingEngine(b *testing.B) {
	log := logger.New(logger.ERROR)
	me := engine.NewMatchingEngine(10000, log)
	me.Start()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		order := engine.NewOrder("TEST", engine.BUY, 2500.0, 10, "user")
		me.GetOrderChan() <- order
	}
}
