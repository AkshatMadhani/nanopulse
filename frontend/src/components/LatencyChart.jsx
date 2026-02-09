import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { Activity } from 'lucide-react';
import { useState, useEffect } from 'react';

const LatencyChart = ({ currentLatency, maxLatency }) => {
  const [data, setData] = useState([]);
  const maxDataPoints = 30;

  useEffect(() => {
    if (currentLatency === undefined || currentLatency === null) {
      return;
    }

    setData(prevData => {
      const newData = [
        ...prevData,
        {
          time: new Date().toLocaleTimeString('en-IN', { 
            hour12: false,
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
          }),
          latency: Math.min(currentLatency, 1000), 
          max: Math.min(maxLatency || 0, 1000),
        }
      ].slice(-maxDataPoints);
      return newData;
    });
  }, [currentLatency, maxLatency]);

  const avgLatency = data.length > 0
    ? (data.reduce((sum, d) => sum + d.latency, 0) / data.length).toFixed(1)
    : 0;

  return (
    <div className="card p-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold flex items-center" style={{ color: 'var(--text-primary)' }}>
          <Activity className="w-5 h-5 mr-2 text-blue-500" />
          Latency Monitor
        </h2>
        <div className="flex items-center space-x-4 text-sm">
          <div>
            <span style={{ color: 'var(--text-secondary)' }}>Avg: </span>
            <span className="font-mono text-green-400">{avgLatency}µs</span>
          </div>
          <div>
            <span style={{ color: 'var(--text-secondary)' }}>Current: </span>
            <span className="font-mono text-blue-400">{currentLatency?.toFixed(1) || '0'}µs</span>
          </div>
          <div>
            <span style={{ color: 'var(--text-secondary)' }}>Threshold: </span>
            <span className="font-mono text-yellow-400">200µs</span>
          </div>
        </div>
      </div>

      <div className="h-64">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={data}>
            <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
            <XAxis 
              dataKey="time" 
              stroke="#94a3b8"
              tick={{ fontSize: 11 }}
              interval="preserveStartEnd"
            />
            <YAxis 
              stroke="#94a3b8"
              tick={{ fontSize: 11 }}
              domain={[0, 500]} 
              label={{ 
                value: 'Latency (µs)', 
                angle: -90, 
                position: 'insideLeft',
                style: { fill: '#94a3b8', fontSize: 12 }
              }}
            />
            <Tooltip 
              contentStyle={{
                backgroundColor: 'var(--bg-secondary)',
                border: '1px solid var(--border-color)',
                borderRadius: '8px',
                fontSize: '12px',
                color: 'var(--text-primary)'
              }}
              labelStyle={{ color: 'var(--text-primary)' }}
            />
            <Line 
              type="monotone" 
              dataKey="latency" 
              stroke="#10b981" 
              strokeWidth={2}
              dot={false}
              name="Current"
            />
            <Line 
              type="monotone" 
              dataKey="max" 
              stroke="#ef4444" 
              strokeWidth={1}
              strokeDasharray="5 5"
              dot={false}
              name="Max"
            />
          </LineChart>
        </ResponsiveContainer>
      </div>

      <div className="flex items-center justify-center space-x-6 mt-4 text-sm">
        <div className="flex items-center space-x-2">
          <div className="w-3 h-3 bg-green-500 rounded-full"></div>
          <span style={{ color: 'var(--text-secondary)' }}>Current Latency</span>
        </div>
        <div className="flex items-center space-x-2">
          <div className="w-3 h-0.5 bg-red-500"></div>
          <span style={{ color: 'var(--text-secondary)' }}>Max Latency</span>
        </div>
      </div>
    </div>
  );
};

export default LatencyChart;