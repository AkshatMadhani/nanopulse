import { useEffect, useState, useCallback, useRef } from 'react';
import wsService from '../service/socket';

export const Socket = () => {
  const [connectionStatus, setConnectionStatus] = useState('disconnected');
  const [systemState, setSystemState] = useState({
    timestamp: null,
    best_bid: null,
    best_ask: null,
    spread: null,
    latency_us: 0,
    max_latency_us: 0,
    mode: 'NORMAL',
    queue_depth: 0,
    mm_profit: 0,
    total_trades: 0,
    order_books: {},
    order_book: null, 
  });
  const [recentTrades, setRecentTrades] = useState([]);
  const [error, setError] = useState(null);
  
  const lastUpdateRef = useRef(Date.now());

  const handleMessage = useCallback((message) => {
    if (message.type === 'connection') {
      setConnectionStatus(message.status);
      
      if (message.status === 'error' || message.status === 'failed') {
        setError(message.message || 'Connection error');
      } else {
        setError(null);
      }
    } else if (message.type === 'data') {
      const now = Date.now();
      
      if (now - lastUpdateRef.current < 1000) {
        return;
      }
      
      lastUpdateRef.current = now;
      
      setSystemState(prevState => ({
        ...prevState,
        ...message.data,
        order_books: message.data.order_books || prevState.order_books || {},
      }));

      if (message.data.recent_trade) {
        setRecentTrades(prev => {
          const newTrades = [message.data.recent_trade, ...prev].slice(0, 20);
          return newTrades;
        });
      }
    }
  }, []);

  useEffect(() => {
    wsService.connect();

    const unsubscribe = wsService.subscribe(handleMessage);

    return () => {
      unsubscribe();
      wsService.disconnect();
    };
  }, [handleMessage]);

  return {
    connectionStatus,
    systemState,
    recentTrades,
    error,
  };
};
