package engine

import (
	"time"

	"github.com/google/uuid"
)

type Side int

const (
	BUY Side = iota
	SELL
)

func (s Side) String() string {
	if s == BUY {
		return "BUY"
	}
	return "SELL"
}

type Order struct {
	ID        uuid.UUID `json:"id"`
	Symbol    string    `json:"symbol"`
	Side      Side      `json:"side"`
	Price     float64   `json:"price"`
	Qty       int       `json:"qty"`
	Timestamp int64     `json:"timestamp"`
	UserID    string    `json:"user_id"`
}

func NewOrder(symbol string, side Side, price float64, qty int, userID string) *Order {
	return &Order{
		ID:        uuid.New(),
		Symbol:    symbol,
		Side:      side,
		Price:     price,
		Qty:       qty,
		Timestamp: time.Now().UnixNano(),
		UserID:    userID,
	}
}

type Trade struct {
	ID        uuid.UUID `json:"id"`
	Symbol    string    `json:"symbol"`
	BuyOrder  uuid.UUID `json:"buy_order"`
	SellOrder uuid.UUID `json:"sell_order"`
	Price     float64   `json:"price"`
	Qty       int       `json:"qty"`
	Timestamp int64     `json:"timestamp"`
	Side      Side      `json:"side"`
}

func NewTrade(symbol string, buyOrder, sellOrder uuid.UUID, price float64, qty int, side Side) *Trade {
	return &Trade{
		ID:        uuid.New(),
		Symbol:    symbol,
		BuyOrder:  buyOrder,
		SellOrder: sellOrder,
		Price:     price,
		Qty:       qty,
		Timestamp: time.Now().UnixNano(),
		Side:      side,
	}
}
