package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/AkshatMadhani/nanopulse/api"
	"github.com/AkshatMadhani/nanopulse/engine"
	"github.com/AkshatMadhani/nanopulse/logger"
	"github.com/AkshatMadhani/nanopulse/market"
	"github.com/AkshatMadhani/nanopulse/monitor"
	"github.com/AkshatMadhani/nanopulse/simulator"
)

type TradeBroadcaster struct {
	input  <-chan *engine.Trade
	output []chan *engine.Trade
	logger *logger.Logger
}

func NewTradeBroadcaster(input <-chan *engine.Trade, count int, log *logger.Logger) *TradeBroadcaster {
	tb := &TradeBroadcaster{
		input:  input,
		output: make([]chan *engine.Trade, count),
		logger: log,
	}

	for i := 0; i < count; i++ {
		tb.output[i] = make(chan *engine.Trade, 100)
	}

	return tb
}

func (tb *TradeBroadcaster) Start() {
	go func() {
		for trade := range tb.input {
			if trade != nil {
				for i, ch := range tb.output {
					select {
					case ch <- trade:
					default:
						tb.logger.Warn("Trade channel full, dropping trade", "channel", i)
					}
				}
			}
		}
		for _, ch := range tb.output {
			close(ch)
		}
	}()
}

func (tb *TradeBroadcaster) GetChannel(index int) <-chan *engine.Trade {
	if index < 0 || index >= len(tb.output) {
		return nil
	}
	return tb.output[index]
}

func main() {
	port := flag.String("port", "8080", "API server port")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	enableSimulator := flag.Bool("simulator", false, "Enable market simulator")
	simRate := flag.Int("sim-rate", 10, "Simulator orders per second")
	flag.Parse()

	level := logger.INFO
	switch *logLevel {
	case "debug":
		level = logger.DEBUG
	case "warn":
		level = logger.WARN
	case "error":
		level = logger.ERROR
	}
	log := logger.New(level)

	log.Info(" Starting NanoPulse Core")
	log.Info("Configuration",
		"port", *port,
		"log_level", *logLevel,
		"simulator_enabled", *enableSimulator,
	)

	matchingEngine := engine.NewMatchingEngine(10000, log)
	matchingEngine.Start()

	tradeBroadcaster := NewTradeBroadcaster(matchingEngine.GetTradeChan(), 3, log)
	tradeBroadcaster.Start()

	monitorConfig := monitor.DefaultConfig()
	systemMonitor := monitor.NewMonitor(
		tradeBroadcaster.GetChannel(0),
		matchingEngine.GetMetricsChan(),
		log,
		monitorConfig,
	)
	systemMonitor.Start()

	selfHealer := monitor.NewSelfHealer(systemMonitor, matchingEngine, log)
	selfHealer.Start()

	marketMaker := market.NewBot(
		matchingEngine,
		tradeBroadcaster.GetChannel(1),
		log,
	)
	marketMaker.Start()

	if *enableSimulator {
		sim := simulator.NewSimulator(matchingEngine, log, *simRate)
		sim.Start()
		log.Info("Market simulator started", "rate", *simRate)
	}

	apiServer := api.NewServer(
		matchingEngine,
		systemMonitor,
		marketMaker,
		selfHealer,
		log,
	)

	go func() {
		log.Info("API server listening", "port", *port)
		if err := apiServer.Start(*port); err != nil {
			log.Error("API server failed", "error", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Info("Final stats", "monitor", systemMonitor.GetStats())
}
