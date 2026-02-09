package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AkshatMadhani/nanopulse/engine"
)

type OrderRequest struct {
	Symbol string
	Side   string
	Price  float64
	Qty    int
	UserID string
}

type OrderResponse struct {
	Status  string
	OrderID string
	Message string
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := s.monitor.GetStats()

	response := map[string]interface{}{
		"status":          "healthy",
		"mode":            stats.CurrentMode.String(),
		"avg_latency_us":  stats.AvgLatencyUs,
		"queue_depth":     s.engine.GetQueueDepth(),
		"injection_count": s.selfHealer.GetInjectionCount(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Symbol == "" || req.Price <= 0 || req.Qty <= 0 {
		s.respondError(w, "Invalid order parameters", http.StatusBadRequest)
		return
	}

	queueDepth := s.engine.GetQueueDepth()
	if s.monitor.ShouldThrottle(queueDepth) {
		s.respondError(w, "System under heavy load - order throttled", http.StatusServiceUnavailable)
		return
	}

	var side engine.Side
	switch strings.ToUpper(req.Side) {
	case "BUY":
		side = engine.BUY
	case "SELL":
		side = engine.SELL
	default:
		s.respondError(w, "Invalid side - must be BUY or SELL", http.StatusBadRequest)
		return
	}

	order := engine.NewOrder(req.Symbol, side, req.Price, req.Qty, req.UserID)

	select {
	case s.engine.GetOrderChan() <- order:
		s.logger.Info("Order received",
			"order_id", order.ID,
			"symbol", order.Symbol,
			"side", order.Side,
			"price", order.Price,
			"qty", order.Qty,
		)

		s.respondJSON(w, OrderResponse{
			Status:  "accepted",
			OrderID: order.ID.String(),
		}, http.StatusAccepted)
	default:
		s.respondError(w, "Order queue full", http.StatusServiceUnavailable)
	}
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	monitorStats := s.monitor.GetStats()
	mmStats := s.marketMaker.GetStats()

	response := map[string]interface{}{
		"monitor":         monitorStats,
		"market_maker":    mmStats,
		"queue_depth":     s.engine.GetQueueDepth(),
		"injection_count": s.selfHealer.GetInjectionCount(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleOrderBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		s.respondError(w, "Symbol required in URL path", http.StatusBadRequest)
		return
	}

	symbol := strings.ToUpper(parts[2])
	book := s.engine.GetBook(symbol)

	if book == nil {
		s.respondError(w, "Order book not found", http.StatusNotFound)
		return
	}

	snapshot := book.GetSnapshot(10)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshot)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("WebSocket upgrade failed", "error", err)
		return
	}

	client := &WebSocketClient{
		hub:  s.wsHub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	s.wsHub.register <- client

	go client.writePump()
	go client.readPump()

	s.logger.Info("WebSocket client connected", "remote", r.RemoteAddr)
}

func (s *Server) respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
func (s *Server) respondError(w http.ResponseWriter, message string, status int) {
	s.respondJSON(w, OrderResponse{
		Status:  "error",
		Message: message,
	}, status)
}
