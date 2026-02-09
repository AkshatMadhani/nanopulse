package api

import (
	"sync"
	"time"

	"github.com/AkshatMadhani/nanopulse/engine"
)

type SystemState struct {
	Timestamp    int64
	BestBid      *float64
	BestAsk      *float64
	Spread       *float64
	LatencyUs    float64
	MaxLatencyUs float64
	Mode         string
	QueueDepth   int
	MMProfit     float64
	TotalTrades  int64
	OrderBooks   map[string]interface{}
	OrderBook    interface{}
	RecentTrade  *TradeEvent
}

type TradeEvent struct {
	ID        string
	Symbol    string
	Price     float64
	Qty       int
	Side      string
	Timestamp int64
	BuyerID   string
	SellerID  string
}

type TradeBuffer struct {
	trades []*engine.Trade
	mu     sync.RWMutex
}

func NewTradeBuffer() *TradeBuffer {
	return &TradeBuffer{
		trades: make([]*engine.Trade, 0, 50),
	}
}

func (tb *TradeBuffer) Add(trade *engine.Trade) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.trades = append([]*engine.Trade{trade}, tb.trades...)
	if len(tb.trades) > 50 {
		tb.trades = tb.trades[:50]
	}
}

func (tb *TradeBuffer) GetRecent() *engine.Trade {
	tb.mu.RLock()
	defer tb.mu.RUnlock()

	if len(tb.trades) > 0 {
		return tb.trades[0]
	}
	return nil
}

func (s *Server) startTradeListener() {
	go func() {
		for trade := range s.tradeChan {
			if trade != nil {
				s.tradeBuffer.Add(trade)
				s.logger.Info("Trade executed",
					"id", trade.ID.String(),
					"symbol", trade.Symbol,
					"price", trade.Price,
					"qty", trade.Qty,
					"side", trade.Side.String(),
					"buy_order", trade.BuyOrder.String(),
					"sell_order", trade.SellOrder.String(),
				)
			}
		}
	}()
}

func (s *Server) broadcastSystemState() {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for range ticker.C {
		state := s.collectSystemState()
		s.wsHub.Broadcast(state)
	}
}

func (s *Server) collectSystemState() SystemState {
	monitorStats := s.monitor.GetStats()
	mmStats := s.marketMaker.GetStats()

	symbols := []string{"RELIANCE", "TCS", "INFY", "HDFC", "ICICI"}

	orderBooks := make(map[string]interface{})
	var primaryOrderBook interface{}
	var bestBid, bestAsk, spread *float64

	for i, symbol := range symbols {
		book := s.engine.GetBook(symbol)
		if book != nil {
			bid := book.GetBestBid()
			ask := book.GetBestAsk()
			spr := book.GetSpread()

			snapshot := book.GetSnapshot(10)

			orderBookData := map[string]interface{}{
				"symbol":    snapshot.Symbol,
				"buy_book":  snapshot.BuyBook,
				"sell_book": snapshot.SellBook,
				"best_bid":  bid,
				"best_ask":  ask,
				"spread":    spr,
			}

			orderBooks[symbol] = orderBookData

			if i == 0 {
				primaryOrderBook = orderBookData
				bestBid = bid
				bestAsk = ask
				spread = spr
			}
		}
	}

	modeStr := "NORMAL"
	switch monitorStats.CurrentMode {
	case 0:
		modeStr = "NORMAL"
	case 1:
		modeStr = "SAFE"
	case 2:
		modeStr = "THROTTLED"
	}

	var recentTrade *TradeEvent
	trade := s.tradeBuffer.GetRecent()
	if trade != nil {
		recentTrade = &TradeEvent{
			ID:        trade.ID.String(),
			Symbol:    trade.Symbol,
			Price:     trade.Price,
			Qty:       trade.Qty,
			Side:      trade.Side.String(),
			Timestamp: trade.Timestamp,
			BuyerID:   trade.BuyOrder.String(),
			SellerID:  trade.SellOrder.String(),
		}
	}

	return SystemState{
		Timestamp:    time.Now().UnixNano(),
		BestBid:      bestBid,
		BestAsk:      bestAsk,
		Spread:       spread,
		LatencyUs:    monitorStats.AvgLatencyUs,
		MaxLatencyUs: monitorStats.MaxLatencyUs,
		Mode:         modeStr,
		QueueDepth:   s.engine.GetQueueDepth(),
		MMProfit:     mmStats.Profit,
		TotalTrades:  monitorStats.TotalTrades,
		OrderBooks:   orderBooks,
		OrderBook:    primaryOrderBook,
		RecentTrade:  recentTrade,
	}
}
