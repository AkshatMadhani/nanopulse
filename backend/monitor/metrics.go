package monitor

import (
	"sync"
	"time"

	"github.com/AkshatMadhani/nanopulse/engine"
	"github.com/AkshatMadhani/nanopulse/logger"
)

type SystemMode int

const (
	NORMAL SystemMode = iota
	SAFE
	THROTTLED
)

func (m SystemMode) String() string {
	switch m {
	case NORMAL:
		return "NORMAL"
	case SAFE:
		return "SAFE"
	case THROTTLED:
		return "THROTTLED"
	default:
		return "UNKNOWN"
	}
}

type Monitor struct {
	latencyWindow    []float64
	latencyWindowMu  sync.RWMutex
	avgLatency       float64
	maxLatency       float64
	queueDepth       int
	currentMode      SystemMode
	modeMu           sync.RWMutex
	tradeChan        <-chan *engine.Trade
	metricsChan      <-chan engine.Metric
	logger           *logger.Logger
	config           Config
	totalTrades      int64
	safeModeTriggers int64
	throttleCount    int64
	mu               sync.RWMutex
}

type Config struct {
	LatencyThresholdUs float64
	QueueThreshold     int
	WindowSize         int
	CheckInterval      time.Duration
	SpreadThreshold    float64
}

func DefaultConfig() Config {
	return Config{
		LatencyThresholdUs: 200.0,
		QueueThreshold:     8000,
		WindowSize:         100,
		CheckInterval:      time.Second * 2,
		SpreadThreshold:    0.1,
	}
}

func NewMonitor(tradeChan <-chan *engine.Trade, metricsChan <-chan engine.Metric, log *logger.Logger, cfg Config) *Monitor {
	return &Monitor{
		latencyWindow: make([]float64, 0, cfg.WindowSize),
		tradeChan:     tradeChan,
		metricsChan:   metricsChan,
		logger:        log,
		config:        cfg,
		currentMode:   NORMAL,
	}
}

func (m *Monitor) Start() {
	m.logger.Info("Starting system monitor")
	go m.collectMetrics()
	go m.collectTrades()
	go m.checkHealth()
}

func (m *Monitor) collectMetrics() {
	for metric := range m.metricsChan {
		switch metric.Type {
		case "latency":
			m.recordLatency(metric.Value)
		}
	}
}
func (m *Monitor) collectTrades() {
	for trade := range m.tradeChan {
		if trade != nil {
			m.IncrementTrades()
			m.logger.Debug("Monitor tracked trade",
				"id", trade.ID.String(),
				"symbol", trade.Symbol,
				"price", trade.Price,
				"qty", trade.Qty,
				"total_trades", m.GetTotalTrades(),
			)
		}
	}
}

func (m *Monitor) recordLatency(latencyUs float64) {
	m.latencyWindowMu.Lock()
	defer m.latencyWindowMu.Unlock()

	m.latencyWindow = append(m.latencyWindow, latencyUs)
	if len(m.latencyWindow) > m.config.WindowSize {
		m.latencyWindow = m.latencyWindow[1:]
	}

	if latencyUs > m.maxLatency {
		m.maxLatency = latencyUs
	}

	sum := 0.0
	for _, l := range m.latencyWindow {
		sum += l
	}
	m.avgLatency = sum / float64(len(m.latencyWindow))
}

func (m *Monitor) checkHealth() {
	ticker := time.NewTicker(m.config.CheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		m.evaluateHealth()
	}
}

func (m *Monitor) evaluateHealth() {
	m.latencyWindowMu.RLock()
	avgLatency := m.avgLatency
	m.latencyWindowMu.RUnlock()

	m.modeMu.Lock()
	defer m.modeMu.Unlock()

	currentMode := m.currentMode

	if avgLatency > m.config.LatencyThresholdUs && currentMode == NORMAL {
		m.currentMode = SAFE
		m.safeModeTriggers++
		m.logger.Warn("Entering SAFE mode due to high latency",
			"avg_latency_us", avgLatency,
			"threshold_us", m.config.LatencyThresholdUs,
		)
	}

	if avgLatency < m.config.LatencyThresholdUs*0.7 && currentMode == SAFE {
		m.currentMode = NORMAL
		m.logger.Info("Returning to NORMAL mode", "avg_latency_us", avgLatency)
	}

	m.logger.Debug("Health check",
		"mode", m.currentMode,
		"avg_latency_us", avgLatency,
		"max_latency_us", m.maxLatency,
		"total_trades", m.totalTrades,
	)
}

func (m *Monitor) GetMode() SystemMode {
	m.modeMu.RLock()
	defer m.modeMu.RUnlock()
	return m.currentMode
}

func (m *Monitor) GetStats() Stats {
	m.latencyWindowMu.RLock()
	avgLatency := m.avgLatency
	maxLatency := m.maxLatency
	m.latencyWindowMu.RUnlock()

	m.mu.RLock()
	defer m.mu.RUnlock()

	return Stats{
		AvgLatencyUs:     avgLatency,
		MaxLatencyUs:     maxLatency,
		TotalTrades:      m.totalTrades,
		SafeModeTriggers: m.safeModeTriggers,
		ThrottleCount:    m.throttleCount,
		CurrentMode:      m.GetMode(),
	}
}

type Stats struct {
	AvgLatencyUs     float64    `json:"avg_latency_us"`
	MaxLatencyUs     float64    `json:"max_latency_us"`
	TotalTrades      int64      `json:"total_trades"`
	SafeModeTriggers int64      `json:"safe_mode_triggers"`
	ThrottleCount    int64      `json:"throttle_count"`
	CurrentMode      SystemMode `json:"current_mode"`
}

func (m *Monitor) IncrementTrades() {
	m.mu.Lock()
	m.totalTrades++
	m.mu.Unlock()
}

func (m *Monitor) GetTotalTrades() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.totalTrades
}

func (m *Monitor) ShouldThrottle(queueDepth int) bool {
	if queueDepth > m.config.QueueThreshold {
		m.mu.Lock()
		m.throttleCount++
		m.mu.Unlock()
		return true
	}
	return false
}
