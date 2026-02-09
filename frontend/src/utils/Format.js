export const formatNumber = (num, decimals = 2) => {
  if (num === null || num === undefined) return '-';
  return Number(num).toLocaleString('en-IN', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  });
};
export const formatCurrency = (num) => {
  if (num === null || num === undefined) return '₹-';
  return '₹' + formatNumber(num, 2);
};

export const formatLatency = (us) => {
  if (us === null || us === undefined) return '-';
  
  if (us < 1000) {
    return `${Math.round(us)}µs`;
  } else {
    return `${(us / 1000).toFixed(2)}ms`;
  }
};

export const getLatencyColor = (latencyUs, threshold = 200) => {
  if (latencyUs < threshold * 0.5) return 'text-green-400';
  if (latencyUs < threshold) return 'text-yellow-400';
  return 'text-red-400';
};

export const getModeBadgeClass = (mode) => {
  switch (mode) {
    case 'NORMAL':
      return 'status-normal';
    case 'SAFE':
      return 'status-safe';
    case 'THROTTLED':
      return 'status-throttled';
    default:
      return 'status-badge bg-gray-900/30 text-gray-400 border border-gray-700';
  }
};

export const getConnectionIndicator = (status) => {
  switch (status) {
    case 'connected':
      return { class: 'pulse-green', text: 'Connected' };
    case 'connecting':
      return { class: 'pulse-yellow', text: 'Connecting...' };
    case 'disconnected':
    case 'error':
    case 'failed':
      return { class: 'pulse-red', text: 'Disconnected' };
    default:
      return { class: 'pulse-indicator bg-gray-500', text: 'Unknown' };
  }
};

export const formatTimestamp = (timestamp) => {
  if (!timestamp) return '-';
  
  const date = new Date(timestamp / 1000000);
  return date.toLocaleTimeString('en-IN', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  });
};

export const calculatePercentage = (value, total) => {
  if (!total) return 0;
  return ((value / total) * 100).toFixed(1);
};

export const debounce = (func, wait) => {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
};

export const generateOrderId = () => {
  return `ORDER-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
};