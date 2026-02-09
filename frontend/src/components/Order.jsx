import { TrendingUp, TrendingDown, BarChart3 } from 'lucide-react';
import { formatCurrency, formatNumber } from '../utils/Format';
import { useState, useEffect } from 'react';

const OrderBook = ({ orderBook, recentUserOrders = [] }) => {
  const [flashBuy, setFlashBuy] = useState(false);
  const [flashSell, setFlashSell] = useState(false);

  useEffect(() => {
    if (orderBook?.buy_book?.length) {
      setFlashBuy(true);
      const timer = setTimeout(() => setFlashBuy(false), 600);
      return () => clearTimeout(timer);
    }
  }, [orderBook?.buy_book]);

  useEffect(() => {
    if (orderBook?.sell_book?.length) {
      setFlashSell(true);
      const timer = setTimeout(() => setFlashSell(false), 600);
      return () => clearTimeout(timer);
    }
  }, [orderBook?.sell_book]);

  if (!orderBook) {
    return (
      <div className="card p-6">
        <h2 className="text-lg font-semibold mb-4 flex items-center">
          <BarChart3 className="w-5 h-5 mr-2 text-blue-500" />
          Order Book
        </h2>
        <div className="text-center py-12" style={{ color: 'var(--text-muted)' }}>
          <BarChart3 className="w-16 h-16 mx-auto mb-3 opacity-20" />
          <p>No order book data available</p>
          <p className="text-sm mt-2">Waiting for market data...</p>
        </div>
      </div>
    );
  }

  const buyBook = orderBook.buy_book || [];
  const sellBook = orderBook.sell_book || [];
  const symbol = orderBook.symbol || 'RELIANCE';

  const totalBuyVolume = buyBook.reduce((sum, order) => sum + order.qty, 0);
  const totalSellVolume = sellBook.reduce((sum, order) => sum + order.qty, 0);
  const maxVolume = Math.max(totalBuyVolume, totalSellVolume) || 1;

  return (
    <div className="card p-6">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-lg font-semibold flex items-center" style={{ color: 'var(--text-primary)' }}>
          <BarChart3 className="w-5 h-5 mr-2 text-blue-500" />
          Order Book
        </h2>
        <div className="flex items-center gap-3">
          <div className="px-3 py-1 rounded-full text-sm font-bold bg-gradient-to-r from-blue-600 to-blue-700 text-white">
            {symbol}
          </div>
          <div className="text-xs px-3 py-1 rounded-full" 
            style={{ 
              backgroundColor: 'var(--bg-tertiary)',
              color: 'var(--text-secondary)',
              border: '1px solid var(--border-color)'
            }}>
            {buyBook.length + sellBook.length} levels
          </div>
        </div>
      </div>

      {recentUserOrders.length > 0 && (
        <div className="mb-4 p-3 rounded-lg" 
          style={{ 
            backgroundColor: 'var(--bg-tertiary)',
            border: '1px solid var(--border-color)' 
          }}>
          <div className="text-xs font-semibold mb-2 flex items-center gap-1" 
            style={{ color: 'var(--text-primary)' }}>
            <span className="text-blue-500">📝</span> Your Recent Orders
          </div>
          {recentUserOrders.slice(0, 3).map((order, idx) => (
            <div key={idx} className="text-xs font-mono" style={{ color: 'var(--text-secondary)' }}>
              {order.side} {order.qty} @ ₹{order.price}
            </div>
          ))}
        </div>
      )}

      <div className="grid grid-cols-2 gap-4">
        <div>
          <div className="flex items-center justify-between mb-3 pb-3" 
            style={{ borderBottom: '2px solid rgba(34, 197, 94, 0.2)' }}>
            <div className="flex items-center gap-2 text-green-400">
              <TrendingUp className="w-5 h-5" />
              <span className="font-bold text-sm">BIDS</span>
            </div>
            <div className="text-xs px-2 py-1 rounded" 
              style={{ 
                backgroundColor: 'rgba(34, 197, 94, 0.1)',
                color: 'var(--text-secondary)' 
              }}>
              Vol: {formatNumber(totalBuyVolume, 0)}
            </div>
          </div>

          <div className="space-y-1.5">
            {buyBook.length > 0 ? (
              buyBook.slice(0, 10).map((order, idx) => (
                <OrderRow
                  key={idx}
                  price={order.price}
                  qty={order.qty}
                  type="buy"
                  percentage={(order.qty / maxVolume) * 100}
                  flash={flashBuy && idx === 0}
                />
              ))
            ) : (
              <div className="text-center text-sm py-8" style={{ color: 'var(--text-muted)' }}>
                No buy orders
              </div>
            )}
          </div>
        </div>

        <div>
          <div className="flex items-center justify-between mb-3 pb-3" 
            style={{ borderBottom: '2px solid rgba(239, 68, 68, 0.2)' }}>
            <div className="flex items-center gap-2 text-red-400">
              <TrendingDown className="w-5 h-5" />
              <span className="font-bold text-sm">ASKS</span>
            </div>
            <div className="text-xs px-2 py-1 rounded" 
              style={{ 
                backgroundColor: 'rgba(239, 68, 68, 0.1)',
                color: 'var(--text-secondary)' 
              }}>
              Vol: {formatNumber(totalSellVolume, 0)}
            </div>
          </div>

          <div className="space-y-1.5">
            {sellBook.length > 0 ? (
              sellBook.slice(0, 10).map((order, idx) => (
                <OrderRow
                  key={idx}
                  price={order.price}
                  qty={order.qty}
                  type="sell"
                  percentage={(order.qty / maxVolume) * 100}
                  flash={flashSell && idx === 0}
                />
              ))
            ) : (
              <div className="text-center text-sm py-8" style={{ color: 'var(--text-muted)' }}>
                No sell orders
              </div>
            )}
          </div>
        </div>
      </div>

      {orderBook.best_bid && orderBook.best_ask && (
        <div className="mt-6 pt-5" style={{ borderTop: '1px solid var(--border-color)' }}>
          <div className="grid grid-cols-3 gap-4 text-sm">
            <div className="text-center p-3 rounded-lg" 
              style={{ backgroundColor: 'rgba(34, 197, 94, 0.1)' }}>
              <div className="text-xs mb-1" style={{ color: 'var(--text-secondary)' }}>
                Best Bid
              </div>
              <div className="font-mono font-bold text-lg text-green-400">
                {formatCurrency(orderBook.best_bid)}
              </div>
            </div>
            
            <div className="text-center p-3 rounded-lg" 
              style={{ backgroundColor: 'var(--bg-tertiary)' }}>
              <div className="text-xs mb-1" style={{ color: 'var(--text-secondary)' }}>
                Spread
              </div>
              <div className="font-mono font-bold text-lg text-blue-400">
                {formatCurrency(orderBook.spread || 0)}
              </div>
            </div>

            <div className="text-center p-3 rounded-lg" 
              style={{ backgroundColor: 'rgba(239, 68, 68, 0.1)' }}>
              <div className="text-xs mb-1" style={{ color: 'var(--text-secondary)' }}>
                Best Ask
              </div>
              <div className="font-mono font-bold text-lg text-red-400">
                {formatCurrency(orderBook.best_ask)}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

const OrderRow = ({ price, qty, type, percentage, flash }) => {
  const isBuy = type === 'buy';
  const bgColor = isBuy ? 'rgba(34, 197, 94, 0.1)' : 'rgba(239, 68, 68, 0.1)';
  const textColor = isBuy ? 'text-green-400' : 'text-red-400';
  const flashClass = flash ? (isBuy ? 'flash-buy' : 'flash-sell') : '';

  return (
    <div
      className={`relative overflow-hidden rounded-lg px-3 py-2 ${flashClass} ${
        isBuy ? 'order-row-buy' : 'order-row-sell'
      }`}
      style={{ backgroundColor: 'var(--bg-tertiary)' }}
    >
      <div
        className="absolute inset-0 transition-all duration-300"
        style={{ 
          backgroundColor: bgColor,
          width: `${percentage}%` 
        }}
      />

      <div className="relative flex items-center justify-between text-sm">
        <span className={`font-mono font-bold ${textColor}`}>
          {formatCurrency(price)}
        </span>
        <span className="font-mono font-semibold" style={{ color: 'var(--text-primary)' }}>
          {formatNumber(qty, 0)}
        </span>
      </div>
    </div>
  );
};

export default OrderBook;