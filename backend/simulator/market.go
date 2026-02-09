package simulator

import (
	"math/rand"
	"time"
	"github.com/AkshatMadhani/nanopulse/engine"
	"github.com/AkshatMadhani/nanopulse/logger"
)

type Simulator struct {
	engine       *engine.MatchingEngine
	logger       *logger.Logger
	enabled      bool
	ordersPerSec int
}

func NewSimulator(eng *engine.MatchingEngine, log *logger.Logger, ordersPerSec int) *Simulator {
	return &Simulator{
		engine:       eng,
		logger:       log,
		enabled:      false,
		ordersPerSec: ordersPerSec,
	}
}

func (s *Simulator) Start() {
	s.enabled = true
	s.logger.Info("Starting market simulator", "orders_per_sec", s.ordersPerSec)
	go s.generateOrders()
}

func (s *Simulator) Stop() {
	s.enabled = false
	s.logger.Info("Stopping market simulator")
}

func (s *Simulator) generateOrders() {
	symbols := []string{"RELIANCE", "TCS", "INFY", "HDFC", "ICICI"}
	interval := time.Second / time.Duration(s.ordersPerSec)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if !s.enabled {
			return
		}

		symbol := symbols[rand.Intn(len(symbols))]
		s.generateRandomOrder(symbol)
	}
}

func (s *Simulator) generateRandomOrder(symbol string) {

	book := s.engine.GetBook(symbol)

	var basePrice float64
	if book != nil {
		bid := book.GetBestBid()
		ask := book.GetBestAsk()

		if bid != nil && ask != nil {
			basePrice = (*bid + *ask) / 2
		} else if bid != nil {
			basePrice = *bid
		} else if ask != nil {
			basePrice = *ask
		} else {
			basePrice = s.getDefaultPrice(symbol)
		}
	} else {
		basePrice = s.getDefaultPrice(symbol)
	}

	side := engine.BUY
	if rand.Float64() > 0.5 {
		side = engine.SELL
	}

	priceVariation := basePrice * 0.02 * (rand.Float64()*2 - 1)
	price := basePrice + priceVariation

	qty := rand.Intn(50) + 1

	order := engine.NewOrder(symbol, side, price, qty, "simulator")
	s.engine.GetOrderChan() <- order

	s.logger.Debug("Simulated order",
		"symbol", symbol,
		"side", side,
		"price", price,
		"qty", qty,
	)
}

func (s *Simulator) getDefaultPrice(symbol string) float64 {
	defaults := map[string]float64{
		"RELIANCE": 2500.0,
		"TCS":      3500.0,
		"INFY":     1500.0,
		"HDFC":     1600.0,
		"ICICI":    900.0,
	}

	if price, exists := defaults[symbol]; exists {
		return price
	}
	return 1000.0
}

func (s *Simulator) SetRate(ordersPerSec int) {
	s.ordersPerSec = ordersPerSec
	s.logger.Info("Simulator rate updated", "orders_per_sec", ordersPerSec)
}
