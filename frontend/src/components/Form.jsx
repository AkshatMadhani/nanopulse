import { useState, useEffect } from 'react';
import { Send, TrendingUp, TrendingDown, AlertCircle, Zap } from 'lucide-react';
import apiService from '../service/api'

const Form = ({ onOrderSubmitted, selectedSymbol = 'RELIANCE' }) => {
  const [formData, setFormData] = useState({
    symbol: selectedSymbol,
    side: 'BUY',
    price: '',
    qty: '',
    user_id: 'web-ui',
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  useEffect(() => {
    setFormData(prev => ({ ...prev, symbol: selectedSymbol }));
  }, [selectedSymbol]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    setError(null);
    setSuccess(null);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.price || !formData.qty) {
      setError('Price and quantity are required');
      return;
    }

    if (parseFloat(formData.price) <= 0) {
      setError('Price must be greater than 0');
      return;
    }

    if (parseInt(formData.qty) <= 0) {
      setError('Quantity must be greater than 0');
      return;
    }

    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      const order = {
        ...formData,
        price: parseFloat(formData.price),
        qty: parseInt(formData.qty),
      };

      console.log('Submitting order:', order);
      const response = await apiService.placeOrder(order);
      
      console.log('Order response:', response);
      
      const successMsg = `${formData.side} ${formData.qty} ${formData.symbol} @ ₹${formData.price}`;
      setSuccess(successMsg);
      
      setFormData(prev => ({ ...prev, price: '', qty: '' }));
      
      if (onOrderSubmitted) {
        onOrderSubmitted(response);
      }

      setTimeout(() => setSuccess(null), 5000);
    } catch (err) {
      console.error('❌Order error:', err);
      setError(err.message || 'Failed to place order');
    } finally {
      setLoading(false);
    }
  };

  const isBuy = formData.side === 'BUY';

  return (
    <div className="card p-6">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-lg font-semibold flex items-center" style={{ color: 'var(--text-primary)' }}>
          <Zap className="w-5 h-5 mr-2 text-blue-500" />
          Quick Trade
        </h2>
        <div className="px-3 py-1 rounded-full text-xs font-semibold"
          style={{ 
            backgroundColor: 'var(--bg-tertiary)',
            color: 'var(--text-secondary)',
            border: '1px solid var(--border-color)'
          }}>
          {selectedSymbol}
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-5">
        <div>
          <label className="block text-xs font-semibold mb-3 uppercase tracking-wider" 
            style={{ color: 'var(--text-secondary)' }}>
            Order Type
          </label>
          <div className="grid grid-cols-2 gap-3">
            <button
              type="button"
              onClick={() => setFormData(prev => ({ ...prev, side: 'BUY' }))}
              className={`py-3 rounded-xl font-semibold transition-all duration-200 flex items-center justify-center gap-2 ${
                isBuy
                  ? 'bg-gradient-to-r from-green-600 to-green-700 text-white shadow-lg shadow-green-900/50 scale-105'
                  : ''
              }`}
              style={{
                backgroundColor: !isBuy ? 'var(--bg-tertiary)' : undefined,
                color: !isBuy ? 'var(--text-secondary)' : undefined,
                border: !isBuy ? '1px solid var(--border-color)' : undefined,
              }}
            >
              <TrendingUp className="w-4 h-4" />
              BUY
            </button>
            <button
              type="button"
              onClick={() => setFormData(prev => ({ ...prev, side: 'SELL' }))}
              className={`py-3 rounded-xl font-semibold transition-all duration-200 flex items-center justify-center gap-2 ${
                !isBuy
                  ? 'bg-gradient-to-r from-red-600 to-red-700 text-white shadow-lg shadow-red-900/50 scale-105'
                  : ''
              }`}
              style={{
                backgroundColor: isBuy ? 'var(--bg-tertiary)' : undefined,
                color: isBuy ? 'var(--text-secondary)' : undefined,
                border: isBuy ? '1px solid var(--border-color)' : undefined,
              }}
            >
              <TrendingDown className="w-4 h-4" />
              SELL
            </button>
          </div>
        </div>

        <div>
          <label className="block text-xs font-semibold mb-2 uppercase tracking-wider" 
            style={{ color: 'var(--text-secondary)' }}>
            Price (₹)
          </label>
          <div className="flex items-center gap-2">
            <div className="flex-shrink-0 w-8 h-12 flex items-center justify-center rounded-lg font-bold text-lg"
              style={{ 
                backgroundColor: 'var(--bg-tertiary)',
                color: 'var(--text-muted)',
                border: '1px solid var(--border-color)'
              }}>
              ₹
            </div>
            <input
              type="number"
              name="price"
              value={formData.price}
              onChange={handleChange}
              placeholder="2501.00"
              step="0.01"
              min="0"
              className="flex-1 rounded-xl px-4 py-3 font-mono text-lg font-semibold focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
              style={{
                backgroundColor: 'var(--bg-tertiary)',
                border: '1px solid var(--border-color)',
                color: 'var(--text-primary)'
              }}
            />
          </div>
        </div>

        <div>
          <label className="block text-xs font-semibold mb-2 uppercase tracking-wider" 
            style={{ color: 'var(--text-secondary)' }}>
            Quantity
          </label>
          <input
            type="number"
            name="qty"
            value={formData.qty}
            onChange={handleChange}
            placeholder="10"
            min="1"
            step="1"
            className="w-full rounded-xl px-4 py-3 font-mono text-lg font-semibold focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
            style={{
              backgroundColor: 'var(--bg-tertiary)',
              border: '1px solid var(--border-color)',
              color: 'var(--text-primary)'
            }}
          />
        </div>

        {formData.price && formData.qty && (
          <div className="p-3 rounded-lg" style={{ backgroundColor: 'var(--bg-tertiary)' }}>
            <div className="flex items-center justify-between text-sm">
              <span style={{ color: 'var(--text-secondary)' }}>Order Value</span>
              <span className="font-mono font-bold text-lg" style={{ color: 'var(--text-primary)' }}>
                ₹{(parseFloat(formData.price) * parseInt(formData.qty)).toFixed(2)}
              </span>
            </div>
          </div>
        )}

        <button
          type="submit"
          disabled={loading}
          className={`w-full py-4 rounded-xl font-bold text-lg transition-all duration-200 flex items-center justify-center gap-2 ${
            isBuy
              ? 'bg-gradient-to-r from-green-600 to-green-700 hover:from-green-700 hover:to-green-800'
              : 'bg-gradient-to-r from-red-600 to-red-700 hover:from-red-700 hover:to-red-800'
          } text-white ${loading ? 'opacity-50 cursor-not-allowed' : 'shadow-lg hover:shadow-xl hover:scale-[1.02]'}`}
        >
          {loading ? (
            <>
              <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
              Submitting...
            </>
          ) : (
            <>
              <Send className="w-5 h-5" />
              Place {formData.side} Order
            </>
          )}
        </button>

        {success && (
          <div className="bg-green-900/30 border-2 border-green-500 rounded-xl p-4 animate-slide-up">
            <div className="flex items-start space-x-3">
              <div className="w-8 h-8 bg-green-500 rounded-full flex items-center justify-center flex-shrink-0">
                <svg className="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                </svg>
              </div>
              <div className="flex-1">
                <div className="font-bold text-green-400 text-lg mb-1">Order Placed! 🎉</div>
                <div className="text-sm text-green-300 font-mono">{success}</div>
                <div className="text-xs text-green-400 mt-2 flex items-center gap-1">
                  <span>✓</span> Check trades below for execution
                </div>
              </div>
            </div>
          </div>
        )}

        {error && (
          <div className="bg-red-900/20 border-2 border-red-700 rounded-xl p-4 flex items-start space-x-3 animate-slide-up">
            <AlertCircle className="w-5 h-5 text-red-400 flex-shrink-0 mt-0.5" />
            <div>
              <div className="font-semibold text-red-400">Error</div>
              <div className="text-sm text-red-300 mt-1">{error}</div>
            </div>
          </div>
        )}
      </form>

      <div className="mt-5 p-4 rounded-xl" 
        style={{ 
          backgroundColor: 'var(--bg-tertiary)',
          border: '1px solid var(--border-color)' 
        }}>
        <div className="text-xs flex items-start gap-2" style={{ color: 'var(--text-secondary)' }}>
          <span className="text-blue-500">💡</span>
          <div>
            <span className="font-semibold" style={{ color: 'var(--text-primary)' }}>Tip:</span> Orders execute instantly when price matches. Check order book for current prices.
          </div>
        </div>
      </div>
    </div>
  );
};

export default Form;
