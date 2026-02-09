package engine

import "container/heap"

type BuyHeap []*Order

func (h BuyHeap) Len() int { return len(h) }

func (h BuyHeap) Less(i, j int) bool {
	if h[i].Price != h[j].Price {
		return h[i].Price > h[j].Price
	}
	return h[i].Timestamp < h[j].Timestamp
}

func (h BuyHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *BuyHeap) Push(x interface{}) {
	*h = append(*h, x.(*Order))
}

func (h *BuyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	order := old[n-1]
	old[n-1] = nil
	*h = old[0 : n-1]
	return order
}

func (h *BuyHeap) Peek() *Order {
	if h.Len() == 0 {
		return nil
	}
	return (*h)[0]
}

func NewBuyHeap() *BuyHeap {
	h := &BuyHeap{}
	heap.Init(h)
	return h
}
