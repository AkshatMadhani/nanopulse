import { ArrowRightLeft, TrendingUp, Activity, Zap } from 'lucide-react';
import { formatCurrency, formatNumber, formatTimestamp } from '../utils/Format';

const Trades = ({ trades }) => {
  const totalVolume = trades.reduce((sum, trade) => sum + (trade.qty || 0), 0);
  const totalValue = trades.reduce((sum, trade) => sum + (trade.price * trade.qty || 0), 0);
  const avgPrice = trades.length > 0 ? totalValue / totalVolume : 0;

  return (
    <div className="card p-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold flex items-center" style={{ color: 'var(--text-primary)' }}>
          <ArrowRightLeft className="w-5 h-5 mr-2 text-blue-500" />
          Recent Trades
        </h2>
        <div className="flex items-center gap-2">
          <div className="px-3 py-1 rounded-full text-xs font-semibold"
            style={{ 
              backgroundColor: 'var(--bg-tertiary)',
              color: 'var(--text-secondary)',
              border: '1px solid var(--border-color)'
            }}>
            {trades.length} trades
          </div>
        </div>
      </div>

      {trades.length > 0 && (
        <div className="grid grid-cols-3 gap-3 mb-4">
          <div className="p-3 rounded-lg" style={{ backgroundColor: 'var(--bg-tertiary)' }}>
            <div className="text-xs mb-1" style={{ color: 'var(--text-secondary)' }}>
              Total Volume
            </div>
            <div className="font-mono font-bold" style={{ color: 'var(--text-primary)' }}>
              {formatNumber(totalVolume, 0)}
            </div>
          </div>
          
          <div className="p-3 rounded-lg" style={{ backgroundColor: 'var(--bg-tertiary)' }}>
            <div className="text-xs mb-1" style={{ color: 'var(--text-secondary)' }}>
              Total Value
            </div>
            <div className="font-mono font-bold text-blue-400">
              {formatCurrency(totalValue)}
            </div>
          </div>

          <div className="p-3 rounded-lg" style={{ backgroundColor: 'var(--bg-tertiary)' }}>
            <div className="text-xs mb-1" style={{ color: 'var(--text-secondary)' }}>
              Avg Price
            </div>
            <div className="font-mono font-bold text-purple-400">
              {formatCurrency(avgPrice)}
            </div>
          </div>
        </div>
      )}

      <div className="space-y-2 max-h-[400px] overflow-y-auto">
        {trades.length > 0 ? (
          trades.map((trade, idx) => (
            <TradeRow key={idx} trade={trade} isNew={idx === 0} />
          ))
        ) : (
          <div className="text-center py-12" style={{ color: 'var(--text-muted)' }}>
            <Activity className="w-16 h-16 mx-auto mb-3 opacity-20" />
            <p className="font-semibold">No trades yet</p>
            <p className="text-xs mt-2">Waiting for orders to match...</p>
            <div className="mt-4 p-3 rounded-lg inline-block" 
              style={{ backgroundColor: 'var(--bg-tertiary)' }}>
              <p className="text-xs flex items-center gap-2">
                <Zap className="w-4 h-4 text-blue-500" />
                <span>Place an order to see trades execute</span>
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

const TradeRow = ({ trade, isNew }) => {
  const {
    id,
    symbol = 'RELIANCE',
    price,
    qty,
    timestamp,
    side = 'BUY',
    buyer_id,
    seller_id,
  } = trade;

  const isBuy = side === 'BUY';
  const tradeValue = price * qty;
  const animationClass = isNew ? 'trade-flash' : '';

  const isMakerTrade = buyer_id === 'market-maker' || seller_id === 'market-maker';
  const isUserTrade = buyer_id === 'web-ui' || seller_id === 'web-ui';

  return (
    <div
      className={`rounded-xl p-4 ${animationClass} hover:scale-[1.01] transition-all duration-200`}
      style={{
        backgroundColor: 'var(--bg-tertiary)',
        border: `2px solid ${isBuy ? 'rgba(34, 197, 94, 0.3)' : 'rgba(239, 68, 68, 0.3)'}`,
      }}
    >
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          <span className="font-bold" style={{ color: 'var(--text-primary)' }}>
            {symbol}
          </span>
          
          <span className={`text-xs px-2 py-1 rounded-full font-bold ${
            isBuy ? 'bg-green-500/20 text-green-400' : 'bg-red-500/20 text-red-400'
          }`}>
            {side}
          </span>

          {isUserTrade && (
            <span className="text-xs px-2 py-1 rounded-full bg-blue-500/20 text-blue-400 font-semibold">
              YOUR TRADE
            </span>
          )}
          {isMakerTrade && (
            <span className="text-xs px-2 py-1 rounded-full font-semibold"
              style={{ 
                backgroundColor: 'rgba(168, 85, 247, 0.2)',
                color: '#a855f7' 
              }}>
              MARKET MAKER
            </span>
          )}
        </div>

        <span className="text-xs" style={{ color: 'var(--text-muted)' }}>
          {formatTimestamp(timestamp)}
        </span>
      </div>

      <div className="grid grid-cols-3 gap-4">
        <div>
          <div className="text-xs mb-1" style={{ color: 'var(--text-secondary)' }}>
            Price
          </div>
          <div className={`font-mono font-bold text-lg ${
            isBuy ? 'text-green-400' : 'text-red-400'
          }`}>
            {formatCurrency(price)}
          </div>
        </div>

        <div>
          <div className="text-xs mb-1" style={{ color: 'var(--text-secondary)' }}>
            Quantity
          </div>
          <div className="font-mono font-bold text-lg" style={{ color: 'var(--text-primary)' }}>
            {formatNumber(qty, 0)}
          </div>
        </div>

        <div>
          <div className="text-xs mb-1" style={{ color: 'var(--text-secondary)' }}>
            Value
          </div>
          <div className="font-mono font-bold text-lg text-blue-400">
            {formatCurrency(tradeValue)}
          </div>
        </div>
      </div>

      {/* Participants (if available) */}
      {(buyer_id || seller_id) && (
        <div className="mt-3 pt-3 grid grid-cols-2 gap-3 text-xs"
          style={{ borderTop: '1px solid var(--border-color)' }}>
          {buyer_id && (
            <div>
              <span style={{ color: 'var(--text-muted)' }}>Buyer: </span>
              <span className="font-mono" style={{ color: 'var(--text-secondary)' }}>
                {buyer_id}
              </span>
            </div>
          )}
          {seller_id && (
            <div>
              <span style={{ color: 'var(--text-muted)' }}>Seller: </span>
              <span className="font-mono" style={{ color: 'var(--text-secondary)' }}>
                {seller_id}
              </span>
            </div>
          )}
        </div>
      )}

      {id && (
        <div className="mt-2 pt-2" style={{ borderTop: '1px solid var(--border-color)' }}>
          <div className="text-xs font-mono truncate" style={{ color: 'var(--text-muted)' }}>
            ID: {id}
          </div>
        </div>
      )}
    </div>
  );
};

export default Trades;
