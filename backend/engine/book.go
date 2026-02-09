package engine

import (
	"container/heap"
	"sync"
)

type OrderBook struct {
	Symbol   string
	BuyHeap  *BuyHeap
	SellHeap *SellHeap
	mu       sync.RWMutex
}

func NewOrderBook(symbol string) *OrderBook {
	return &OrderBook{
		Symbol:   symbol,
		BuyHeap:  NewBuyHeap(),
		SellHeap: NewSellHeap(),
	}
}

func (ob *OrderBook) AddOrder(order *Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if order.Side == BUY {
		heap.Push(ob.BuyHeap, order)
	} else {
		heap.Push(ob.SellHeap, order)
	}
}

func (ob *OrderBook) GetBestBid() *float64 {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if ob.BuyHeap.Len() == 0 {
		return nil
	}
	price := ob.BuyHeap.Peek().Price
	return &price
}

func (ob *OrderBook) GetBestAsk() *float64 {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if ob.SellHeap.Len() == 0 {
		return nil
	}
	price := ob.SellHeap.Peek().Price
	return &price
}

func (ob *OrderBook) GetSpread() *float64 {
	bid := ob.GetBestBid()
	ask := ob.GetBestAsk()

	if bid == nil || ask == nil {
		return nil
	}

	spread := *ask - *bid
	return &spread
}

type BookSnapshot struct {
	Symbol   string       `json:"symbol"`
	BuyBook  []PriceLevel `json:"buy_book"`
	SellBook []PriceLevel `json:"sell_book"`
	BestBid  *float64     `json:"best_bid"`
	BestAsk  *float64     `json:"best_ask"`
	Spread   *float64     `json:"spread"`
}

type PriceLevel struct {
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

func (ob *OrderBook) GetSnapshot(depth int) BookSnapshot {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	snapshot := BookSnapshot{
		Symbol:   ob.Symbol,
		BuyBook:  make([]PriceLevel, 0),
		SellBook: make([]PriceLevel, 0),
		BestBid:  ob.GetBestBid(),
		BestAsk:  ob.GetBestAsk(),
		Spread:   ob.GetSpread(),
	}

	priceMap := make(map[float64]int)
	for _, order := range *ob.BuyHeap {
		priceMap[order.Price] += order.Qty
	}
	for price, qty := range priceMap {
		snapshot.BuyBook = append(snapshot.BuyBook, PriceLevel{Price: price, Qty: qty})
		if len(snapshot.BuyBook) >= depth {
			break
		}
	}

	priceMap = make(map[float64]int)
	for _, order := range *ob.SellHeap {
		priceMap[order.Price] += order.Qty
	}
	for price, qty := range priceMap {
		snapshot.SellBook = append(snapshot.SellBook, PriceLevel{Price: price, Qty: qty})
		if len(snapshot.SellBook) >= depth {
			break
		}
	}

	return snapshot
}
