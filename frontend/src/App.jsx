import { Socket } from './hooks/socket';
import Header from './components/Header';
import MetricsPanel from './components/Metrics';
import OrderBook from './components/Order';
import Trades from './components/Trades';
import LatencyChart from './components/LatencyChart';
import Form from './components/Form';
import { AlertCircle, TrendingUp, Info } from 'lucide-react';
import { useState, useMemo, useEffect } from 'react';

function App() {
  const { connectionStatus, systemState, recentTrades, error } = Socket();
  const [selectedSymbol, setSelectedSymbol] = useState('RELIANCE');

  const handleOrderSubmitted = (response) => {
    console.log('Order submitted:', response);
  };

  const selectedOrderBook = useMemo(() => {
    if (systemState.order_books && Object.keys(systemState.order_books).length > 0) {
      const book = systemState.order_books[selectedSymbol];
      if (book) {
        console.log(`Order book for ${selectedSymbol}:`, book);
        return book;
      }
    }
    
    if (systemState.order_book && systemState.order_book.symbol === selectedSymbol) {
      console.log(`Using primary order book for ${selectedSymbol}`);
      return systemState.order_book;
    }
    
    console.log(`No order book data for ${selectedSymbol}`);
    return null;
  }, [systemState.order_books, systemState.order_book, selectedSymbol]);

  const symbolStats = useMemo(() => {
    const book = selectedOrderBook;
    if (book) {
      return {
        best_bid: book.best_bid,
        best_ask: book.best_ask,
        spread: book.spread,
      };
    }
    if (systemState.order_book?.symbol === selectedSymbol) {
      return {
        best_bid: systemState.best_bid,
        best_ask: systemState.best_ask,
        spread: systemState.spread,
      };
    }
    return {
      best_bid: null,
      best_ask: null,
      spread: null,
    };
  }, [selectedOrderBook, systemState, selectedSymbol]);

  const symbols = ['RELIANCE', 'TCS', 'INFY', 'HDFC', 'ICICI'];
  useEffect(() => {
    console.log('System State:', {
      order_books: systemState.order_books,
      selectedSymbol,
      selectedOrderBook,
      total_trades: systemState.total_trades
    });
  }, [systemState, selectedSymbol, selectedOrderBook]);

  return (
    <div className="min-h-screen p-4 md:p-6" style={{ backgroundColor: 'var(--bg-primary)' }}>
      <div className="max-w-[1600px] mx-auto">
        <Header 
          connectionStatus={connectionStatus} 
          systemMode={systemState.mode}
        />

        {error && (
          <div className="mb-6 bg-red-900/20 border border-red-700 rounded-lg p-4 flex items-start space-x-3 animate-slide-up">
            <AlertCircle className="w-6 h-6 text-red-400 flex-shrink-0 mt-0.5" />
            <div>
              <div className="font-semibold text-red-400 mb-1">Connection Error</div>
              <div className="text-sm text-red-300">{error}</div>
              <div className="text-xs text-red-400 mt-2">
                Make sure the backend is running on http://localhost:8080
              </div>
            </div>
          </div>
        )}
        <div className="mb-6">
          <div className="card p-4 md:p-6">
            <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
              <div>
                <div className="text-sm font-semibold mb-1" style={{ color: 'var(--text-primary)' }}>
                  Select Trading Symbol
                </div>
                <div className="text-xs" style={{ color: 'var(--text-muted)' }}>
                  Choose which stock to trade
                </div>
              </div>
              <div className="flex gap-2 flex-wrap">
                {symbols.map(symbol => (
                  <button
                    key={symbol}
                    onClick={() => {
                      console.log(`Switching to ${symbol}`);
                      setSelectedSymbol(symbol);
                    }}
                    className={`px-4 md:px-6 py-2 md:py-3 rounded-lg font-semibold transition-all duration-200 flex items-center gap-2 text-sm ${
                      selectedSymbol === symbol
                        ? 'bg-gradient-to-r from-blue-600 to-blue-700 text-white shadow-lg shadow-blue-900/50 scale-105'
                        : ''
                    }`}
                    style={{
                      backgroundColor: selectedSymbol !== symbol ? 'var(--bg-tertiary)' : undefined,
                      color: selectedSymbol !== symbol ? 'var(--text-secondary)' : undefined,
                      border: selectedSymbol !== symbol ? '1px solid var(--border-color)' : undefined,
                    }}
                  >
                    {selectedSymbol === symbol && <TrendingUp className="w-4 h-4" />}
                    {symbol}
                  </button>
                ))}
              </div>
            </div>

            {symbolStats.best_bid !== null && symbolStats.best_ask !== null && (
              <div className="mt-4 pt-4 grid grid-cols-3 gap-4 text-sm"
                style={{ borderTop: `1px solid var(--border-color)` }}>
                <div className="text-center p-2 rounded-lg" style={{ backgroundColor: 'rgba(34, 197, 94, 0.1)' }}>
                  <div className="text-xs mb-1" style={{ color: 'var(--text-secondary)' }}>Best Bid</div>
                  <div className="font-mono font-bold text-green-400">
                    ₹{symbolStats.best_bid.toFixed(2)}
                  </div>
                </div>
                <div className="text-center p-2 rounded-lg" style={{ backgroundColor: 'var(--bg-tertiary)' }}>
                  <div className="text-xs mb-1" style={{ color: 'var(--text-secondary)' }}>Spread</div>
                  <div className="font-mono font-bold text-blue-400">
                    ₹{symbolStats.spread?.toFixed(2) || '0.00'}
                  </div>
                </div>
                <div className="text-center p-2 rounded-lg" style={{ backgroundColor: 'rgba(239, 68, 68, 0.1)' }}>
                  <div className="text-xs mb-1" style={{ color: 'var(--text-secondary)' }}>Best Ask</div>
                  <div className="font-mono font-bold text-red-400">
                    ₹{symbolStats.best_ask.toFixed(2)}
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
        {systemState.total_trades === 0 && (
          <div className="mb-6 p-4 rounded-lg flex items-start gap-3"
            style={{ 
              backgroundColor: 'rgba(59, 130, 246, 0.1)',
              border: '1px solid rgba(59, 130, 246, 0.3)'
            }}>
            <Info className="w-5 h-5 text-blue-400 flex-shrink-0 mt-0.5" />
            <div className="flex-1">
              <div className="font-semibold text-blue-400 text-sm mb-1">
                How to Execute Your First Trade
              </div>
              <div className="text-xs" style={{ color: 'var(--text-secondary)' }}>
                1. Look at the <strong>Order Book</strong> below for current prices<br/>
                2. To BUY: Enter the <strong className="text-red-400">Best Ask</strong> price (₹{symbolStats.best_ask?.toFixed(2) || '2502.00'})<br/>
                3. To SELL: Enter the <strong className="text-green-400">Best Bid</strong> price (₹{symbolStats.best_bid?.toFixed(2) || '2500.00'})<br/>
                4. Click "Place Order" and watch it execute instantly! ⚡
              </div>
            </div>
          </div>
        )}

        <MetricsPanel systemState={systemState} />
        <div className="grid grid-cols-1 xl:grid-cols-3 gap-6 mb-6">
          <div className="xl:col-span-2">
            <OrderBook orderBook={selectedOrderBook} />
          </div>
          <div className="xl:col-span-1">
            <Form 
              onOrderSubmitted={handleOrderSubmitted}
              selectedSymbol={selectedSymbol}
            />
          </div>
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
          <LatencyChart 
            currentLatency={systemState.latency_us}
            maxLatency={systemState.max_latency_us}
          />

          <Trades trades={recentTrades} />
        </div>

        <footer className="mt-6 text-center text-xs" style={{ color: 'var(--text-muted)' }}>
          <div className="card p-3">
            <div className="flex items-center justify-center flex-wrap gap-3">
              <span className="flex items-center gap-1">
                ⚡NanoPulse Core
              </span>
              <span>•</span>
              <span className="flex items-center gap-1">
                Trading: <span className="text-blue-500 font-semibold">{selectedSymbol}</span>
              </span>
              <span>•</span>
              <span className="flex items-center gap-1">
                {connectionStatus === 'connected' ? '🟢' : '🔴'} {connectionStatus}
              </span>
              <span>•</span>
              <span className="flex items-center gap-1">
                Queue: <span className="font-mono">{systemState.queue_depth}</span>
              </span>
              <span>•</span>
              <span className="flex items-center gap-1">
                Trades: <span className="font-mono">{systemState.total_trades}</span>
              </span>
            </div>
          </div>
        </footer>
      </div>
    </div>
  );
}

export default App;