package api

import (
	"net/http"

	"github.com/AkshatMadhani/nanopulse/engine"
	"github.com/AkshatMadhani/nanopulse/logger"
	"github.com/AkshatMadhani/nanopulse/market"
	"github.com/AkshatMadhani/nanopulse/monitor"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	engine      *engine.MatchingEngine
	monitor     *monitor.Monitor
	marketMaker *market.Bot
	selfHealer  *monitor.SelfHealer
	logger      *logger.Logger
	wsHub       *WebSocketHub
	tradeChan   <-chan *engine.Trade
	tradeBuffer *TradeBuffer
}

func NewServer(
	eng *engine.MatchingEngine,
	mon *monitor.Monitor,
	mm *market.Bot,
	sh *monitor.SelfHealer,
	log *logger.Logger,
) *Server {
	hub := NewWebSocketHub(log)

	return &Server{
		engine:      eng,
		monitor:     mon,
		marketMaker: mm,
		selfHealer:  sh,
		logger:      log,
		wsHub:       hub,
		tradeChan:   eng.GetTradeChan(),
		tradeBuffer: NewTradeBuffer(),
	}
}

func (s *Server) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/order", s.handleOrder)
	mux.HandleFunc("/stats", s.handleStats)
	mux.HandleFunc("/book/", s.handleOrderBook)
	mux.HandleFunc("/ws", s.handleWebSocket)

	return mux
}

func (s *Server) Start(port string) error {
	s.logger.Info("Starting API server", "port", port)

	go s.wsHub.Run()
	go s.startTradeListener()
	go s.broadcastSystemState()

	mux := s.SetupRoutes()
	return http.ListenAndServe(":"+port, s.corsMiddleware(mux))
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
