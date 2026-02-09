import { TrendingUp, Clock, Layers, DollarSign, Activity, AlertTriangle } from 'lucide-react';
import { formatLatency, formatCurrency, formatNumber, getLatencyColor } from '../utils/Format';
const MetricsPanel = ({ systemState }) => {
  const {
    best_bid,
    best_ask,
    spread,
    latency_us,
    max_latency_us,
    queue_depth,
    mm_profit,
    total_trades,
  } = systemState;

  const metrics = [
    {
      icon: TrendingUp,
      label: 'Best Bid',
      value: best_bid ? formatCurrency(best_bid) : '-',
      color: 'text-green-400',
      bgColor: 'bg-green-900/20',
    },
    {
      icon: TrendingUp,
      label: 'Best Ask',
      value: best_ask ? formatCurrency(best_ask) : '-',
      color: 'text-red-400',
      bgColor: 'bg-red-900/20',
    },
    {
      icon: Activity,
      label: 'Spread',
      value: spread !== null ? formatCurrency(spread) : '-',
      color: 'text-blue-400',
      bgColor: 'bg-blue-900/20',
    },
    {
      icon: Clock,
      label: 'Avg Latency',
      value: formatLatency(latency_us),
      color: getLatencyColor(latency_us),
      bgColor: 'bg-slate-900/20',
    },
    {
      icon: AlertTriangle,
      label: 'Max Latency',
      value: formatLatency(max_latency_us),
      color: getLatencyColor(max_latency_us),
      bgColor: 'bg-slate-900/20',
    },
    {
      icon: Layers,
      label: 'Queue Depth',
      value: formatNumber(queue_depth, 0),
      color: queue_depth > 8000 ? 'text-red-400' : 'text-slate-300',
      bgColor: 'bg-slate-900/20',
    },
    {
      icon: DollarSign,
      label: 'MM Profit',
      value: formatCurrency(mm_profit),
      color: mm_profit > 0 ? 'text-green-400' : 'text-slate-300',
      bgColor: mm_profit > 0 ? 'bg-green-900/20' : 'bg-slate-900/20',
    },
    {
      icon: Activity,
      label: 'Total Trades',
      value: formatNumber(total_trades, 0),
      color: 'text-purple-400',
      bgColor: 'bg-purple-900/20',
    },
  ];

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
      {metrics.map((metric, index) => (
        <MetricCard key={index} {...metric} />
      ))}
    </div>
  );
};

const MetricCard = ({ icon: Icon, label, value, color, bgColor }) => {
  return (
    <div className={`metric-card ${bgColor}`}>
      <div className="flex items-center justify-between mb-2">
        <span className="metric-label">{label}</span>
        <Icon className={`w-4 h-4 ${color}`} />
      </div>
      <div className={`metric-value ${color} number-transition`}>
        {value}
      </div>
    </div>
  );
};

export default MetricsPanel;