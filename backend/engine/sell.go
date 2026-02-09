package engine

import "container/heap"

type SellHeap []*Order

func (h SellHeap) Len() int { return len(h) }

func (h SellHeap) Less(i, j int) bool {
	if h[i].Price != h[j].Price {
		return h[i].Price < h[j].Price
	}
	return h[i].Timestamp < h[j].Timestamp
}

func (h SellHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *SellHeap) Push(x interface{}) {
	*h = append(*h, x.(*Order))
}

func (h *SellHeap) Pop() interface{} {
	old := *h
	n := len(old)
	order := old[n-1]
	old[n-1] = nil
	*h = old[0 : n-1]
	return order
}

func (h *SellHeap) Peek() *Order {
	if h.Len() == 0 {
		return nil
	}
	return (*h)[0]
}

func NewSellHeap() *SellHeap {
	h := &SellHeap{}
	heap.Init(h)
	return h
}
